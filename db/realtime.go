package db

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"sync"
	"time"
)

type AqiRealtime struct {
	Idx      int     `json:"idx"`
	Sid      string  `json:"sid"`
	Pol      string  `json:"pol"`
	Data     float64 `json:"data"`
	Daily    string  `json:"daily"`
	Forecast string  `json:"forecast"`
	Tz       string  `json:"tz"`
	Tm       int64   `json:"tm"`
	Tms      string  `json:"tms"`
}

type ForecastItem struct {
	Avg float64 `json:"avg"`
	Day string  `json:"day"`
	Max float64 `json:"max"`
	Min float64 `json:"min"`
}

type RealtimeInfo struct {
	Pol  string  `json:"pol"`
	Data float64 `json:"data"`
}

type ForecastInfo struct {
	Daily map[string][]ForecastItem `json:"daily"`
}

type RealtimeGetResponse struct {
	EsGetResponse
	Source AqiRealtime `json:"_source"`
}

type ForecastResp struct {
	Idx      int                       `json:"idx"`
	Sid      string                    `json:"sid"`
	Name     string                    `json:"name"`
	Loc      GeoPoint                  `json:"loc"`
	CityName string                    `json:"city_name"`
	Forecast map[string][]ForecastItem `json:"forecast"`
	Tz       string                    `json:"tz"`
	Tm       int64                     `json:"tm"`
	Tms      string                    `json:"tms"`
}

type RealtimeResp struct {
	Idx      int            `json:"idx"`
	Sid      string         `json:"sid"`
	Name     string         `json:"name"`
	Loc      GeoPoint       `json:"loc"`
	CityName string         `json:"city_name"`
	Realtime []RealtimeInfo `json:"realtime"`
	MainPol  string         `json:"main_pol"`
	Tz       string         `json:"tz"`
	Tm       int64          `json:"tm"`
	Tms      string         `json:"tms"`
}

type RealtimeStMap struct {
	RealTimeMap map[string]float64
}

type RealtimeItem struct {
	EsSearchItem
	Source AqiRealtime `json:"_source"`
}

type BucketItem struct {
	Key      string `json:"key"`
	DocCount int    `json:"doc_count"`
	Data     struct {
		Value float64 `json:"value"`
	} `json:"data"`
}

type RealtimeAggResponse struct {
	EsSearchRespMeta
	Hits struct {
		Total    EsRespTotal    `json:"total"`
		MaxScore float64        `json:"max_score"`
		Hits     []RealtimeItem `json:"hits"`
	} `json:"hits"`
	Aggregations struct {
		Buckets struct {
			Buckets []BucketItem `json:"buckets"`
		} `json:"buckets"`
	} `json:"aggregations"`
}

type RealtimeSearchResponse struct {
	EsSearchRespMeta
	Hits struct {
		Total    EsRespTotal    `json:"total"`
		MaxScore float64        `json:"max_score"`
		Hits     []RealtimeItem `json:"hits"`
	} `json:"hits"`
}

func (db *DB) GetAllAqiRealtime() (*RealtimeStMap, error) {
	response := &RealtimeStMap{
		RealTimeMap: map[string]float64{},
	}
	var wg sync.WaitGroup
	resCh := make(chan BucketItem)
	for _, s := range []int{0, 8000} {
		wg.Add(1)
		go func(st int) {
			realResp, err := db.getHalfRealtimeStation(st, st+8000)
			if err != nil {
				return
			}
			for _, item := range realResp.Aggregations.Buckets.Buckets {
				resCh <- item
			}
			wg.Done()
		}(s)
	}
	go func() {
		for item := range resCh {
			response.RealTimeMap[item.Key] = item.Data.Value
		}
	}()
	wg.Wait()
	close(resCh)
	fmt.Println(len(resCh))
	return response, nil
}

func (db *DB) GetAqiRealtimeById(sid string) (*RealtimeResp, error) {
	st, err := db.getStationFromCache(sid)
	if err != nil || st == nil {
		return nil, err
	}
	response := &RealtimeResp{
		Idx:      st.Idx,
		Sid:      st.Sid,
		Name:     st.Name,
		Loc:      st.Loc,
		CityName: st.CityName,
	}

	size := 10
	query := `{
        "query": {
            "match": {
                "sid": ` + sid + `
            }
       }
    }`
	search := &esapi.SearchRequest{
		Index:          []string{db.Conf.RealtimeIndex},
		Body:           strings.NewReader(query),
		Size:           &size,
		SourceExcludes: []string{"forecast", "daily"},
		Timeout:        20 * time.Second,
	}
	resp, err := db.api.ProcessRespWithCli(search)
	var esSearchResp RealtimeSearchResponse
	if err != nil {
		if strings.HasPrefix(err.Error(), "404") {
			return response, nil
		}
		db.log.Error("GetAqiRealtimeById(). es.ProcessRespWithCli(). err:", zap.String("query", strings.ReplaceAll(query, " ", "")), zap.Error(err))
		return nil, err
	}
	defer func() {
		resp = nil
	}()
	err = json.Unmarshal(resp, &esSearchResp)
	if err != nil {
		db.log.Error("GetAqiRealtimeById(). json.Unmarshal(). err:", zap.Error(err))
		return nil, err
	}
	var rts []RealtimeInfo
	maxVal := -1.0
	mainPol := ""
	if esSearchResp.Hits.Total.Value > 0 {
		for _, item := range esSearchResp.Hits.Hits {
			info := RealtimeInfo{
				Pol:  item.Source.Pol,
				Data: item.Source.Data,
			}
			if item.Source.Data > maxVal {
				mainPol = item.Source.Pol
			}
			rts = append(rts, info)
		}
		response.Realtime = rts
		response.Tz = esSearchResp.Hits.Hits[0].Source.Tz
		response.Tm = esSearchResp.Hits.Hits[0].Source.Tm
		response.Tms = esSearchResp.Hits.Hits[0].Source.Tms
		response.MainPol = mainPol
		return response, nil
	}
	return response, nil
}

func (db *DB) GetAqiRealtimeByIdAndPol(sid string, pol string) (*RealtimeResp, error) {
	st, err := db.getStationFromCache(sid)
	if err != nil || st == nil {
		return nil, err
	}
	infoResp := &RealtimeResp{
		Idx:      st.Idx,
		Sid:      st.Sid,
		Name:     st.Name,
		Loc:      st.Loc,
		CityName: st.CityName,
		Realtime: []RealtimeInfo{},
	}
	search := &esapi.GetRequest{
		Index:          db.Conf.RealtimeIndex,
		DocumentID:     "rt_" + sid + "$" + pol,
		SourceExcludes: []string{"forecast", "daily"},
	}
	resp, err := db.api.ProcessRespWithCli(search)
	defer func() {
		resp = nil
	}()
	if err != nil {
		if strings.HasPrefix(err.Error(), "404") {
			return infoResp, nil
		}
		db.log.Error("GetAqiRealtimeByIdAndPol(). es.ProcessRespWithCli(). err:", zap.Error(err))
		return nil, err
	}
	var response RealtimeGetResponse
	err = json.Unmarshal(resp, &response)
	if err != nil {
		db.log.Error("GetAqiRealtimeByIdAndPol(). json.Unmarshal(). err:", zap.Error(err))
		return nil, err
	}
	if response.Found {
		info := RealtimeInfo{
			Pol:  response.Source.Pol,
			Data: response.Source.Data,
		}
		infoResp.Realtime = []RealtimeInfo{info}
		infoResp.Tz = response.Source.Tz
		infoResp.Tm = response.Source.Tm
		infoResp.Tms = response.Source.Tms
	}
	return infoResp, nil
}

func (db *DB) GetForecast(sid string, pol string) (*ForecastResp, error) {
	st, err := db.getStationFromCache(sid)
	if err != nil || st == nil {
		return nil, err
	}
	response := &ForecastResp{
		Idx:      st.Idx,
		Sid:      st.Sid,
		Name:     st.Name,
		Loc:      st.Loc,
		CityName: st.CityName,
	}

	size := 10
	query := `{
        "query": {
            "match": {
                "sid": ` + sid + `
            }
       }
    }`
	search := &esapi.SearchRequest{
		Index:          []string{db.Conf.RealtimeIndex},
		Body:           strings.NewReader(query),
		Size:           &size,
		SourceExcludes: []string{"data", "pol", "daily"},
		Timeout:        20 * time.Second,
	}
	resp, err := db.api.ProcessRespWithCli(search)
	var esSearchResp RealtimeSearchResponse
	if err != nil {
		if strings.HasPrefix(err.Error(), "404") {
			return response, nil
		}
		db.log.Error("GetForecast(). es.ProcessRespWithCli(). err:", zap.String("query", strings.ReplaceAll(query, " ", "")), zap.Error(err))
		return nil, err
	}
	defer func() {
		resp = nil
	}()
	err = json.Unmarshal(resp, &esSearchResp)
	if err != nil {
		db.log.Error("GetForecast(). json.Unmarshal(). err:", zap.Error(err))
		return nil, err
	}

	if esSearchResp.Hits.Total.Value > 0 {
		forecastStr := esSearchResp.Hits.Hits[0].Source.Forecast
		var forecastSource ForecastInfo
		err = json.Unmarshal([]byte(forecastStr), &forecastSource)
		if err != nil {
			db.log.Error("GetForecast(). json.Unmarshal(). err:", zap.Error(err))
			return nil, err
		}
		if pol == "all" {
			response.Forecast = forecastSource.Daily
		} else {
			response.Forecast = map[string][]ForecastItem{pol: forecastSource.Daily[pol]}
		}
		response.Tz = esSearchResp.Hits.Hits[0].Source.Tz
		response.Tm = esSearchResp.Hits.Hits[0].Source.Tm
		response.Tms = esSearchResp.Hits.Hits[0].Source.Tms
		return response, nil
	}
	return response, nil
}

func (db *DB) getHalfRealtimeStation(from int, to int) (*RealtimeAggResponse, error) {
	query := `{
        "query": {
            "bool": {
                "must": {
                    "range": {
                        "idx": {
                            "lt": ` + strconv.Itoa(to) + `,
                            "gte": ` + strconv.Itoa(from) + `
                        }
                    }
                }
            }
        },
        "aggs": {
            "buckets": {
                "terms": {
                    "field": "sid",
                    "order": {
                        "data": "desc"
                    },
                    "size": 20000
                },
                "aggs": {
                    "data": {
                        "max": {
                            "field": "data"
                        }
                    }
                }
            }
        }
    }`
	size := 0
	search := &esapi.SearchRequest{
		Index:          []string{db.Conf.RealtimeIndex},
		Body:           strings.NewReader(query),
		Size:           &size,
		SourceExcludes: []string{"forecast", "daily"},
		Timeout:        20 * time.Second,
	}
	resp, err := db.api.ProcessRespWithCli(search)
	if err != nil {
		return nil, err
	}
	var respEs RealtimeAggResponse
	err = json.Unmarshal(resp, &respEs)
	if err != nil {
		return nil, err
	}
	return &respEs, nil
}
