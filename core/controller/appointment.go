package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/service"

	"github.com/gofiber/fiber/v2"
)

type appointment_controller struct {
	service.Base[model.Company, DTO.Company]
}

// CreateAppointment creates an appointment
//
//	@Summary		Create appointment
//	@Description	Create an appointment
//	@Tags			Appointment
//	@Accept			json
//	@Produce		json
//	@Param			appointment	body		DTO.CreateAppointment	true	"Appointment"
//	@Success		200		{object}	DTO.Appointment
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/appointment [post]
func (ac *appointment_controller) CreateAppointment(c *fiber.Ctx) error {
	return ac.CreateOne(c)
}

// GetAppointmentByID gets an appointment by ID
//
//	@Summary		Get appointment
//  @Description	Get an appointment by ID
//	@Tags			Appointment
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		200		{object}	DTO.Appointment
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/appointment/{id} [get]
func (ac *appointment_controller) GetAppointmentByID(c *fiber.Ctx) error {
	return ac.GetOneById(c)
}

// UpdateAppointmentByID updates an appointment by ID
//
//	@Summary		Update appointment
//	@Description	Update an appointment by ID
//	@Tags			Appointment
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Param			appointment	body		DTO.UpdateAppointment	true	"Appointment"
//	@Success		200		{object}	DTO.Appointment
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/appointment/{id} [patch]
func (ac *appointment_controller) UpdateAppointmentByID(c *fiber.Ctx) error {
	return ac.UpdateOneById(c)
}

// DeleteAppointmentByID deletes an appointment by ID
//
//	@Summary		Delete appointment
//	@Description	Delete an appointment by ID
//	@Tags			Appointment
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		200		{object}	DTO.Appointment
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/appointment/{id} [delete]
func (ac *appointment_controller) DeleteAppointmentByID(c *fiber.Ctx) error {
	return ac.DeleteOneById(c)
}

// Constructor for appointment_controller
func Appointment(Gorm *handler.Gorm) *appointment_controller {
	return &appointment_controller{
		Base: service.Base[model.Company, DTO.Company]{
			Name:         namespace.CompanyKey.Name,
			Request:      handler.Request(Gorm),
			Associations: []string{"Sector", "Branches", "Employees", "Services"},
		},
	}
}
