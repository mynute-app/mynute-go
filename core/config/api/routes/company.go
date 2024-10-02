package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/services"

	"github.com/gofiber/fiber/v3"
)

func Company(DB *services.Postgres, App *fiber.App) {
	cc := controllers.Company{DB: DB, App: App}
	r := App.Group("/company")

	r.Post("/", cc.Create) // ok
	r.Get("/", cc.GetAll) // ok
	r.Get("/:id", cc.GetOneById) // ok
	r.Get("/name/:name", cc.GetOneByName) // ok
	r.Get("/tax_id/:tax_id", cc.GetOneByTaxId) // ok
	r.Patch("/:id", cc.UpdateById) 
	r.Delete("/:id", cc.DeleteById)
}
