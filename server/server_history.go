package server

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type HistoryRequest struct {
	QType string `json:"qType" validate:"required,oneof=_get"`
	PType string `json:"pType" validate:"required,oneof=time"`
	Sid   string `json:"sid" validate:"required_if=QType _get,number"`
	Pol   string `json:"pol" validate:"required_if=QType _get,oneof=all no2 pm25 pm10 o3 so2 co"`
	Range string `json:"range" validate:"required_if=QType _get PType time,omitempty,oneof=lastDay lastWeek lastMonth lastQuarter lastYear"`
}

func (app *AQIServer) HistoryGet(ctx *fiber.Ctx) error {
	var query HistoryRequest
	err := ctx.QueryParser(&query)
	if err != nil {
		return FailWithMessage(http.StatusBadRequest, "can't parser params", ctx)
	}
	errResp := ValidateStruct(query)
	if errResp != nil {
		return FailWithDetailed(http.StatusBadRequest, errResp, "", ctx)
	}
	if query.PType == "time" {
		if query.Range == "lastDay" {
			return app.GetHistoryYesterday(query.Sid, query.Pol, ctx)
		} else if query.Range == "lastWeek" {
			return app.GetHistoryWeek(query.Sid, query.Pol, ctx)
		} else if query.Range == "lastMonth" {
			return app.GetHistoryMonth(query.Sid, query.Pol, ctx)
		} else if query.Range == "lastQuarter" {
			return app.GetHistoryQuarter(query.Sid, query.Pol, ctx)
		} else {
			return app.GetHistoryYear(query.Sid, query.Pol, ctx)
		}
	} else {
		return nil
	}
}

func (app *AQIServer) GetHistoryYesterday(sid string, pol string, ctx *fiber.Ctx) error {
	rt, err := app.DB.GetHistoryYesterday(sid, pol)
	if err != nil {
		return FailWithMessage(http.StatusInternalServerError, err.Error(), ctx)
	}
	if rt == nil {
		return OkWithNotFound(fiber.MIMEApplicationJSON, ctx)
	}
	return OkWithData(rt, ctx)
}

func (app *AQIServer) GetHistoryWeek(sid string, pol string, ctx *fiber.Ctx) error {
	rt, err := app.DB.GetHistoryLastWeek(sid, pol)
	if err != nil {
		return FailWithMessage(http.StatusInternalServerError, err.Error(), ctx)
	}
	if rt == nil {
		return OkWithNotFound(fiber.MIMEApplicationJSON, ctx)
	}
	return OkWithData(rt, ctx)
}

func (app *AQIServer) GetHistoryMonth(sid string, pol string, ctx *fiber.Ctx) error {
	rt, err := app.DB.GetHistoryLastMonth(sid, pol)
	if err != nil {
		return FailWithMessage(http.StatusInternalServerError, err.Error(), ctx)
	}
	if rt == nil {
		return OkWithNotFound(fiber.MIMEApplicationJSON, ctx)
	}
	return OkWithData(rt, ctx)
}

func (app *AQIServer) GetHistoryQuarter(sid string, pol string, ctx *fiber.Ctx) error {
	rt, err := app.DB.GetHistoryLastQuarter(sid, pol)
	if err != nil {
		return FailWithMessage(http.StatusInternalServerError, err.Error(), ctx)
	}
	if rt == nil {
		return OkWithNotFound(fiber.MIMEApplicationJSON, ctx)
	}
	return OkWithData(rt, ctx)
}

func (app *AQIServer) GetHistoryYear(sid string, pol string, ctx *fiber.Ctx) error {
	rt, err := app.DB.GetHistoryYear(sid, pol)
	if err != nil {
		return FailWithMessage(http.StatusInternalServerError, err.Error(), ctx)
	}
	if rt == nil {
		return OkWithNotFound(fiber.MIMEApplicationJSON, ctx)
	}
	return OkWithData(rt, ctx)
}
