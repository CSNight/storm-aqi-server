package server

import "github.com/gofiber/fiber/v2"

func (app *AQIServer) GetStationById(ctx *fiber.Ctx) error {
	sid := ctx.Query("sid")
	st, err := app.DB.GetStationById(sid)
	if err != nil {
		return err
	}
	return ctx.JSON(st)
}

func (app *AQIServer) GetStationByName(ctx *fiber.Ctx) error {
	name := ctx.Query("name")
	sts, err := app.DB.SearchStationByName(name)
	if err != nil {
		return err
	}
	return ctx.JSON(sts)
}
