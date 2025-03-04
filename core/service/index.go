package service

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v2"
)

var _ IService = (*Base[IService, IService])(nil)

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

type Base[MODEL any, DTO any] struct {
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

func (b *Base[MODEL, DTO]) SetAction(c *fiber.Ctx) {
	b.saveLocals(c)
	b.AutoReqActions = b.Request.SetAutomatedActions(c)
	b.runMiddlewares(c)
}

func (b *Base[MODEL, DTO]) runMiddlewares(c *fiber.Ctx) {
	mdws := b.Middleware.GetActions(b.Name, c.Method())
	for _, mdw := range mdws {
		if s, err := mdw(c); err != nil {
			b.AutoReqActions.ActionFailed(s, err)
			return
		}
	}
}

func (b *Base[MODEL, DTO]) saveLocals(c *fiber.Ctx) {
	var modelArr []MODEL
	var dtoArr []DTO
	var model MODEL
	var dto DTO
	var changes map[string]any
	keys := namespace.GeneralKey
	if s, err := middleware.ParseBodyToContext(c, keys.Model, &model); err != nil {
		b.AutoReqActions.ActionFailed(s, err)
		return
	}
	c.Locals(keys.ModelArr, &modelArr)
	c.Locals(keys.Dto, &dto)
	c.Locals(keys.DtoArr, &dtoArr)
	c.Locals(keys.Changes, changes)
	c.Locals(keys.Associations, b.Associations)
}

func (b *Base[MODEL, DTO]) GetBy(paramKey string, c *fiber.Ctx) error {
	b.SetAction(c)
	b.AutoReqActions.GetBy(paramKey)
	return nil
}

func (b *Base[MODEL, DTO]) ForceGetBy(paramKey string, c *fiber.Ctx) error {
	b.SetAction(c)
	b.AutoReqActions.ForceGetBy(paramKey)
	return nil
}

func (b *Base[MODEL, DTO]) DeleteOneById(c *fiber.Ctx) error {
	b.SetAction(c)
	b.AutoReqActions.DeleteOneById()
	return nil
}

func (b *Base[MODEL, DTO]) ForceDeleteOneById(c *fiber.Ctx) error {
	b.SetAction(c)
	b.AutoReqActions.ForceDeleteOneById()
	return nil
}

func (b *Base[MODEL, DTO]) UpdateOneById(c *fiber.Ctx) error {
	b.SetAction(c)
	b.AutoReqActions.UpdateOneById()
	return nil
}

func (b *Base[MODEL, DTO]) CreateOne(c *fiber.Ctx) error {
	b.SetAction(c)
	b.AutoReqActions.CreateOne()
	return nil
}

func (b *Base[MODEL, DTO]) GetAll(c *fiber.Ctx) error {
	b.SetAction(c)
	return b.GetBy("", c)
}

func (b *Base[MODEL, DTO]) GetOneById(c *fiber.Ctx) error {
	b.SetAction(c)
	return b.GetBy("id", c)
}

func (b *Base[MODEL, DTO]) ForceGetOneById(c *fiber.Ctx) error {
	b.SetAction(c)
	return b.ForceGetBy("id", c)
}

func (b *Base[MODEL, DTO]) ForceGetAll(c *fiber.Ctx) error {
	b.SetAction(c)
	return b.ForceGetBy("", c)
}
