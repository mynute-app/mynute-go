package lib

import "github.com/gofiber/fiber/v3"

func FiberError(s int, c fiber.Ctx, err error) error {
	return c.Status(s).JSON(fiber.Map{"error": err.Error()})
}