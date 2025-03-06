package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/service"
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"
)

// EmployeeController embeds Base in order to extend it with the functions below
type auth_controller struct {
	service.Base[model.User, DTO.User]
}

func Auth(Gorm *handler.Gorm) *auth_controller {
	return &auth_controller{
		Base: service.Base[model.User, DTO.User]{
			Name:         namespace.UserKey.Name,
			Request:      handler.Request(Gorm),
			Associations: []string{"Branches", "Services", "Appointment", "Company"},
		},
	}
}

// Custom extension method to login an user
func (cc *auth_controller) Login(c *fiber.Ctx) error {
	cc.SetAction(c)
	body := c.Locals(namespace.GeneralKey.Model).(*model.User)
	var userDatabase model.User
	if err := cc.Request.Gorm.GetOneBy("email", body.Email, &userDatabase, []string{}); err != nil {
		cc.AutoReqActions.ActionFailed(404, err)
		return nil
	}

	if !handler.ComparePassword(userDatabase.Password, body.Password) {
		cc.AutoReqActions.ActionFailed(401, errors.New("invalid password"))
		return nil
	}
	claims := handler.JWT(c).CreateClaims(userDatabase.Email)
	token, err := handler.JWT(c).CreateToken(claims)
	if err != nil {
		cc.AutoReqActions.ActionFailed(500, err)
	}
	log.Println("User logged in: ", userDatabase.Email)
	c.Response().Header.Set("Authorization", token)

	return nil
}

func (cc *auth_controller) Register(c *fiber.Ctx) error {
	cc.SetAction(c)
	body := c.Locals(namespace.GeneralKey.Model).(*model.User)
	body.Password, _ = handler.HashPassword(body.Password)
	if err := cc.Request.Gorm.Create(body); err != nil {
		cc.AutoReqActions.ActionFailed(500, err)
		return nil
	}
	return nil
}

func (cc *auth_controller) VerifyEmail(c *fiber.Ctx) error {
	cc.SetAction(c)
	userId := c.Query("id")
	validationCode := c.Query("code")

	email := handler.NewGoMailService(&handler.Configs{})

	err := email.SendEmail(handler.EmailRequestTest)
	if err != nil {
		cc.AutoReqActions.ActionFailed(500, err)
	}
	return nil
	var userDatabase model.User
	if err := cc.Request.Gorm.GetOneBy("id", userId, &userDatabase, []string{}); err != nil {
		cc.AutoReqActions.ActionFailed(500, err)
	}
	if userDatabase.VerificationCode != validationCode {
		cc.AutoReqActions.ActionFailed(401, errors.New("invalid validation code"))
	}
	if userDatabase.Verified {
		cc.AutoReqActions.ActionFailed(409, errors.New("account already verified"))
	}
	userDatabase.Verified = true
	if err := cc.Request.Gorm.UpdateOneById(userId, model.User{}, &userDatabase, []string{}); err != nil {
		cc.AutoReqActions.ActionFailed(500, err)
	}
	cc.AutoReqActions.ActionSuccess(200, nil, nil)
	return nil
}

func (cc *auth_controller) VerifyExistingAccount(c *fiber.Ctx) error {
	cc.SetAction(c)
	body := c.Locals(namespace.GeneralKey.Model).(*model.User)
	var userDatabase model.User
	if err := cc.Request.Gorm.GetOneBy("email", body.Email, &userDatabase, []string{}); err != nil {
		if err.Error() == "record not found" {
			cc.AutoReqActions.ActionFailed(200, err)
			return nil
		}
		cc.AutoReqActions.ActionFailed(404, err)
		return nil
	}
	if userDatabase.Email == body.Email {
		cc.AutoReqActions.ActionFailed(409, errors.New("email already registered"))
		return nil
	}
	return nil
}

// OAUTH logics
func (cc *auth_controller) BeginAuthProviderCallback(c *fiber.Ctx) error {
	if err := goth_fiber.BeginAuthHandler(c); err != nil {
		cc.AutoReqActions.ActionFailed(500, err)
		return nil
	}
	return nil
}

func (cc *auth_controller) GetAuthCallbackFunction(c *fiber.Ctx) error {
	user, err := goth_fiber.CompleteUserAuth(c)
	if err != nil {
		cc.AutoReqActions.ActionFailed(500, err)
		return nil
	}
	handler.Auth(c).StoreUserSession(user)
	c.Redirect("/")
	return nil
}

func (cc *auth_controller) LogoutProvider(c *fiber.Ctx) error {
	goth_fiber.Logout(c)
	c.Redirect("/")
	return nil
}
