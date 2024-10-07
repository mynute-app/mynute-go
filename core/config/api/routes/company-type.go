package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v3"
)

func CompanyType(Gorm *handlers.Gorm, App *fiber.App) {
	mdw := middleware.CompanyType{Gorm: Gorm}
	handler := handlers.HTTP{Gorm: Gorm}
	cct := controllers.CompanyType{Gorm: Gorm, Middleware: &mdw, HttpHandler: &handler}
	r := App.Group("/companyType")

	r.Post("/", cct.Create) // ok
	r.Get("/", cct.GetAll) // ok
	r.Get("/:id", cct.GetOneById) // ok
	r.Get("/name/:name", cct.GetOneByName) // ok
	r.Delete("/:id", cct.DeleteOneById)
	r.Patch("/:id", cct.UpdateOneById)
	
}
