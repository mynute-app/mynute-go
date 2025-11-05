package routes

import (
	"log"
	"mynute-go/services/core/api/controller"
	"mynute-go/services/core/api/handler"
	authClient "mynute-go/services/core/api/lib/auth_client"
	"mynute-go/services/core/api/middleware"
	"os"

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

	// Fetch endpoints from auth service API and register routes dynamically
	client := authClient.NewAuthClient()
	var endpoints []*authClient.EndPoint
	var err error

	// Check if auth service is available
	if !client.IsAvailable() {
		appEnv := os.Getenv("APP_ENV")
		if appEnv == "test" {
			// In test mode, use hardcoded test endpoints
			log.Println("Warning: Auth service is not available in test mode.")
			log.Println("Using hardcoded test endpoints for route registration.")
			endpoints = GetTestEndpoints()
		} else {
			// In production/dev, auth service is required
			log.Println("Error: Auth service is not available. Routes will not be registered.")
			log.Println("Make sure auth service is running at", client.BaseURL)
			return
		}
	} else {
		// Auth service is available, fetch endpoints from API
		endpoints, err = client.FetchEndpoints()
		if err != nil {
			log.Printf("Error fetching endpoints from auth service: %v\n", err)
			log.Println("Routes will not be registered. Please check auth service.")
			return
		}
		log.Printf("Successfully fetched %d endpoints from auth service\n", len(endpoints))
	}

	// Register routes using the endpoint middleware
	ep := &middleware.Endpoint{DB: Gorm}
	if err := ep.BuildFromAPI(apiRoutes, endpoints); err != nil {
		log.Printf("Error building routes: %v\n", err)
	}
}
