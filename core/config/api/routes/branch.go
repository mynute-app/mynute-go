package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/handlers"

	"github.com/gofiber/fiber/v3"
)

func Branch(Gorm *handlers.Gorm, App *fiber.App) {
	cb := controllers.Branch(Gorm)
	r := App.Group("/company/:companyId/branch")
	controllers.CreateRoutes(r, cb)
}