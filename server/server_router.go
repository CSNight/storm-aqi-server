package server

import "github.com/gofiber/fiber/v2"

func (app *AQIServer) Register(root fiber.Router) {
	root.Get("/aqi/station", app.StationGet)
	root.Get("/aqi/stations", app.StationSearch)
	root.Get("/aqi/realtime", app.RealtimeGet)
	root.Get("/aqi/forecast", app.ForecastGet)
}
