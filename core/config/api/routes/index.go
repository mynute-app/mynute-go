package routes

import (
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/middleware"
	"log/slog"

	
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Builds all available routes for the application
func Build(DB *gorm.DB, App *fiber.App, logger *slog.Logger) {
	Gorm := &handler.Gorm{DB: DB}
	auth := middleware.Auth(Gorm)
	App.Use(middleware.Logger(logger))
	Prometheus(App)
	Swagger(App)
	r := App.Group("/", auth.WhoAreYou)
	Auth(Gorm, r)
	Holidays(Gorm, r)
	Sector(Gorm, r)
	Company(Gorm, r)
	Client(Gorm, r)
	Branch(Gorm, r)
	Service(Gorm, r)
	Employee(Gorm, r)
	Appointment(Gorm, r)
}
