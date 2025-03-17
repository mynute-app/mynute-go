package routes

import (
	"agenda-kaki-go/core/controller"
	"agenda-kaki-go/core/handler"

	"github.com/gofiber/fiber/v2"
)

func Appointment(Gorm *handler.Gorm, r fiber.Router) {
	cb := controller.Appointment(Gorm)
	b := r.Group("/appointment")
	b.Post("/", cb.CreateAppointment)
	b.Get("/:id", cb.GetAppointmentByID)
	b.Patch("/:id", cb.UpdateAppointmentByID)
}