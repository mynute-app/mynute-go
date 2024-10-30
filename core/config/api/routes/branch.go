package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/handlers"

	"github.com/gofiber/fiber/v3"
)

func Branch(Gorm *handlers.Gorm, r fiber.Router) {
	cb := controllers.Branch(Gorm)
	b := r.Group("/company/:companyId/branch")
	controllers.CreateRoutes(b, cb)
}