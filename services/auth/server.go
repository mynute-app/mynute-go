package auth

import (
	"fmt"
	"log"
	"log/slog"
	"mynute-go/services/auth/api/routes"
	database "mynute-go/services/auth/config/db"
	"mynute-go/services/auth/config/db/model"
	"mynute-go/services/auth/lib"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	fiberSwagger "github.com/gofiber/swagger"
)

type Server struct {
	App *fiber.App
	Db  *database.Database
}

// Creates a new auth server instance
func NewServer() *Server {
	slogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			slogger.Error("Request error", "error", err.Error(), "path", c.Path(), "method", c.Method())
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
		BodyLimit: 1 * 1024 * 1024, // 1 MB
	})

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, X-Auth-Token, X-Company-ID",
	}))

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))

	lib.LoadEnv()
	db := database.Connect()

	// Only run migrations and seeding in dev/test environments
	// Production migrations must be run manually before deployment
	app_env := os.Getenv("APP_ENV")
	if app_env == "dev" || app_env == "test" {
		// Migrate auth database (endpoints, policies, resources, users, roles)
		log.Println("Migrating auth database...")
		db.WithDB(db.Gorm).Migrate(model.AuthDBModels)

		// Run initial seeding (auth-related data only)
		log.Println("Seeding auth database...")
		db.InitialSeed()
	} else {
		log.Println("Production environment detected - skipping automatic migrations and seeding")
		log.Println("Run migrations manually with: make migrate-up (see docs/MIGRATIONS.md)")
	}

	// Setup auth routes
	routes.SetupAuthRoutes(app, db.Gorm)

	// Swagger documentation
	app.Get("/swagger/*", fiberSwagger.HandlerDefault)

	return &Server{App: app, Db: db}
}

func (s *Server) Shutdown() {
	// Check if server is already running
	if s.App.Handler() == nil {
		return
	}
	if err := s.App.Shutdown(); err != nil {
		fmt.Printf("Auth server did not shutdown gracefully: %v", err)
	}
	s.Db.Test().Clear()
	s.Db.Disconnect()
	fmt.Printf("Finished auth server shutdown procedure. \n")
}

func (s *Server) parallel() *Server {
	go func() {
		s.listen()
	}()
	return s
}

func (s *Server) listen() *Server {
	auth_port := os.Getenv("AUTH_SERVICE_PORT")
	if auth_port == "" {
		auth_port = "4001"
	}
	log.Printf("Auth Service is starting at http://localhost:%s\n", auth_port)
	if err := s.App.Listen(":" + auth_port); err != nil {
		log.Fatalf("Auth Service failed to start: %v", err)
	}
	return s
}

// Runs the server in two modes: test and listen
//
//	@parallel:	starts the server in a goroutine. This is useful for unit testing.
//	@listen:	starts the server and listens for incoming requests. This is useful for production or normal dev.
func (s *Server) Run(in string) *Server {
	log.Printf("Starting auth server in '%s' mode...\n", in)
	switch in {
	case "parallel":
		app_env := os.Getenv("APP_ENV")
		if app_env == "prod" {
			log.Fatal("Auth server run for production can not be in parallel. For parallel running set APP_ENV=test or APP_ENV=dev at .env file")
		} else if app_env != "test" && app_env != "dev" {
			log.Fatal("Auth server run for parallel can only be in test or dev environment. For parallel running set APP_ENV=test or APP_ENV=dev at .env file")
		}
		s.parallel()
	case "listen":
		s.listen()
	default:
		log.Fatal("Auth server run mode not recognized. Please, provide a valid argument")
	}
	return s
}
