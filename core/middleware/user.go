package middleware

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"

	"github.com/gofiber/fiber/v2"
)

type UserMiddlewareActions struct {
	Gorm *handler.Gorm
}

func User(Gorm *handler.Gorm) *Registry {
	// user := &UserMiddlewareActions{Gorm: Gorm}
	registry := NewRegistry()
	user := &UserMiddlewareActions{Gorm: Gorm}
	registry.RegisterAction(namespace.UserKey.Name, "POST", user.Create)

	return registry
}

func (em *UserMiddlewareActions) Create(c *fiber.Ctx) (int, error) {
	user, err := lib.GetFromCtx[*model.User](c, namespace.GeneralKey.Model)
	if err != nil {
		return 500, err
	}
	// Perform validation
	if err := lib.ValidateName(user.Name, "user"); err != nil {
		return 400, err
	}
	// Proceed to the next middleware or handler
	return 0, nil
}

func FindUserByEmail(Gorm *handler.Gorm) fiber.Handler {
	return func(c *fiber.Ctx) error {
		body, err := lib.GetFromCtx[DTO.CreateEmployee](c, namespace.RequestKey.Body)
		if err != nil {
			return err
		}
		user := &model.User{}
		if err := Gorm.DB.Where("email = ?", body.Email).First(user).Error; err != nil {
			return err
		}
		c.Locals(namespace.UserKey.Model, user)
		return c.Next()
	}
}
