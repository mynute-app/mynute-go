package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/service"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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

	var new_appointment model.Appointment

	if err := c.BodyParser(&new_appointment); err != nil {
		return lib.Error.General.UpdatedError.WithError(err)
	}

	var old_appointment model.Appointment

	if err := ac.Request.Gorm.DB.Find(&old_appointment, "id = ?", appointment_id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.Appointment.NotFound
		}
		return err
	}

	if old_appointment.CompanyID != new_appointment.CompanyID {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("company_id cannot be changed"))
	} else if old_appointment.ClientID != new_appointment.ClientID {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("client_id cannot be changed"))
	} else if old_appointment.Cancelled {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("appointment is already cancelled"))
	} else if old_appointment.MovedToID != nil {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("appointment has already been moved"))
	}

	hasChanges := false

	if old_appointment.ServiceID != new_appointment.ServiceID {
		old_appointment.ChangedService = true
		hasChanges = true
	}
	if old_appointment.StartTime != new_appointment.StartTime {
		old_appointment.ChangedTime = true
		hasChanges = true
	}
	if old_appointment.EmployeeID != new_appointment.EmployeeID {
		old_appointment.ChangedEmployee = true
		hasChanges = true
	}
	if old_appointment.BranchID != new_appointment.BranchID {
		old_appointment.ChangedBranch = true
		hasChanges = true
	}

	if !hasChanges {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("no changes detected"))
	}

	

}

// DeleteAppointmentByID deletes an appointment by ID
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
func (ac *appointment_controller) DeleteAppointmentByID(c *fiber.Ctx) error {
	return ac.DeleteOneById(c)
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
		ac.DeleteAppointmentByID,
	})
	return ac
}
