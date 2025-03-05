package middleware

import (
	"agenda-kaki-go/core/handler"

	"github.com/gofiber/fiber/v2"
)

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
