package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/handlers"

	"github.com/gofiber/fiber/v2"
)

func Service(Gorm *handlers.Gorm, r fiber.Router) {
	cs := controllers.Service(Gorm)
	s := r.Group("/service")
	s.Post("/", cs.CreateService)             // ok
	s.Get("/:id", cs.GetServiceById)           // ok
	s.Get("/name/:name", cs.GetServiceByName)  // ok
	s.Patch("/:id", cs.UpdateServiceById)      // ok
	s.Delete("/:id", cs.DeleteServiceById)     // ok
}
