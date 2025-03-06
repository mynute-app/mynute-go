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
	Routine *auth_mdw_routines
}

func Auth(Gorm *handler.Gorm) *auth_middleware {
	return &auth_middleware{Gorm: Gorm}
}

type auth_mdw_routines struct {
	Login *func () []fiber.Handler
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
	login, err := lib.GetFromCtx[*DTO.LoginUser](c, namespace.RequestKey.Body_Parsed)
	if err != nil {
		return err
	}
	// Get user from email
	user := &[]model.User{}
	am.Gorm.DB.Where("email = ?", login.Email).Find(user)
	fmt.Printf("User: %+v\n", *user)
	if len(*user) == 0 {
		return lib.MyErrors.InvalidLogin.SendToClient(c)
	}
	// Check if user is verified
	if !(*user)[0].Verified {
		return lib.MyErrors.UserNotVerified.SendToClient(c)
	}
	fmt.Println("User is verified")
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
