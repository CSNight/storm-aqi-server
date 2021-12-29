package server

import (
	"aqi-server/conf"
	"aqi-server/db"
	"aqi-server/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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
	Tag         string
	Value       interface{}
}

func ValidateStruct(params interface{}) []*ErrorResponse {
	var errorResponses []*ErrorResponse
	err := validate.Struct(params)
	if err != nil {
		for _, errItem := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = errItem.StructNamespace()
			element.Tag = errItem.Tag()
			element.Value = errItem.Value()
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
			element.Tag = errItem.Tag()
			element.Value = errItem.Value()
			errorResponses = append(errorResponses, &element)
		}
	}
	return errorResponses
}
