package middleware

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"

	"github.com/gofiber/fiber/v2"
)

type EmployeeMiddlewareActions struct {
	Gorm *handler.Gorm
}

func Employee(Gorm *handler.Gorm) *Registry {
	registry := NewRegistry()
	// employee := &EmployeeMiddlewareActions{Gorm: Gorm}
	return registry
}

func SetEmployeeUserAccount(gorm *handler.Gorm) fiber.Handler {
	return func(c *fiber.Ctx) error {
		body, err := lib.GetFromCtx[*DTO.CreateEmployee](c, namespace.RequestKey.Body)
		if err != nil {
			return err
		}
		if body.UserID != 0 {
			return c.Next() // Skip if user ID is already set
		}
		user, err := lib.GetFromCtx[*model.User](c, namespace.GeneralKey.Model)
		if err != nil {
			return err
		}
		if user.ID != 0 {
			body.UserID = user.ID // Set user ID if user is already created
			c.Locals(namespace.RequestKey.Body, body)
			return c.Next()
		}
		// Create an user account in case it doesn't exist
		lib.ParseToDTO(body, user)
		if err := gorm.DB.Create(user).Error; err != nil {
			return err
		}
		body.UserID = user.ID
		return c.Next()
	}
}
