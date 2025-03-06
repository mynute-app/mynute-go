package routes

import (
	"agenda-kaki-go/core/controller"
	"agenda-kaki-go/core/handler"

	"github.com/gofiber/fiber/v2"
)

func Auth(Gorm *handler.Gorm, r fiber.Router) {
	ce := controller.Auth(Gorm)
	e := r.Group("/auth")
	e.Post("/login", ce.Login) // ok
	e.Post("/register", ce.Register)
	e.Post("/verify-existing-account", ce.VerifyExistingAccount)
	e.Post("/verify-email", ce.VerifyEmail)
	e.Get("/oauth/logout", ce.LogoutProvider)
	e.Get("/oauth/:provider", ce.BeginAuthProviderCallback)
	e.Get("/oauth/:provider/callback", ce.GetAuthCallbackFunction)
}
