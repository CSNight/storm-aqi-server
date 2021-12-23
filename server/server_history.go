package server

type HistoryRequest struct {
	QType string `json:"qType" validate:"required,oneof=_get"`
	PType string `json:"pType" validate:"required,oneof=time"`
	Sid   string `json:"sid" validate:"required_if=QType _get,number"`
	Pol   string `json:"pol" validate:"required_if=QType _get,oneof=all no2 pm25 pm10 o3 so2 co"`
	Range string `json:"range" validate:"required_if=QType _get PType time,omitempty,oneof="`
}
