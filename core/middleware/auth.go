package middleware

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type auth_middleware struct {
	Gorm *handler.Gorm
}

func Auth(Gorm *handler.Gorm) *auth_middleware {
	return &auth_middleware{Gorm: Gorm}
}

func (am *auth_middleware) DenyUnauthorized(c *fiber.Ctx) error {
	auth_claims := c.Locals(namespace.RequestKey.Auth_Claims)
	claim, ok := auth_claims.(*DTO.Claims)
	if !ok {
		return lib.Error.Auth.InvalidToken.SendToClient(c)
	}
	if claim.ID == 0 {
		return lib.Error.Auth.InvalidToken.SendToClient(c)
	}
	if !claim.Verified {
		return lib.Error.Client.NotVerified.SendToClient(c)
	}
	return c.Next()
}

func (am *auth_middleware) Login() []fiber.Handler {
	return []fiber.Handler{
		lib.SaveBodyOnCtx[DTO.LoginClient],
		am.DenyLoginFromUnverified,
	}
}

func (am *auth_middleware) DenyUnverified(c *fiber.Ctx) error {
	auth_claims := c.Locals(namespace.RequestKey.Auth_Claims)
	user, ok := auth_claims.(*DTO.ClientPopulated)
	if !ok {
		return lib.Error.Auth.InvalidToken.SendToClient(c)
	}
	if !user.Verified {
		return lib.Error.Client.NotVerified.SendToClient(c)
	}
	return c.Next()
}

func (am *auth_middleware) DenyClaimless(c *fiber.Ctx) error {
	res := lib.SendResponse{Ctx: c}
	auth_claims, ok := c.Locals(namespace.RequestKey.Auth_Claims).(*model.Client)
	if !ok {
		return res.Http401(nil)
	}
	if auth_claims.ID == 0 {
		return res.Http401(nil)
	}
	return c.Next()
}

func (am *auth_middleware) DenyLoginFromUnverified(c *fiber.Ctx) error {
	// Get login body from context
	login, err := lib.GetBodyFromCtx[*DTO.LoginClient](c)
	if err != nil {
		return err
	}
	// Get user from email
	user := &[]model.Client{}
	am.Gorm.DB.Where("email = ?", login.Email).Find(user)
	if len(*user) == 0 {
		fmt.Printf("Client %v not found\n", login.Email)
		return lib.Error.Auth.InvalidLogin.SendToClient(c)
	}
	// Check if user is verified
	if !(*user)[0].Verified {
		fmt.Printf("Client %v is not verified\n", login.Email)
		return lib.Error.Client.NotVerified.SendToClient(c)
	}
	return c.Next()
}

func (am *auth_middleware) WhoAreYou(c *fiber.Ctx) error {
	authorization := c.Get("Authorization")
	if authorization == "" {
		return c.Next()
	}
	user, err := handler.JWT(c).WhoAreYou()
	if err != nil {
		return err
	} else if user == nil {
		return c.Next()
	}
	c.Locals(namespace.RequestKey.Auth_Claims, user)
	return c.Next()
}

// func (am *auth_middleware) WhoAreYou(c *fiber.Ctx) error {
// 	res := lib.SendResponse{Ctx: c}
// 	// if c.Get("Authorization") == "" {
// 	// 	err := handler.Auth(c).WhoAreYou()
// 	// 	if err != nil {
// 	// 		return res.Http401(err).Next()
// 	// 	}
// 	// 	return nil
// 	// }
// 	if c.Get("Authorization") != "" {
// 		err := handler.JWT(c).WhoAreYou()
// 		if err != nil {
// 			return res.Http401(nil)
// 		}
// 	}
// 	return c.Next()
// }
