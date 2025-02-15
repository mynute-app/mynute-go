package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func Swagger(App *fiber.App) {
	App.Get("/swagger/*", swagger.HandlerDefault)
}
