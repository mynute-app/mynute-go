package routes

import (
	"agenda-kaki-go/core/controller"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v2"
)

func Auth(Gorm *handler.Gorm, r fiber.Router) {
	ce := controller.Auth(Gorm)
	e := r.Group("/auth")
	mdw := middleware.Auth(Gorm)
	LoginRoutine := append(mdw.Login(), ce.Login)
	e.Post("/login", LoginRoutine...) // ok
	e.Post("/register", ce.Register)
	e.Post("/verify-existing-account", ce.VerifyExistingAccount)
	e.Get("/verifyemail", ce.VerifyEmail)
	e.Get("/oauth/logout", ce.LogoutProvider)
	e.Get("/oauth/:provider", ce.BeginAuthProviderCallback)
	e.Get("/oauth/:provider/callback", ce.GetAuthCallbackFunction)
}
