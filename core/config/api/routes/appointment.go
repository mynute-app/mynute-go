package routes

import (
	"agenda-kaki-go/core/controller"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v2"
)

func Appointment(Gorm *handler.Gorm, r fiber.Router) {
	cb := controller.Appointment(Gorm)
	b := r.Group("/appointment")
	auth := b.Group("/", middleware.Auth(Gorm).DenyUnauthorized)
	auth.Post("/", cb.CreateAppointment)
	auth.Get("/:id", cb.GetAppointmentByID)
	auth.Patch("/:id", cb.UpdateAppointmentByID)
}