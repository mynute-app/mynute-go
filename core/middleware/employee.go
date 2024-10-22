package middleware

import (
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/lib"

	"github.com/gofiber/fiber/v3"
)

type Employee struct {
	Gorm *handlers.Gorm
}

type EmployeeMiddlewareActions struct {}

func (em *EmployeeMiddlewareActions) Create(c fiber.Ctx) (int, error) {
	employee, err := lib.GetFromCtx[*models.Employee](c, namespace.GeneralKey.Model)
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

var employeeActs = EmployeeMiddlewareActions{}

func (e *Employee) POST() []func(fiber.Ctx) (int, error) {
	return []func(fiber.Ctx) (int, error){employeeActs.Create}
}

func (e *Employee) PUT() []func(fiber.Ctx) (int, error) {
	return []func(fiber.Ctx) (int, error){}
}

func (e *Employee) DELETE() []func(fiber.Ctx) (int, error) {
	return []func(fiber.Ctx) (int, error){}
}

func (e *Employee) GET() []func(fiber.Ctx) (int, error) {
	return []func(fiber.Ctx) (int, error){}
}

func (e *Employee) PATCH() []func(fiber.Ctx) (int, error) {
	return []func(fiber.Ctx) (int, error){}
}

func (e *Employee) ForceDELETE() []func(fiber.Ctx) (int, error) {
	return []func(fiber.Ctx) (int, error){}
}

func (e *Employee) ForceGET() []func(fiber.Ctx) (int, error) {
	return []func(fiber.Ctx) (int, error){}
}