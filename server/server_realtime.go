package server

type RealtimeRequest struct {
	QType string `json:"qType" validate:"required,oneof=_get"`
	PType string `json:"pType" validate:"required,oneof=all single"`
	Sid   string `json:"sid" validate:"required_if=QType _get,number"`
	Pol   string `json:"pol" validate:"required_if=QType _get PType single,omitempty,oneof=no2 pm25 pm10 o3 so2 co"`
}

type ForecastRequest struct {
}
