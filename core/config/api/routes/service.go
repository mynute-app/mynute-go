package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/handlers"

	"github.com/gofiber/fiber/v3"
)

func Service(Gorm *handlers.Gorm, r fiber.Router) {
	cs := controllers.Service(Gorm)
	s := r.Group("/service")

	controllers.CreateRoutes(s, cs)
}