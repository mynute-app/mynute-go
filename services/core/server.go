package core

import (
	"fmt"
	"log"
	"log/slog"
	"mynute-go/debug"
	"mynute-go/services/core/api/lib"
	myUploader "mynute-go/services/core/api/lib/cloud_uploader"
	"mynute-go/services/core/api/middleware"
	"mynute-go/services/core/api/routes"
	database "mynute-go/services/core/config/db"
	"mynute-go/services/core/config/db/model"
	"os"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/gofiber/swagger"
)

type Server struct {
	App *fiber.App
	Db  *database.Database
}

// Creates a new server instance
func NewServer() *Server {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	app := fiber.New(fiber.Config{
		ErrorHandler:          middleware.ErrorV13(logger),
		BodyLimit:             2 * 1024 * 1024, // 2 MB
		DisableStartupMessage: true,
	})
	app.Use(middleware.LogV13(logger))
	lib.LoadEnv()
	db := database.Connect()
	// session := handler.NewCookieStore(handler.SessionOpts())
	// handler.NewAuth(session)

	// Only run migrations and seeding in dev/test environments
	// Production migrations must be run manually before deployment
	// See: docs/MIGRATIONS.md or run `make migrate-help`
	app_env := os.Getenv("APP_ENV")
	if app_env == "dev" || app_env == "test" {
		// NOTE: Auth database migrations are handled by the auth service
		// Core service connects to auth DB in read-only mode for permission checks

		// Migrate main database (business data)
		log.Println("Migrating main database...")
		db.Migrate(model.MainDBModels)

		// Run initial seeding (core business data only)
		db.InitialSeed()
	} else {
		log.Println("Production environment detected - skipping automatic migrations and seeding")
		log.Println("Run migrations manually with: make migrate-up (see docs/MIGRATIONS.md)")
	}

	app.Static("/", "./static")

	// Transpile TypeScript files on-the-fly (only for .ts files)
	app.Use("/admin", middleware.TypeScriptTranspiler())

	// Serve admin panel static files
	app.Static("/admin", "./admin")

	routes.Build(db.Gorm, app)

	// Swagger documentation
	app.Get("/swagger/*", fiberSwagger.HandlerDefault)

	// Setup live reload for development (admin panel)
	middleware.SetupLiveReload(app)

	if err := myUploader.StartProvider(); err != nil {
		panic(err)
	}
	debug.Clear()
	return &Server{App: app, Db: db}
}

func (s *Server) Shutdown() {
	// Check if server is already running
	if s.App.Handler() == nil {
		return
	}
	if err := s.App.Shutdown(); err != nil {
		fmt.Printf("Server did not shutdown gracefully: %v", err)
	}
	s.Db.Test().Clear()
	s.Db.Disconnect()
	fmt.Printf("Finished server shutdown procedure. \n")
}

func (s *Server) parallel() *Server {
	go func() {
		s.listen()
	}()
	return s
}

func (s *Server) listen() *Server {
	app_port := os.Getenv("APP_PORT")
	log.Printf("Server is starting at http://localhost:%s\n", app_port)
	if err := s.App.Listen(":" + app_port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
	return s
}

// Runs the server in two modes: test and listen
//
//	@parallel:	starts the server in a goroutine. This is useful for unit testing.
//	@listen:	starts the server and listens for incoming requests. This is useful for production or normal dev.
func (s *Server) Run(in string) *Server {
	log.Printf("Starting server in '%s' mode...\n", in)
	switch in {
	case "parallel":
		app_env := os.Getenv("APP_ENV")
		if app_env == "prod" {
			log.Fatal("Server run for production can not be in parallel. For parallel running set APP_ENV=test or APP_ENV=dev at .env file")
		} else if app_env != "test" && app_env != "dev" {
			log.Fatal("Server run for parallel can only be in test or dev environment. For parallel running set APP_ENV=test or APP_ENV=dev at .env file")
		}
		s.parallel()
	case "listen":
		s.listen()
	default:
		log.Fatal("Server run mode not recognized. Please, provide a valid argument")
	}
	return s
}
