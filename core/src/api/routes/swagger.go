package routes

import (
	"mynute-go/core/src/api/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func Swagger(App *fiber.App) {
	App.Get("/swagger/*", middleware.SwaggerAuth(), swagger.HandlerDefault)
}

