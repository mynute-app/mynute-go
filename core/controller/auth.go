package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/service"
	"errors"
	"fmt"

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

// func (cc *user_controller) Login(c *fiber.Ctx) error {
// 	cc.SetAction(c)
// 	body := c.Locals(namespace.GeneralKey.Model).(*model.User)
// 	var user model.User
// 	if err := cc.Request.Gorm.GetOneBy("email", body.Email, &user, []string{}); err != nil {
// 		cc.AutoReqActions.ActionFailed(404, err)
// 		return nil
// 	}
// 	if handler.ComparePassword(user.Password, body.Password) && user.Verified {
// 		cc.AutoReqActions.Status = 401
// 		return nil
// 	}
// 	claims := handler.JWT(c).CreateClaims(user.Email)
// 	token, err := handler.JWT(c).CreateToken(claims)
// 	if err != nil {
// 		cc.AutoReqActions.ActionFailed(500, err)
// 	}
// 	log.Println("User logged in")
// 	c.Response().Header.Set("Authorization", token)

// 	return nil
// }

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
//	@Router			/auth/login [post]
func (cc *auth_controller) Login(c *fiber.Ctx) error {
	cc.SetAction(c)
	body, err := lib.GetBodyFromCtx[*DTO.LoginUser](c)
	if err != nil {
		return err
	}
	var user model.User
	if err := cc.Request.Gorm.GetOneBy("email", body.Email, &user, []string{}); err != nil {
		// cc.AutoReqActions.ActionFailed(404, err)
		return lib.MyErrors.InvalidLogin.SendToClient(c)
	}
	if !handler.ComparePassword(user.Password, body.Password) {
		return lib.MyErrors.InvalidLogin.SendToClient(c)
	}
	jwt := handler.JWT(c)
	claims := jwt.CreateClaims(user)
	token, err := jwt.CreateToken(claims)
	if err != nil {
		return err
	}
	fmt.Printf("User %v logged in\n", user.Email)
	c.Response().Header.Set("Authorization", token)
	return nil
}

// func (cc *auth_controller) Login(c *fiber.Ctx) error {
// 	cc.SetAction(c)
// 	body := c.Locals(namespace.GeneralKey.Model).(*model.User)
// 	var user model.User
// 	if err := cc.Request.Gorm.GetOneBy("email", body.Email, &user, []string{}); err != nil {
// 		cc.AutoReqActions.ActionFailed(404, err)
// 		return nil
// 	}

// 	if !handler.ComparePassword(user.Password, body.Password) {
// 		cc.AutoReqActions.ActionFailed(401, errors.New("invalid password"))
// 		return nil
// 	}
// 	claims := handler.JWT(c).CreateClaims(user.Email)
// 	token, err := handler.JWT(c).CreateToken(claims)
// 	if err != nil {
// 		cc.AutoReqActions.ActionFailed(500, err)
// 	}
// 	c.Response().Header.Set("Authorization", token)

// 	return nil
// }

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

// Verify the user's email
//
//	@Summary		Verify email
//	@Description	Verify an user's email
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			id	query	string	true	"User ID"
//	@Param			code	query	string	true	"Validation code"
//	@Success		200
//	@Failure		401	{object}	DTO.ErrorResponse
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/auth/verify-email/{id}/{code} [get]
func (cc *auth_controller) VerifyEmail(c *fiber.Ctx) error {
	cc.SetAction(c)
	id := c.Params("id")
	// validationCode := c.Params("code")
	var user model.User
	if err := cc.Request.Gorm.GetOneBy("id", id, &user, UserAssociations); err != nil {
		return lib.MyErrors.UserNotFoundById.SendToClient(c)
	}
	// if userDatabase.VerificationCode != validationCode {
	// 	return lib.MyErrors.InvalidLogin.SendToClient(c)
	// }
	if user.Verified {
		return lib.MyErrors.UserNotVerified.SendToClient(c)
	}
	user.Verified = true
	if err := cc.Request.Gorm.UpdateOneById(id, &user, &user, UserAssociations); err != nil {
		return err
	}
	return nil
}

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
