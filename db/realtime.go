package db

type AqiRealtime struct {
	Idx      int     `json:"idx"`
	Sid      string  `json:"sid"`
	Pol      string  `json:"pol"`
	Data     float64 `json:"data"`
	Daily    string  `json:"daily"`
	Forecast string  `json:"forecast"`
	Tz       string  `json:"tz"`
	Tm       int64   `json:"tm"`
	Tms      string  `json:"tms"`
}
