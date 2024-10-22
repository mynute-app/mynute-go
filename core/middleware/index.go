package middleware

import "github.com/gofiber/fiber/v3"

type IMiddleware interface {
	GET() []func(fiber.Ctx) (int, error)
	POST() []func(fiber.Ctx) (int, error)
	PATCH() []func(fiber.Ctx) (int, error)
	DELETE() []func(fiber.Ctx) (int, error)
	ForceGET() []func(fiber.Ctx) (int, error)
	ForceDELETE() []func(fiber.Ctx) (int, error)
}