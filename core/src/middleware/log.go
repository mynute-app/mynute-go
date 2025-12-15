package middleware

import (
	"log/slog"
	"mynute-go/core/src/lib"
	myLogger "mynute-go/core/src/lib/metrics"
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

		var http_status int

		if myError, ok := err.(lib.ErrorStruct); ok {
			http_status = myError.HTTPStatus
		} else {
			http_status = c.Response().StatusCode()
		}

		ResLabels := labels.GetResLabels(start, http_status)

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
	const maxBodySize = 1024 * 1024 // 1MB - max size for request/response bodies

	return func(c *fiber.Ctx) error {
		start := time.Now()

		// --- Request log ---
		reqStreamLabels := map[string]string{
			"app":   "main-api",
			"level": "info",
			"type":  "request",
		}

		host := string(c.Request().Host())
		ip := c.IP()
		path := c.Path()
		method := c.Method()

		reqBody := myLogger.MaskSensibleInformation(string(c.Request().Body()))
		reqHeaders := myLogger.MaskSensibleInformation(string(c.Request().Header.Header()))

		// Truncate large request bodies to prevent Loki errors
		reqTruncated := false
		if len(reqBody) > maxBodySize {
			reqBody = reqBody[:maxBodySize] + "... [TRUNCATED - request too large]"
			reqTruncated = true
		}

		reqBodyLabels := map[string]any{
			"method":        method,
			"path":          path,
			"ip":            ip,
			"host":          host,
			"req_headers":   reqBody,
			"req_body":      reqHeaders,
			"req_truncated": reqTruncated,
		}

		// Log to Loki asynchronously to avoid blocking
		go func() {
			if err := loki.LogV13(reqStreamLabels, reqBodyLabels); err != nil {
				logger.Error("Failed to log request to Loki", slog.String("error", err.Error()))
			}
		}()

		// --- Proceed ---
		err := c.Next()

		// --- Response log ---
		http_status := c.Response().StatusCode()
		if myError, ok := err.(lib.ErrorStruct); ok {
			http_status = myError.HTTPStatus
		}

		level := "info"
		switch {
		case http_status >= 500:
			level = "error"
		case http_status >= 400:
			level = "warning"
		case http_status >= 200:
			level = "success"
		}

		resHeaders := myLogger.MaskSensibleInformation(string(c.Response().Header.Header()))
		resBody := myLogger.MaskSensibleInformation(string(c.Response().Body()))

		// Truncate large response bodies to prevent Loki errors
		resTruncated := false
		if len(resBody) > maxBodySize {
			resBody = resBody[:maxBodySize] + "... [TRUNCATED - response too large]"
			resTruncated = true
		}

		resStreamLabels := map[string]string{
			"app":   "main-api",
			"level": level,
			"type":  "response",
		}

		resBodyLabels := map[string]any{
			"method":        method,
			"path":          path,
			"ip":            ip,
			"host":          host,
			"http_status":   http_status,
			"duration":      time.Since(start).String(),
			"req_body":      reqBody,
			"req_headers":   reqHeaders,
			"res_headers":   resHeaders,
			"res_body":      resBody,
			"req_truncated": reqTruncated,
			"res_truncated": resTruncated,
		}

		// Log to Loki asynchronously to avoid blocking
		go func() {
			if err := loki.LogV13(resStreamLabels, resBodyLabels); err != nil {
				logger.Error("Failed to log response to Loki", slog.String("error", err.Error()))
			}
		}()

		return err
	}
}
