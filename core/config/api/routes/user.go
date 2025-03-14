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
	e.Post("/", ce.CreateUser)
	e.Post("/verify-email/:email/:code", ce.VerifyUserEmail)
	e.Post("/login", ce.LoginUser)
	auth := e.Group("/", middleware.Auth(Gorm).DenyUnauthorized)
	auth.Get("/email/:email", ce.GetUserByEmail)
	auth.Patch("/:id", ce.UpdateUserById)       
	auth.Delete("/:id", ce.DeleteUserById)
}
