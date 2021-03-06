package db

import (
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"sort"
	"strconv"
	"strings"
	"sync"
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
	et := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	st := et.AddDate(0, 0, -1)
	hisList, err := db.getHistoryByRange(sid, pol, st, et)
	if err != nil {
		db.log.Error("GetHistoryYesterday(). db.getHistoryByRange(). err:", zap.Error(err))
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
	et := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	st := et.AddDate(0, 0, -7)
	hisList, err := db.getHistoryByRange(sid, pol, st, et)
	if err != nil {
		db.log.Error("GetHistoryLastWeek(). db.getHistoryByRange(). err:", zap.Error(err))
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
	et := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	st := et.AddDate(0, -1, 0)
	hisList, err := db.getHistoryByRange(sid, pol, st, et)
	if err != nil {
		db.log.Error("GetHistoryLastMonth(). db.getHistoryByRange(). err:", zap.Error(err))
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

func (db *DB) GetHistoryLastQuarter(sid string, pol string) (*AqiHistoryResp, error) {
	station, err := db.getStationFromCache(sid)
	if err != nil || station == nil {
		return nil, err
	}
	now := time.Now().UTC()
	et := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	st := et.AddDate(0, -3, 0)
	hisList, err := db.getHistoryByRange(sid, pol, st, et)
	if err != nil {
		db.log.Error("GetHistoryLastQuarter(). db.getHistoryByRange(). err:", zap.Error(err))
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
	et := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	st := et.AddDate(-1, 0, 0)
	hisList, err := db.getHistoryByRange(sid, pol, st, et)
	if err != nil {
		db.log.Error("GetHistoryYear(). db.getHistoryByRange(). err:", zap.Error(err))
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

func (db *DB) GetHistoryRange(sid string, pol string, st time.Time, et time.Time) (*AqiHistoryResp, error) {
	station, err := db.getStationFromCache(sid)
	if err != nil || station == nil {
		return nil, err
	}
	hisList, err := db.getHistoryByRange(sid, pol, st, et)
	if err != nil {
		db.log.Error("GetHistoryRange(). db.getHistoryByRange(). err:", zap.Error(err))
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

func (db *DB) GetNoneStation() []string {
	stations, _ := db.GetAllStations()
	var wg = sync.WaitGroup{}
	ch := make(chan bool, 20)
	defer close(ch)
	var indexes []string
	for i := 2014; i <= 2021; i++ {
		indexes = append(indexes, strings.Replace(db.Conf.HisIndex, "$year", strconv.Itoa(i), -1))
	}
	empty := make(chan string, 600)
	for _, st := range stations {
		wg.Add(1)
		ch <- true
		go func(sd string) {
			defer func() {
				<-ch
				wg.Done()
			}()
			req := esapi.CountRequest{
				Index: indexes,
				Body:  strings.NewReader(`{"query":{"match":{"sid":"` + sd + `"}}}`),
			}
			resp, err := db.api.ProcessRespWithCli(req)
			if err != nil {
				return
			}
			result := gjson.ParseBytes(resp)
			if result.Get("count").Int() == 0 {
				empty <- sd
			}
		}(st.Sid)
	}
	var sidx []string
	go func() {
		for sid := range empty {
			sidx = append(sidx, sid)
		}
	}()
	wg.Wait()
	close(empty)
	sort.Slice(sidx, func(i, j int) bool {
		a, _ := strconv.ParseInt(sidx[i], 10, 64)
		b, _ := strconv.ParseInt(sidx[j], 10, 64)
		return a < b
	})
	return sidx
}

func BuildResp(list []AqiHistory) map[string][]AqiHisItem {
	items := map[string][]AqiHisItem{}
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
	for _, v := range items {
		sort.Slice(v, func(i, j int) bool {
			return v[i].Tm > v[j].Tm
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
		query += `,{"match": {"pol": "` + pol + `"}}`
	}
	query += "]}}}"
	size := 10000
	request := esapi.SearchRequest{
		Index:   indexes,
		Body:    strings.NewReader(query),
		Size:    &size,
		Timeout: 20 * time.Second,
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
