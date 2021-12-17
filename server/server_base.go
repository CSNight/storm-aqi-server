package server

import (
	"aqi-server/conf"
	"aqi-server/db"
	"aqi-server/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

type AQIServer struct {
	App *fiber.App
	Log *zap.Logger
	DB  *db.DB
	Oss *minio.Client
}

var json = jsoniter.Config{
	EscapeHTML:             false,
	SortMapKeys:            true,
	ValidateJsonRawMessage: true,
}.Froze()

var validate = validator.New()

func New(conf *conf.GConfig) *AQIServer {
	server := fiber.New(fiber.Config{
		CaseSensitive:     true,
		AppName:           "AQI-SERVER",
		ReduceMemoryUsage: true,
		JSONEncoder:       json.Marshal,
		JSONDecoder:       json.Unmarshal,
	})
	logger := middleware.Use(server, conf)

	dbEs := db.Init(conf, logger)

	ossCli, err := minio.New(conf.OssConf.Server, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.OssConf.Account, conf.OssConf.Secret, ""),
		Secure: false,
	})
	if err != nil {
		return nil
	}
	api := server.Group("/api")
	v1 := api.Group("v1")

	app := &AQIServer{
		App: server,
		Log: logger,
		DB:  dbEs,
		Oss: ossCli,
	}
	app.Register(v1)
	return app
}

func (app *AQIServer) Close() {
	app.DB.Close()
	app.Log.Info(`elasticsearch api closed`)
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func ValidateStruct(params interface{}) []*ErrorResponse {
	var errorResponses []*ErrorResponse
	err := validate.Struct(params)
	if err != nil {
		for _, errItem := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = errItem.StructNamespace()
			element.Tag = errItem.Tag()
			element.Value = errItem.Param()
			errorResponses = append(errorResponses, &element)
		}
	}
	return errorResponses
}
