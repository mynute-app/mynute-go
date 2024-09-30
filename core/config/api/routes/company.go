package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/services"

	"github.com/gofiber/fiber/v3"
)

func Company(DB *services.Postgres, App *fiber.App) {
	cc := controllers.Company{DB: DB, App: App}
	r := App.Group("/company")

	r.Post("/", cc.Create)
	r.Get("/:id", cc.GetOneById)
	r.Put("/:id", cc.UpdateById)
	r.Delete("/:id", cc.DeleteById)
	r.Get("/name/:name", cc.GetOneByName)
}
