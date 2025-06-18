package middleware

import (
	database "agenda-kaki-go/core/config/db"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/lib/Logger"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Deprecated: Use ErrorV13 instead.
func ErrorV12(logger *slog.Logger) fiber.ErrorHandler {
	loki := &myLogger.Loki{}

	return func(c *fiber.Ctx, err error) error {
		status := fiber.StatusInternalServerError
		var responseBody []byte

		switch e := err.(type) {
		case *fiber.Error:
			status = e.Code
			responseBody = []byte(e.Message)

		case lib.ErrorStruct:
			status = e.HTTPStatus
			var marshalErr error
			responseBody, marshalErr = json.Marshal(e)
			if marshalErr != nil {
				responseBody = []byte(`{"error":"failed to marshal custom error"}`)
			}

		default:
			responseBody = []byte(err.Error())
		}

		if sendErr := c.Status(status).Send(responseBody); sendErr != nil {
			labels := myLogger.Labels{
				App:    "main-api",
				Method: c.Method(),
				Path:   c.Path(),
				IP:     c.IP(),
				Host:   string(c.Request().Host()),
			}
			ResLabels := labels.GetResLabels(time.Now(), status)
			ResLabels["level"] = "critical"

			ResMsg := myLogger.GetResMessage(c)
			ResMsg += "\nFailed to send error to client:\n" + err.Error() + "\n\n" + string(responseBody)

			if err := loki.LogV12(ResMsg, ResLabels); err != nil {
				logger.Error("Failed to log error to Loki", slog.String("error", err.Error()))
			}
		}

		return nil
	}
}

// ErrorV13 is the new version of the error handler that uses the Loki Schema V13 updated logging system.
func ErrorV13(logger *slog.Logger) fiber.ErrorHandler {
	loki := &myLogger.Loki{}

	return func(c *fiber.Ctx, err error) error {
		status := fiber.StatusInternalServerError
		var responseBody []byte

		switch e := err.(type) {
		case *fiber.Error:
			status = e.Code
			responseBody = []byte(e.Message)

		case lib.ErrorStruct:
			status = e.HTTPStatus
			var marshalErr error
			responseBody, marshalErr = json.Marshal(e)
			if marshalErr != nil {
				responseBody = []byte(`{"error":"failed to marshal custom error"}`)
			}

		default:
			responseBody = []byte(err.Error())
		}

		// Tenta enviar a resposta
		sendErr := c.Status(status).Send(responseBody)
		if sendErr != nil {
			// 🔥 Falha ao enviar resposta — log crítico
			resLabels := map[string]string{
				"app":   "main-api",
				"level": "critical",
				"type":  "response_error",
			}

			resBody := map[string]any{
				"method":   c.Method(),
				"path":     c.Path(),
				"ip":       c.IP(),
				"host":     string(c.Request().Host()),
				"status":   status,
				"error":    err.Error(),
				"response": string(responseBody),
			}

			if err := loki.LogV13(resLabels, resBody); err != nil {
				logger.Error("Failed to log critical error to Loki", slog.String("error", err.Error()))
			}
		}

		tx, _, _ := database.ContextTransaction(c)
		if tx != nil {
			if err := tx.Rollback().Error; err != nil {
				logger.Error("Failed to rollback transaction", slog.String("error", err.Error()))
			}
		}

		return nil
	}
}

