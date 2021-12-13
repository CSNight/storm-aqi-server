package db

import (
	"aqi-server/elastic"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"strconv"
	"strings"
	"time"
)

func (db *DB) GetStationById(idx string) *AqiStation {
	query := `{
       "query":{"match":{"sid":"` + idx + `"}}
    }`
	search := &esapi.SearchRequest{
		Index: []string{db.Conf.StationIndex},
		Body:  strings.NewReader(query),
	}
	elastic.ProcessResp(search, nil)
	return nil
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
