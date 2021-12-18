package db

import (
	"aqi-server/conf"
	"aqi-server/elastic"
	"context"
	pool "github.com/jolestar/go-commons-pool/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/oschwald/geoip2-golang"
	"go.uber.org/zap"
	"os"
)

type DB struct {
	Conf *conf.AQIConfig
	api  *elastic.EsAPI
	ipDB *geoip2.Reader
	pool *pool.ObjectPool
	log  *zap.SugaredLogger
	ctx  context.Context
}

var json = jsoniter.Config{
	EscapeHTML:             false,
	SortMapKeys:            true,
	ValidateJsonRawMessage: true,
}.Froze()

func Init(conf *conf.GConfig, logger *zap.Logger) (*DB, error) {
	var ctx = context.Background()
	poolEs := elastic.InitEsPool(ctx, conf.ESConf)

	elasticApi := &elastic.EsAPI{
		Log:       logger.Sugar().Named("\u001B[33m[ES]\u001B[0m"),
		EsPool:    poolEs,
		FailQueue: []elastic.BulkIndexerItem{},
	}
	elasticApi.Init()
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	ipDB, err := geoip2.Open(pwd + string(os.PathSeparator) + "assets" + string(os.PathSeparator) + "GeoLite2-City.mmdb")
	if err != nil {
		return nil, err
	}
	return &DB{
		Conf: conf.AQIConf,
		api:  elasticApi,
		ipDB: ipDB,
		pool: poolEs,
		log:  logger.Sugar().Named("[DB]"),
		ctx:  ctx,
	}, nil
}

func (db *DB) Close() {
	db.api.Close()
	db.pool.Close(db.ctx)
}
