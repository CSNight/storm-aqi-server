package middleware

import (
	"aqi-server/conf"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	rcp "github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

func Use(server *fiber.App, config *conf.GConfig) *zap.Logger {

	logger := InitLogger(config.LogConf)

	server.Use(rcp.New())

	server.Use(New(LogConfig{
		Next:     nil,
		Logger:   logger,
		Fields:   []string{"ip", "port", "url", "method", "status", "latency", "queryParams", "body"},
		Messages: []string{"Server error", "Client error", "Success"},
	}))

	server.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token,X-Token,X-User-Id",
		AllowCredentials: true,
	}))

	server.Use(NewCache(CacheConfig{
		Expiration:  300,
		Compress:    config.AppConf.EnableCompress,
		CacheHeader: "X-Cache-Storm",
	}))

	if config.AppConf.EnableCompress {
		server.Use(compress.New(compress.Config{
			Level: compress.LevelBestSpeed, // 1
		}))
	}

	server.Get("/monitor", monitor.New())

	return logger
}
