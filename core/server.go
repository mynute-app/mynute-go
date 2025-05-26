package core

import (
	"agenda-kaki-go/core/config/api/routes"
	database "agenda-kaki-go/core/config/db"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/middleware"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Server struct {
	App *fiber.App
	Db  *database.Database
}

// Creates a new server instance
func NewServer() *Server {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	fiberConfig := fiber.Config{
		ErrorHandler:          middleware.Error(logger),
		BodyLimit:             2 * 1024 * 1024, // 2 MB
		DisableStartupMessage: true,
	}
	app := fiber.New(fiberConfig)
	app.Use(middleware.Log(logger))
	lib.LoadEnv()
	db := database.Connect()
	session := handler.NewCookieStore(handler.SessionOpts())
	handler.NewAuth(session)
	app_env := os.Getenv("APP_ENV")
	if app_env == "test" {
		db.Test().Clear()
	}
	if app_env == "dev" || app_env == "test" {
		app.Static(namespace.StaticServerFolder, namespace.UploadsFolder)
	}
	db.Migrate(model.GeneralModels)
	if err := Seed(db.Gorm); err != nil {
		panic(err)
	}
	routes.Build(db.Gorm, app)
	return &Server{App: app, Db: db}
}

func Seed(db *gorm.DB) error {
	tx, end, err := database.Transaction(db)
	defer end()
	if err != nil {
		return err
	}

	Database := &database.Database{Gorm: tx}

	if err := Database.
		Seed("Resources", model.Resources, `"table" = ?`, []string{"Table"}).
		Seed("Roles", model.Roles, "name = ? AND company_id IS NULL", []string{"Name"}).
		Error; err != nil {
		return err
	}

	if err := model.LoadSystemRoleIDs(tx); err != nil {
		return fmt.Errorf("failed to load system role IDs: %w", err)
	}

	endpoints, deferEndpoint, err := model.EndPoints(&model.EndpointCfg{AllowCreation: true}, tx)
	if err != nil {
		return err
	}
	defer deferEndpoint()

	if err := Database.
		Seed("Endpoints", endpoints, "method = ? AND path = ?", []string{"Method", "Path"}).
		Error; err != nil {
		return err
	}

	policies, deferPolicies := model.Policies(&model.PolicyCfg{AllowNilCompanyID: true, AllowNilCreatedBy: true})
	defer deferPolicies()
	if err := Database.
		Seed("Policies", policies, "name = ?", []string{"Name"}).
		Error; err != nil {
		return err
	}

	return nil
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
	if err := s.App.Listen(":" + app_port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
	return s
}

// Runs the server in two modes: test and listen
//
//	@test:		starts the server in a goroutine. This is useful for unit testing.
//	@listen:	starts the server and listens for incoming requests. This is useful for production or normal dev.
func (s *Server) Run(in string) *Server {
	if in == "test" {
		app_env := os.Getenv("APP_ENV")
		if app_env != "test" {
			log.Fatalf("Server run for tests must have APP_ENV as 'test'. Currently is '%s'.\nPlease, set APP_ENV=test at .env file", app_env)
		}
		s.parallel()
	} else if in == "listen" {
		s.listen()
	} else {
		log.Fatalf("Server run mode not recognized. Please, use 'test' or 'listen' as argument.")
	}
	return s
}
