package server

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"
)

type HistoryRequest struct {
	QType  string `json:"qType" validate:"required,oneof=_get"`
	PType  string `json:"pType" validate:"required,oneof=recent range"`
	Sid    string `json:"sid" validate:"required_if=QType _get,number"`
	Pol    string `json:"pol" validate:"required_if=QType _get,oneof=all no2 pm25 pm10 o3 so2 co"`
	Recent string `json:"recent" validate:"required_if=QType _get PType recent,omitempty,oneof=lastDay lastWeek lastMonth lastQuarter lastYear"`
	Start  string `json:"start" validate:"required_if=QType _get PType range,omitempty,datetime=2006-01-02"`
	End    string `json:"end" validate:"required_if=QType _get PType range,omitempty,datetime=2006-01-02"`
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
	if query.PType == "recent" {
		switch query.Recent {
		case "lastDay":
			return app.GetHistoryYesterday(query.Sid, query.Pol, ctx)
		case "lastWeek":
			return app.GetHistoryWeek(query.Sid, query.Pol, ctx)
		case "lastMonth":
			return app.GetHistoryMonth(query.Sid, query.Pol, ctx)
		case "lastQuarter":
			return app.GetHistoryQuarter(query.Sid, query.Pol, ctx)
		case "lastYear":
			return app.GetHistoryYear(query.Sid, query.Pol, ctx)
		default:
			return nil
		}
	} else {
		return app.GetHistoryRange(query.Sid, query.Pol, query.Start, query.End, ctx)
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

func (app *AQIServer) GetNoneStation(ctx *fiber.Ctx) error {
	idx := app.DB.GetNoneStation()
	return OkWithData(idx, ctx)
}

func (app *AQIServer) GetHistoryRange(sid string, pol string, st string, et string, ctx *fiber.Ctx) error {
	stTime, err := time.ParseInLocation("2006-01-02", st, time.UTC)
	if err != nil {
		return FailWithMessage(http.StatusBadRequest, "bad start time", ctx)
	}
	etTime, err := time.ParseInLocation("2006-01-02", et, time.UTC)
	if err != nil {
		return FailWithMessage(http.StatusBadRequest, "bad end time", ctx)
	}
	if etTime.Before(stTime) {
		return FailWithMessage(http.StatusBadRequest, "end time can't less then start time", ctx)
	}
	rt, err := app.DB.GetHistoryRange(sid, pol, stTime, etTime)
	if err != nil {
		return FailWithMessage(http.StatusInternalServerError, err.Error(), ctx)
	}
	if rt == nil {
		return OkWithNotFound(fiber.MIMEApplicationJSON, ctx)
	}
	return OkWithData(rt, ctx)
}
