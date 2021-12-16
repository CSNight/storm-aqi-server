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
	Source      string `json:"_source"`
}

type EsSearchResponse struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float64  `json:"max_score"`
		Hits     []string `json:"hits"`
	} `json:"hits"`
}
