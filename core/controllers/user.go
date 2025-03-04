package controllers

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"
	"agenda-kaki-go/core/service"
	"log"

	"github.com/gofiber/fiber/v2"
)

// EmployeeController embeds service.Base in order to extend it with the functions below
type userController struct {
	service.Base[models.User, DTO.User]
}

func User(Gorm *handlers.Gorm) *userController {
	return &userController{
		Base: service.Base[models.User, DTO.User]{
			Name:         namespace.UserKey.Name,
			Request:      handlers.Request(Gorm),
			Middleware:   middleware.User(Gorm),
			Associations: []string{"Branches", "Services", "Appointments", "Company"},
		},
	}
}

// Custom extension method to get an employee by email
func (cc *userController) GetOneByEmail(c *fiber.Ctx) error {
	return cc.GetBy("email", c)
}

// Custom extension method to login an user
func (cc *userController) Login(c *fiber.Ctx) error {
	cc.SetAction(c)
	body := c.Locals(namespace.GeneralKey.Model).(*models.User)
	var userDatabase models.User
	if err := cc.Request.Gorm.GetOneBy("email", body.Email, &userDatabase, []string{}); err != nil {
		cc.AutoReqActions.ActionFailed(404, err)
		return nil
	}
	if handlers.ComparePassword(userDatabase.Password, body.Password) && userDatabase.Verified {
		cc.AutoReqActions.Status = 401
		return nil
	}
	claims := handlers.JWT(c).CreateClaims(userDatabase.Email)
	token, err := handlers.JWT(c).CreateToken(claims)
	if err != nil {
		cc.AutoReqActions.ActionFailed(500, err)
	}
	log.Println("User logged in")
	c.Response().Header.Set("Authorization", token)

	return nil
}
