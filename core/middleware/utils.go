package middleware

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"

	"github.com/gofiber/fiber/v2"
)

func SaveReqInfo(c *fiber.Ctx) error {
	// Save request information here
	c.Locals(namespace.RequestKey.Body, c.Body())
	return c.Next()
}