package routes

import (
	"agenda-kaki-go/core/controller"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/middleware"

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

	a := middleware.Auth(Gorm)

	r := App.Group("/")
	mdwPub := []fiber.Handler{a.WhoAreYou}
	mdwPrv := []fiber.Handler{a.WhoAreYou, a.DenyUnauthorized}

	endpoint := &handler.Endpoint{DB: DB}
	if err := endpoint.Build(r, r, mdwPub, mdwPrv); err != nil {
		panic(err)
	}
}
