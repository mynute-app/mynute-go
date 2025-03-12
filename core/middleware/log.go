package middleware

import (
	"bytes"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

// LoggerMiddleware logs request and response details without modifying the response.
func LoggerMiddleware(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Capture request headers and body
		reqHeaders := c.GetReqHeaders()
		reqBody := string(c.Body())

		// Capture response body using a buffer
		resBodyBuffer := new(bytes.Buffer)

		// Clone response writer to capture response body without modifying it
		originalBody := c.Response().Body()
		resBodyBuffer.Write(originalBody)

		// Process the request
		err := c.Next()

		// Measure execution time
		duration := time.Since(start)

		// Capture response headers and body again after the request is processed
		resHeaders := getResponseHeaders(c)
		resBody := resBodyBuffer.String()

		// Log structured data
		logger.Info("API Request",
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"duration_ms", duration.Milliseconds(),
			"ip", c.IP(),
			"error", errorString(err),
			"request_headers", reqHeaders,
			"request_body", reqBody,
			"response_headers", resHeaders,
			"response_body", resBody,
		)

		return err
	}
}

// getResponseHeaders extracts response headers.
func getResponseHeaders(c *fiber.Ctx) map[string]string {
	headers := make(map[string]string)
	c.Response().Header.VisitAll(func(key, value []byte) {
		headers[string(key)] = string(value)
	})
	return headers
}

// errorString converts an error to a string (or returns "none" if nil).
func errorString(err error) string {
	if err != nil {
		return err.Error()
	}
	return "none"
}