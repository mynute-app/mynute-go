package routes

import (
	"agenda-kaki-go/core/controller"
	"agenda-kaki-go/core/handler"

	"github.com/gofiber/fiber/v2"
)

func Branch(Gorm *handler.Gorm, r fiber.Router) {
	cb := controller.Branch(Gorm)
	b := r.Group("/branch")
	b.Post("/", cb.CreateBranch)          // ok
	b.Get("/:id", cb.GetBranchById)       // ok
	b.Get("/name/:name", cb.GetBranchByName) // ok
	b.Patch("/:id", cb.UpdateBranchById)  // ok
	b.Delete("/:id", cb.DeleteBranchById) // ok
}
