package routes

import (
	"mynute-go/core/controller"
	"mynute-go/core/handler"
	"mynute-go/core/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Builds all available routes for the application
func Build(DB *gorm.DB, App *fiber.App) {
	Prometheus(App)
	Swagger(App)

	Gorm := &handler.Gorm{DB: DB}

	controller.Appointment(Gorm)
	controller.Auth(Gorm)
	controller.Branch(Gorm)
	controller.Client(Gorm)
	controller.Company(Gorm)
	controller.Employee(Gorm)
	controller.Holiday(Gorm)
	controller.Sector(Gorm)
	controller.Service(Gorm)

	r := App.Group("/")

	r.Get("/", controller.Home)

	endpoints := &middleware.Endpoint{DB: Gorm}
	if err := endpoints.Build(r); err != nil {
		panic(err)
	}
}
