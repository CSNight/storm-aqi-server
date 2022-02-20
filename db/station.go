package db

import (
	"github.com/csnight/storm-aqi-server/tools"
	"github.com/elastic/go-elasticsearch/v8/esapi"
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

type Source struct {
	Logo string   `json:"logo,omitempty"`
	Name string   `json:"name,omitempty"`
	Url  string   `json:"url,omitempty"`
	Pols []string `json:"pols,omitempty"`
}

type AqiStationResp struct {
	Sid      string   `json:"sid"`
	Idx      int      `json:"idx"`
	Name     string   `json:"name"`
	Loc      GeoPoint `json:"loc"`
	UpTime   int64    `json:"up_time"`
	Tms      string   `json:"tms"`
	Tz       string   `json:"tz"`
	CityName string   `json:"city_name,omitempty"`
	HisRange string   `json:"his_range,omitempty"`
	Sources  []Source `json:"sources,omitempty"`
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

func (db *DB) GetStationById(idx string) (*AqiStationResp, error) {
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
		db.log.Error("GetStationById(). es.ProcessRespWithCli(). err:", zap.Error(err))
		return nil, err
	}
	var response StationGetResponse
	err = json.Unmarshal(resp, &response)
	if err != nil {
		db.log.Error("GetStationById(). json.Unmarshal(). err:", zap.Error(err))
		return nil, err
	}
	if response.Found {
		return buildResponse(&response.Source)
	}
	return nil, nil
}

func (db *DB) GetStationByName(name string) (*AqiStationResp, error) {
	sts, err := db.SearchStationsByName(name, 1)
	if err != nil {
		db.log.Error("GetStationByName(). SearchStationsByName("+name+"). err:", zap.Error(err))
		return nil, err
	}
	if len(sts) > 0 {
		return &sts[0], nil
	}
	return nil, nil
}

func (db *DB) GetStationByCityName(name string) (*AqiStationResp, error) {
	sts, err := db.SearchStationsByCityName(name, 1)
	if err != nil {
		db.log.Error("GetStationByCityName(). SearchStationsByCityName("+name+"). err:", zap.Error(err))
		return nil, err
	}
	if len(sts) > 0 {
		return &sts[0], nil
	}
	return nil, nil
}

func (db *DB) SearchStationsByName(name string, size int) ([]AqiStationResp, error) {
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
		db.log.Error("SearchStationsByName(). es.ProcessRespWithCli(). err:", zap.String("query", strings.ReplaceAll(query, " ", "")), zap.Error(err))
		return nil, err
	}
	err = json.Unmarshal(resp, &esSearchResp)
	if err != nil {
		db.log.Error("SearchStationsByName(). json.Unmarshal(). err:", zap.Error(err))
		return nil, err
	}
	var sts []AqiStation
	if esSearchResp.Hits.Total.Value > 0 {
		for _, item := range esSearchResp.Hits.Hits {
			sts = append(sts, item.Source)
		}
		return buildResponses(sts), nil
	}
	return []AqiStationResp{}, nil
}

func (db *DB) SearchStationsByCityName(name string, size int) ([]AqiStationResp, error) {
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
		db.log.Error("SearchStationsByCityName(). es.ProcessRespWithCli(). err:", zap.String("query", strings.ReplaceAll(query, " ", "")), zap.Error(err))
		return nil, err
	}
	err = json.Unmarshal(resp, &esSearchResp)
	if err != nil {
		db.log.Error("SearchStationsByCityName(). json.Unmarshal(). err:", zap.Error(err))
		return nil, err
	}
	var sts []AqiStation
	if esSearchResp.Hits.Total.Value > 0 {
		for _, item := range esSearchResp.Hits.Hits {
			sts = append(sts, item.Source)
		}
		return buildResponses(sts), nil
	}
	return []AqiStationResp{}, nil
}

func (db *DB) SearchStationByRadius(x string, y string, dis float64, unit string, size int) ([]AqiStationResp, error) {
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
		db.log.Error("SearchStationByRadius(). es.ProcessRespWithCli(). err:", zap.String("query", strings.ReplaceAll(query, " ", "")), zap.Error(err))
		return nil, err
	}
	err = json.Unmarshal(resp, &esSearchResp)
	if err != nil {
		db.log.Error("SearchStationByRadius(). json.Unmarshal(). err:", zap.Error(err))
		return nil, err
	}
	var sts []AqiStation
	if esSearchResp.Hits.Total.Value > 0 {
		for _, item := range esSearchResp.Hits.Hits {
			sts = append(sts, item.Source)
		}
		return buildResponses(sts), nil
	}
	return []AqiStationResp{}, nil

}

func (db *DB) GetStationByIp(ip string) (*AqiStationResp, error) {
	ipNet := net.ParseIP(ip)
	city, err := db.ipDB.City(ipNet.To4())
	if err != nil {
		db.log.Error("GetStationByIp(). ipDB.City(). err:", zap.Error(err))
		return nil, err
	}
	loc := city.Location
	x := strconv.FormatFloat(loc.Longitude, 'f', 10, 64)
	y := strconv.FormatFloat(loc.Latitude, 'f', 10, 64)
	sts, err := db.SearchStationByRadius(x, y, 10, "km", 10)
	if err != nil {
		db.log.Error("GetStationByIp(). db.SearchStationByRadius(). err:", zap.Error(err))
		return nil, err
	}
	if len(sts) == 0 {
		return nil, nil
	}
	return &sts[0], nil
}

func (db *DB) SearchStationsByArea(bounds Bounds, size int) ([]AqiStationResp, error) {
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
		db.log.Error("SearchStationsByArea(). es.ProcessRespWithCli(). err:", zap.String("query", strings.ReplaceAll(query, " ", "")), zap.Error(err))
		return nil, err
	}
	err = json.Unmarshal(resp, &esSearchResp)
	if err != nil {
		db.log.Error("SearchStationsByArea(). json.Unmarshal(). err:", zap.Error(err))
		return nil, err
	}
	var sts []AqiStation
	if esSearchResp.Hits.Total.Value > 0 {
		for _, item := range esSearchResp.Hits.Hits {
			sts = append(sts, item.Source)
		}
		return buildResponses(sts), nil
	}
	return []AqiStationResp{}, nil
}

func (db *DB) GetAllStations() ([]AqiStationResp, error) {
	query := `{
       "query":{"match_all":{}}
    }`
	return db.ScrollSearchStation(query)
}

func (db *DB) GetStationsByRange(st int, et int) ([]AqiStationResp, error) {
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

func (db *DB) ScrollSearchStation(query string) ([]AqiStationResp, error) {
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
		if strings.HasPrefix(err.Error(), "404") {
			return nil, nil
		}
		db.log.Error("ScrollSearchStation(). es.ScrollSearch(). err:", zap.String("query", strings.ReplaceAll(query, " ", "")), zap.Error(err))
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
	return buildResponses(sts), nil
}

func (db *DB) GetStationLogo(logo string) ([]byte, error) {
	return tools.GetObject(db.oss, logo)
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
		if len(st.Sources) > 0 {
			for _, source := range st.Sources {
				logoName := source.Logo
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
					db.log.Error("download station logo failed, err:", zap.Error(err))
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
func buildResponse(st *AqiStation) (*AqiStationResp, error) {
	var sources []Source
	err := json.UnmarshalFromString(st.Sources, &sources)
	if err != nil {
		return nil, err
	}
	return &AqiStationResp{
		Sid:      st.Sid,
		Idx:      st.Idx,
		Name:     st.Name,
		Loc:      st.Loc,
		UpTime:   st.UpTime,
		Tms:      st.Tms,
		Tz:       st.Tz,
		CityName: st.CityName,
		HisRange: st.HisRange,
		Sources:  sources,
	}, nil
}
func buildResponses(sts []AqiStation) []AqiStationResp {
	var stResp []AqiStationResp
	for _, st := range sts {
		stn, err := buildResponse(&st)
		if err != nil {
			continue
		}
		stResp = append(stResp, *stn)
	}
	return stResp
}
