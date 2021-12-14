package db

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
