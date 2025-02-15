package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/handlers"

	"github.com/gofiber/fiber/v2"
)

func CompanyType(Gorm *handlers.Gorm, r fiber.Router) {
	cct := controllers.CompanyType(Gorm)
	c := r.Group("/companyType")
	c.Get("/name/:name", cct.GetOneByName) // ok

	controllers.CreateRoutes(c, cct)
}
