package middleware

import (
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/lib"

	"github.com/gofiber/fiber/v3"
)

type EmployeeMiddlewareActions struct {
	Gorm *handlers.Gorm
}

func Employee(Gorm *handlers.Gorm) *Registry {
	employee := &EmployeeMiddlewareActions{Gorm: Gorm}
	registry := NewRegistry()

	registry.RegisterAction(namespace.EmployeeKey.Name, "POST", employee.Create)

	return registry
}

func (em *EmployeeMiddlewareActions) Create(c fiber.Ctx) (int, error) {
	employee, err := lib.GetFromCtx[*models.User](c, namespace.GeneralKey.Model)
	if err != nil {
		return 500, err
	}
	// Perform validation
	if err := lib.ValidateName(employee.Name, "employee"); err != nil {
		return 400, err
	}
	// Proceed to the next middleware or handler
	return 0, nil
}