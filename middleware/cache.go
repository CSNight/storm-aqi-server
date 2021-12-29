package middleware

import (
	"github.com/coocood/freecache"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"net/http"
	"time"
)

type CacheConfig struct {
	Expiration  int
	CacheHeader string
	Compress    bool
}

// NewCache creates a new cache handler
func NewCache(cfg CacheConfig) fiber.Handler {

	// Nothing to cache
	if cfg.Expiration*int(time.Second) < 0 {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}
	manager := freecache.NewCache(100 * 1024 * 1024)
	// Return new handler
	return func(c *fiber.Ctx) error {
		// Only cache GET methods
		if c.Method() != fiber.MethodGet {
			c.Set(cfg.CacheHeader, "unreachable")
			return c.Next()
		}

		// Get key from request
		key := utils.CopyString(c.OriginalURL())

		// Get entry from pool
		e, err := manager.Get([]byte(key))
		if err == nil {
			// Set response headers from cache
			c.Response().SetBodyRaw(e)
			c.Response().SetStatusCode(http.StatusOK)
			if cfg.Compress {
				c.Response().Header.SetBytesV(fiber.HeaderContentEncoding, []byte("br"))
			}
			c.Response().Header.SetContentTypeBytes([]byte(fiber.MIMEApplicationJSON))
			c.Set(cfg.CacheHeader, "hit")
			return nil
		}
		// Continue stack, return err to Fiber if exist
		if err = c.Next(); err != nil {
			return err
		}
		if c.Response().StatusCode() == http.StatusOK {
			body := utils.CopyBytes(c.Response().Body())
			if len(body) > 0 {
				_ = manager.Set([]byte(key), body, cfg.Expiration)
			}
		}
		c.Set(cfg.CacheHeader, "miss")
		// Finish response
		return nil
	}
}
