package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/handlers"

	"github.com/gofiber/fiber/v2"
)

func CompanyType(Gorm *handlers.Gorm, r fiber.Router) {
	cct := controllers.CompanyType(Gorm)
	c := r.Group("/company_type")
	c.Post("/", cct.CreateCompanyType)             // ok
	c.Get("/:id", cct.GetCompanyTypeById)          // ok
	c.Get("/name/:name", cct.GetCompanyTypeByName) // ok
	c.Patch("/:id", cct.UpdateCompanyTypeById)     // ok
	c.Delete("/:id", cct.DeleteCompanyTypeById)    // ok
}
