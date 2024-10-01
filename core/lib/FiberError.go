package lib

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

func FiberError(s int, c fiber.Ctx, err error) error {
	if s == 500 {
		log.Printf("An internal error occurred! \n Error: %v", err)
	}
	return c.Status(s).JSON(fiber.Map{"error": err.Error()})
}