package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/service"
	"errors"

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

// Login just logs an user in case the password is correct
//
//	@Summary		Login
//	@Description	Log in an user
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			user	body	DTO.LoginUser	true	"User"
//	@Success		200
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Failure		401	{object}	DTO.ErrorResponse
//	@Router			/auth/verify-existing-account/ [post]
func (cc *auth_controller) VerifyExistingAccount(c *fiber.Ctx) error {
	cc.SetAction(c)
	body := c.Locals(namespace.GeneralKey.Model).(*model.User)
	var user model.User
	if err := cc.Request.Gorm.GetOneBy("email", body.Email, &user, []string{}); err != nil {
		if err.Error() == "record not found" {
			cc.AutoReqActions.ActionFailed(200, err)
			return nil
		}
		cc.AutoReqActions.ActionFailed(404, err)
		return nil
	}
	if user.Email == body.Email {
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
