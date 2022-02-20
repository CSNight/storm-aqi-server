package server

import (
	"github.com/csnight/storm-aqi-server/conf"
	"github.com/csnight/storm-aqi-server/db"
	"github.com/csnight/storm-aqi-server/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type AQIServer struct {
	App *fiber.App
	Log *zap.Logger
	DB  *db.DB
}

var json = jsoniter.Config{
	EscapeHTML:             false,
	SortMapKeys:            true,
	ValidateJsonRawMessage: true,
}.Froze()

var validate = validator.New()

func New(conf *conf.GConfig) (*AQIServer, error) {
	server := fiber.New(fiber.Config{
		Views:             html.New("./assets", ".html"),
		CaseSensitive:     true,
		AppName:           "AQI-SERVER",
		ReduceMemoryUsage: true,
		JSONEncoder:       json.Marshal,
		JSONDecoder:       json.Unmarshal,
	})
	logger := middleware.Use(server, conf)

	dbEs, err := db.Init(conf, logger)
	if err != nil {
		return nil, err
	}
	dbEs.RefreshCache()
	server.Static("/static", "./assets/static")
	api := server.Group("/api")
	v1 := api.Group("v1")

	app := &AQIServer{
		App: server,
		Log: logger,
		DB:  dbEs,
	}
	app.Register(v1)
	return app, nil
}

func (app *AQIServer) Close() {
	app.DB.Close()
	app.Log.Info(`elasticsearch api closed`)
}

type ErrorResponse struct {
	FailedField string
	Rule        string
	ErrValue    interface{}
}

func ValidateStruct(params interface{}) []*ErrorResponse {
	var errorResponses []*ErrorResponse
	err := validate.Struct(params)
	if err != nil {
		for _, errItem := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = errItem.StructNamespace()
			if errItem.Param() == "" {
				element.Rule = errItem.Tag()
			} else {
				element.Rule = errItem.Tag() + "=" + errItem.Param()
			}
			element.ErrValue = errItem.Value()
			errorResponses = append(errorResponses, &element)
		}
	}
	return errorResponses
}

func ValidateVar(param interface{}, tag string) []*ErrorResponse {
	var errorResponses []*ErrorResponse
	err := validate.Var(param, tag)
	if err != nil {
		for _, errItem := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = errItem.StructNamespace()
			if errItem.Param() == "" {
				element.Rule = errItem.Tag()
			} else {
				element.Rule = errItem.Tag() + "=" + errItem.Param()
			}
			element.ErrValue = errItem.Value()
			errorResponses = append(errorResponses, &element)
		}
	}
	return errorResponses
}
