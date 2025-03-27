package service

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"bytes"

	"reflect"

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
	Request        *handler.Req
	AutoReqActions *handler.AutoReqActions
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

func (b *Base[MODEL, DTO]) SetAction(c *fiber.Ctx) error {
	if err := b.saveLocals(c); err != nil {
		return err
	}
	b.AutoReqActions = b.Request.SetAutomatedActions(c)
	return nil
}

func (b *Base[MODEL, DTO]) saveLocals(c *fiber.Ctx) error {
	var modelArr []MODEL
	var dtoArr []DTO
	var model MODEL
	var dto DTO
	var changes map[string]any
	keys := namespace.GeneralKey
	method := c.Method()

	if method == "PATCH" {
		if err := c.BodyParser(&changes); err != nil {
			return err
		}
	} else {
		body := c.Request().Body()
		if len(body) != 0 && string(bytes.TrimSpace(body)) != "" {
			if body[0] == '[' {
				if err := c.BodyParser(&modelArr); err != nil {
					return err
				}
			} else {
				if err := c.BodyParser(&model); err != nil {
					return err
				}
			}
		}
	}

	if hasDto := c.Locals(keys.Dto); hasDto == nil {
		c.Locals(keys.Dto, &dto)
	}
	if hasDtoArr := c.Locals(keys.DtoArr); hasDtoArr == nil {
		c.Locals(keys.DtoArr, &dtoArr)
	}
	c.Locals(keys.ModelArr, &modelArr)
	c.Locals(keys.Model, &model)
	c.Locals(keys.Changes, changes)
	c.Locals(keys.Associations, b.Associations)
	return nil
}

func (b *Base[MODEL, DTO]) SetDTO(c *fiber.Ctx, newDTO any) *Base[MODEL, DTO] {
	keys := namespace.GeneralKey

	// Store single DTO instance
	c.Locals(keys.Dto, newDTO)

	// Dynamically create a slice of the same type as newDTO
	newDtoType := reflect.TypeOf(newDTO)

	// Handle pointer types properly
	if newDtoType.Kind() == reflect.Ptr {
		newDtoType = newDtoType.Elem()
	}

	// Create an empty slice of the same type
	newDtoArr := reflect.MakeSlice(reflect.SliceOf(newDtoType), 0, 0).Interface()

	// Store the slice in Fiber locals
	c.Locals(keys.DtoArr, &newDtoArr)

	return b
}

func (b *Base[MODEL, DTO]) GetBy(paramKey string, c *fiber.Ctx) error {
	if err := b.SetAction(c); err != nil {
		return err
	}
	return b.AutoReqActions.GetBy(paramKey)
}

func (b *Base[MODEL, DTO]) ForceGetBy(paramKey string, c *fiber.Ctx) error {
	if err := b.SetAction(c); err != nil {
		return err
	}
	return b.AutoReqActions.ForceGetBy(paramKey)
}

func (b *Base[MODEL, DTO]) DeleteOneById(c *fiber.Ctx) error {
	if err := b.SetAction(c); err != nil {
		return err
	}
	return b.AutoReqActions.DeleteOneById()
}

func (b *Base[MODEL, DTO]) ForceDeleteOneById(c *fiber.Ctx) error {
	if err := b.SetAction(c); err != nil {
		return err
	}
	return b.AutoReqActions.ForceDeleteOneById()
}

func (b *Base[MODEL, DTO]) UpdateOneById(c *fiber.Ctx) error {
	if err := b.SetAction(c); err != nil {
		return err
	}
	return b.AutoReqActions.UpdateOneById()
}

func (b *Base[MODEL, DTO]) CreateOne(c *fiber.Ctx) error {
	if err := b.SetAction(c); err != nil {
		return err
	}
	return b.AutoReqActions.CreateOne()
}

func (b *Base[MODEL, DTO]) GetAll(c *fiber.Ctx) error {
	if err := b.SetAction(c); err != nil {
		return err
	}
	return b.GetBy("", c)
}

func (b *Base[MODEL, DTO]) GetOneById(c *fiber.Ctx) error {
	if err := b.SetAction(c); err != nil {
		return err
	}
	return b.GetBy("id", c)
}

func (b *Base[MODEL, DTO]) ForceGetOneById(c *fiber.Ctx) error {
	if err := b.SetAction(c); err != nil {
		return err
	}
	return b.ForceGetBy("id", c)
}

func (b *Base[MODEL, DTO]) ForceGetAll(c *fiber.Ctx) error {
	if err := b.SetAction(c); err != nil {
		return err
	}
	return b.ForceGetBy("", c)
}
