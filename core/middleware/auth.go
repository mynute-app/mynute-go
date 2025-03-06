package middleware

import (
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
	res := lib.SendResponse{Ctx: c}
	if c.Get("Authorization") == "" {
		fmt.Printf("Access Denied!\n")
		return res.Http401(nil)
	}
	return c.Next()
}

func (am *auth_middleware) WhoAreYou(c *fiber.Ctx) error {
	res := lib.SendResponse{Ctx: c}
	fmt.Printf("Authorization: %s\n", c.Get("Authorization"))
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
			return res.Http401(err)
		}
	}
	return c.Next()
}

func WhoAreYou(c *fiber.Ctx) (int, error) {
	//verify if is oauth or jwt
	if c.Get("Authorization") == "" {
		err := handler.Auth(c).WhoAreYou()
		if err != nil {
			return 401, err
		}
		return 0, nil
	}
	err := handler.JWT(c).WhoAreYou()
	if err != nil {
		return 401, err
	}
	return 0, nil
}
