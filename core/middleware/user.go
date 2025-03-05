package middleware

import (
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"

	"github.com/gofiber/fiber/v2"
)

type EmployeeMiddlewareActions struct {
	Gorm *handler.Gorm
}

func User(Gorm *handler.Gorm) *Registry {
	// user := &EmployeeMiddlewareActions{Gorm: Gorm}
	registry := NewRegistry()
	user := &EmployeeMiddlewareActions{Gorm: Gorm}
	registry.RegisterAction(namespace.UserKey.Name, "POST", user.Create)

	return registry
}

func (em *EmployeeMiddlewareActions) Create(c *fiber.Ctx) (int, error) {
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
