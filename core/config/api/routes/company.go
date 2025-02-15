package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/handlers"

	"github.com/gofiber/fiber/v2"
)

func Company(Gorm *handlers.Gorm, r fiber.Router) {
	cc := controllers.Company(Gorm)
	c := r.Group("/company")
	c.Get("/name/:name", cc.GetOneByName)      // ok
	c.Get("/tax_id/:tax_id", cc.GetOneByTaxId) // ok

	controllers.CreateRoutes(c, cc)
}
