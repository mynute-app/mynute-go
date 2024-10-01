package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/services"

	"github.com/gofiber/fiber/v3"
)

func CompanyType(DB *services.Postgres, App *fiber.App) {
	cct := controllers.CompanyType{DB: DB, App: App}
	r := App.Group("/companyType")

	r.Post("/", cct.Create)
	r.Get("/", cct.GetAll)
	r.Get("/:id", cct.GetOneById)
	r.Get("/name/:name", cct.GetOneByName)
	
}
