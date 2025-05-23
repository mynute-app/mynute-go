package middleware

import (
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/lib/Logger"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Error(logger *slog.Logger) fiber.ErrorHandler {
	loki := &myLogger.Loki{}
	return func(c *fiber.Ctx, err error) error {
		if e, ok := err.(lib.ErrorStruct); ok {
			if err := e.SendToClient(c); err != nil {
				ResMsg := myLogger.GetResMessage(c)
				labels := myLogger.Labels{
					App:    "main-api",
					Method: c.Method(),
					Path:   c.Path(),
					IP:     c.IP(),
					Host:   string(c.Request().Host()),
				}
				ResLabels := labels.GetResLabels(time.Now(), e.HTTPStatus)
				ResLabels["level"] = "critical"
				ResMsg += "Failed to send error response to client!" + ResMsg + "\n\n" + e.WithError(err).Error()
				if err := loki.Log(ResMsg, ResLabels); err != nil {
					logger.Error("Failed to even send error response log to Loki", slog.String("error", err.Error()))
					return nil
				}
				return nil
			}
			return nil
		}

		MyErr := lib.ErrorStruct{
			DescriptionEn: "Internal Server Error",
			DescriptionBr: "Erro interno do servidor",
			HTTPStatus:    500,
			InnerError: map[int]string{
				1: err.Error(),
			},
		}

		MyErrJson := MyErr.ToJSON()

		// fallback for unknown errors
		return c.Status(500).Send([]byte(MyErrJson))
	}
}
