package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/middleware"
	"agenda-kaki-go/core/services"

	"github.com/gofiber/fiber/v3"
)

func Company(DB *services.Postgres, App *fiber.App) {
	mdw := &middleware.Company{DB: DB}
	cc := &controllers.Company{DB: DB, Middleware: mdw}
	r := App.Group("/company")

	r.Post("/", cc.Create) // ok
	r.Get("/", cc.GetAll) // ok
	r.Get("/:id", cc.GetOneById) // ok
	r.Get("/name/:name", cc.GetOneByName) // ok
	r.Get("/tax_id/:tax_id", cc.GetOneByTaxId) // ok
	r.Patch("/:id", cc.UpdateById) // ok
	r.Delete("/:id", cc.DeleteById)
}
