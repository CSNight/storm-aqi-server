package server

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type RealtimeRequest struct {
	QType string `json:"qType" validate:"required,oneof=_get"`
	PType string `json:"pType" validate:"required,oneof=all single"`
	Sid   string `json:"sid" validate:"required_if=QType _get,number"`
	Pol   string `json:"pol" validate:"required_if=QType _get PType single,omitempty,oneof=no2 pm25 pm10 o3 so2 co"`
}

type ForecastRequest struct {
	QType string `json:"qType" validate:"required,oneof=_get"`
	PType string `json:"pType" validate:"required,oneof=all single"`
	Sid   string `json:"sid" validate:"required_if=QType _get,number"`
	Pol   string `json:"pol" validate:"required_if=QType _get PType single,omitempty,oneof=no2 pm25 pm10 o3 so2 co"`
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
		return app.GetAllRealtime(query.Sid, ctx)
	} else {
		return app.GetSingleRealtime(query.Sid, query.Pol, ctx)
	}
}

func (app *AQIServer) GetAllRealtime(sid string, ctx *fiber.Ctx) error {
	rt, err := app.DB.GetAqiRealtimeById(sid)
	if err != nil {
		return FailWithMessage(http.StatusInternalServerError, err.Error(), ctx)
	}
	if rt == nil {
		return OkWithNotFound(fiber.MIMEApplicationJSON, ctx)
	}
	return OkWithData(rt, ctx)
}

func (app *AQIServer) GetSingleRealtime(sid string, pol string, ctx *fiber.Ctx) error {
	rt, err := app.DB.GetAqiRealtimeByIdAndPol(sid, pol)
	if err != nil {
		return FailWithMessage(http.StatusInternalServerError, err.Error(), ctx)
	}
	if rt == nil {
		return OkWithNotFound(fiber.MIMEApplicationJSON, ctx)
	}
	return OkWithData(rt, ctx)
}

func (app *AQIServer) GetAllForecast(sid string, ctx *fiber.Ctx) error {
	rt, err := app.DB.GetForecast(sid)
	if err != nil {
		return FailWithMessage(http.StatusInternalServerError, err.Error(), ctx)
	}
	if rt == nil {
		return OkWithNotFound(fiber.MIMEApplicationJSON, ctx)
	}
	return OkWithData(rt, ctx)
}
