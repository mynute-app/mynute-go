package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/handlers"

	"github.com/gofiber/fiber/v3"
)

func Employee(Gorm *handlers.Gorm, r fiber.Router) {
	ce := controllers.Employee(Gorm)
	e := r.Group("/employee")

	e.Get("/email/:email", ce.GetOneByEmail) // ok

	controllers.CreateRoutes(e, ce)
}
