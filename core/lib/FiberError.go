package lib

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

func FiberError(s int, c fiber.Ctx, err error) error {
 	return c.Status(s).JSON(fiber.Map{"error": err.Error()})
}

func Fiber400(c fiber.Ctx, err error) error {
	return c.Status(400).JSON(fiber.Map{"error": err.Error()})
}

func Fiber404(c fiber.Ctx) error {
	return c.SendStatus(404)
}

func Fiber500(c fiber.Ctx, err error) error {
	log.Printf("An internal error occurred! \n Error: %v", err)
	return c.SendStatus(500)
}

func Fiber201(c fiber.Ctx) error {
	return c.SendStatus(201)
}

func Fiber204(c fiber.Ctx) error {
	return c.SendStatus(204)
}