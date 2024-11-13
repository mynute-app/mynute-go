package middleware

import (
	"agenda-kaki-go/core/handlers"

	"github.com/gofiber/fiber/v3"
)

func WhoAreYou(c fiber.Ctx) error {
	return handlers.JWT(c).WhoAreYou()
}

