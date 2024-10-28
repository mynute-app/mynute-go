package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/handlers"

	"github.com/gofiber/fiber/v3"
)

func CompanyType(Gorm *handlers.Gorm, App *fiber.App) {
	cct := controllers.CompanyType(Gorm)
	r := App.Group("/companyType")
	r.Get("/name/:name", cct.GetOneByName) // ok

	controllers.CreateRoutes(r, cct)
	
}
