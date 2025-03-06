package middleware

import (
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"

	"github.com/gofiber/fiber/v2"
)

type ServiceMiddlewareActions struct {
	Gorm *handler.Gorm
}

func Service(Gorm *handler.Gorm) *Registry {
	service := &ServiceMiddlewareActions{Gorm: Gorm}
	registry := NewRegistry()

	registry.RegisterAction(namespace.ServiceKey.Name, "POST", service.Create)

	return registry
}

func (sm *ServiceMiddlewareActions) Create(c *fiber.Ctx) (int, error) {
	service, err := lib.GetFromCtx[*model.Service](c, namespace.GeneralKey.Model)
	if err != nil {
		return 500, err
	}
	// Perform validation
	if err := lib.ValidateName(service.Name, "service"); err != nil {
		return 400, err
	}
	// Proceed to the next middleware or handler
	return 0, nil
}
