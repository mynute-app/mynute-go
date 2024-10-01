package services

import (
	"agenda-kaki-go/core/lib"
	"log"

	"github.com/gofiber/fiber/v3"
)

type Controller struct {
	Ctx fiber.Ctx
	DB  *Postgres
}

func (c *Controller) Create(model interface{}, dto interface{}, associations []string) error {
	if err := lib.BodyParser(c.Ctx.Body(), model); err != nil {
		return lib.FiberError(400, c.Ctx, err)
	}
	if err := c.DB.Create(model, associations); err != nil {
		return lib.FiberError(400, c.Ctx, err)
	}
	if err := ConvertToDTO(model, dto); err != nil {
		return lib.FiberError(500, c.Ctx, err)
	}
	log.Printf("Created on Database! \n %+v", dto)
	return c.Ctx.JSON(dto)
}
