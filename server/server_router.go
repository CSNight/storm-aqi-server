package server

import "github.com/gofiber/fiber/v2"

func (app *AQIServer) Register(root fiber.Router) {
	root.Get("/stationById", app.GetStationById)
	root.Get("/stationByName", app.GetStationByName)
	root.Get("/stationByCity", app.GetStationByName)
}
