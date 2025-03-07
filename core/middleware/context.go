package middleware

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"

	"github.com/gofiber/fiber/v2"
)

func ParseBodyToContext[T any](c *fiber.Ctx, model T) (int, error) {
	// Ignore in case the body has been already parsed into context.
	if b := c.Locals(namespace.RequestKey.Body_Parsed); b != nil {
		return 0, nil
	}
	bodyKey := namespace.RequestKey.Body_Parsed
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
