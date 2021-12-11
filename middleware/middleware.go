package middleware

import (
	"aqi-server/conf"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/utils"
	"go.uber.org/zap"
	"time"
)

func Use(server *fiber.App, config *conf.GConfig) *zap.Logger {

	logger := InitLogger(config.LogConf)

	server.Use(New(LogConfig{
		Next:     nil,
		Logger:   logger,
		Fields:   []string{"ip", "port", "path", "method", "status", "latency"},
		Messages: []string{"Server error", "Client error", "Success"},
	}))

	server.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-User-Id",
		AllowCredentials: true,
	}))

	server.Use(csrf.New(csrf.Config{
		KeyLookup:      "header:X-Csrf-Token",
		CookieName:     "csrf_",
		CookieSameSite: "Lax",
		Expiration:     1 * time.Minute,
		KeyGenerator:   utils.UUID,
	}))

	server.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))

	return logger
}
