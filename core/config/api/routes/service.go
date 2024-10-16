package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/middleware"
	"agenda-kaki-go/core/handlers"

	"github.com/gofiber/fiber/v3"
)

func Service(Gorm *handlers.Gorm, App *fiber.App) {
	Middleware := &middleware.Service{Gorm: Gorm}
	HTTP := &handlers.HTTP{Gorm: Gorm}
	RequestHandler := &handlers.Request{HTTP: HTTP}
	Associations := []string{"ServiceType"}
	cs := controllers.Service{
		Request: RequestHandler,
		Middleware: Middleware,
		Associations: Associations,
	}
	r := App.Group("/service")

	r.Post("/", cs.CreateOne) // ok
	r.Get("/", cs.GetAll) // ok
	r.Get("/:id", cs.GetOneById) // ok
	// r.Get("/name/:name", cs.GetOneByName) // ok
	r.Delete("/:id", cs.DeleteOneById) // ok
	r.Delete("/:id/force", cs.ForceDeleteOneById) // ok
	r.Patch("/:id", cs.UpdateOneById) // ok

}