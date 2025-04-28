package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/service"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type appointment_controller struct {
	service.Base[model.Appointment, DTO.Appointment]
}

// CreateAppointment creates an appointment
//
//	@Summary		Create appointment
//	@Description	Create an appointment
//	@Tags			Appointment
//	@Accept			json
//	@Produce		json
//	@Param			appointment	body		DTO.CreateAppointment	true	"Appointment"
//	@Success		200			{object}	DTO.Appointment
//	@Failure		400			{object}	DTO.ErrorResponse
//	@Router			/appointment [post]
func (ac *appointment_controller) CreateAppointment(c *fiber.Ctx) error {
	return ac.CreateOne(c)
}

// GetAppointmentByID gets an appointment by ID
//
//	@Summary		Get appointment
//	@Description	Get an appointment by ID
//	@Tags			Appointment
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	DTO.Appointment
//	@Failure		400	{object}	DTO.ErrorResponse
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
//	@Param			id			path		string					true	"ID"
//	@Param			appointment	body		DTO.CreateAppointment	true	"Appointment"
//	@Success		200			{object}	DTO.Appointment
//	@Failure		400			{object}	DTO.ErrorResponse
//	@Router			/appointment/{id} [patch]
func (ac *appointment_controller) UpdateAppointmentByID(c *fiber.Ctx) error {
	res := &lib.SendResponse{Ctx: c}

	tx := ac.Request.Gorm.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			if err, ok := r.(error); ok {
				res.Http500(err)
			} else {
				res.Http500(lib.Error.General.InternalError)
			}
		} else if tx.Error != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	appointment_id := c.Params("id")
	if appointment_id == "" {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("missing appointment's id in the url"))
	}

	var appointment model.Appointment

	tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", appointment_id).Find(&appointment)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return lib.Error.Appointment.NotFound
		}
		return lib.Error.General.UpdatedError.WithError(tx.Error)
	}

	if appointment.Cancelled {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("appointment is cancelled"))
	}

	var updated_appointment model.Appointment

	if err := c.BodyParser(&updated_appointment); err != nil {
		return lib.Error.General.UpdatedError.WithError(err)
	}

	if updated_appointment.ID != uuid.Nil {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("appointment update can not have pre defined ID"))
	}

	tx.Model(&model.Appointment{}).Where("id = ?", appointment_id).Updates(updated_appointment)
	if tx.Error != nil {
		return lib.Error.General.UpdatedError.WithError(tx.Error)
	}

	return nil
}

// CancelAppointmentByID deletes an appointment by ID
//
//	@Summary		Delete appointment
//	@Description	Delete an appointment by ID
//	@Tags			Appointment
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	DTO.Appointment
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/appointment/{id} [delete]
func (ac *appointment_controller) CancelAppointmentByID(c *fiber.Ctx) error {
	// Set the appointment to cancelled
	// Check if the appointment is already cancelled
	appointment_id := c.Params("id")
	if appointment_id == "" {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("missing appointment's id in the url"))
	}
	var appointment model.Appointment
	if err := ac.Request.Gorm.DB.Where("id = ?", appointment_id).First(&appointment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.Appointment.NotFound
		}
		return lib.Error.General.UpdatedError.WithError(err)
	}
	if appointment.Cancelled {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("appointment is already cancelled"))
	}
	if appointment.StartTime.Before(time.Now()) {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("cannot cancel appointment as it already happened"))
	}
	
	return nil
}

// Constructor for appointment_controller
func Appointment(Gorm *handler.Gorm) *appointment_controller {
	ac := &appointment_controller{
		Base: service.Base[model.Appointment, DTO.Appointment]{
			Name:    namespace.CompanyKey.Name,
			Request: handler.Request(Gorm),
		},
	}
	endpoint := &handler.Endpoint{DB: Gorm.DB}
	endpoint.BulkRegisterHandler([]fiber.Handler{
		ac.CreateAppointment,
		ac.GetAppointmentByID,
		ac.UpdateAppointmentByID,
		ac.CancelAppointmentByID,
	})
	return ac
}
