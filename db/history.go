package db

import (
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"strconv"
	"strings"
	"time"
)

type AqiHistory struct {
	Idx   int     `json:"idx"`
	Sid   string  `json:"sid"`
	Pol   string  `json:"pol"`
	Name  string  `json:"name"`
	Data  float64 `json:"data"`
	Tz    string  `json:"tz"`
	Month int     `json:"month"`
	Tm    int64   `json:"tm"`
	Tms   string  `json:"tms"`
}

type AqiHisItem struct {
	Pol   string  `json:"pol"`
	Name  string  `json:"name"`
	Data  float64 `json:"data"`
	Tz    string  `json:"tz"`
	Month int     `json:"month"`
	Year  int     `json:"year"`
	Tm    int64   `json:"tm"`
	Tms   string  `json:"tms"`
}

type AqiHistoryResp struct {
	Idx      int                     `json:"idx"`
	Sid      string                  `json:"sid"`
	Name     string                  `json:"name"`
	Loc      GeoPoint                `json:"loc"`
	CityName string                  `json:"city_name"`
	History  map[string][]AqiHisItem `json:"history"`
}

type HistoryItem struct {
	EsSearchItem
	Source AqiHistory `json:"_source"`
}

type HistorySearchResponse struct {
	EsSearchRespMeta
	Hits struct {
		Total    EsRespTotal   `json:"total"`
		MaxScore float64       `json:"max_score"`
		Hits     []HistoryItem `json:"hits"`
	} `json:"hits"`
}

var pols = []string{"no2", "pm25", "pm10", "o3", "so2", "co"}

func (db *DB) GetHistoryYesterday(sid string, pol string) (*AqiHistoryResp, error) {
	station, err := db.getStationFromCache(sid)
	if err != nil || station == nil {
		return nil, err
	}
	now := time.Now().UTC()
	st := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	et := st.AddDate(0, 0, -1)
	hisList, err := db.getHistoryByRange(sid, pol, st, et)
	if err != nil {
		return nil, err
	}
	items := BuildResp(hisList)
	return &AqiHistoryResp{
		Idx:      station.Idx,
		Sid:      station.Sid,
		Name:     station.Name,
		Loc:      station.Loc,
		CityName: station.CityName,
		History:  items,
	}, nil
}

func (db *DB) GetHistoryLastWeek(sid string, pol string) (*AqiHistoryResp, error) {
	station, err := db.getStationFromCache(sid)
	if err != nil || station == nil {
		return nil, err
	}
	now := time.Now().UTC()
	st := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	et := st.AddDate(0, 0, -7)
	hisList, err := db.getHistoryByRange(sid, pol, st, et)
	if err != nil {
		return nil, err
	}
	items := BuildResp(hisList)
	return &AqiHistoryResp{
		Idx:      station.Idx,
		Sid:      station.Sid,
		Name:     station.Name,
		Loc:      station.Loc,
		CityName: station.CityName,
		History:  items,
	}, nil
}

func (db *DB) GetHistoryLastMonth(sid string, pol string) (*AqiHistoryResp, error) {
	station, err := db.getStationFromCache(sid)
	if err != nil || station == nil {
		return nil, err
	}
	now := time.Now().UTC()
	st := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	et := st.AddDate(0, -1, 0)
	hisList, err := db.getHistoryByRange(sid, pol, st, et)
	if err != nil {
		return nil, err
	}
	items := BuildResp(hisList)
	return &AqiHistoryResp{
		Idx:      station.Idx,
		Sid:      station.Sid,
		Name:     station.Name,
		Loc:      station.Loc,
		CityName: station.CityName,
		History:  items,
	}, nil
}

func (db *DB) GetHistoryLastSeason(sid string, pol string) (*AqiHistoryResp, error) {
	station, err := db.getStationFromCache(sid)
	if err != nil || station == nil {
		return nil, err
	}
	now := time.Now().UTC()
	st := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	et := st.AddDate(0, -3, 0)
	hisList, err := db.getHistoryByRange(sid, pol, st, et)
	if err != nil {
		return nil, err
	}
	items := BuildResp(hisList)
	return &AqiHistoryResp{
		Idx:      station.Idx,
		Sid:      station.Sid,
		Name:     station.Name,
		Loc:      station.Loc,
		CityName: station.CityName,
		History:  items,
	}, nil
}

func (db *DB) GetHistoryYear(sid string, pol string) (*AqiHistoryResp, error) {
	station, err := db.getStationFromCache(sid)
	if err != nil || station == nil {
		return nil, err
	}
	now := time.Now().UTC()
	st := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	et := st.AddDate(-1, 0, 0)
	hisList, err := db.getHistoryByRange(sid, pol, st, et)
	if err != nil {
		return nil, err
	}
	items := BuildResp(hisList)
	return &AqiHistoryResp{
		Idx:      station.Idx,
		Sid:      station.Sid,
		Name:     station.Name,
		Loc:      station.Loc,
		CityName: station.CityName,
		History:  items,
	}, nil
}

func BuildResp(list []AqiHistory) map[string][]AqiHisItem {
	var items map[string][]AqiHisItem
	for _, pol := range pols {
		items[pol] = []AqiHisItem{}
	}
	for _, item := range list {
		year := time.Unix(item.Tm/1000, 0).UTC().Year()
		items[item.Pol] = append(items[item.Pol], AqiHisItem{
			Pol:   item.Pol,
			Name:  item.Name,
			Data:  item.Data,
			Tz:    item.Tz,
			Month: item.Month,
			Year:  year,
			Tm:    item.Tm,
			Tms:   item.Tms,
		})
	}
	list = nil
	return items
}

func (db *DB) getHistoryByRange(sid string, pol string, st time.Time, et time.Time) ([]AqiHistory, error) {
	stYear := st.Year()
	etYear := et.Year()
	var indexes []string
	for i := stYear; i <= etYear; i++ {
		indexes = append(indexes, strings.Replace(db.Conf.HisIndex, "$year", strconv.Itoa(i), -1))
	}
	query := `{
        "query": {
            "bool": {
                "must": [
                    {"match": {"sid": "` + sid + `"}},
                    {"range": {"tm":{"gte": ` + strconv.Itoa(int(st.UnixMilli())) + `,"lte": ` + strconv.Itoa(int(et.UnixMilli())) + `}}}`
	if pol != "all" {
		query += `{"match": {"pol": "` + pol + `"}}`
	}
	query += "]}}}"
	size := 10000
	request := esapi.SearchRequest{
		Index:   indexes,
		Body:    strings.NewReader(query),
		Size:    &size,
		Timeout: 20 * time.Second,
		Sort:    []string{"{tm:desc"},
	}
	resp, err := db.api.ProcessRespWithCli(request)
	var esSearchResp HistorySearchResponse
	if err != nil {
		if strings.HasPrefix(err.Error(), "404") {
			return nil, nil
		}
		return nil, err
	}
	err = json.Unmarshal(resp, &esSearchResp)
	if err != nil {
		return nil, err
	}
	var hisList []AqiHistory
	if esSearchResp.Hits.Total.Value > 0 {
		for _, item := range esSearchResp.Hits.Hits {
			hisList = append(hisList, item.Source)
		}
		return hisList, nil
	}
	defer func() {
		resp = nil
	}()
	return nil, nil
}
