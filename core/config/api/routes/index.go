package routes

import (
	"agenda-kaki-go/core/config/db/model"
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

	auth := middleware.Auth(Gorm)

	router_pub := App.Group("/", auth.WhoAreYou)
	router_auth := router_pub.Group("/", auth.DenyUnauthorized)

	route := &handler.Route{DB: DB}
	route.Build(router_pub, router_auth)
}
