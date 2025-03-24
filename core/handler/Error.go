package handler

import (
	"agenda-kaki-go/core/lib"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func Error(logger *slog.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		if e, ok := err.(lib.ErrorStruct); ok {
			if err := e.SendToClient(c); err != nil {
				if err := lib.SendLogToLoki("Failed to send error response to client!", map[string]string{
					"app":        "main-api",
					"level":      "error",
					"type":       "response",
					"method":     c.Method(),
					"path":       c.Path(),
					"ip":         c.IP(),
					"req_body":   string(c.Request().Body()),
					"req_header": string(c.Request().Header.Header()),
					"res_body":   string(c.Response().Body()),
					"res_header": string(c.Response().Header.Header()),
					"error":      e.WithError(err).Error(),
				}); err != nil {
					logger.Error("Failed to send error response log to Loki", slog.String("error", err.Error()))
					return nil
				}
				return nil
			}
		}

		MyErr := lib.ErrorStruct{
			DescriptionEn: "Internal Server Error",
			DescriptionBr: "Erro interno do servidor",
			HTTPStatus:    500,
			InnerError:    err.Error(),
		}

		MyErrJson := MyErr.ToJSON()

		// fallback for unknown errors
		return c.Status(500).Send([]byte(MyErrJson))
	}
}
