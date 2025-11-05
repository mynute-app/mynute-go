package routes

import (
	"log"
	"mynute-go/services/core/api/controller"
	"mynute-go/services/core/api/handler"
	"mynute-go/services/core/api/middleware"

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

	// Fetch endpoints from auth service API and register routes
	// TODO: Implement HTTP client to fetch endpoints from http://localhost:4001/api/endpoints
	// For now, log that routes need to be registered via auth service
	log.Println("TODO: Fetch endpoints from auth service API at http://localhost:4001/api/endpoints")
	log.Println("Routes will be registered dynamically after auth service API integration")

	_ = apiRoutes             // Prevent unused variable error
	_ = middleware.Endpoint{} // Prevent unused import error
}
