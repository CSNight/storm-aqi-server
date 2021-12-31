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
	manager := freecache.NewCache(500 * 1024 * 1024)
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
			c.Response().SetBodyRaw(e[2:])
			c.Response().SetStatusCode(http.StatusOK)
			if cfg.Compress && e[1] > 0 {
				c.Response().Header.SetBytesV(fiber.HeaderContentEncoding, []byte(getEncoding(e[1])))
			}
			c.Response().Header.SetContentTypeBytes([]byte(getContentType(e[0])))
			c.Set(cfg.CacheHeader, "hit")
			return nil
		}
		if err = c.Next(); err != nil {
			return err
		}
		if c.Response().StatusCode() == http.StatusOK {
			body := utils.CopyBytes(c.Response().Body())
			if len(body) > 0 {
				ct := utils.CopyBytes(c.Response().Header.ContentType())
				encoding := utils.CopyBytes(c.Response().Header.Peek(fiber.HeaderContentEncoding))
				var cacheBytes = make([]byte, 2)
				ctB := getContentTypeByte(ct)
				if ctB != 0 {
					cacheBytes[0] = ctB
					cacheBytes[1] = getEncodingByte(encoding)
					cacheBytes = append(cacheBytes, body...)
					_ = manager.Set([]byte(key), cacheBytes, cfg.Expiration)
					cacheBytes = nil
				}
			}
			body = nil
		}
		return nil
	}
}

func getContentType(ct uint8) string {
	switch ct {
	case 1:
		return "text/xml"
	case 2:
		return "text/html"
	case 3:
		return "text/plain"
	case 4:
		return "application/xml"
	case 5:
		return "application/json"
	case 6:
		return "application/javascript"
	case 7:
		return "application/x-www-form-urlencoded"
	case 8:
		return "application/octet-stream"
	case 9:
		return "multipart/form-data"
	case 11:
		return "text/xml; charset=utf-8"
	case 12:
		return "text/html; charset=utf-8"
	case 13:
		return "text/plain; charset=utf-8"
	case 14:
		return "application/xml; charset=utf-8"
	case 15:
		return "application/json; charset=utf-8"
	case 16:
		return "application/javascript; charset=utf-8"
	case 17:
		return "application/x-www-form-urlencoded; charset=utf-8"
	default:
		return "text/plain"
	}

}

func getEncoding(encoding uint8) string {
	switch encoding {
	case 1:
		return "br"
	case 2:
		return "gzip"
	case 3:
		return "deflate"
	default:
		return ""
	}
}

func getContentTypeByte(ct []byte) uint8 {
	cts := string(ct)
	switch cts {
	case "text/xml":
		return 1
	case "text/html":
		return 2
	case "text/plain":
		return 3
	case "application/xml":
		return 4
	case "application/json":
		return 5
	case "application/javascript":
		return 6
	case "application/x-www-form-urlencoded":
		return 7
	case "application/octet-stream":
		return 8
	case "multipart/form-data":
		return 9
	case "text/xml; charset=utf-8":
		return 11
	case "text/html; charset=utf-8":
		return 12
	case "text/plain; charset=utf-8":
		return 13
	case "application/xml; charset=utf-8":
		return 14
	case "application/json; charset=utf-8":
		return 15
	case "application/javascript; charset=utf-8":
		return 16
	case "application/x-www-form-urlencoded; charset=utf-8":
		return 17
	default:
		return 0
	}
}

func getEncodingByte(encoding []byte) uint8 {
	cts := string(encoding)
	switch cts {
	case "br":
		return 1
	case "gzip":
		return 2
	case "deflate":
		return 3
	default:
		return 0
	}
}
