package server

import (
	"github.com/gofiber/fiber/v2"
	"reflect"
	"strconv"
)

func (app *AQIServer) GetStationById(ctx *fiber.Ctx) error {
	sid := ctx.Params("sid")
	valid, err := paramsCheck(sid, ctx, reflect.String)
	if !valid {
		return err
	}
	st, err := app.DB.GetStationById(sid)
	if err != nil {
		return FailWithMessage(500, err.Error(), ctx)
	}
	return OkWithData(st, ctx)
}

func (app *AQIServer) GetStationByName(ctx *fiber.Ctx) error {
	name := ctx.Query("name")
	valid, err := paramsCheck(name, ctx, reflect.String)
	if !valid {
		return err
	}
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
	city := ctx.Query("city")
	valid, err := paramsCheck(city, ctx, reflect.String)
	if !valid {
		return err
	}
	st, err := app.DB.GetStationByCityName(city)
	if err != nil {
		return FailWithMessage(500, err.Error(), ctx)
	}
	if st == nil {
		return OkWithNotFound("application/json", ctx)
	}
	return OkWithData(st, ctx)
}

func (app *AQIServer) SearchStationByName(ctx *fiber.Ctx) error {
	name := ctx.Query("name")
	size := ctx.Query("size")
	valid, err := paramsCheck(name, ctx, reflect.String)
	if !valid {
		return err
	}
	validSize, err := paramsCheck(size, ctx, reflect.Int)
	if !validSize {
		return err
	}
	sizeInt, _ := strconv.ParseInt(size, 10, 64)
	st, err := app.DB.SearchStationByName(name, int(sizeInt))
	if err != nil {
		return FailWithMessage(500, err.Error(), ctx)
	}
	if st == nil {
		return OkWithNotFound("application/json", ctx)
	}
	return OkWithData(st, ctx)
}

func (app *AQIServer) SearchStationByCityName(ctx *fiber.Ctx) error {
	city := ctx.Query("city")
	size := ctx.Query("size")
	valid, err := paramsCheck(city, ctx, reflect.String)
	if !valid {
		return err
	}
	validSize, err := paramsCheck(size, ctx, reflect.Int)
	if !validSize {
		return err
	}
	sizeInt, _ := strconv.ParseInt(size, 10, 64)
	st, err := app.DB.SearchStationByCityName(city, int(sizeInt))
	if err != nil {
		return FailWithMessage(500, err.Error(), ctx)
	}
	if st == nil {
		return OkWithNotFound("application/json", ctx)
	}
	return OkWithData(st, ctx)
}
