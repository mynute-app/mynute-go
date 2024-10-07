package middleware

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"

	"github.com/gofiber/fiber/v3"
)

func ParseBodyToContext[T any](key namespace.ContextKey, model *T) func(c fiber.Ctx) (int, error) {
	return func(c fiber.Ctx) (int, error) {
		if err := lib.BodyParser(c.Body(), model); err != nil {
			return 500, err
		}
		c.Locals(key, model)
		return 0, nil
	}
}

func AddToContext[T any](key namespace.ContextKey, model *T) func(c fiber.Ctx) (int, error) {
	return func(c fiber.Ctx) (int, error) {
		c.Locals(key, model)
		return 0, nil
	}
}