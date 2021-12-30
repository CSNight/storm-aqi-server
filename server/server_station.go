package server

import (
	"aqi-server/db"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

type StationGetRequest struct {
	QType string `json:"qType" validate:"required,oneof=_get"`
	PType string `json:"pType" validate:"required,oneof=sid ip name city loc"`
	Sid   string `json:"sid" validate:"required_if=QType _get PType sid,omitempty,number"`
	Name  string `json:"name" validate:"required_if=QType _get PType name,omitempty,excludesall=@?*%"`
	City  string `json:"city" validate:"required_if=QType _get PType city,omitempty,excludesall=@?*%"`
	Ip    string `json:"ip" validate:"required_if=QType _get PType ip,omitempty,ip4_addr"`
	Lon   string `json:"lon" validate:"required_if=QType _get PType loc,omitempty,longitude"`
	Lat   string `json:"lat" validate:"required_if=QType _get PType loc,omitempty,latitude"`
}

type StationSearchRequest struct {
	QType       string    `json:"qType" validate:"required,oneof=_search _all"`
	PType       string    `json:"pType" validate:"required_if=QType _search,omitempty,oneof=name city area radius"`
	Size        int       `json:"size" validate:"required_if=QType _search,omitempty,number,min=1,max=10000"`
	Name        string    `json:"name" validate:"required_if=QType _search PType name,omitempty,excludesall=@?*%"`
	City        string    `json:"city" validate:"required_if=QType _search PType city,omitempty,excludesall=@?*%"`
	TopLeft     []float64 `json:"topLeft" validate:"required_if=QType _search PType area,omitempty,len=2"`
	BottomRight []float64 `json:"bottomRight" validate:"required_if=QType _search PType area,omitempty,len=2"`
	Center      []float64 `json:"center" validate:"required_if=QType _search PType radius,omitempty,len=2"`
	Radius      float64   `json:"radius" validate:"required_if=QType _search PType radius,omitempty,gt=0,max=10000"`
	Unit        string    `json:"unit" validate:"required_if=QType _search PType radius,omitempty,oneof=kilometers miles meters"`
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
		return app.GetStationByLoc(query.Lon, query.Lat, ctx)
	} else {
		return app.GetStationByIp(query.Ip, ctx)
	}
}

func (app *AQIServer) StationSearch(ctx *fiber.Ctx) error {
	var query StationSearchRequest
	err := ctx.QueryParser(&query)
	if err != nil {
		return FailWithMessage(http.StatusBadRequest, "can't parser params", ctx)
	}
	errResp := ValidateStruct(query)
	if errResp != nil {
		return FailWithDetailed(http.StatusBadRequest, errResp, "", ctx)
	}
	if query.QType == "_all" {
		return app.SearchAllStations(ctx)
	}
	if query.PType == "name" {
		return app.SearchStationsByName(query.Name, query.Size, ctx)
	} else if query.PType == "city" {
		return app.SearchStationsByCityName(query.City, query.Size, ctx)
	} else if query.PType == "area" {
		errResp = ValidateVar(query.TopLeft[0], "longitude")
		if errResp != nil {
			return FailWithDetailed(http.StatusBadRequest, errResp, "", ctx)
		}
		errResp = ValidateVar(query.TopLeft[1], "latitude")
		if errResp != nil {
			return FailWithDetailed(http.StatusBadRequest, errResp, "", ctx)
		}
		errResp = ValidateVar(query.BottomRight[0], "longitude")
		if errResp != nil {
			return FailWithDetailed(http.StatusBadRequest, errResp, "", ctx)
		}
		errResp = ValidateVar(query.BottomRight[1], "latitude")
		if errResp != nil {
			return FailWithDetailed(http.StatusBadRequest, errResp, "", ctx)
		}
		return app.SearchStationsByArea(query.TopLeft, query.BottomRight, ctx)
	} else {
		errResp = ValidateVar(query.Center[0], "longitude")
		if errResp != nil {
			return FailWithDetailed(http.StatusBadRequest, errResp, "", ctx)
		}
		errResp = ValidateVar(query.Center[1], "latitude")
		if errResp != nil {
			return FailWithDetailed(http.StatusBadRequest, errResp, "", ctx)
		}
		return app.SearchStationsByRadius(query.Center, query.Unit, query.Radius, query.Size, ctx)
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

func (app *AQIServer) GetStationByLoc(x string, y string, ctx *fiber.Ctx) error {
	st, err := app.DB.SearchStationByRadius(x, y, 10, "km", 10)
	if err != nil {
		return FailWithMessage(http.StatusInternalServerError, err.Error(), ctx)
	}
	if len(st) == 0 {
		return OkWithNotFound(fiber.MIMEApplicationJSON, ctx)
	}
	return OkWithData(st[0], ctx)
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

func (app *AQIServer) SearchStationsByName(name string, size int, ctx *fiber.Ctx) error {
	sts, err := app.DB.SearchStationsByName(name, size)
	if err != nil {
		return FailWithMessage(http.StatusInternalServerError, err.Error(), ctx)
	}
	return OkWithData(sts, ctx)
}

func (app *AQIServer) SearchStationsByCityName(city string, size int, ctx *fiber.Ctx) error {
	sts, err := app.DB.SearchStationsByCityName(city, size)
	if err != nil {
		return FailWithMessage(http.StatusInternalServerError, err.Error(), ctx)
	}
	return OkWithData(sts, ctx)
}

func (app *AQIServer) SearchStationsByArea(topLeft []float64, bottomRight []float64, ctx *fiber.Ctx) error {
	sts, err := app.DB.SearchStationsByArea(db.Bounds{
		TopLeft: db.GeoPoint{
			Lon: topLeft[0],
			Lat: topLeft[1],
		},
		BottomRight: db.GeoPoint{
			Lon: bottomRight[0],
			Lat: bottomRight[1],
		},
	})
	if err != nil {
		return FailWithMessage(http.StatusInternalServerError, err.Error(), ctx)
	}
	return OkWithData(sts, ctx)
}

func (app *AQIServer) SearchStationsByRadius(center []float64, unit string, radius float64, size int, ctx *fiber.Ctx) error {
	x := strconv.FormatFloat(center[0], 'f', 8, 64)
	y := strconv.FormatFloat(center[1], 'f', 8, 64)
	unitMark := "km"
	switch unit {
	default:
	case "kilometers":
		unitMark = "km"
		break
	case "miles":
		unitMark = "mi"
		break
	case "meters":
		unitMark = "m"
		break
	}
	sts, err := app.DB.SearchStationByRadius(x, y, radius, unitMark, size)
	if err != nil {
		return FailWithMessage(http.StatusInternalServerError, err.Error(), ctx)
	}
	return OkWithData(sts, ctx)
}

func (app *AQIServer) SearchAllStations(ctx *fiber.Ctx) error {
	sts, err := app.DB.GetAllStations()
	if err != nil {
		return FailWithMessage(http.StatusInternalServerError, err.Error(), ctx)
	}
	return OkWithData(sts, ctx)
}

func (app *AQIServer) StationLogoGet(ctx *fiber.Ctx) error {
	logo := ctx.Params("logo")
	if logo == "" {
		return FailWithMessage(http.StatusNotFound, "empty logo", ctx)
	}
	img, err := app.DB.GetStationLogo(logo)
	if logo == "" {
		return FailWithMessage(http.StatusBadRequest, err.Error(), ctx)
	}
	return OkWithRaw("image/png", img, ctx)
}

func (app *AQIServer) SyncStationLog(ctx *fiber.Ctx) error {
	err := app.DB.SyncStationLogos()
	if err != nil {
		return FailWithMessage(http.StatusBadRequest, err.Error(), ctx)
	}
	return Ok(ctx)
}
