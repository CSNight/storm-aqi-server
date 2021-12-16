package server

import "github.com/gofiber/fiber/v2"

func (app *AQIServer) GetStationById(ctx *fiber.Ctx) error {
	sid := ctx.Query("sid")
	st, err := app.DB.GetStationById(sid)
	if err != nil {
		return FailWithMessage(500, err.Error(), ctx)
	}
	return OkWithData(st, ctx)
}

func (app *AQIServer) GetStationByName(ctx *fiber.Ctx) error {
	name := ctx.Query("name")
	st, err := app.DB.GetStationByName(name)
	if err != nil {
		return FailWithMessage(500, err.Error(), ctx)
	}
	if st == nil {
		return OkWithNotFound("application/json", ctx)
	}
	return OkWithData(st, ctx)
}

func (app *AQIServer) GetStationByCityName(ctx *fiber.Ctx) error {
	name := ctx.Query("name")
	st, err := app.DB.GetStationByCityName(name)
	if err != nil {
		return FailWithMessage(500, err.Error(), ctx)
	}
	if st == nil {
		return OkWithNotFound("application/json", ctx)
	}
	return OkWithData(st, ctx)
}
