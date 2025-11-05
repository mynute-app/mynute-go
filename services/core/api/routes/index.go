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
	
	// TODO: In production, endpoint authorization should be handled by the auth service API
	// For now, we'll set up basic routes manually for testing
	
	// Company routes (public schema)
	apiRoutes.Post("/company", middleware.SavePublicSession(DB), middleware.ChangeToPublicSchema, controller.CreateCompany)
	apiRoutes.Get("/company/:id", middleware.SavePublicSession(DB), middleware.ChangeToPublicSchema, controller.GetCompanyById)
	apiRoutes.Get("/company/name/:name", middleware.SavePublicSession(DB), middleware.ChangeToPublicSchema, controller.GetCompanyByName)
	apiRoutes.Get("/company/tax/:taxId", middleware.SavePublicSession(DB), middleware.ChangeToPublicSchema, controller.GetCompanyByTaxId)
	apiRoutes.Get("/company/subdomain/:subdomain", middleware.SavePublicSession(DB), middleware.ChangeToPublicSchema, controller.GetCompanyBySubdomain)
	apiRoutes.Put("/company/:id", middleware.SavePublicSession(DB), middleware.ChangeToPublicSchema, controller.UpdateCompanyById)
	apiRoutes.Delete("/company/:id", middleware.SavePublicSession(DB), middleware.ChangeToPublicSchema, controller.DeleteCompanyById)
	
	log.Println("Routes registered manually (endpoint middleware disabled)")
}
