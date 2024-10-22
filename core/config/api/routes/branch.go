package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v3"
)

func Branch(Gorm *handlers.Gorm, App *fiber.App) {
	Middleware := &middleware.Branch{Gorm: Gorm}
	HTTP := &handlers.HTTP{Gorm: Gorm}
	RequestHandler := &handlers.Request{HTTP: HTTP}
	cb := controllers.NewBranchController(RequestHandler, Middleware)
	r := App.Group("/company/:companyId/branch")
	controllers.CreateRoutes(r, cb)
}