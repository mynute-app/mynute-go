package middleware

import (
	"agenda-kaki-go/core/handlers"

	"github.com/gofiber/fiber/v2"
)

func WhoAreYou(c *fiber.Ctx) error {
	//verify if is oauth or jwt
	if c.Get("Authorization") == "" {
		return handlers.Auth(c).WhoAreYou()
	}
	return handlers.JWT(c).WhoAreYou()
}
