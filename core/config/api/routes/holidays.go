package routes

import (
	"agenda-kaki-go/core/controller"
	"agenda-kaki-go/core/handler"

	"github.com/gofiber/fiber/v2"
)

func Holidays(Gorm *handler.Gorm, App fiber.Router) {
	ce := controller.Holidays(Gorm)
	r := App.Group("/holidays")
	r.Post("/", ce.CreateHoliday)
	r.Get("/:id", ce.GetHolidayById)
	r.Patch("/:id", ce.UpdateHolidayById)
	r.Delete("/:id", ce.DeleteHolidayById)
}
