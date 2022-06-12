package db

import (
	"context"
	"github.com/coocood/freecache"
	"github.com/csnight/storm-aqi-server/conf"
	"github.com/csnight/storm-aqi-server/elastic"
	pool "github.com/jolestar/go-commons-pool/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
	"runtime"
	"time"
)

type DB struct {
	Conf  *conf.AQIConfig
	api   *elastic.EsAPI
	pool  *pool.ObjectPool
	log   *zap.Logger
	cache *freecache.Cache
	ctx   context.Context
	oss   *minio.Client
}

var json = jsoniter.Config{
	EscapeHTML:             false,
	SortMapKeys:            true,
	ValidateJsonRawMessage: true,
}.Froze()

var tick = time.NewTicker(time.Minute * 9)

func Init(conf *conf.GConfig, logger *zap.Logger) (*DB, error) {
	var ctx = context.Background()
	poolEs := elastic.InitEsPool(ctx, conf.ESConf)

	elasticApi := &elastic.EsAPI{
		Log:       logger.Sugar().Named("\u001B[33m[ES]\u001B[0m"),
		EsPool:    poolEs,
		FailQueue: []elastic.BulkIndexerItem{},
	}
	elasticApi.Init()

	ossCli, err := minio.New(conf.OssConf.Server, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.OssConf.Account, conf.OssConf.Secret, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	cache := freecache.NewCache(20 * 1024 * 1024)

	return &DB{
		Conf:  conf.AQIConf,
		api:   elasticApi,
		pool:  poolEs,
		cache: cache,
		log:   logger.Named("\u001B[33m[db]\u001B[0m"),
		ctx:   ctx,
		oss:   ossCli,
	}, nil
}

func (db *DB) RefreshCache() {
	db.loadStations()
	go func() {
		for {
			select {
			case <-tick.C:
				db.loadStations()
			}
		}
	}()
}

func (db *DB) loadStations() {
	stations, err := db.GetAllStations()
	if err != nil {
		db.log.Error("refresh stations cache error:", zap.String("err", err.Error()))
		return
	}
	db.cache.Clear()
	for _, st := range stations {
		stb, err := json.Marshal(st)
		if err != nil {
			continue
		}
		_ = db.cache.Set([]byte(st.Sid), stb, 600)
	}
	defer func() {
		stations = nil
		runtime.GC()
	}()
	db.log.Info("refresh stations cache success")
}

func (db *DB) Close() {
	tick.Stop()
	db.api.Close()
	db.pool.Close(db.ctx)
}
