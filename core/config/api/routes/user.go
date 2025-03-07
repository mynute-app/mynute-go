package routes

import (
	"agenda-kaki-go/core/controller"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v2"
)

func User(Gorm *handler.Gorm, r fiber.Router) {
	ce := controller.User(Gorm)
	e := r.Group("/user")
	mu := middleware.User(Gorm)
	e.Post("/", append(mu.Create(), ce.CreateUser)...)               // ok
	e.Get("/email/:email", ce.GetOneByEmail) // ok
	e.Patch("/:id", ce.UpdateUserById)       // ok
	e.Delete("/:id", ce.DeleteUserById)      // ok
}
