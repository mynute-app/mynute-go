package routes

import (
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/middleware"
	"log/slog"
	"os"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Builds all available routes for the application
func Build(DB *gorm.DB, App *fiber.App) {
	Gorm := &handler.Gorm{DB: DB}
	auth := middleware.Auth(Gorm)
	prometheus := fiberprometheus.New("fiber_app")
	prometheus.RegisterAt(App, "/metrics")
	App.Use(prometheus.Middleware)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	App.Use(middleware.Logger(logger))
	r := App.Group("/", auth.WhoAreYou)
	Auth(Gorm, r)
	Holidays(Gorm, r)
	Sector(Gorm, r)
	Company(Gorm, r)
	User(Gorm, r)
	Branch(Gorm, r)
	Service(Gorm, r)
	Swagger(App)
	Employee(Gorm, r)
}
