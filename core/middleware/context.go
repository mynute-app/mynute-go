package middleware

import (
	"agenda-kaki-go/core/lib"

	"github.com/gofiber/fiber/v3"
)

func ParseBodyToContext[T any](c fiber.Ctx, key string, model T) (int, error) {
	method := c.Method()
	if method != "POST" && method != "PUT" && method != "PATCH" {
		c.Locals(key, model)
		return 0, nil
	}
	if err := lib.BodyParser(c.Body(), &model); err != nil {
		return 500, err
	}
	c.Locals(key, model)
	return 0, nil
}