package db

import (
	"aqi-server/conf"
	"aqi-server/elastic"
	"aqi-server/tools"
	"context"
	"github.com/coocood/freecache"
	pool "github.com/jolestar/go-commons-pool/v2"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"os"
	"runtime"
	"time"
)

type DB struct {
	Conf  *conf.AQIConfig
	api   *elastic.EsAPI
	ipDB  *tools.Reader
	pool  *pool.ObjectPool
	log   *zap.Logger
	cache *freecache.Cache
	ctx   context.Context
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
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	ipDB, err := tools.Open(pwd + string(os.PathSeparator) + "assets" + string(os.PathSeparator) + "GeoLite2-City.mmdb")
	if err != nil {
		return nil, err
	}

	cache := freecache.NewCache(20 * 1024 * 1024)

	return &DB{
		Conf:  conf.AQIConf,
		api:   elasticApi,
		ipDB:  ipDB,
		pool:  poolEs,
		cache: cache,
		log:   logger.Named("\u001B[33m[DB]\u001B[0m"),
		ctx:   ctx,
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
