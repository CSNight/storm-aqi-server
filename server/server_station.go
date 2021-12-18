package server

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type StationGetRequest struct {
	QType string `json:"qtype" validate:"required,oneof=_get"`
	PType string `json:"ptype" validate:"required,oneof=sid ip name city loc"`
	Sid   string `json:"sid" validate:"required_if=PType sid,omitempty,number"`
	Name  string `json:"name" validate:"required_if=PType name,omitempty,excludesall=@?*%"`
	City  string `json:"city" validate:"required_if=PType city,omitempty,excludesall=@?*%"`
	Ip    string `json:"ip" validate:"required_if=PType ip,omitempty,ip4_addr"`
	Lng   string `json:"lng" validate:"required_if=PType loc,omitempty,longitude"`
	Lat   string `json:"lat" validate:"required_if=PType loc,omitempty,latitude"`
}

type StationSearchRequest struct {
	QType string `json:"qtype" validate:"required,oneof=_search"`
	PType string `json:"ptype" validate:"required,oneof=name city area"`
	Size  int    `json:"size" validate:"required,number,min=1,max=10000"`
	Name  string `json:"name" validate:"required_if=PType name,omitempty,excludesall=@?*%"`
	City  string `json:"city" validate:"required_if=PType city,omitempty,excludesall=@?*%"`
}

func (app *AQIServer) StationGet(ctx *fiber.Ctx) error {
	var query StationGetRequest
	err := ctx.QueryParser(&query)
	if err != nil {
		return FailWithMessage(http.StatusBadRequest, "can't parser params", ctx)
	}
	errResp := ValidateStruct(query)
	if errResp != nil {
		return FailWithDetailed(http.StatusBadRequest, errResp, "", ctx)
	}
	if query.PType == "sid" {
		return app.GetStationById(query.Sid, ctx)
	} else if query.PType == "name" {
		return app.GetStationByName(query.Name, ctx)
	} else if query.PType == "city" {
		return app.GetStationByCity(query.City, ctx)
	} else if query.PType == "loc" {
		return app.GetStationByCity(query.City, ctx)
	} else {
		return app.GetStationByIp(query.Ip, ctx)
	}
}

func (app *AQIServer) StationSearch(ctx *fiber.Ctx) error {
	var query StationSearchRequest
	err := ctx.QueryParser(&query)
	if err != nil {
		return FailWithMessage(http.StatusBadRequest, "can't parser paramsOh6ChfVOSqPq2IgQ", ctx)
	}
	errResp := ValidateStruct(query)
	if errResp != nil {
		return FailWithDetailed(http.StatusBadRequest, errResp, "", ctx)
	}
	if query.PType == "name" {
		return app.SearchStationByName(query.Name, query.Size, ctx)
	} else if query.PType == "city" {
		return app.SearchStationByName(query.City, query.Size, ctx)
	} else {
		return nil
	}
}

func (app *AQIServer) GetStationById(sid string, ctx *fiber.Ctx) error {
	st, err := app.DB.GetStationById(sid)
	if err != nil {
		return FailWithMessage(http.StatusInternalServerError, err.Error(), ctx)
	}
	if st == nil {
		return OkWithNotFound(fiber.MIMEApplicationJSON, ctx)
	}
	return OkWithData(st, ctx)
}

func (app *AQIServer) GetStationByName(name string, ctx *fiber.Ctx) error {
	st, err := app.DB.GetStationByName(name)
	if err != nil {
		return FailWithMessage(http.StatusInternalServerError, err.Error(), ctx)
	}
	if st == nil {
		return OkWithNotFound(fiber.MIMEApplicationJSON, ctx)
	}
	return OkWithData(st, ctx)
}

func (app *AQIServer) GetStationByCity(city string, ctx *fiber.Ctx) error {
	st, err := app.DB.GetStationByCityName(city)
	if err != nil {
		return FailWithMessage(http.StatusInternalServerError, err.Error(), ctx)
	}
	if st == nil {
		return OkWithNotFound(fiber.MIMEApplicationJSON, ctx)
	}
	return OkWithData(st, ctx)
}

func (app *AQIServer) GetStationByIp(ip string, ctx *fiber.Ctx) error {
	st, err := app.DB.GetStationByIp(ip)
	if err != nil {
		return FailWithMessage(http.StatusInternalServerError, err.Error(), ctx)
	}
	if st == nil {
		return OkWithNotFound(fiber.MIMEApplicationJSON, ctx)
	}
	return OkWithData(st, ctx)
}

func (app *AQIServer) SearchStationByName(name string, size int, ctx *fiber.Ctx) error {
	st, err := app.DB.SearchStationByName(name, size)
	if err != nil {
		return FailWithMessage(http.StatusInternalServerError, err.Error(), ctx)
	}
	return OkWithData(st, ctx)
}

func (app *AQIServer) SearchStationByCityName(city string, size int, ctx *fiber.Ctx) error {
	st, err := app.DB.SearchStationByCityName(city, size)
	if err != nil {
		return FailWithMessage(http.StatusInternalServerError, err.Error(), ctx)
	}
	return OkWithData(st, ctx)
}
