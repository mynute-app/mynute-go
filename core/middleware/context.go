package middleware

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"

	"github.com/gofiber/fiber/v2"
)

func ParseBodyToContext[T any](c *fiber.Ctx, model T) (int, error) {
	bodyKey := namespace.GeneralKey.Model
	// Ignore in case the body has been already parsed into context.
	if b := c.Locals(bodyKey); b != nil {
		return 0, nil
	}
	method := c.Method()
	if method != "POST" && method != "PUT" && method != "PATCH" {
		c.Locals(bodyKey, model)
		return 0, nil
	}
	if err := lib.BodyParser(c.Body(), model); err != nil {
		return 500, err
	}
	c.Locals(bodyKey, model)
	return 0, nil
}
