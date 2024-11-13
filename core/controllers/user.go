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
type userController struct {
	BaseController[models.User, DTO.User]
}

func User(Gorm *handlers.Gorm) *userController {
	return &userController{
		BaseController: BaseController[models.User, DTO.User]{
			Name:         namespace.UserKey.Name,
			Request:      handlers.Request(Gorm),
			Middleware:   middleware.User(Gorm),
			Associations: []string{"Branches", "Services", "Appointment", "Company"},
		},
	}
}

// Custom extension method to get an employee by email
func (cc *userController) GetOneByEmail(c fiber.Ctx) error {
	return cc.GetBy("email", c)
}

// Custom extension method to login an user
func (cc *userController) Login(c fiber.Ctx) error {
	cc.init(c)
	body := c.Locals(namespace.GeneralKey.Model).(DTO.User)
	var user models.User
	if err := cc.Request.Gorm.GetOneBy("email", body.Email, &user, []string{}); err != nil {
		cc.reqActions.SendError(404, err)
		return nil
	}

	return middleware.WhoAreYou(c)
}
