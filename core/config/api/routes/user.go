package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/handlers"

	"github.com/gofiber/fiber/v2"
)

func User(Gorm *handlers.Gorm, r fiber.Router) {
	ce := controllers.User(Gorm)
	e := r.Group("/user")
	e.Post("/", ce.CreateUser)                    // ok
	e.Post("/login", ce.Login)               // ok
	e.Get("/email/:email", ce.GetOneByEmail) // ok
	e.Patch("/:id", ce.UpdateUserById)         // ok
	e.Delete("/:id", ce.DeleteUserById)        // ok
}
