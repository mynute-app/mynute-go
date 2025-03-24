package routes

import (
	"agenda-kaki-go/core/controller"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v2"
)

func Client(Gorm *handler.Gorm, r fiber.Router) {
	ce := controller.Client(Gorm)
	e := r.Group("/client")
	e.Post("/", ce.CreateClient)
	e.Post("/verify-email/:email/:code", ce.VerifyClientEmail)
	e.Post("/login", ce.LoginClient)
	auth := e.Group("/", middleware.Auth(Gorm).DenyUnauthorized)
	auth.Get("/email/:email", ce.GetClientByEmail)
	auth.Patch("/:id", ce.UpdateClientById)
	auth.Delete("/:id", ce.DeleteClientById)
}
