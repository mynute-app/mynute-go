package routes

import (
	"agenda-kaki-go/core/controller"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v2"
)

func Company(Gorm *handler.Gorm, r fiber.Router) {
	cc := controller.Company(Gorm)
	c := r.Group("/company")
	c.Post("/", cc.CreateCompany) // ok
	c.Get("/:id", cc.GetCompanyById)                             // ok
	c.Get("/name/:name", cc.GetCompanyByName)                    // ok
	c.Get("/tax_id/:tax_id", cc.GetCompanyByTaxId)               // ok
	auth := c.Group("/", middleware.Auth(Gorm).DenyUnauthorized) // ok
	auth.Patch("/:id", cc.UpdateCompanyById)                        // ok
	auth.Delete("/:id", cc.DeleteCompanyById)                       // ok
}
