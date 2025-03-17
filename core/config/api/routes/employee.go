package routes

import (
	"agenda-kaki-go/core/controller"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v2"
)

func Employee(Gorm *handler.Gorm, r fiber.Router) {
	ce := controller.Employee(Gorm)
	e := r.Group("/employee")
	AuthMdw := middleware.Auth(Gorm)
	e.Post("/login", ce.LoginEmployee)                                   // ok
	e.Post("/verify-email/:email/:code", ce.VerifyEmployeeEmail)         // ok
	auth := e.Group("/", AuthMdw.DenyUnauthorized)                       // ok
	auth.Post("/", ce.CreateEmployee)                                    // ok
	auth.Get("/:id", ce.GetEmployeeById)                                 // ok
	auth.Patch("/:id", ce.UpdateEmployeeById)                            // ok
	auth.Delete("/:id", ce.DeleteEmployeeById)                           // ok
	auth.Post("/:employee_id/branch/:branch_id", ce.AddEmployeeToBranch) // ok
}
