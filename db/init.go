package db

import (
	"aqi-server/conf"
	"aqi-server/elastic"
	"context"
	pool "github.com/jolestar/go-commons-pool/v2"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type DB struct {
	Conf *conf.AQIConfig
	api  *elastic.EsAPI
	pool *pool.ObjectPool
	log  *zap.Logger
}

var ctx = context.Background()
var json = jsoniter.Config{
	EscapeHTML:             false,
	SortMapKeys:            true,
	ValidateJsonRawMessage: true,
}.Froze()

func Init(conf *conf.GConfig, logger *zap.Logger) *DB {
	poolEs := elastic.InitEsPool(ctx, conf.ESConf)

	elasticApi := &elastic.EsAPI{
		Log:       logger.Sugar().Named("\u001B[33m[DB]\u001B[0m"),
		EsPool:    poolEs,
		FailQueue: []elastic.BulkIndexerItem{},
	}
	elasticApi.Init()
	return &DB{
		Conf: conf.AQIConf,
		api:  elasticApi,
		pool: poolEs,
		log:  logger,
	}
}

func (db *DB) Close() {
	db.api.Close()
	db.pool.Close(ctx)
}