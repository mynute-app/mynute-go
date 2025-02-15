package controllers

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	gothfiber "github.com/luigiazoreng/goth_fiber"
)

// EmployeeController embeds BaseController in order to extend it with the functions below
type authController struct {
	BaseController[models.User, DTO.User]
}

func Auth(Gorm *handlers.Gorm) *authController {
	return &authController{
		BaseController: BaseController[models.User, DTO.User]{
			Name:         namespace.UserKey.Name,
			Request:      handlers.Request(Gorm),
			Middleware:   middleware.User(Gorm),
			Associations: []string{"Branches", "Services", "Appointment", "Company"},
		},
	}
}

// Custom extension method to login an user
func (cc *authController) Login(c *fiber.Ctx) error {
	cc.init(c)
	body := c.Locals(namespace.GeneralKey.Model).(*models.User)
	var userDatabase models.User
	if err := cc.Request.Gorm.GetOneBy("email", body.Email, &userDatabase, []string{}); err != nil {
		cc.reqActions.SendError(404, err)
		return nil
	}

	if !handlers.ComparePassword(userDatabase.Password, body.Password) {
		cc.reqActions.SendError(401, errors.New("invalid password"))
		return nil
	}
	claims := handlers.JWT(c).CreateClaims(userDatabase.Email)
	token, err := handlers.JWT(c).CreateToken(claims)
	if err != nil {
		cc.reqActions.SendError(500, err)
	}
	log.Println("User logged in: ", userDatabase.Email)
	c.Response().Header.Set("Authorization", token)

	return nil
}

func (cc *authController) Register(c *fiber.Ctx) error {
	cc.init(c)
	body := c.Locals(namespace.GeneralKey.Model).(*models.User)
	body.Password, _ = handlers.HashPassword(body.Password)
	if err := cc.Request.Gorm.Create(body); err != nil {
		cc.reqActions.SendError(500, err)
		return nil
	}
	return nil
}

func (cc *authController) VerifyEmail(c *fiber.Ctx) error {

	return nil
}

func (cc *authController) VerifyExistingAccount(c *fiber.Ctx) error {
	cc.init(c)
	body := c.Locals(namespace.GeneralKey.Model).(*models.User)
	var userDatabase models.User
	if err := cc.Request.Gorm.GetOneBy("email", body.Email, &userDatabase, []string{}); err != nil {
		if err.Error() == "record not found" {
			cc.reqActions.SendError(200, err)
			return nil
		}
		cc.reqActions.SendError(404, err)
		return nil
	}
	if userDatabase.Email == body.Email {
		cc.reqActions.SendError(409, errors.New("email already registered"))
		return nil
	}
	return nil
}

// OAUTH logics
func (cc *authController) BeginAuthProviderCallback(c *fiber.Ctx) error {
	if err := gothfiber.BeginAuthHandler(c); err != nil {
		cc.reqActions.SendError(500, err)
		return nil
	}
	return nil
}

func (cc *authController) GetAuthCallbackFunction(c *fiber.Ctx) error {
	user, err := gothfiber.CompleteUserAuth(c)
	if err != nil {
		cc.reqActions.SendError(500, err)
		return nil
	}
	handlers.Auth(c).StoreUserSession(user)
	c.Redirect().To("/")
	return nil
}

func (cc *authController) LogoutProvider(c *fiber.Ctx) error {
	gothfiber.Logout(c)
	c.Redirect().To("/")
	return nil
}
