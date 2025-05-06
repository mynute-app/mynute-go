package controller

import (
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/service"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Create(c *fiber.Ctx, model any) error {
	if err := c.BodyParser(model); err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}
	Service := service.Factory(c).Model(model)
	if err := Service.Create(); err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}
	m, ok := model.(interface{ GetID() string })
	if !ok {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("model does not implement GetID method"))
	}
	if err := Service.MyGorm.GetOneBy("id", m.GetID(), model); err != nil {
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
