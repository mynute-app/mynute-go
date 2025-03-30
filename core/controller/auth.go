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
	service.Base[model.Client, DTO.Client]
}

// Login just logs an client in case the password is correct
//
//	@Summary		Login
//	@Description	Log in an client
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			client	body	DTO.LoginClient	true	"Client"
//	@Success		200
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Failure		401	{object}	DTO.ErrorResponse
//	@Router			/auth/verify-existing-account/ [post]
func (cc *auth_controller) VerifyExistingAccount(c *fiber.Ctx) error {
	if err := cc.SetAction(c); err != nil {
		return err
	}
	body := c.Locals(namespace.GeneralKey.Model).(*model.Client)
	var client model.Client
	if err := cc.Request.Gorm.GetOneBy("email", body.Email, &client, []string{}); err != nil {
		if err.Error() == "record not found" {
			return cc.AutoReqActions.ActionFailed(200, err)
		}
		return cc.AutoReqActions.ActionFailed(404, err)
	}
	if client.Email == body.Email {
		return cc.AutoReqActions.ActionFailed(409, errors.New("email already registered"))
	}
	return nil
}

// OAUTH logics
func (cc *auth_controller) BeginAuthProviderCallback(c *fiber.Ctx) error {
	if err := goth_fiber.BeginAuthHandler(c); err != nil {
		return err
	}
	return nil
}

func (cc *auth_controller) GetAuthCallbackFunction(c *fiber.Ctx) error {
	client, err := goth_fiber.CompleteUserAuth(c)
	if err != nil {
		return err
	}
	if err := handler.Auth(c).StoreClientSession(client); err != nil {
		return err
	}
	if err := c.Redirect("/"); err != nil {
		return err
	}
	return nil
}

func (cc *auth_controller) LogoutProvider(c *fiber.Ctx) error {
	if err := goth_fiber.Logout(c); err != nil {
		return err
	}
	if err := c.Redirect("/"); err != nil {
		return err
	}
	return nil
}

func Auth(Gorm *handler.Gorm) *auth_controller {
	ac := &auth_controller{
		Base: service.Base[model.Client, DTO.Client]{
			Name:         namespace.ClientKey.Name,
			Request:      handler.Request(Gorm),
			Associations: []string{"Branches", "Services", "Appointment", "Company"},
		},
	}
	route := &handler.Route{DB: Gorm.DB}
	AuthResources := []*handler.ResourceRoute{
		{
			Path:        "/auth/verify-existing-account",
			Method:      "POST",
			Handler:     ac.VerifyExistingAccount,
			Description: "Verify if an account exists",
			Access:      "public",
		},
		{
			Path:        "/auth/oauth/:provider",
			Method:      "GET",
			Handler:     ac.BeginAuthProviderCallback,
			Description: "Begin auth provider callback",
			Access:      "public",
		},
		{
			Path:        "/auth/oauth/:provider/callback",
			Method:      "GET",
			Handler:     ac.GetAuthCallbackFunction,
			Description: "Get auth callback function",
			Access:      "public",
		},
		{
			Path:        "/auth/oauth/logout",
			Method:      "GET",
			Handler:     ac.LogoutProvider,
			Description: "Logout provider",
			Access:      "public",
		},
	}
	route.BulkRegisterAndSave(AuthResources)
	return ac
}
