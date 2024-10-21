package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v3"
)

func Company(Gorm *handlers.Gorm, App *fiber.App) {
	Middleware := &middleware.Company{Gorm: Gorm}
	HTTP := &handlers.HTTP{Gorm: Gorm}
	RequestHandler := &handlers.Request{HTTP: HTTP}
	Associations := []string{"CompanyTypes", "Branches"}
	cc := &controllers.Company{
		Request: RequestHandler, 
		Middleware: Middleware,
		Associations: Associations,
	}

	r := App.Group("/company")
	r.Get("/name/:name", cc.GetOneByName) // ok
	r.Get("/tax_id/:tax_id", cc.GetOneByTaxId) // ok

	controllers.CreateRoutes(r, cc)
}
