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
	cct := controllers.CompanyType{Request: RequestHandler, Middleware: Middleware}
	r := App.Group("/companyType")

	r.Post("/", cct.CreateOne) // ok
	r.Get("/", cct.GetAll) // ok
	r.Get("/:id", cct.GetOneById) // ok
	r.Get("/name/:name", cct.GetOneByName) // ok
	r.Delete("/:id", cct.DeleteOneById) // ok
	r.Patch("/:id", cct.UpdateOneById) // ok
	
}
