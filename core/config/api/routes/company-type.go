package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/middleware"
	"agenda-kaki-go/core/services"

	"github.com/gofiber/fiber/v3"
)

func CompanyType(Gorm *services.Gorm, App *fiber.App) {
	mdw := middleware.CompanyType{Gorm: Gorm}
	cct := controllers.CompanyType{Gorm: Gorm, Middleware: &mdw}
	r := App.Group("/companyType")

	r.Post("/", cct.Create)
	r.Get("/", cct.GetAll)
	r.Get("/:id", cct.GetOneById)
	r.Get("/name/:name", cct.GetOneByName)
	r.Delete("/:id", cct.DeleteById)
	
}
