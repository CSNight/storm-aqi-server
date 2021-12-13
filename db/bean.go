package db

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

type GeoPoint struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

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
