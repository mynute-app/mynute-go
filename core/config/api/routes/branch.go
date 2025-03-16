package routes

import (
	"agenda-kaki-go/core/controller"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v2"
)

func Branch(Gorm *handler.Gorm, r fiber.Router) {
	cb := controller.Branch(Gorm)
	b := r.Group("/branch")
	b.Post("/", cb.CreateBranch)                                                   // ok
	b.Get("/:id", cb.GetBranchById)                                                // ok
	b.Get("/name/:name", cb.GetBranchByName)                                       // ok
	b.Get("/:branch_id/employee/:employee_id/services", cb.GetEmployeeServicesByBranchId) // ok
	auth := b.Group("/", middleware.Auth(Gorm).DenyUnauthorized)                   // ok
	auth.Patch("/:id", cb.UpdateBranchById)                                        // ok
	auth.Delete("/:id", cb.DeleteBranchById)                                       // ok
	auth.Post("/:branch_id/service/:service_id", cb.AddServiceToBranch)                   // ok
	auth.Delete("/:branch_id/service/:service_id", cb.RemoveServiceFromBranch)            // ok
	auth.Post("/:branch_id/employee/:employee_id", cb.AddEmployeeToBranch)                // ok
	auth.Delete("/:branch_id/employee/:employee_id", cb.RemoveEmployeeFromBranch)         // ok
}
