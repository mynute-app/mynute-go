package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/services"

	"github.com/gofiber/fiber/v3"
)

func Company(DB *services.Postgres, App *fiber.App) {
	Controller := controllers.Company{DB: DB, App: App}

	App.Post("/company", Controller.Create)
	App.Get("/company/:id", Controller.GetOneById)
	App.Put("/company/:id", Controller.UpdateById)
	App.Delete("/company/:id", Controller.DeleteById)
}
