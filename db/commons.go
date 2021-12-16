package db

type GeoPoint struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type Bounds struct {
	TopLeft     GeoPoint
	BottomRight GeoPoint
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
