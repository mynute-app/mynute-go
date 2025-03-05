package routes

import (
	"agenda-kaki-go/core/controller"
	"agenda-kaki-go/core/handler"

	"github.com/gofiber/fiber/v2"
)

func Sector(Gorm *handler.Gorm, r fiber.Router) {
	cct := controller.Sector(Gorm)
	c := r.Group("/sector")
	c.Post("/", cct.CreateCompanyType)             // ok
	c.Get("/:id", cct.GetCompanyTypeById)          // ok
	c.Get("/name/:name", cct.GetCompanyTypeByName) // ok
	c.Patch("/:id", cct.UpdateCompanyTypeById)     // ok
	c.Delete("/:id", cct.DeleteCompanyTypeById)    // ok
}
