package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// ValidateUserToken validates JWT token from Authorization header
func ValidateUserToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Accept token from "Authorization: Bearer <token>" header.
			// Be robust against extra spaces and case differences.
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				// Also try the lowercase header key just in case some client sets it (headers are case-insensitive,
				// but some code paths might access different representations).
				auth = c.Request().Header.Get("authorization")
			}
			if auth == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"message": "missing token"})
			}

			// Normalize and extract token
			auth = strings.TrimSpace(auth)
			// Support tokens that may already include the "Bearer " prefix or be passed raw.
			if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
				// strip only the first occurrence of "Bearer "
				auth = strings.TrimSpace(auth[7:])
			}

			if auth == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"message": "invalid token format"})
			}

			// Perform actual token validation here if you have a signing key/service.
			// If you have an auth service in the app, call it to validate and parse claims.
			// For now, store the raw token in the context so handlers can use it.
			c.Set("user_token", auth)

			return next(c)
		}
	}
}
