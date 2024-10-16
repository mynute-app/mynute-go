package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/middleware"
	"agenda-kaki-go/core/handlers"

	"github.com/gofiber/fiber/v3"
)

func Employee(Gorm *handlers.Gorm, App *fiber.App) {
	Middleware := &middleware.Employee{Gorm: Gorm}
	HTTP := &handlers.HTTP{Gorm: Gorm}
	RequestHandler := &handlers.Request{HTTP: HTTP}
	Associations := []string{}
	ce := controllers.Employee{
		Request: RequestHandler,
		Middleware: Middleware,
		Associations: Associations,
	}
	r := App.Group("/employee")

	r.Post("/", ce.CreateOne) // ok
	r.Get("/", ce.GetAll) // ok
	r.Get("/:id", ce.GetOneById) // ok
	// r.Get("/name/:name", ce.GetOneByName) // ok
	r.Get("/email/:email", ce.GetOneByEmail) // ok
	r.Patch("/:id", ce.UpdateOneById) // ok
	r.Delete("/:id", ce.DeleteOneById) // ok
	r.Delete("/:id/force", ce.ForceDeleteOneById) // ok
}
