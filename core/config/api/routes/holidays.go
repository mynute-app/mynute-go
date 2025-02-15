package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/handlers"

	"github.com/gofiber/fiber/v2"
)

func Holidays(Gorm *handlers.Gorm, App *fiber.App) {
	ce := controllers.Holidays(Gorm)
	r := App.Group("/holidays")

	controllers.CreateRoutes(r, ce)
}
