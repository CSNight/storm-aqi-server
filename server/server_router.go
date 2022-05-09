package server

import "github.com/gofiber/fiber/v2"

func (app *AQIServer) Register(root fiber.Router) {
	root.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})
	root.Get("/station", app.StationGet)
	root.Get("/stations", app.StationSearch)
	root.Get("/realtime", app.RealtimeGet)
	root.Get("/forecast", app.ForecastGet)
	root.Get("/image", app.ImageGet)
	root.Get("/history", app.HistoryGet)
	root.Get("/none_his", app.GetNoneStation)
	root.Get("/logo/:logo", app.StationLogoGet)
	root.Post("/sync_logo", app.SyncStationLog)
}
