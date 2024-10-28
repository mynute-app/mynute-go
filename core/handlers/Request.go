package handlers

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"

	"github.com/gofiber/fiber/v3"
)

type Request struct {
	HTTP *HTTP
}

func (req *Request) saveLocals(c fiber.Ctx, model interface{}, dto interface{}, assocs []string, changes map[string]interface{}) {
	keys := namespace.GeneralKey
	c.Locals(keys.Model, model)
	c.Locals(keys.Dto, dto)
	c.Locals(keys.Associations, assocs)
	c.Locals(keys.Changes, changes)
}

func (req *Request) CreateOne(c fiber.Ctx, model interface{}, dto interface{}, assocs []string, mdws []func(fiber.Ctx) (int, error)) {
	req.saveLocals(c, model, dto, assocs, nil)
	actions := req.HTTP.FiberCtx(c)
	if err := lib.BodyParser(c.Body(), &model); err != nil {
		actions.sendResponse.Http500(err)
	}
	actions.RunMiddlewares(mdws).CreateOne()
}

func (req *Request) GetBy(c fiber.Ctx, paramKey string, model interface{}, dto interface{}, assocs []string, mdws []func(fiber.Ctx) (int, error)) {
	req.saveLocals(c, model, dto, assocs, nil)
	actions := req.HTTP.FiberCtx(c)
	actions.RunMiddlewares(mdws).GetBy(paramKey)
}

func (req *Request) ForceGetBy(c fiber.Ctx, paramKey string, model interface{}, dto interface{}, assocs []string, mdws []func(fiber.Ctx) (int, error)) {
	req.saveLocals(c, model, dto, assocs, nil)
	actions := req.HTTP.FiberCtx(c)
	actions.RunMiddlewares(mdws).ForceGetBy(paramKey)
}

func (req *Request) DeleteOneById(c fiber.Ctx, model interface{}, mdws []func(fiber.Ctx) (int, error)) {
	req.saveLocals(c, model, nil, nil, nil)
	actions := req.HTTP.FiberCtx(c)
	actions.RunMiddlewares(mdws).DeleteOneById()
}

func (req *Request) ForceDeleteOneById(c fiber.Ctx, model interface{}, mdws []func(fiber.Ctx) (int, error)) {
	req.saveLocals(c, model, nil, nil, nil)
	actions := req.HTTP.FiberCtx(c)
	actions.RunMiddlewares(mdws).ForceDeleteOneById()
}

func (req *Request) UpdateOneById(c fiber.Ctx, model interface{}, dto interface{}, changes map[string]interface{}, assocs []string, mdws []func(fiber.Ctx) (int, error)) {
	req.saveLocals(c, model, dto, assocs, changes)
	actions := req.HTTP.FiberCtx(c)
	if err := lib.BodyParser(c.Body(), &changes); err != nil {
		actions.sendResponse.HttpError(500, err)
	}
	actions.RunMiddlewares(mdws).UpdateOneById()
}
