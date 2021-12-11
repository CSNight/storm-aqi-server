package elastic

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	pool "github.com/jolestar/go-commons-pool/v2"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
	"time"
)

type EsAPI struct {
	Log         *zap.SugaredLogger
	EsPool      *pool.ObjectPool
	Bulk        BulkIndexer
	sucTotal    uint64
	sucPre      uint64
	errTotal    int64
	globalCli   *elasticsearch.Client
	ctxCli      context.Context
	ctxBulk     context.Context
	FailQueue   []BulkIndexerItem
	lock        sync.RWMutex
	checkTicker *time.Ticker
	isReachable bool
}

type BulkItem struct {
	Index      string
	Action     string
	DocumentID string
	Body       []byte
}

func (t *EsAPI) Init() {
	t.ctxCli = context.Background()
	t.ctxBulk = context.Background()
	_ = t.initClient()
	go t.checkConn()
}

func (t *EsAPI) Close() {
	t.checkTicker.Stop()
	err := t.Bulk.Close(t.ctxBulk)
	if err != nil {
		t.Log.Errorf("EsAPI bulker close \u001B[31merr: %v\u001B[0m", err)
	}
	t.lock.Lock()
	defer t.lock.Unlock()
	t.FailQueue = []BulkIndexerItem{}
}

func (t *EsAPI) GetClient(ctx context.Context) (*elasticsearch.Client, error) {
	cli, err := t.EsPool.BorrowObject(ctx)
	if err != nil {
		return nil, err
	}
	return cli.(*elasticsearch.Client), nil
}

func (t *EsAPI) CloseClient(ctx context.Context, cli *elasticsearch.Client) error {
	err := t.EsPool.ReturnObject(ctx, cli)
	if err != nil {
		return err
	}
	return nil
}

func (t *EsAPI) initClient() error {
	t.isReachable = false
	cli, err := t.EsPool.BorrowObject(t.ctxCli)
	if err != nil {
		t.globalCli = nil
		return err
	}
	t.globalCli = cli.(*elasticsearch.Client)
	t.errTotal = 0
	t.sucTotal = 0
	t.sucPre = 0
	t.Bulk, err = t.NewBulkProcessor()
	if err != nil {
		return err
	}
	t.isReachable = true
	t.Log.Infof("Elasticsearch client connect successfully!")
	return nil
}

func (t *EsAPI) isConnected() bool {
	if t.globalCli == nil {
		return false
	}
	res, err := t.globalCli.Ping()
	if err != nil || res.IsError() {
		t.Log.Errorf("Elasticsearch global client is disconnected \u001B[31merr: %v\u001B[0m", err)
		return false
	}
	return true
}

func (t *EsAPI) checkConn() {
	t.checkTicker = time.NewTicker(time.Second * 3)
	for {
		select {
		case <-t.checkTicker.C:
			isCon := t.isConnected()
			if !isCon {
				if t.Bulk != nil {
					err := t.Bulk.Close(t.ctxBulk)
					if err != nil {
						t.Log.Errorf("EsAPI bulker close \u001B[31merr: %v\u001B[0m", err)
					}
					t.Bulk = nil
				}
				_ = t.initClient()
			} else if isCon && t.isReachable {
				t.lock.Lock()
				count := 0
				cF := make([]BulkIndexerItem, len(t.FailQueue))
				copy(cF, t.FailQueue)
				t.FailQueue = t.FailQueue[:0]
				t.lock.Unlock()
				for _, item := range cF {
					_ = t.AddToBulk(context.Background(), item)
					count++
				}
				atomic.AddInt64(&t.errTotal, -int64(len(cF)))
				cF = cF[:0]
				cF = nil
				if count > 0 {
					t.Log.Infof("retry failure bulk request count: %d", count)
				}
			}
		}
	}
}
