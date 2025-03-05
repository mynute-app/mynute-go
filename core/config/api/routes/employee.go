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
	e.Post("/", middleware.SaveReqInfo, middleware.FindUserByEmail(Gorm), middleware.SetEmployeeUserAccount(Gorm), ce.CreateEmployee)          // ok
	e.Get("/:id", ce.GetEmployeeById)       // ok
	e.Patch("/:id", ce.UpdateEmployeeById)  // ok
	e.Delete("/:id", ce.DeleteEmployeeById) // ok
}