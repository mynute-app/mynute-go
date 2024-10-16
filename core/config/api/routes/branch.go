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
	r := App.Group("/branch")

	r.Post("/", cb.CreateOne) // ok
	r.Get("/", cb.GetAll) // ok
	r.Get("/:id/:companyId", cb.GetOneById) // ok
	// r.Get("/name/:name", cb.GetOneByName) // ok
	r.Delete("/:id/:companyId", cb.DeleteOneById) // ok
	r.Delete("/:id/:companyId/force", cb.ForceDeleteOneById) // ok
	r.Patch("/:id/:companyId", cb.UpdateOneById) // ok
	
}