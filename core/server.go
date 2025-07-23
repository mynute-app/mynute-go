package core

import (
	"fmt"
	"log"
	"log/slog"
	"mynute-go/core/config/api/routes"
	database "mynute-go/core/config/db"
	"mynute-go/core/config/db/model"
	"mynute-go/core/lib"
	myUploader "mynute-go/core/lib/cloud_uploader"
	"mynute-go/core/middleware"
	"os"

	"github.com/gofiber/fiber/v2"
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
	db.Migrate(model.GeneralModels)
	db.InitialSeed()
	routes.Build(db.Gorm, app)
	if err := myUploader.StartProvider(); err != nil {
		panic(err)
	}
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
	if in == "parallel" {
		app_env := os.Getenv("APP_ENV")
		if app_env == "prod" {
			log.Fatal("Server run for production can not be in parallel. For parallel running set APP_ENV=test or APP_ENV=dev at .env file")
		}
		s.parallel()
	} else if in == "listen" {
		s.listen()
	} else {
		log.Fatal("Server run mode not recognized. Please, provide a valid argument")
	}
	return s
}
