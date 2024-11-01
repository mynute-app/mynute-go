package controllers

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v3"
)

// EmployeeController embeds BaseController in order to extend it with the functions below
type employeeController struct {
	BaseController[models.User, DTO.User]
}

func Employee(Gorm *handlers.Gorm) *employeeController {
	return &employeeController{
		BaseController: BaseController[models.User, DTO.User]{
			Name: namespace.EmployeeKey.Name,
			Request:      handlers.Request(Gorm),
			Middleware:   middleware.Employee(Gorm),
			Associations: []string{"Branches", "Services", "Appointment", "Company"},
		},
	}
}

// Custom extension method to get an employee by email
func (cc *employeeController) GetOneByEmail(c fiber.Ctx) error {
	return cc.GetBy("email", c)
}