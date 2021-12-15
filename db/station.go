package db

import (
	"errors"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
	"time"
)

type AqiStation struct {
	Sid      string   `json:"sid"`
	Idx      int      `json:"idx"`
	Name     string   `json:"name"`
	Loc      GeoPoint `json:"loc"`
	UpTime   int64    `json:"up_time"`
	Tms      string   `json:"tms"`
	Tz       string   `json:"tz"`
	CityName string   `json:"city_name,omitempty"`
	HisRange string   `json:"his_range,omitempty"`
	Sources  string   `json:"sources,omitempty"`
}

func (db *DB) GetStationById(idx string) (*AqiStation, error) {
	search := &esapi.GetRequest{
		Index:      db.Conf.StationIndex,
		DocumentID: idx,
	}
	resp, err := db.api.ProcessRespWithCli(search)
	if err != nil {
		return nil, err
	}
	res := gjson.ParseBytes(resp)
	if res.Get("found").Bool() {
		var station AqiStation
		err = json.UnmarshalFromString(res.Get("_source").String(), &station)
		if err != nil {
			return nil, err
		}
		return &station, nil
	}
	return nil, errors.New("record not found")
}

func (db *DB) GetStationByName(name string) ([]AqiStation, error) {
	size := 10000
	search := &esapi.SearchRequest{
		Index:   []string{db.Conf.StationIndex},
		Body:    nil,
		Size:    &size,
		Sort:    nil,
		Timeout: 20 * time.Second,
	}
	_, err := db.api.ProcessRespWithCli(search)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (db *DB) GetStationByCityName(name string) *AqiStation {
	return nil
}

func (db *DB) GetStationByLoc(x float64, y float64) *AqiStation {
	return nil
}

func (db *DB) GetStationByIp(ip string) *AqiStation {
	return nil
}

func (db *DB) GetStationByArea(bounds Bounds) {

}

func (db *DB) GetAllStations() []AqiStation {
	query := `{
       "query":{"match_all":{}}
    }`
	return db.ScrollSearchStation(query)
}

func (db *DB) GetStationByRange(st int, et int) []AqiStation {
	query := `{
       "query": {
           "range" : {
               "idx" : {
                   "gte" : ` + strconv.Itoa(st) + `,
                   "lte" : ` + strconv.Itoa(et) + `
               }
           }
       }
    }`
	return db.ScrollSearchStation(query)
}

func (db *DB) ScrollSearchStation(query string) []AqiStation {
	size := 10000
	search := &esapi.SearchRequest{
		Index:  []string{db.Conf.StationIndex},
		Body:   strings.NewReader(query),
		Scroll: time.Second * 20,
		Size:   &size,
		Sort:   []string{"idx"},
	}
	results := db.api.ScrollSearch(search)

	var sts []AqiStation
	for _, hit := range results {
		var station AqiStation
		err := json.UnmarshalFromString(hit.Raw, &station)
		if err != nil {
			continue
		}
		sts = append(sts, station)
	}
	return sts
}
