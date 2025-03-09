package core

import (
	"agenda-kaki-go/core/config/api/routes"
	database "agenda-kaki-go/core/config/db"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	App *fiber.App
	Db  *database.Database
}

// Creates a new server instance
func NewServer() *Server {
	app := fiber.New()
	lib.LoadEnv()
	db := database.Connect()
	session := handler.NewCookieStore(handler.SessionOpts())
	handler.NewAuth(session)
	db.Migrate()
	routes.Build(db.Gorm, app)
	return &Server{App: app, Db: db}
}

func (s *Server) Shutdown() {
	if err := s.App.Shutdown(); err != nil {
		fmt.Printf("Server did not shutdown gracefully: %v", err)
	}
	s.Db.Test().Clear()
	s.Db.Disconnect()
	fmt.Printf("Finished server shutdown procedure. \n")
}

func (s *Server) parallel() *Server {
	go func() {
		if err := s.App.Listen(":" + namespace.AppPort); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()
	return s
}
// Runs the server in two modes: test and listen
// @test: starts the server in a goroutine. This is useful for unit testing.
// @listen: starts the server and listens for incoming requests. This is useful for production or normal dev.
func (s *Server) Run(in string) *Server {
	if in == "test" {
		app_env := os.Getenv("APP_ENV")
		if app_env != "test" {
			log.Fatalf("Server run for tests must have APP_ENV as 'test'. Currently is '%s'.\nPlease, set APP_ENV=test at .env file", app_env)
		}
		s.parallel()
	} else if in == "listen" {
		log.Fatal(s.App.Listen(":" + namespace.AppPort))
	}
	return s
}
