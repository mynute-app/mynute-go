package middleware

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
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

func (am *auth_middleware) Login() []fiber.Handler {
	return []fiber.Handler{
		lib.SaveBodyOnCtx[DTO.LoginUser],
		am.DenyLoginFromUnverified,
	}
}

func (am *auth_middleware) DenyUnauthorized(c *fiber.Ctx) error {
	res := lib.SendResponse{Ctx: c}
	if c.Get("Authorization") == "" {
		fmt.Printf("Access Denied!\n")
		return res.Http401(nil)
	}
	return c.Next()
}

func (am *auth_middleware) DenyLoginFromUnverified(c *fiber.Ctx) error {
	// Get login body from context
	login, err := lib.GetBodyFromCtx[*DTO.LoginUser](c)
	if err != nil {
		return err
	}
	// Get user from email
	user := &[]model.User{}
	am.Gorm.DB.Where("email = ?", login.Email).Find(user)
	if len(*user) == 0 {
		fmt.Printf("User %v not found\n", login.Email)
		return lib.Error.Auth.InvalidLogin.SendToClient(c)
	}
	// Check if user is verified
	if !(*user)[0].Verified {
		fmt.Printf("User %v is not verified\n", login.Email)
		return lib.Error.User.NotVerified.SendToClient(c)
	}
	return c.Next()
}

func (am *auth_middleware) WhoAreYou(c *fiber.Ctx) error {
	res := lib.SendResponse{Ctx: c}
	// if c.Get("Authorization") == "" {
	// 	err := handler.Auth(c).WhoAreYou()
	// 	if err != nil {
	// 		return res.Http401(err).Next()
	// 	}
	// 	return nil
	// }
	if c.Get("Authorization") != "" {
		err := handler.JWT(c).WhoAreYou()
		if err != nil {
			return res.Http401(nil)
		}
	}
	return c.Next()
}
