package routes

import (
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Builds all available routes for the application
func Build(DB *gorm.DB, App *fiber.App) {
	Gorm := &handler.Gorm{DB: DB}
	auth := middleware.Auth(Gorm)
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
