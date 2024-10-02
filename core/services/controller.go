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

func (c *Controller) UpdateOneBy(param string, model interface{}, dto interface{}, associations []string) error {
	paramValue := c.Ctx.Params(param)
	var changes map[string]interface{}
	if err := lib.BodyParser(c.Ctx.Body(), &changes); err != nil {
		return lib.FiberError(400, c.Ctx, err)
	}
	log.Printf("Changes: %+v", changes)
	if err := c.DB.UpdateOneBy(param, paramValue, model, changes, associations); err != nil {
		return lib.FiberError(400, c.Ctx, err)
	}
	if err := ConvertToDTO(changes, dto); err != nil {
		return lib.FiberError(500, c.Ctx, err)
	}
	log.Printf("Updated on Database using '%s'! \n %+v", param, dto)
	return c.Ctx.JSON(dto)
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

func (c *Controller) GetAll(model interface{}, dto interface{}, associations []string) error {
	if err := c.DB.GetAll(model, associations); err != nil {
		return lib.FiberError(404, c.Ctx, err)
	}
	if err := ConvertToDTO(model, dto); err != nil {
		return lib.FiberError(500, c.Ctx, err)
	}
	log.Printf("Retrieved from Database! \n %+v", dto)
	return c.Ctx.JSON(dto)
}

func (c *Controller) GetOneBy(param string, model interface{}, dto interface{}, associations []string) error {
	paramValue := c.Ctx.Params(param)
	if err := c.DB.GetOneBy(param, paramValue, model, associations); err != nil {
		lib.FiberError(404, c.Ctx, err)
	}
	if err := ConvertToDTO(model, dto); err != nil {
		lib.FiberError(500, c.Ctx, err)
	}
	log.Printf("Got from Database using '%s'! \n %+v", param, dto)
	return c.Ctx.JSON(dto)
}

func (c *Controller) DeleteOneBy(param string, model interface{}) error {
	paramValue := c.Ctx.Params(param)
	if err := c.DB.DeleteOneBy(param, paramValue, model); err != nil {
		return lib.FiberError(404, c.Ctx, err)
	}
	log.Printf("Deleted from Database using '%s'! \n %+v", param, model)
	return c.Ctx.JSON(model)
}
