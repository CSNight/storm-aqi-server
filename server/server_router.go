package server

import "github.com/gofiber/fiber/v2"

func (app *AQIServer) Register(root fiber.Router) {
	root.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})
	root.Static("/static", "./assets/static")
	root.Get("/aqi/station", app.StationGet)
	root.Get("/aqi/stations", app.StationSearch)
	root.Get("/aqi/realtime", app.RealtimeGet)
	root.Get("/aqi/forecast", app.ForecastGet)
	root.Get("/aqi/history", app.HistoryGet)
	root.Get("/aqi/none", app.GetNoneStation)
	root.Get("/aqi/logo/:logo", app.StationLogoGet)
	root.Post("/aqi/sync_logo", app.SyncStationLog)
}
