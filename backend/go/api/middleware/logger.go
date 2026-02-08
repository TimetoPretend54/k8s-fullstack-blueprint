package middleware

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

// Logger returns a middleware that logs HTTP requests
func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			latency := time.Since(start)

			status := c.Response().Status
			method := c.Request().Method
			path := c.Path()
			userAgent := c.Request().UserAgent()

			log.Printf("[%s] %s %d %s %s", method, path, status, latency, userAgent)

			return err
		}
	}
}
