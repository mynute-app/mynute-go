package routes

import (
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Builds all available routes for the application
func Build(DB *gorm.DB, App *fiber.App) {
	Gorm := &handler.Gorm{DB: DB}
	Router := App.Group("/", lib.ApiLog)
	Auth(Gorm, Router)
	Holidays(Gorm, Router)
	Sector(Gorm, Router)
	Company(Gorm, Router)
	User(Gorm, Router)
	Branch(Gorm, Router)
	Service(Gorm, Router)
	Swagger(App)
	Employee(Gorm, Router)
}
