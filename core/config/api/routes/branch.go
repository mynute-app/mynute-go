package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v3"
)

func Branch(Gorm *handlers.Gorm, App *fiber.App) {
	Middleware := &middleware.Branch{Gorm: Gorm}
	HTTP := &handlers.HTTP{Gorm: Gorm}
	RequestHandler := &handlers.Request{HTTP: HTTP}
	Associations := []string{}
	cb := controllers.Branch{
		Request: RequestHandler,
		Middleware: Middleware,
		Associations: Associations,
	}
	r := App.Group("/company/:companyId/branch")

	r.Post("/", cb.CreateOne) // ok
	r.Get("/", cb.GetAll) // ok
	r.Get("/:id", cb.GetOneById) // ok
	r.Delete("/:id", cb.DeleteOneById) // ok
	r.Delete("/:id/force", cb.ForceDeleteOneById) // ok
	r.Patch("/:id", cb.UpdateOneById) // ok
	
}