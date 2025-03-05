package routes

import (
	"agenda-kaki-go/core/controller"
	"agenda-kaki-go/core/handler"

	"github.com/gofiber/fiber/v2"
)

func Service(Gorm *handler.Gorm, r fiber.Router) {
	cs := controller.Service(Gorm)
	s := r.Group("/service")
	s.Post("/", cs.CreateService)             // ok
	s.Get("/:id", cs.GetServiceById)          // ok
	s.Get("/name/:name", cs.GetServiceByName) // ok
	s.Patch("/:id", cs.UpdateServiceById)     // ok
	s.Delete("/:id", cs.DeleteServiceById)    // ok
}
