package middleware

import (
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/lib/Logger"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Error(logger *slog.Logger) fiber.ErrorHandler {
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

			if err := loki.Log(ResMsg, ResLabels); err != nil {
				logger.Error("Failed to log error to Loki", slog.String("error", err.Error()))
			}
		}

		return nil
	}
}
