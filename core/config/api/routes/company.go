package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/handlers"

	"github.com/gofiber/fiber/v3"
)

func Company(Gorm *handlers.Gorm, App *fiber.App) {
	cc := controllers.Company(Gorm)
	r := App.Group("/company")
	r.Get("/name/:name", cc.GetOneByName) // ok
	r.Get("/tax_id/:tax_id", cc.GetOneByTaxId) // ok

	controllers.CreateRoutes(r, cc)
}
