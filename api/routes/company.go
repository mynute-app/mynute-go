package routes

import (
	"agenda-kaki-company-go/api/controllers"
	"agenda-kaki-company-go/api/services"

	"github.com/gofiber/fiber/v3"
)

func Company(DB *services.Postgres, App *fiber.App) {
	Controller := controllers.Company{ DB: DB, App: App }

	App.Post("/companies", Controller.Create)
	App.Get("/companies/:id", Controller.GetOneById)
	App.Put("/companies/:id", Controller.UpdateById)
	App.Delete("/companies/:id", Controller.DeleteById)
}