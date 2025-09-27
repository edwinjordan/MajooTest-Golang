package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func SecurityHeadersMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path

			// Longgarkan aturan kalau akses Swagger
			if strings.HasPrefix(path, "/swagger/") {
				c.Response().Header().Set("Content-Security-Policy",
					"default-src 'self'; "+
						"script-src 'self' 'unsafe-inline' 'unsafe-eval'; "+
						"style-src 'self' 'unsafe-inline'; "+ // 🔥 ini fix inline style
						"img-src 'self' data:; "+
						"font-src 'self' data:;")
			} else {
				// Ketat untuk API biasa
				c.Response().Header().Set("Content-Security-Policy", "default-src 'self'")
			}

			// Header keamanan tambahan
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")
			c.Response().Header().Set("X-Frame-Options", "DENY")
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")

			return next(c)
		}
	}
}
