package server

import (
	"aqi-server/conf"
	"aqi-server/elastic"
	"aqi-server/middleware"
	"context"
	"github.com/gofiber/fiber/v2"
	pool "github.com/jolestar/go-commons-pool/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

type AQIServer struct {
	App  *fiber.App
	Log  *zap.Logger
	Es   *elastic.EsAPI
	pool *pool.ObjectPool
	Oss  *minio.Client
}

var ctx = context.Background()

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
	poolEs := elastic.InitEsPool(ctx, conf.ESConf)

	elasticApi := &elastic.EsAPI{
		Log:       logger.Sugar().Named("\u001B[33m[ESClient]\u001B[0m"),
		EsPool:    poolEs,
		FailQueue: []elastic.BulkIndexerItem{},
	}
	elasticApi.Init()

	ossCli, err := minio.New(conf.OssConf.Server, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.OssConf.Account, conf.OssConf.Secret, ""),
		Secure: false,
	})

	if err != nil {
		return nil
	}

	return &AQIServer{
		App:  server,
		Log:  logger,
		Es:   elasticApi,
		pool: poolEs,
		Oss:  ossCli,
	}
}

func (app *AQIServer) Close() {
	app.Es.Close()
	app.pool.Close(ctx)
	app.Log.Info(`elasticsearch api closed`)
}
