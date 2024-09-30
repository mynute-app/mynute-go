package routes

import (
	"agenda-kaki-company-go/api/controllers"
	"agenda-kaki-company-go/api/services"

	"github.com/gofiber/fiber/v3"
)

func CompanyType(DB *services.Postgres, App *fiber.App) {
	Controller := controllers.CompanyType{}

	App.Post("/companyType", Controller.Create)
	App.Get("/companyType/:id", Controller.GetOneById)
}