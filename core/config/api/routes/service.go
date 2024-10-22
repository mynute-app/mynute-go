package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/middleware"
	"agenda-kaki-go/core/handlers"

	"github.com/gofiber/fiber/v3"
)

func Service(Gorm *handlers.Gorm, App *fiber.App) {
	Middleware := &middleware.Service{Gorm: Gorm}
	HTTP := &handlers.HTTP{Gorm: Gorm}
	RequestHandler := &handlers.Request{HTTP: HTTP}
	cs := controllers.NewServiceController(RequestHandler, Middleware)
	r := App.Group("/service")

	controllers.CreateRoutes(r, cs)
}