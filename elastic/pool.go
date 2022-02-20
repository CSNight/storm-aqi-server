package elastic

import (
	"context"
	"errors"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/jolestar/go-commons-pool/v2"
	"net/http"
	"storm-aqi-server/conf"
	"time"
)

type PoolFactory struct {
	Uris              []string
	Username          string
	Password          string
	EnableDebugLogger bool
	MaxRetries        int
}

func (f *PoolFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	cfg := elasticsearch.Config{
		Addresses:         f.Uris,
		Username:          f.Username,
		Password:          f.Password,
		EnableDebugLogger: f.EnableDebugLogger,
		MaxRetries:        f.MaxRetries,
		Transport: &http.Transport{
			DisableKeepAlives:     false,
			DisableCompression:    false,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   10,
			MaxConnsPerHost:       0,
			IdleConnTimeout:       time.Second * 20,
			ResponseHeaderTimeout: 0,
			ExpectContinueTimeout: 0,
			Proxy:                 http.ProxyFromEnvironment,
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return pool.NewPooledObject(es), nil
}

func (f *PoolFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	object.Object = nil
	return nil
}

func (f *PoolFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	if object.Object == nil {
		return false
	}
	cli := object.Object.(*elasticsearch.Client)
	if cli != nil {
		ping, err := cli.Ping()
		if err != nil || ping.IsError() {
			return false
		}
	}
	return true
}

func (f *PoolFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	if object.Object == nil {
		return errors.New("empty pool object")
	}
	return nil
}

func (f *PoolFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}

func InitEsPool(ctx context.Context, conf *conf.ESConfig) *pool.ObjectPool {
	factory := &PoolFactory{
		Uris:              conf.Uri,
		Username:          conf.Username,
		Password:          conf.Password,
		EnableDebugLogger: conf.EnableDebugLogger,
		MaxRetries:        conf.MaxRetries,
	}
	config := &pool.ObjectPoolConfig{
		LIFO:                    conf.LIFO,
		MaxTotal:                conf.MaxTotal,
		MaxIdle:                 conf.MaxIdle,
		MinIdle:                 conf.MinIdle,
		TestOnCreate:            conf.TestOnCreate,
		TestOnBorrow:            conf.TestOnBorrow,
		TestOnReturn:            conf.TestOnReturn,
		TestWhileIdle:           conf.TestWhileIdle,
		BlockWhenExhausted:      conf.BlockWhenExhausted,
		NumTestsPerEvictionRun:  conf.NumTestsPerEvictionRun,
		TimeBetweenEvictionRuns: time.Duration(conf.TimeBetweenEvictionRuns * int(time.Second)),
		EvitionContext:          context.Background(),
	}
	abConfig := &pool.AbandonedConfig{
		RemoveAbandonedOnBorrow:      conf.RemoveAbandonedOnBorrow,
		RemoveAbandonedOnMaintenance: conf.RemoveAbandonedOnMaintenance,
		RemoveAbandonedTimeout:       time.Duration(conf.RemoveAbandonedTimeout * int(time.Second)),
	}

	objectPool := pool.NewObjectPoolWithAbandonedConfig(ctx, factory, config, abConfig)
	return objectPool
}
