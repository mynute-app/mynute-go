package middleware

import (
	"agenda-kaki-go/core/lib"
	"fmt"
	"log/slog"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
)

// LoggerMiddleware logs request and response details without modifying the response.
func Logger(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		method := c.Method()
		path := c.Path()
		ip := c.IP()
		reqBody := maskSensibleInformation(string(c.Request().Body()))
		reqHeaders := maskSensibleInformation(string(c.Request().Header.Header()))

		loggerDefaultMap := map[string]any{
			"slog": []slog.Attr{
				slog.String("method", method),
				slog.String("path", path),
				slog.String("ip", ip),
				slog.String("request body", reqBody),
				slog.String("request header", reqHeaders),
			},
		}

		labelMsg := "Incoming request"

		lokiDefaultMap := map[string]string{
			"app":        "main-api",
			"level":      "info",
			"type":       "request",
			"method":     method,
			"path":       path,
			"ip":         ip,
			"req_body":   reqBody,
			"req_header": reqHeaders,
		}

		logger.Info(labelMsg, loggerDefaultMap["slog"])

		_ = lib.SendLogToLoki(labelMsg, lokiDefaultMap)

		err := c.Next()

		duration := time.Since(start)

		resStatus := c.Response().Header.StatusCode()
		resBody := maskSensibleInformation(string(c.Response().Body()))
		resHeaders := maskSensibleInformation(string(c.Response().Header.Header()))

		loggerDefaultMap["slog"] = append(loggerDefaultMap["slog"].([]slog.Attr),
			slog.Int("status", resStatus),
			slog.String("response body", resBody),
			slog.String("response header", resHeaders),
			slog.Duration("duration", duration),
		)

		lokiDefaultMap["status"] = fmt.Sprintf("%d", resStatus)
		lokiDefaultMap["res_body"] = resBody
		lokiDefaultMap["res_header"] = resHeaders
		lokiDefaultMap["duration"] = duration.String()

		if err != nil {
			loggerDefaultMap["slog"] = append(loggerDefaultMap["slog"].([]slog.Attr),
				slog.String("error", err.Error()),
				slog.String("stack", fmt.Sprintf("%+v", err)),
			)

			lokiDefaultMap["type"] = "response"
			lokiDefaultMap["level"] = "error"
			lokiDefaultMap["error"] = err.Error()
			lokiDefaultMap["stack"] = fmt.Sprintf("%+v", err)

			labelMsg = "Request error!"
		} else {
			lokiDefaultMap["level"] = "info"
			lokiDefaultMap["type"] = "response"

			labelMsg = "Resquest success!"
		}

		logger.Error(labelMsg, loggerDefaultMap["slog"])
		_ = lib.SendLogToLoki(labelMsg, lokiDefaultMap)

		return err
	}
}

func maskSensibleInformation(body string) string {
	patterns := []string{
		`("password"\s*:\s*")([^"]*)(")`,
		`("token"\s*:\s*")([^"]*)(")`,
		`("secret"\s*:\s*")([^"]*)(")`,
		`("access_token"\s*:\s*")([^"]*)(")`,
		`("refresh_token"\s*:\s*")([^"]*)(")`,
		`("authorization"\s*:\s*")([^"]*)(")`,
	}

	masked := body
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		masked = re.ReplaceAllString(masked, `$1***masked***$3`)
	}

	return masked
}
