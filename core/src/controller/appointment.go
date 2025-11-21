package controller

import (
	"context"
	"fmt"
	"log"
	DTO "mynute-go/core/src/config/api/dto"
	database "mynute-go/core/src/config/db"
	"mynute-go/core/src/config/db/model"
	"mynute-go/core/src/handler"
	"mynute-go/core/src/lib"
	"mynute-go/core/src/lib/email"
	"mynute-go/core/src/middleware"
	"mynute-go/debug"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateAppointment creates an appointment
//
//	@Summary		Create appointment
//	@Description	Create an appointment
//	@Tags			Appointment
//	@Accept			json
//	@Produce		json
//	@Param			X-Company-ID	header		string					true	"X-Company-ID"
//	@Param			appointment		body		DTO.CreateAppointment	true	"Appointment"
//	@Success		200				{object}	DTO.Appointment
//	@Failure		400				{object}	DTO.ErrorResponse
//	@Router			/appointment [post]
func CreateAppointment(c *fiber.Ctx) error {
	var appointment model.Appointment
	if err := Create(c, &appointment); err != nil {
		return err
	}
	if err := debug.Output("controller_CreateAppointment", appointment); err != nil {
		return err
	}

	// Send appointment created emails
	go func() {
		ctx := context.Background()
		tx, err := lib.Session(c)
		if err != nil {
			log.Printf("Failed to get database session for email: %v", err)
			return
		}

		// Initialize email service
		sender, err := email.NewProvider(nil)
		if err != nil {
			log.Printf("Failed to create email provider: %v", err)
			return
		}

		templateDir := filepath.Join("static", "email")
		translationDir := filepath.Join("translation", "email")
		emailService := email.NewAppointmentEmailService(sender, templateDir, translationDir)

		if err := emailService.SendAppointmentCreatedEmails(ctx, tx, &appointment); err != nil {
			log.Printf("Failed to send appointment created emails: %v", err)
		}
	}()

	if err := lib.ResponseFactory(c).SendDTO(200, &appointment, &DTO.Appointment{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// GetAppointmentByID gets an appointment by ID
//
//	@Summary		Get appointment
//	@Description	Get an appointment by ID
//	@Tags			Appointment
//	@Accept			json
//	@Produce		json
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			id				path		string	true	"ID"
//	@Success		200				{object}	DTO.Appointment
//	@Failure		400				{object}	DTO.ErrorResponse
//	@Router			/appointment/{id} [get]
func GetAppointmentByID(c *fiber.Ctx) error {
	var appointment model.Appointment
	if err := GetOneBy("id", c, &appointment, nil, nil); err != nil {
		return err
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &appointment, &DTO.Appointment{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// UpdateAppointmentByID updates an appointment by ID
//
//	@Summary		Update appointment
//	@Description	Update an appointment by ID
//	@Tags			Appointment
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string					true	"X-Company-ID"
//	@Param			id				path		string					true	"ID"
//	@Param			appointment		body		DTO.CreateAppointment	true	"Appointment"
//	@Success		200				{object}	DTO.Appointment
//	@Failure		400				{object}	DTO.ErrorResponse
//	@Router			/appointment/{id} [patch]
func UpdateAppointmentByID(c *fiber.Ctx) error {
	var err error

	appointment_id := c.Params("id")
	if appointment_id == "" {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("missing appointment's id in the url"))
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var appointment model.Appointment
	if err = database.LockForUpdate(tx, &appointment, "id", appointment_id); err != nil {
		return err
	}

	if appointment.IsCancelled {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("appointment is cancelled"))
	}

	var updated_appointment model.Appointment

	if err = c.BodyParser(&updated_appointment); err != nil {
		return lib.Error.General.UpdatedError.WithError(err)
	}

	if updated_appointment.ID != uuid.Nil {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("appointment update can not have pre defined ID"))
	}

	tx.Model(appointment).Where("id = ?", appointment_id).Updates(updated_appointment)
	if tx.Error != nil {
		return lib.Error.General.UpdatedError.WithError(tx.Error)
	}

	tx.Model(appointment).Where("id = ?", appointment_id).First(&appointment)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return lib.Error.Appointment.NotFound
		}
		return lib.Error.General.UpdatedError.WithError(tx.Error)
	}

	// Send appointment updated emails
	go func() {
		ctx := context.Background()
		// Initialize email service
		sender, err := email.NewProvider(nil)
		if err != nil {
			log.Printf("Failed to create email provider: %v", err)
			return
		}

		templateDir := filepath.Join("static", "email")
		translationDir := filepath.Join("translation", "email")
		emailService := email.NewAppointmentEmailService(sender, templateDir, translationDir)

		if err := emailService.SendAppointmentUpdatedEmails(ctx, tx, &appointment); err != nil {
			log.Printf("Failed to send appointment updated emails: %v", err)
		}
	}()

	if err = lib.ResponseFactory(c).SendDTO(200, &appointment, &DTO.Appointment{}); err != nil {
		return lib.Error.General.UpdatedError.WithError(err)
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
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			id				path		string	true	"ID"
//	@Success		200				{object}	DTO.Appointment
//	@Failure		400				{object}	DTO.ErrorResponse
//	@Router			/appointment/{id} [delete]
func CancelAppointmentByID(c *fiber.Ctx) error {
	var err error
	appointment_id := c.Params("id")
	if appointment_id == "" {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("missing appointment's id in the url"))
	}
	uuid, err := uuid.Parse(appointment_id)
	if err != nil {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("invalid appointment's id in the url"))
	}
	var appointment model.Appointment
	appointment.ID = uuid
	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	// Load appointment data before cancelling for email
	if err := tx.Where("id = ?", uuid).First(&appointment).Error; err != nil {
		return lib.Error.Appointment.NotFound.WithError(err)
	}

	if err := appointment.Cancel(tx); err != nil {
		return err
	}

	// Send appointment cancelled emails
	go func() {
		ctx := context.Background()
		// Initialize email service
		sender, err := email.NewProvider(nil)
		if err != nil {
			log.Printf("Failed to create email provider: %v", err)
			return
		}

		templateDir := filepath.Join("static", "email")
		translationDir := filepath.Join("translation", "email")
		emailService := email.NewAppointmentEmailService(sender, templateDir, translationDir)

		if err := emailService.SendAppointmentCancelledEmails(ctx, tx, &appointment); err != nil {
			log.Printf("Failed to send appointment cancelled emails: %v", err)
		}
	}()

	return nil
}

// Constructor for appointment_controller
func Appointment(Gorm *handler.Gorm) {
	endpoint := &middleware.Endpoint{DB: Gorm}
	endpoint.BulkRegisterHandler([]fiber.Handler{
		CreateAppointment,
		GetAppointmentByID,
		UpdateAppointmentByID,
		CancelAppointmentByID,
	})
}
