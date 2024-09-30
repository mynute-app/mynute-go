package routes

import (
	"agenda-kaki-go/api/controllers"
	"agenda-kaki-go/api/services"

	"github.com/gofiber/fiber/v3"
)

func CompanyType(DB *services.Postgres, App *fiber.App) {
	Controller := controllers.CompanyType{DB: DB, App: App}

	App.Post("/companyType", Controller.Create)
	App.Get("/companyType/:id", Controller.GetOneById)
}
