package server

import (
	"net/http"

	"github.com/csnight/storm-aqi-server/db"
	"github.com/gofiber/fiber/v2"
)

type RealtimeRequest struct {
	QType string `json:"qType" validate:"required,oneof=_get"`
	PType string `json:"pType" validate:"required,oneof=all single"`
	Sid   string `json:"sid" validate:"required_if=PType single,omitempty,number"`
	Pol   string `json:"pol" validate:"required_if=PType single,omitempty,oneof=all no2 pm25 pm10 o3 so2 co"`
}

type ForecastRequest struct {
	QType string `json:"qType" validate:"required,oneof=_get"`
	PType string `json:"pType" validate:"required,oneof=single"`
	Sid   string `json:"sid" validate:"required,number"`
	Pol   string `json:"pol" validate:"required,oneof=all no2 pm25 pm10 o3 so2 co"`
}

func (app *AQIServer) RealtimeGet(ctx *fiber.Ctx) error {
	var query RealtimeRequest
	err := ctx.QueryParser(&query)
	if err != nil {
		return FailWithMessage(http.StatusBadRequest, "can't parser params", ctx)
	}
	errResp := ValidateStruct(query)
	if errResp != nil {
		return FailWithDetailed(http.StatusBadRequest, errResp, "", ctx)
	}
	if query.PType == "all" {
		return app.GetAllRealtime(ctx)
	} else {
		return app.GetSingleRealtime(query.Sid, query.Pol, ctx)
	}
}

func (app *AQIServer) ForecastGet(ctx *fiber.Ctx) error {
	var query ForecastRequest
	err := ctx.QueryParser(&query)
	if err != nil {
		return FailWithMessage(http.StatusBadRequest, "can't parser params", ctx)
	}
	errResp := ValidateStruct(query)
	if errResp != nil {
		return FailWithDetailed(http.StatusBadRequest, errResp, "", ctx)
	}
	return app.GetForecastByPol(query.Sid, query.Pol, ctx)
}

func (app *AQIServer) GetAllRealtime(ctx *fiber.Ctx) error {
	rt, err := app.db.GetAllAqiRealtime()
	if err != nil {
		return FailWithMessage(http.StatusInternalServerError, err.Error(), ctx)
	}
	if rt == nil {
		return OkWithNotFound(fiber.MIMEApplicationJSON, ctx)
	}
	return OkWithData(rt, ctx)
}

func (app *AQIServer) GetSingleRealtime(sid string, pol string, ctx *fiber.Ctx) error {
	var rt *db.RealtimeResp
	var err error
	if pol == "all" {
		rt, err = app.db.GetAqiRealtimeById(sid)
	} else {
		rt, err = app.db.GetAqiRealtimeByIdAndPol(sid, pol)
	}
	if err != nil {
		return FailWithMessage(http.StatusInternalServerError, err.Error(), ctx)
	}
	if rt == nil {
		return OkWithNotFound(fiber.MIMEApplicationJSON, ctx)
	}
	return OkWithData(rt, ctx)
}

func (app *AQIServer) GetAllForecast(sid string, ctx *fiber.Ctx) error {
	fore, err := app.db.GetForecast(sid, "all")
	if err != nil {
		return FailWithMessage(http.StatusInternalServerError, err.Error(), ctx)
	}
	if fore == nil {
		return OkWithNotFound(fiber.MIMEApplicationJSON, ctx)
	}
	return OkWithData(fore, ctx)
}

func (app *AQIServer) GetForecastByPol(sid string, pol string, ctx *fiber.Ctx) error {
	fore, err := app.db.GetForecast(sid, pol)
	if err != nil {
		return FailWithMessage(http.StatusInternalServerError, err.Error(), ctx)
	}
	if fore == nil {
		return OkWithNotFound(fiber.MIMEApplicationJSON, ctx)
	}
	return OkWithData(fore, ctx)
}
