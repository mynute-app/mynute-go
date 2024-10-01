package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/services"

	"github.com/gofiber/fiber/v3"
)

func CompanyType(DB *services.Postgres, App *fiber.App) {
	Controller := controllers.CompanyType{DB: DB, App: App}

	App.Post("/companyType", Controller.Create)
	App.Get("/companyType/:id", Controller.GetOneById)
	App.Get("/companyType/name/:name", Controller.GetOneByName)
	App.Get("/companyType", Controller.GetAll)
}
