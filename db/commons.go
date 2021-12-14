package db

type GeoPoint struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type Bounds struct {
	TopLeft     GeoPoint
	BottomRight GeoPoint
}
