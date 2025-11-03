package main

import (
	"log"
	"mynute-go/auth/api/routes"
	database "mynute-go/core/src/config/db"
	"mynute-go/core/src/config/db/model"
	"mynute-go/core/src/lib"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// @title						Auth Service API
// @version					1.0
// @description				Authentication and Authorization Service
// @termsOfService				http://swagger.io/terms/
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						X-Auth-Token
// @description				Enter the token in the format: <token>
// @contact.name				API Support
// @contact.email				auth@mynute.com
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
// @host						localhost:4001
// @BasePath					/
func main() {
	log.Println("Starting Auth Service...")

	// Load environment variables
	lib.LoadEnv()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, X-Auth-Token, X-Company-ID",
	}))
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))

	// Connect to auth database
	db := database.Connect()

	// Run migrations and seeding in dev/test environments
	app_env := os.Getenv("APP_ENV")
	if app_env == "dev" || app_env == "test" {
		log.Println("Migrating auth database...")
		db.WithDB(db.AuthDB).Migrate(model.AuthDBModels)

		// Run initial seeding (auth-related data only)
		log.Println("Seeding auth database...")
		db.InitialSeed()
	} else {
		log.Println("Production environment - skipping automatic migrations")
	}

	// Setup routes
	routes.SetupAuthRoutes(app, db.AuthDB)

	// Get port from env or default to 4001
	port := os.Getenv("AUTH_SERVICE_PORT")
	if port == "" {
		port = "4001"
	}

	// Start server
	log.Printf("Auth Service running on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start auth service: %v", err)
	}
}

