package service

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v2"
)

var _ IService = (*BaseService[IService, IService])(nil)

type IService interface {
	GetBy(paramKey string, c *fiber.Ctx) error
	ForceGetBy(paramKey string, c *fiber.Ctx) error
	CreateOne(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	GetOneById(c *fiber.Ctx) error
	UpdateOneById(c *fiber.Ctx) error
	DeleteOneById(c *fiber.Ctx) error
	ForceDeleteOneById(c *fiber.Ctx) error
	ForceGetOneById(c *fiber.Ctx) error
	ForceGetAll(c *fiber.Ctx) error
}

type BaseService[MODEL any, DTO any] struct {
	Name           string
	Request        *handlers.Req
	AutoReqActions *handlers.AutoReqActions
	Middleware     *middleware.Registry
	Associations   []string
}

func CreateRoutes(r fiber.Router, ci IService) {
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

func (bs *BaseService[MODEL, DTO]) SetAction(c *fiber.Ctx) {
	bs.saveLocals(c)
	bs.runMiddlewares(c)
	bs.AutoReqActions = bs.Request.SetAutomatedActions(c)
}

func (bc *BaseService[MODEL, DTO]) runMiddlewares(c *fiber.Ctx) {
	mdws := bc.Middleware.GetActions(bc.Name, c.Method())
	for _, mdw := range mdws {
		if s, err := mdw(c); err != nil {
			bc.AutoReqActions.ActionFailed(s, err)
			return
		}
	}
}

func (bc *BaseService[MODEL, DTO]) saveLocals(c *fiber.Ctx) {
	var modelArr []MODEL
	var dtoArr []DTO
	var model MODEL
	var dto DTO
	var changes map[string]any
	keys := namespace.GeneralKey
	if s, err := middleware.ParseBodyToContext(c, keys.Model, &model); err != nil {
		bc.AutoReqActions.ActionFailed(s, err)
		return
	}
	c.Locals(keys.ModelArr, &modelArr)
	c.Locals(keys.Dto, &dto)
	c.Locals(keys.DtoArr, &dtoArr)
	c.Locals(keys.Changes, changes)
	c.Locals(keys.Associations, bc.Associations)
}

func (bc *BaseService[MODEL, DTO]) GetBy(paramKey string, c *fiber.Ctx) error {
	bc.AutoReqActions.GetBy(paramKey)
	return nil
}

func (bc *BaseService[MODEL, DTO]) ForceGetBy(paramKey string, c *fiber.Ctx) error {
	bc.AutoReqActions.ForceGetBy(paramKey)
	return nil
}

func (bc *BaseService[MODEL, DTO]) DeleteOneById(c *fiber.Ctx) error {
	bc.AutoReqActions.DeleteOneById()
	return nil
}

func (bc *BaseService[MODEL, DTO]) ForceDeleteOneById(c *fiber.Ctx) error {
	bc.AutoReqActions.ForceDeleteOneById()
	return nil
}

func (bc *BaseService[MODEL, DTO]) UpdateOneById(c *fiber.Ctx) error {
	bc.AutoReqActions.UpdateOneById()
	return nil
}

func (bc *BaseService[MODEL, DTO]) CreateOne(c *fiber.Ctx) error {
	bc.AutoReqActions.CreateOne()
	return nil
}

func (bc *BaseService[MODEL, DTO]) GetAll(c *fiber.Ctx) error {
	return bc.GetBy("", c)
}

func (bc *BaseService[MODEL, DTO]) GetOneById(c *fiber.Ctx) error {
	return bc.GetBy("id", c)
}

func (bc *BaseService[MODEL, DTO]) ForceGetOneById(c *fiber.Ctx) error {
	return bc.ForceGetBy("id", c)
}

func (bc *BaseService[MODEL, DTO]) ForceGetAll(c *fiber.Ctx) error {
	return bc.ForceGetBy("", c)
}
