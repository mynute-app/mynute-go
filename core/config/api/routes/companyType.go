package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v3"
)

func CompanyType(Gorm *handlers.Gorm, App *fiber.App) {
	Middleware := &middleware.CompanyType{Gorm: Gorm}
	HTTP := &handlers.HTTP{Gorm: Gorm}
	RequestHandler := &handlers.Request{HTTP: HTTP}
	cct := controllers.NewCompanyTypeController(RequestHandler, Middleware)
	r := App.Group("/companyType")
	r.Get("/name/:name", cct.GetOneByName) // ok

	controllers.CreateRoutes(r, cct)
	
}
