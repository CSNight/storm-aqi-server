package server

import "github.com/gofiber/fiber/v2"

func (app *AQIServer) Register(root fiber.Router) {
	root.Get("/station/@:sid", app.GetStationById)
	root.Get("/station/name/_get", app.GetStationByName)
	root.Get("/station/city/_get", app.GetStationByCityName)
	root.Get("/station/name/_search", app.SearchStationByName)
	root.Get("/station/city/_search", app.SearchStationByCityName)
	root.Get("/station/fuzzy/_search", app.SearchStationByCityName)
}
