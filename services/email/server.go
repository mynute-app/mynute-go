package email

import (
	"fmt"
	"log"
	"log/slog"
	"mynute-go/services/core/src/lib"
	_ "mynute-go/services/email/docs" // Import swagger docs
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	swagger "github.com/gofiber/swagger"
)

type Server struct {
	App *fiber.App
}

// Creates a new email server instance
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
		BodyLimit: 2 * 1024 * 1024, // 2 MB for email attachments
	})

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, X-API-Key",
	}))

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))

	lib.LoadEnv()

	// Initialize email services
	if err := initEmailServices(); err != nil {
		log.Fatalf("Failed to initialize email services: %v", err)
	}

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "email",
		})
	})

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Setup email routes
	setupRoutes(app)

	return &Server{App: app}
}

func (s *Server) Shutdown() {
	if s.App.Handler() == nil {
		return
	}
	if err := s.App.Shutdown(); err != nil {
		fmt.Printf("Email server did not shutdown gracefully: %v", err)
	}
	fmt.Printf("Finished email server shutdown procedure. \n")
}

func (s *Server) parallel() *Server {
	go func() {
		s.listen()
	}()
	return s
}

func (s *Server) listen() *Server {
	email_port := os.Getenv("EMAIL_SERVICE_PORT")
	if email_port == "" {
		email_port = "4002"
	}
	log.Printf("Email Service is starting at http://localhost:%s\n", email_port)
	if err := s.App.Listen(":" + email_port); err != nil {
		log.Fatalf("Email Service failed to start: %v", err)
	}
	return s
}

// Runs the server in two modes: test and listen
func (s *Server) Run(in string) *Server {
	log.Printf("Starting email server in '%s' mode...\n", in)
	switch in {
	case "parallel":
		app_env := os.Getenv("APP_ENV")
		if app_env == "prod" {
			log.Fatal("Email server run for production can not be in parallel. For parallel running set APP_ENV=test or APP_ENV=dev at .env file")
		} else if app_env != "test" && app_env != "dev" {
			log.Fatal("Email server run for parallel can only be in test or dev environment. For parallel running set APP_ENV=test or APP_ENV=dev at .env file")
		}
		s.parallel()
	case "listen":
		s.listen()
	default:
		log.Fatal("Email server run mode not recognized. Please, provide a valid argument")
	}
	return s
}
