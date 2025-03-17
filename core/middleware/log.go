package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

// LoggerMiddleware logs request and response details without modifying the response.
func Logger(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now() // Capture start time

		// Log the request details
		logger.Info("Incoming request",
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.String("ip", c.IP()),
		)

		// Process the request
		err := c.Next()

		// Log the response details
		if err != nil {
			logger.Error("Error processing request",
				slog.String("method", c.Method()),
				slog.String("path", c.Path()),
				slog.String("ip", c.IP()),
				slog.String("error", err.Error()),
				slog.Duration("duration", time.Since(start)),
			)
		} else {
			logger.Info("Response sent",
				slog.String("method", c.Method()),
				slog.String("path", c.Path()),
				slog.String("ip", c.IP()),
				slog.Int("status", c.Response().StatusCode()),
				slog.Duration("duration", time.Since(start)),
			)
		}
		return err
	}
}
