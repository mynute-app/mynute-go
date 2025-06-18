package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	database "agenda-kaki-go/core/config/db"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/lib"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CreateWorkSchedule creates a work schedule for an employee
//
//	@Summary		Create work schedule
//	@Description	Create a work schedule for an employee
//	@Tags			WorkSchedule
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil	"Unauthorized"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Accept			json
//	@Produce		json
//	@Param			work_schedule	body		DTO.CreateWorkSchedule	true	"Work Schedule"
//	@Success		201		{object}	model.WorkSchedule
//	@Failure		400		{object}	lib.ErrorResponse
//	@Router			/work_schedule [post]
func CreateWorkSchedule(c *fiber.Ctx) error {
	var workSchedule DTO.CreateWorkSchedule

	if err := c.BodyParser(&workSchedule); err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("failed to parse request body: %w", err))
	}

	tx, end, err := database.ContextTransaction(c)
	if err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	defer end()

	var employee model.Employee
	if err := tx.First(&employee, "id = ?", workSchedule.EmployeeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee with ID %s not found", workSchedule.EmployeeID))
		}
		return lib.Error.General.DatabaseError.WithError(err)
	}

	return nil
}
