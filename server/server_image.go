package server

import (
	"net/http"
	"time"

	"github.com/csnight/storm-aqi-server/tools"
	"github.com/gofiber/fiber/v2"
)

type ImageRequest struct {
	Time string `json:"time" validate:"required,datetime=2006-01-02T15:04:05Z"`
	Pol  string `json:"pol" validate:"required,oneof=no2 pm25 pm10 co so2 o3 dust pmFRP"`
}

func (app *AQIServer) ImageGet(ctx *fiber.Ctx) error {
	var query ImageRequest
	err := ctx.QueryParser(&query)
	if err != nil {
		return FailWithMessage(http.StatusBadRequest, "can't parser params", ctx)
	}
	errResp := ValidateStruct(query)
	if errResp != nil {
		return FailWithDetailed(http.StatusBadRequest, errResp, "", ctx)
	}
	resp, err := app.db.GetImage(query.Time, query.Pol)
	if err != nil {
		return FailWithDetailed(http.StatusBadRequest, err, "", ctx)
	}

	return OkWithData(resp, ctx)
}

func (app *AQIServer) ImageDownload(ctx *fiber.Ctx) error {
	resp, err := app.db.DownloadImage(ctx.Params("dir"), ctx.Params("file"))
	if err != nil {
		return err
	}
	etag, err := tools.NewNanoId()
	if err != nil {
		return err
	}
	ctx.Set("Cache-Control", "max-age=7200")
	ctx.Set("ETag", etag)
	ctx.Set("Last-Modified", time.Now().Format(time.RFC1123))
	return OkWithRaw("image/png", resp, ctx)
}
