package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/middleware"
	"agenda-kaki-go/core/handlers"

	"github.com/gofiber/fiber/v3"
)

func Company(Gorm *handlers.Gorm, App *fiber.App) {
	mdw := &middleware.Company{Gorm: Gorm}
	cc := &controllers.Company{Gorm: Gorm, Middleware: mdw}
	r := App.Group("/company")

	r.Post("/", cc.Create) // ok
	r.Get("/", cc.GetAll) // ok
	r.Get("/:id", cc.GetOneById) // ok
	r.Get("/name/:name", cc.GetOneByName) // ok
	r.Get("/tax_id/:tax_id", cc.GetOneByTaxId) // ok
	r.Patch("/:id", cc.UpdateById) // ok
	r.Delete("/:id", cc.DeleteById)
}
