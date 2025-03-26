package routes

import (
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/controller"
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
	controller.Company(Gorm)
	router_pub := App.Group("/", auth.WhoAreYou)
	router_auth := router_pub.Group("/", auth.DenyUnauthorized)
	var dbRoutes []model.Route
	route := handler.Route{}
	// Get all routes from database and assign them to the router
	DB.Find(&dbRoutes)
	for _, dbRoute := range dbRoutes {
		dbRouteHandler := route.GetHandler(dbRoute.Path, dbRoute.Method)
		if dbRoute.IsPublic {
			router_pub.Add(dbRoute.Method, dbRoute.Path, dbRouteHandler)
		} else {
			router_auth.Add(dbRoute.Method, dbRoute.Path, dbRouteHandler)
		}
	}
}
