package elastic

import (
	"context"
	"errors"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/tidwall/gjson"
	"io"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

func (t *EsAPI) NewBulkProcessor() (BulkIndexer, error) {
	bulkProcessor, err := NewBulkIndexer(BulkIndexerConfig{
		Client:        t.globalCli,     // The Elasticsearch client
		NumWorkers:    8,               // The number of worker goroutines
		FlushInterval: 5 * time.Second, // The periodic flush interval
		Timeout:       time.Second * 60,
		OnFlushStart: func(ctx context.Context, count int) context.Context {
			stats := t.Bulk.Stats()
			t.Log.Debugf("Executing bulk \u001B[36m[%d]\u001B[0m requests \u001B[36m%d\u001B[0m", stats.NumRequests+1, count)
			ctx = context.WithValue(ctx, "startTm", time.Now())
			ctx = context.WithValue(ctx, "seq", stats.NumRequests+1)
			return ctx
		},
		OnError: func(ctx context.Context, items []BulkIndexerItem, err error) {
			if items != nil && len(items) > 0 {
				t.AddToFailure(items...)
				atomic.AddInt64(&t.errTotal, int64(len(items)))
			}
			stats := t.Bulk.Stats()
			t.Log.Errorf("\u001B[31mExecuting bulk %d \u001B[31merr: %v\u001B[0m failed count: %d\u001B[0m", stats.NumRequests, err, t.errTotal)
		},
		OnFlushEnd: func(ctx context.Context) {
			cost := time.Since(ctx.Value("startTm").(time.Time)).Milliseconds()
			seq := ctx.Value("seq").(uint64)
			success := t.sucTotal - t.sucPre
			t.Log.Debugf("Bulk \u001B[36m[%d]\u001B[0m \u001B[32mcompleted\u001B[0m in \u001B[36m%dms\u001B[0m requests %d", seq, cost, success)
			t.sucPre = t.sucTotal
		},
	})
	if err != nil {
		return nil, err
	}
	return bulkProcessor, err
}

func (t *EsAPI) AddToFailure(item ...BulkIndexerItem) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.FailQueue = append(t.FailQueue, item...)
}

func (t *EsAPI) AddToBulk(ctx context.Context, req BulkIndexerItem) error {
	if !t.isReachable {
		return errors.New("elasticsearch is not reachable")
	}
	req.OnFailure = func(ctx context.Context, item BulkIndexerItem, res BulkIndexerResponseItem, err error) {
		t.AddToFailure(item)
	}
	req.OnSuccess = func(ctx context.Context, item BulkIndexerItem, resp BulkIndexerResponseItem) {
		atomic.AddUint64(&t.sucTotal, 1)
	}
	err := t.Bulk.Add(ctx, req)
	if err != nil {
		t.AddToFailure(req)
		t.Log.Errorf("Add to bulker failure \u001B[31merr: %v\u001B[0m", err)
		return err
	}
	return nil
}

func (t *EsAPI) ScrollSearch(req *esapi.SearchRequest) []gjson.Result {
	ctx := context.Background()
	cli, err := t.GetClient(ctx)
	if err != nil {
		t.Log.Errorf("ScrollSearch(). GetClient(). \u001B[31merr: %v\u001B[0m", err)
		return nil
	}
	defer func() {
		err = t.CloseClient(ctx, cli)
		if err != nil {
			t.Log.Errorf("ScrollSearch(). CloseClient(). \u001B[31merr: %v\u001B[0m", err)
			return
		}
	}()
	respBytes, err := ProcessResp(req, cli)
	if err != nil {
		t.Log.Errorf("ScrollSearch(). ProcessResp(). \u001B[31merr: %v\u001B[0m", err)
		return nil
	}
	root := gjson.ParseBytes(respBytes)
	rootHits := root.Get("hits")
	var results []gjson.Result
	for {
		getSources(rootHits, &results)
		if root.Get("_scroll_id").Exists() {
			scroll := esapi.ScrollRequest{
				ScrollID: root.Get("_scroll_id").String(),
			}
			respBytes, err = ProcessResp(scroll, cli)
			if err != nil {
				t.Log.Errorf("ScrollSearch(). ProcessResp(). \u001B[31merr: %v\u001B[0m", err)
				return results
			}
			root = gjson.ParseBytes(respBytes)
			rootHits = root.Get("hits")
		} else {
			break
		}
	}
	return results
}

func (t *EsAPI) CreateIndex(index string, mappings string, args string) bool {
	ctx := context.Background()
	cli, err := t.GetClient(ctx)
	if err != nil {
		t.Log.Errorf("CreateIndex(). GetClient(). \u001B[31merr: %v\u001B[0m", err)
		return true
	}
	defer func() {
		err = t.CloseClient(ctx, cli)
		if err != nil {
			t.Log.Errorf("CreateIndex(). CloseClient(). \u001B[31merr: %v\u001B[0m", err)
			return
		}
	}()
	indices := index
	if strings.Contains(index, "$") && args != "" {
		indices = strings.Split(index, "$")[0] + args
	} else if strings.Contains(index, "$") && args == "" {
		return false
	}
	if t.ExistIndex(indices) {
		return true
	}
	request := esapi.IndicesCreateRequest{
		Index:         indices,
		Body:          strings.NewReader(mappings),
		MasterTimeout: 0,
		Timeout:       0,
		Pretty:        true,
		ErrorTrace:    true,
	}
	_, err = ProcessResp(request, cli)
	if err != nil {
		t.Log.Errorf("CreateIndex(). \u001B[31merr: %v\u001B[0m", err)
		return false
	}
	t.Log.Infof("CreateIndex(). success create indices %s with mappings: %s", indices, mappings)
	return true
}

func (t *EsAPI) ExistIndex(index string) bool {
	ctx := context.Background()
	cli, err := t.GetClient(ctx)
	if err != nil {
		t.Log.Errorf("ExistIndex(). GetClient(). \u001B[31merr: %v\u001B[0m", err)
		return true
	}
	defer func() {
		err = t.CloseClient(ctx, cli)
		if err != nil {
			t.Log.Errorf("ExistIndex(). CloseClient(). \u001B[31merr: %v\u001B[0m", err)
			return
		}
	}()
	request := esapi.IndicesExistsRequest{
		Index:      []string{index},
		Pretty:     true,
		Human:      true,
		ErrorTrace: true,
	}
	_, err = ProcessResp(request, cli)
	if err != nil {
		if err.Error() == "404" {
			return false
		}
		return true
	}
	return true
}

func ProcessResp(req esapi.Request, cli *elasticsearch.Client) ([]byte, error) {
	resp, err := req.Do(context.Background(), cli)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(strconv.Itoa(resp.StatusCode))
	}
	if resp.Body != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBytes, nil
}

func getSources(rootHits gjson.Result, results *[]gjson.Result) {
	if rootHits.Get("total.value").Exists() && rootHits.Get("total.value").Int() > 0 {
		hits := rootHits.Get("hits").Array()
		for _, hit := range hits {
			*results = append(*results, hit.Get("_source"))
		}
	}
}
