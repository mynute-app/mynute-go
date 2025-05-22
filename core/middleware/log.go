package middleware

import (
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/lib/Logger"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Log(logger *slog.Logger) fiber.Handler {
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

		loki := myLogger.New("loki")

		if err := loki.Log(labelMsg, ReqLabels); err != nil {
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

		if err := loki.Log(labelMsg, ResLabels); err != nil {
			logger.Error("Failed to send log to Loki", slog.String("error", err.Error()))
		}

		return err
	}
}
