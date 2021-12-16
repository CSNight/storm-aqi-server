package db

import (
	"github.com/elastic/go-elasticsearch/v8/esapi"
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

type StationGetResponse struct {
	EsGetResponse
	Source AqiStation `json:"_source"`
}

type StationItem struct {
	EsSearchItem
	Source AqiStation `json:"_source"`
}

type StationSearchResponse struct {
	EsSearchRespMeta
	Hits struct {
		Total    EsRespTotal   `json:"total"`
		MaxScore float64       `json:"max_score"`
		Hits     []StationItem `json:"hits"`
	} `json:"hits"`
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
	var response StationGetResponse
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return nil, err
	}
	if response.Found {
		return &response.Source, nil
	}
	return nil, nil
}

func (db *DB) GetStationByName(name string) (*AqiStation, error) {
	sts, err := db.SearchStationByName(name, 1)
	if err != nil {
		return nil, err
	}
	if len(sts) > 0 {
		return &sts[0], nil
	}
	return nil, nil
}

func (db *DB) GetStationByCityName(name string) (*AqiStation, error) {
	sts, err := db.SearchStationByCityName(name, 1)
	if err != nil {
		return nil, err
	}
	if len(sts) > 0 {
		return &sts[0], nil
	}
	return nil, nil
}

func (db *DB) SearchStationByName(name string, size int) ([]AqiStation, error) {
	query := `{
        "query": {
            "wildcard": {
                "name": {
                    "case_insensitive": true,
                    "value": "*` + name + `*"
                }
            }
       }
    }`
	search := &esapi.SearchRequest{
		Index:   []string{db.Conf.StationIndex},
		Body:    strings.NewReader(query),
		Size:    &size,
		Sort:    nil,
		Timeout: 20 * time.Second,
	}
	resp, err := db.api.ProcessRespWithCli(search)
	var esSearchResp StationSearchResponse
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp, &esSearchResp)
	if err != nil {
		return nil, err
	}
	var sts []AqiStation
	if esSearchResp.Hits.Total.Value > 0 {
		for _, item := range esSearchResp.Hits.Hits {
			sts = append(sts, item.Source)
		}
		return sts, nil
	}
	return nil, nil
}

func (db *DB) SearchStationByCityName(name string, size int) ([]AqiStation, error) {
	query := `{
        "query": {
            "wildcard": {
                "city_name": {
                    "case_insensitive": true,
                    "value": "*` + name + `*"
                }
            }
       }
    }`
	search := &esapi.SearchRequest{
		Index:   []string{db.Conf.StationIndex},
		Body:    strings.NewReader(query),
		Size:    &size,
		Sort:    nil,
		Timeout: 20 * time.Second,
	}
	resp, err := db.api.ProcessRespWithCli(search)
	var esSearchResp StationSearchResponse
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp, &esSearchResp)
	if err != nil {
		return nil, err
	}
	var sts []AqiStation
	if esSearchResp.Hits.Total.Value > 0 {
		for _, item := range esSearchResp.Hits.Hits {
			sts = append(sts, item.Source)
		}
		return sts, nil
	}
	return nil, nil
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
