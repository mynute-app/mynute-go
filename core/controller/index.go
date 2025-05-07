package controller

import (
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/service"
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Create(c *fiber.Ctx, model any) error {
	if err := c.BodyParser(model); err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}
	Service := service.Factory(c).Model(model)
	if err := Service.Create(); err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}
	HasID := func (model any) (uuid.UUID, bool) {
		val := reflect.ValueOf(model)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		field := val.FieldByName("ID")
		if !field.IsValid() || field.Type() != reflect.TypeOf(uuid.UUID{}) {
			return uuid.Nil, false
		}
		return field.Interface().(uuid.UUID), true
	}
	id, ok := HasID(model)
	if !ok {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("model does have ID field"))
	}
	if err := Service.MyGorm.GetOneBy("id", id.String(), model); err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}
	return nil
}

func GetOneBy(param string, c *fiber.Ctx, model any) error {
	Service := service.Factory(c).Model(model)
	if err := Service.GetBy(param); err != nil {
		return err
	}
	return nil
}

func UpdateOneById(c *fiber.Ctx, model any) error {
	return service.Factory(c).Model(model).UpdateOneById()
}

func DeleteOneById(c *fiber.Ctx, model any) error {
	return service.Factory(c).Model(model).DeleteOneById()
}
