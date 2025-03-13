package lib

import (
	"bytes"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// Setup logger with JSON formatting
var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true, // Makes JSON logs more readable
	})
	log.SetLevel(logrus.InfoLevel)
}

func ApiLog(c *fiber.Ctx) error {
	start := time.Now()

	// Call next middleware
	err := c.Next()

	// Calculate response time
	duration := time.Since(start).Milliseconds()

	// Read response body (if needed)
	body := c.Response().Body()

	// Read request headers
	requestHeaders := extractHeaders(c.Request().Header.String())

	// Create a structured log entry
	log.WithFields(logrus.Fields{
		"time":            time.Now().Format(time.RFC3339),
		"method":          c.Method(),
		"path":            c.OriginalURL(),
		"ip":              c.IP(),
		"duration_ms":     duration,
		"status":          c.Response().StatusCode(),
		"request_headers": requestHeaders,
		"request_body":    maskSensitiveData(string(c.Request().Body())),
		"response_body":   string(body),
	}).Info("API Request Processed")

	return err
}

// Mask sensitive fields in logs
func maskSensitiveData(body string) string {
	if bytes.Contains([]byte(body), []byte("password")) {
		return "{masked sensitive data}"
	}
	return body
}

// Extract headers cleanly
func extractHeaders(headerStr string) string {
	lines := bytes.Split([]byte(headerStr), []byte("\n"))
	var filtered [][]byte
	for _, line := range lines {
		if len(bytes.TrimSpace(line)) > 0 {
			filtered = append(filtered, bytes.TrimSpace(line))
		}
	}
	return string(bytes.Join(filtered, []byte(", ")))
}
