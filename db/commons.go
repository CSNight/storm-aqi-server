package db

import "go.uber.org/zap"

type GeoPoint struct {
	Lon float64 `json:"lon" validate:"longitude"`
	Lat float64 `json:"lat" validate:"latitude"`
}

type Bounds struct {
	TopLeft     GeoPoint `json:"top_left" validate:"required"`
	BottomRight GeoPoint `json:"bottom_right" validate:"required"`
}

type EsGetResponse struct {
	Index       string `json:"_index"`
	Type        string `json:"_type"`
	Id          string `json:"_id"`
	Version     int    `json:"_version"`
	SeqNo       int    `json:"_seq_no"`
	PrimaryTerm int    `json:"_primary_term"`
	Found       bool   `json:"found"`
}

type EsSearchRespMeta struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
}

type EsSearchItem struct {
	Index string  `json:"_index"`
	Type  string  `json:"_type"`
	Id    string  `json:"_id"`
	Score float64 `json:"_score"`
}

type EsRespTotal struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

func (db *DB) getStationFromCache(sid string) (*AqiStationResp, error) {
	var st AqiStationResp
	stb, err := db.cache.Get([]byte(sid))
	if err != nil {
		stp, err := db.GetStationById(sid)
		if err != nil {
			db.log.Error("getStationFromCache(). GetStationById(). err:", zap.Error(err))
			return nil, err
		}
		return stp, nil
	}
	err = json.Unmarshal(stb, &st)
	if err != nil {
		db.log.Error("getStationFromCache(). json.Unmarshal(). err:", zap.Error(err))
		return nil, err
	}
	return &st, nil
}
