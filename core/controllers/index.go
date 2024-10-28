package controllers

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v3"
)

var _ IController = (*BaseController[IController, IController])(nil)

type BaseController[MODEL any, DTO any] struct {
	Name         string
	Request      *handlers.Request
	reqActions   *handlers.ReqActions
	Middleware   *middleware.Registry
	Associations []string
}

type IController interface {
	GetBy(paramKey string, c fiber.Ctx) error
	ForceGetBy(paramKey string, c fiber.Ctx) error
	CreateOne(c fiber.Ctx) error
	GetAll(c fiber.Ctx) error
	GetOneById(c fiber.Ctx) error
	UpdateOneById(c fiber.Ctx) error
	DeleteOneById(c fiber.Ctx) error
	ForceDeleteOneById(c fiber.Ctx) error
	ForceGetOneById(c fiber.Ctx) error
	ForceGetAll(c fiber.Ctx) error
}

func CreateRoutes(r fiber.Router, ci IController) {
	r.Post("/", ci.CreateOne)       // ok
	r.Get("/", ci.GetAll)           // ok
	r.Get("/force", ci.ForceGetAll) // ok
	id := r.Group("/:id")
	id.Get("/", ci.GetOneById)                 // ok
	id.Patch("/", ci.UpdateOneById)            // ok
	id.Delete("/", ci.DeleteOneById)           // ok
	id.Delete("/force", ci.ForceDeleteOneById) // ok
	id.Get("/force", ci.ForceGetOneById)       // ok
}

func (bc *BaseController[MODEL, DTO]) init(c fiber.Ctx) {
	bc.saveLocals(c)
	bc.reqActions = bc.Request.FiberCtx(c)
	bc.runMiddlewares(c)
}

func (bc *BaseController[MODEL, DTO]) runMiddlewares(c fiber.Ctx) {
	bc.saveLocals(c)
	mdws := bc.Middleware.GetActions(bc.Name, c.Method())
	for _, mdw := range mdws {
		if s, err := mdw(c); err != nil {
			bc.reqActions.SendError(s, err)
			return
		}
	}
}

func (bc *BaseController[MODEL, DTO]) saveLocals(c fiber.Ctx) {
	var modelArr []MODEL
	var dtoArr []DTO
	var model MODEL
	var dto DTO
	var changes map[string]interface{}
	keys := namespace.GeneralKey
	c.Locals(keys.Model, model)
	c.Locals(keys.Dto, dto)
	c.Locals(keys.ModelArr, modelArr)
	c.Locals(keys.DtoArr, dtoArr)
	c.Locals(keys.Changes, changes)
	c.Locals(keys.Associations, bc.Associations)
}

func (bc *BaseController[MODEL, DTO]) GetBy(paramKey string, c fiber.Ctx) error {
	bc.init(c)
	bc.reqActions.GetBy(paramKey)
	return nil
}

func (bc *BaseController[MODEL, DTO]) ForceGetBy(paramKey string, c fiber.Ctx) error {
	bc.init(c)
	bc.reqActions.ForceGetBy(paramKey)
	return nil
}

func (bc *BaseController[MODEL, DTO]) DeleteOneById(c fiber.Ctx) error {
	bc.init(c)
	bc.reqActions.DeleteOneById()
	return nil
}

func (bc *BaseController[MODEL, DTO]) ForceDeleteOneById(c fiber.Ctx) error {
	bc.init(c)
	bc.reqActions.ForceDeleteOneById()
	return nil
}

func (bc *BaseController[MODEL, DTO]) UpdateOneById(c fiber.Ctx) error {
	bc.init(c)
	bc.reqActions.UpdateOneById()
	return nil
}

func (bc *BaseController[MODEL, DTO]) CreateOne(c fiber.Ctx) error {
	bc.init(c)
	bc.reqActions.CreateOne()
	return nil
}

func (bc *BaseController[MODEL, DTO]) GetAll(c fiber.Ctx) error {
	return bc.GetBy("", c)
}

func (bc *BaseController[MODEL, DTO]) GetOneById(c fiber.Ctx) error {
	return bc.GetBy("id", c)
}

func (bc *BaseController[MODEL, DTO]) ForceGetOneById(c fiber.Ctx) error {
	return bc.ForceGetBy("id", c)
}

func (bc *BaseController[MODEL, DTO]) ForceGetAll(c fiber.Ctx) error {
	return bc.ForceGetBy("", c)
}
