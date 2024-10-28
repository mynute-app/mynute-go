package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/handlers"

	"github.com/gofiber/fiber/v3"
)

func Employee(Gorm *handlers.Gorm, App *fiber.App) {
	ce := controllers.Employee(Gorm)
	r := App.Group("/employee")

	r.Get("/email/:email", ce.GetOneByEmail) // ok

	controllers.CreateRoutes(r, ce)
}
