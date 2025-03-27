package middleware

import (
	"agenda-kaki-go/core/lib"
	"fmt"
	"log/slog"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Logs request and response details without modifying them.
func Log(logger *slog.Logger) fiber.Handler {
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

		// logger.Info(labelMsg, slog.Any("req_info", loggerDefaultMap["slog"]))

		if err := lib.SendLogToLoki(labelMsg, lokiDefaultMap); err != nil {
			logger.Error("Failed to send log to Loki", slog.String("error", err.Error()))
		}

		origin_err := c.Next()

		resStatus := 0

		duration := time.Since(start)

		resBody := maskSensibleInformation(string(c.Response().Body()))
		resHeaders := maskSensibleInformation(string(c.Response().Header.Header()))

		loggerDefaultMap["slog"] = append(loggerDefaultMap["slog"].([]slog.Attr),
			slog.String("response body", resBody),
			slog.String("response header", resHeaders),
			slog.Duration("duration", duration),
		)

		lokiDefaultMap["res_header"] = resHeaders
		lokiDefaultMap["duration"] = duration.String()

		if origin_err != nil {
			loggerDefaultMap["slog"] = append(loggerDefaultMap["slog"].([]slog.Attr),
				slog.String("error", origin_err.Error()),
				slog.String("stack", fmt.Sprintf("%+v", origin_err)),
			)

			lokiDefaultMap["type"] = "response"
			lokiDefaultMap["level"] = "error"
			lokiDefaultMap["error"] = origin_err.Error()

			labelMsg = "Request error!"
			if e, ok := origin_err.(lib.ErrorStruct); ok {
				loggerDefaultMap["slog"] = append(loggerDefaultMap["slog"].([]slog.Attr),
					slog.String("inner_error", fmt.Sprintf("%+v", e.InnerError)),
				)
				lokiDefaultMap["inner_error"] = fmt.Sprintf("%+v", e.InnerError)
				resStatus = e.HTTPStatus
				resBody = e.ToJSON()
			} else {
				resStatus = fiber.ErrInternalServerError.Code
			}
		} else {
			lokiDefaultMap["level"] = "info"
			lokiDefaultMap["type"] = "response"
			resStatus = c.Response().Header.StatusCode()
			if resStatus == 401 {
				lokiDefaultMap["level"] = "warning"
				labelMsg = "Request unauthorized!"
			} else {
				labelMsg = "Request success!"
			}
		}

		lokiDefaultMap["res_body"] = resBody

		lokiDefaultMap["status_code"] = fmt.Sprintf("%d", resStatus)
		loggerDefaultMap["slog"] = append(loggerDefaultMap["slog"].([]slog.Attr),
			slog.Int("status_code", resStatus),
		)

		// logger.Info(labelMsg, slog.Any("req_info", loggerDefaultMap["slog"]))
		if err := lib.SendLogToLoki(labelMsg, lokiDefaultMap); err != nil {
			logger.Error("Failed to send log to Loki", slog.String("error", err.Error()))
		}

		return origin_err
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
