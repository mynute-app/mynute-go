package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/handlers"

	"github.com/gofiber/fiber/v2"
)

func Company(Gorm *handlers.Gorm, r fiber.Router) {
	cc := controllers.Company(Gorm)
	c := r.Group("/company")
	c.Post("/", cc.CreateCompany) 	               // ok
	c.Get("/:id", cc.GetCompanyById)               // ok
	c.Get("/name/:name", cc.GetCompanyByName)      // ok
	c.Get("/tax_id/:tax_id", cc.GetCompanyByTaxId) // ok
	c.Patch("/:id", cc.UpdateCompanyById)          // ok
	c.Delete("/:id", cc.DeleteCompanyById)         // ok
}
