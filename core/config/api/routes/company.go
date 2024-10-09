package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/middleware"
	"agenda-kaki-go/core/handlers"

	"github.com/gofiber/fiber/v3"
)

func Company(Gorm *handlers.Gorm, App *fiber.App) {
	Middleware := &middleware.Company{Gorm: Gorm}
	HTTP := &handlers.HTTP{Gorm: Gorm}
	RequestHandler := &handlers.Request{HTTP: HTTP}
	cc := &controllers.Company{Request: RequestHandler, Middleware: Middleware}
	r := App.Group("/company")

	r.Post("/", cc.CreateOne) // ok
	r.Get("/", cc.GetAll) // ok
	r.Get("/:id", cc.GetOneById) // ok
	r.Get("/name/:name", cc.GetOneByName) // ok
	r.Get("/tax_id/:tax_id", cc.GetOneByTaxId) // ok
	r.Patch("/:id", cc.UpdateOneById) // ok
	r.Delete("/:id", cc.DeleteOneById) // ok
}
