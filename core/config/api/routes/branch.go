package routes

import (
	"agenda-kaki-go/core/controllers"
	"agenda-kaki-go/core/handlers"

	"github.com/gofiber/fiber/v2"
)

func Branch(Gorm *handlers.Gorm, r fiber.Router) {
	cb := controllers.Branch(Gorm)
	b := r.Group("/branch")
	b.Post("/", cb.CreateBranch)          // ok
	b.Get("/:id", cb.GetBranchById)       // ok
	b.Patch("/:id", cb.UpdateBranchById)  // ok
	b.Delete("/:id", cb.DeleteBranchById) // ok
}
