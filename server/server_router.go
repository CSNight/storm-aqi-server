package server

import "github.com/gofiber/fiber/v2"

func (app *AQIServer) Register(root fiber.Router) {
	root.Get("/station", app.StationGet)
	root.Get("/stations", app.StationSearch)
	root.Get("/realtime", app.RealtimeGet)
}
