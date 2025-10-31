package routes

import (
	"mynute-go/core/src/controller"
	"mynute-go/core/src/handler"
	"mynute-go/core/src/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Builds all available routes for the application
func Build(DB *gorm.DB, App *fiber.App) {
	Prometheus(App)
	Swagger(App)

	Gorm := &handler.Gorm{DB: DB}

	controller.Admin(Gorm)
	controller.AdminAuth(Gorm)
	controller.Appointment(Gorm)
	controller.Auth(Gorm)
	controller.Branch(Gorm)
	controller.Client(Gorm)
	controller.Company(Gorm)
	controller.Employee(Gorm)
	controller.Holiday(Gorm)
	controller.Sector(Gorm)
	controller.Service(Gorm)

	r := App.Group("/")

	r.Get("/", controller.Home)
	r.Get("/verify-email", controller.VerifyEmailPage)
	r.Get("/translations/page/:page", controller.GetPageTranslations)

	endpoints := &middleware.Endpoint{DB: Gorm}
	if err := endpoints.Build(r); err != nil {
		panic(err)
	}
}
