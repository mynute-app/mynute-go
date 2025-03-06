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
	mdw := middleware.Employee(Gorm)
	CreateEmployeeRoutine := append(mdw.CreateEmployee(), ce.CreateEmployee)
	e.Post("/", CreateEmployeeRoutine...) // ok
	e.Get("/:id", ce.GetEmployeeById)                                                    // ok
	e.Patch("/:id", ce.UpdateEmployeeById)                                               // ok
	e.Delete("/:id", ce.DeleteEmployeeById)                                              // ok
}
