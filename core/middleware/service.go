package middleware

import (
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/lib"

	"github.com/gofiber/fiber/v3"
)

type Service struct {
	Gorm *handlers.Gorm
}

type ServiceMiddlewareActions struct{}

func (sm *ServiceMiddlewareActions) Create(c fiber.Ctx) (int, error) {
	service, err := lib.GetFromCtx[*models.Service](c, namespace.GeneralKey.Model)
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

var serviceActs = ServiceMiddlewareActions{}

func (s *Service) POST() []func(fiber.Ctx) (int, error) {
	return []func(fiber.Ctx) (int, error){serviceActs.Create}
}

func (s *Service) PUT() []func(fiber.Ctx) (int, error) {
	return []func(fiber.Ctx) (int, error){}
}

func (s *Service) DELETE() []func(fiber.Ctx) (int, error) {
	return []func(fiber.Ctx) (int, error){}
}

func (s *Service) GET() []func(fiber.Ctx) (int, error) {
	return []func(fiber.Ctx) (int, error){}
}

func (s *Service) PATCH() []func(fiber.Ctx) (int, error) {
	return []func(fiber.Ctx) (int, error){}
}

func (s *Service) ForceDELETE() []func(fiber.Ctx) (int, error) {
	return []func(fiber.Ctx) (int, error){}
}

func (s *Service) ForceGET() []func(fiber.Ctx) (int, error) {
	return []func(fiber.Ctx) (int, error){}
}