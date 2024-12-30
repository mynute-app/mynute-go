package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/handlers"

	"github.com/gofiber/fiber/v3"
)

func Auth(Gorm *handlers.Gorm, r fiber.Router) {
	ce := controllers.Auth(Gorm)
	e := r.Group("/auth")
	e.Post("/login", ce.Login) // ok
	e.Post("/register", ce.Register)
	e.Post("/verify-email", ce.VerifyEmail)
	e.Get("/logout", ce.LogoutProvider)
	e.Get("/oauth/:provider", ce.BeginAuthProviderCallback)
	e.Get("/oauth/:provider/callback", ce.GetAuthCallbackFunction)

	controllers.CreateRoutes(e, ce)
}
