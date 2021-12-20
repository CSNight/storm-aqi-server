package db

import (
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"strings"
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
	Tz       string         `json:"tz"`
	Tm       int64          `json:"tm"`
	Tms      string         `json:"tms"`
}

type RealtimeItem struct {
	EsSearchItem
	Source AqiRealtime `json:"_source"`
}

type RealtimeSearchResponse struct {
	EsSearchRespMeta
	Hits struct {
		Total    EsRespTotal    `json:"total"`
		MaxScore float64        `json:"max_score"`
		Hits     []RealtimeItem `json:"hits"`
	} `json:"hits"`
}

func (db *DB) GetAqiRealtimeById(sid string) (*RealtimeResp, error) {
	st, err := db.getStationFromCache(sid)
	if err != nil {
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
		return nil, err
	}
	err = json.Unmarshal(resp, &esSearchResp)
	if err != nil {
		return nil, err
	}
	var rts []RealtimeInfo
	if esSearchResp.Hits.Total.Value > 0 {
		for _, item := range esSearchResp.Hits.Hits {
			info := RealtimeInfo{
				Pol:  item.Source.Pol,
				Data: item.Source.Data,
			}
			rts = append(rts, info)
		}
		response.Realtime = rts
		response.Tz = esSearchResp.Hits.Hits[0].Source.Tz
		response.Tm = esSearchResp.Hits.Hits[0].Source.Tm
		response.Tms = esSearchResp.Hits.Hits[0].Source.Tms
		return response, nil
	}
	return response, nil
}
