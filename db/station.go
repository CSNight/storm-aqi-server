package db

import (
	"aqi-server/tools"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"net"
	"strconv"
	"strings"
	"sync"
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
	defer func() {
		resp = nil
	}()
	if err != nil {
		if strings.HasPrefix(err.Error(), "404") {
			return nil, nil
		}
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
	sts, err := db.SearchStationsByName(name, 1)
	if err != nil {
		return nil, err
	}
	if len(sts) > 0 {
		return &sts[0], nil
	}
	return nil, nil
}

func (db *DB) GetStationByCityName(name string) (*AqiStation, error) {
	sts, err := db.SearchStationsByCityName(name, 1)
	if err != nil {
		return nil, err
	}
	if len(sts) > 0 {
		return &sts[0], nil
	}
	return nil, nil
}

func (db *DB) SearchStationsByName(name string, size int) ([]AqiStation, error) {
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
	defer func() {
		resp = nil
	}()
	var esSearchResp StationSearchResponse
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
	var sts []AqiStation
	if esSearchResp.Hits.Total.Value > 0 {
		for _, item := range esSearchResp.Hits.Hits {
			sts = append(sts, item.Source)
		}
		return sts, nil
	}
	return []AqiStation{}, nil
}

func (db *DB) SearchStationsByCityName(name string, size int) ([]AqiStation, error) {
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
	defer func() {
		resp = nil
	}()
	var esSearchResp StationSearchResponse
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
	var sts []AqiStation
	if esSearchResp.Hits.Total.Value > 0 {
		for _, item := range esSearchResp.Hits.Hits {
			sts = append(sts, item.Source)
		}
		return sts, nil
	}
	return []AqiStation{}, nil
}

func (db *DB) SearchStationByRadius(x string, y string, dis float64, unit string, size int) ([]AqiStation, error) {
	disStr := strconv.FormatFloat(dis, 'f', 8, 64) + unit
	query := `{
      "query": {
        "bool": {
          "must": {
            "match_all": {}
          },
          "filter": {
            "geo_distance": {
              "distance": "` + disStr + `",
              "loc": {
                "lat": ` + y + `,
                "lon": ` + x + `
              }
            }
          }
        }
      }
    }`
	search := &esapi.SearchRequest{
		Index: []string{db.Conf.StationIndex},
		Body:  strings.NewReader(query),
		Size:  &size,
		Sort: []string{`{"_geo_distance": {
        "loc": {
          "lat": ` + y + `,
          "lon": ` + x + `
        }, 
        "order": "asc",
        "unit": ` + unit + `
      }}`},
	}
	resp, err := db.api.ProcessRespWithCli(search)
	defer func() {
		resp = nil
	}()
	var esSearchResp StationSearchResponse
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
	var sts []AqiStation
	if esSearchResp.Hits.Total.Value > 0 {
		for _, item := range esSearchResp.Hits.Hits {
			sts = append(sts, item.Source)
		}
		return sts, nil
	}
	return []AqiStation{}, nil

}

func (db *DB) GetStationByIp(ip string) (*AqiStation, error) {
	ipNet := net.ParseIP(ip)
	city, err := db.ipDB.City(ipNet.To4())
	if err != nil {
		return nil, err
	}
	loc := city.Location
	x := strconv.FormatFloat(loc.Longitude, 'f', 10, 64)
	y := strconv.FormatFloat(loc.Latitude, 'f', 10, 64)
	sts, err := db.SearchStationByRadius(x, y, 10, "km", 10)
	if err != nil {
		return nil, err
	}
	if len(sts) == 0 {
		return nil, nil
	}
	return &sts[0], nil
}

func (db *DB) SearchStationsByArea(bounds Bounds) ([]AqiStation, error) {
	boundsBytes, err := json.Marshal(bounds)
	if err != nil {
		return nil, err
	}
	query := `{
      "query": {
        "bool": {
          "must": {
            "match_all": {}
          },
          "filter": {
            "geo_bounding_box": {
              "loc": ` + string(boundsBytes) + `
            }
          }
        }
      }
    }`
	size := 1000
	search := &esapi.SearchRequest{
		Index:   []string{db.Conf.StationIndex},
		Body:    strings.NewReader(query),
		Size:    &size,
		Sort:    nil,
		Timeout: 20 * time.Second,
	}
	resp, err := db.api.ProcessRespWithCli(search)
	defer func() {
		resp = nil
	}()
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
	return []AqiStation{}, nil
}

func (db *DB) GetAllStations() ([]AqiStation, error) {
	query := `{
       "query":{"match_all":{}}
    }`
	return db.ScrollSearchStation(query)
}

func (db *DB) GetStationsByRange(st int, et int) ([]AqiStation, error) {
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

func (db *DB) ScrollSearchStation(query string) ([]AqiStation, error) {
	size := 10000
	search := &esapi.SearchRequest{
		Index:  []string{db.Conf.StationIndex},
		Body:   strings.NewReader(query),
		Scroll: time.Second * 20,
		Size:   &size,
		Sort:   []string{"idx"},
	}
	results, err := db.api.ScrollSearch(search)
	if err != nil {
		return nil, err
	}
	var sts []AqiStation
	for _, hit := range results {
		var station AqiStation
		err = json.UnmarshalFromString(hit.Raw, &station)
		if err != nil {
			continue
		}
		sts = append(sts, station)
	}
	defer func() {
		results = nil
	}()
	return sts, nil
}

func (db *DB) SyncStationLogos() error {
	stations, err := db.GetAllStations()
	if err != nil {
		return err
	}
	queue := make(chan string, 50)
	wg := sync.WaitGroup{}
	var logos = map[string]string{}
	for _, st := range stations {
		if st.Sources != "" {
			root := gjson.Parse(st.Sources).Array()
			for _, source := range root {
				logoName := source.Get("logo").String()
				if logoName != "" {
					logos[logoName] = logoName
				}
			}
		}
	}
	for _, logo := range logos {
		if !tools.ExistObject(db.oss, "aqi/"+logo) {
			wg.Add(1)
			queue <- logo
			go func(logoImg string) {
				defer func() {
					<-queue
					wg.Done()
				}()
				image, err := tools.DownloadImage(db.Conf.ImageOss + logoImg)
				if err != nil {
					db.log.Sugar().Error(err)
					return
				}
				status := tools.PutObject(db.oss, image, "aqi/"+logoImg)
				if status {
					db.log.Info("save to oss success", zap.String("object", logoImg))
				} else {
					db.log.Error("save to oss failed", zap.String("object", logoImg))
				}
			}(logo)
		}
	}
	wg.Wait()
	return nil
}
