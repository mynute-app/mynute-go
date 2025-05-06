package controller

import (
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"
)

// OAUTH logics
func BeginAuthProviderCallback(c *fiber.Ctx) error {
	if err := goth_fiber.BeginAuthHandler(c); err != nil {
		return err
	}
	return nil
}

func GetAuthCallbackFunction(c *fiber.Ctx) error {
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

func LogoutProvider(c *fiber.Ctx) error {
	if err := goth_fiber.Logout(c); err != nil {
		return err
	}
	if err := c.Redirect("/"); err != nil {
		return err
	}
	return nil
}

func Auth(Gorm *handler.Gorm) {
	endpoint := &middleware.Endpoint{DB: Gorm}
	endpoint.BulkRegisterHandler([]fiber.Handler{
		BeginAuthProviderCallback,
		GetAuthCallbackFunction,
		LogoutProvider,
	})
}
