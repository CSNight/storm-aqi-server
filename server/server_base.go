package server

import (
	"aqi-server/conf"
	"aqi-server/db"
	"aqi-server/middleware"
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

	return &AQIServer{
		App: server,
		Log: logger,
		DB:  dbEs,
		Oss: ossCli,
	}
}

func (app *AQIServer) Close() {
	app.DB.Close()
	app.Log.Info(`elasticsearch api closed`)
}
