package handlers

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"

	"github.com/gofiber/fiber/v3"
)

type Request struct {
	HTTP *HTTP
}

func (req *Request) CreateOne(c fiber.Ctx, model interface{}, dto interface{}, assocs []string, mdws []func(fiber.Ctx) (int, error)) {
	keys := namespace.GeneralKey
	actions := req.HTTP.FiberCtx(c)
	if err := lib.BodyParser(c.Body(), &model); err != nil {
		actions.sendResponse.HttpError(500, err)
	}
	c.Locals(keys.Model, model)
	c.Locals(keys.Dto, dto)
	c.Locals(keys.Associations, assocs)
	actions.RunMiddlewares(mdws).CreateOne()
}

func (req *Request) GetBy(c fiber.Ctx, paramKey string, model interface{}, dto interface{}, assocs []string, mdws []func(fiber.Ctx) (int, error)) {
	keys := namespace.GeneralKey
	actions := req.HTTP.FiberCtx(c)
	c.Locals(keys.Model, model)
	c.Locals(keys.Dto, dto)
	c.Locals(keys.Associations, assocs)

	actions.RunMiddlewares(mdws).GetBy(paramKey)
}

func (req *Request) DeleteOneById(c fiber.Ctx, model interface{}, mdws []func(fiber.Ctx) (int, error)) {
	keys := namespace.GeneralKey
	actions := req.HTTP.FiberCtx(c)
	c.Locals(keys.Model, model)

	actions.RunMiddlewares(mdws).DeleteOneById()
}

func (req *Request) UpdateOneById(c fiber.Ctx, model interface{}, dto interface{}, changes map[string]interface{}, assocs []string, mdws []func(fiber.Ctx) (int, error)) {
	keys := namespace.GeneralKey
	actions := req.HTTP.FiberCtx(c)
	if err := lib.BodyParser(c.Body(), &changes); err != nil {
		actions.sendResponse.HttpError(500, err)
	}
	c.Locals(keys.Model, model)
	c.Locals(keys.Dto, dto)
	c.Locals(keys.Changes, changes)
	c.Locals(keys.Associations, assocs)

	actions.RunMiddlewares(mdws).UpdateOneById()
}
