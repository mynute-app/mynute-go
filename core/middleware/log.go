package middleware

import (
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/lib/Logger"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)
// Deprecated: Use LogV13 instead.
func LogV12(logger *slog.Logger) fiber.Handler {
	loki := &myLogger.Loki{}
	return func(c *fiber.Ctx) error {
		start := time.Now()
		method := c.Method()
		path := c.Path()
		ip := c.IP()
		host := string(c.Request().Host())

		labels := myLogger.Labels{
			App:    "main-api",
			Method: method,
			Path:   path,
			IP:     ip,
			Host:   host,
		}

		ReqLabels := labels.GetReqLabels()

		ReqMsg := myLogger.GetReqMessage(c)

		// Mensagem rica vai fora dos labels
		labelMsg := ReqLabels["title"] + ReqMsg

		if err := loki.LogV12(labelMsg, ReqLabels); err != nil {
			logger.Error("Failed to send log to Loki", slog.String("error", err.Error()))
		}

		err := c.Next()

		var status int

		if myError, ok := err.(lib.ErrorStruct); ok {
			status = myError.HTTPStatus
		} else {
			status = c.Response().StatusCode()
		}

		ResLabels := labels.GetResLabels(start, status)

		ResMsg := myLogger.GetResMessage(c)

		labelMsg = ResLabels["title"] + ReqMsg + ResMsg

		if err := loki.LogV12(labelMsg, ResLabels); err != nil {
			logger.Error("Failed to send log to Loki", slog.String("error", err.Error()))
		}

		return err
	}
}

// LogV13 is a middleware for logging requests and responses to Loki using the new Schema V13 format.
func LogV13(logger *slog.Logger) fiber.Handler {
	loki := &myLogger.Loki{}
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// --- Request log ---
		reqLabels := map[string]string{
			"app":   "main-api",
			"level": "info",
			"type":  "request",
		}

		reqBody := map[string]any{
			"method":  c.Method(),
			"path":    c.Path(),
			"ip":      c.IP(),
			"host":    string(c.Request().Host()),
			"headers": myLogger.MaskSensibleInformation(string(c.Request().Header.Header())),
			"body":    myLogger.MaskSensibleInformation(string(c.Request().Body())),
		}

		if err := loki.LogV13(reqLabels, reqBody); err != nil {
			logger.Error("Failed to log request to Loki", slog.String("error", err.Error()))
		}

		// --- Proceed ---
		err := c.Next()

		// --- Response log ---
		status := c.Response().StatusCode()
		if myError, ok := err.(lib.ErrorStruct); ok {
			status = myError.HTTPStatus
		}

		level := "info"
		switch {
		case status >= 500:
			level = "error"
		case status >= 400:
			level = "warning"
		case status >= 200:
			level = "success"
		}

		resLabels := map[string]string{
			"app":   "main-api",
			"level": level,
			"type":  "response",
		}

		resBody := map[string]any{
			"method":   c.Method(),
			"path":     c.Path(),
			"ip":       c.IP(),
			"host":     string(c.Request().Host()),
			"status":   status,
			"duration": time.Since(start).String(),
			"headers":  myLogger.MaskSensibleInformation(string(c.Response().Header.Header())),
			"body":     myLogger.MaskSensibleInformation(string(c.Response().Body())),
		}

		if err := loki.LogV13(resLabels, resBody); err != nil {
			logger.Error("Failed to log response to Loki", slog.String("error", err.Error()))
		}

		return err
	}
}
