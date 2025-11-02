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
	controller.Appointment(Gorm)
	controller.Auth(Gorm)
	controller.Branch(Gorm)
	controller.Client(Gorm)
	controller.Company(Gorm)
	controller.Employee(Gorm)
	controller.Holiday(Gorm)
	controller.Sector(Gorm)
	controller.Service(Gorm)

	// Public routes (no /api prefix)
	publicRoutes := App.Group("/")
	publicRoutes.Get("/", controller.Home)
	publicRoutes.Get("/verify-email", controller.VerifyEmailPage)
	publicRoutes.Get("/admin/verify-email/:email/:code", controller.VerifyEmailPage)
	publicRoutes.Get("/translations/page/:page", controller.GetPageTranslations)

	// API routes (with /api prefix)
	apiRoutes := App.Group("/api")
	endpoints := &middleware.Endpoint{DB: Gorm}
	if err := endpoints.Build(apiRoutes); err != nil {
		panic(err)
	}
}
