package handler

import (
	"agenda-kaki-go/core/lib"

	"github.com/gofiber/fiber/v2"
)

func Error(c *fiber.Ctx, err error) error {
	if e, ok := err.(lib.ErrorStruct); ok {
		return e.SendToClient(c)
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