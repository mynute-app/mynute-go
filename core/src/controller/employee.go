package controller

import (
	"encoding/json"
	"fmt"
	DTO "mynute-go/core/src/config/api/dto"
	dJSON "mynute-go/core/src/config/api/dto/json"
	database "mynute-go/core/src/config/db"
	"mynute-go/core/src/config/db/model"
	"mynute-go/core/src/config/namespace"
	"mynute-go/core/src/handler"
	"mynute-go/core/src/lib"
	"mynute-go/core/src/middleware"
	"mynute-go/debug"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreateEmployee creates an employee
//
//	@Summary		Create employee
//	@Description	Create an employee
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			employee	body		DTO.CreateEmployee	true	"Employee"
//	@Success		200			{object}	DTO.EmployeeFull
//	@Failure		400			{object}	DTO.ErrorResponse
//	@Router			/employee [post]
func CreateEmployee(c *fiber.Ctx) error {
	var employee model.Employee
	if err := Create(c, &employee); err != nil {
		return err
	}

	if err := debug.Output("controller_CreateEmployee", employee); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &employee, &DTO.EmployeeFull{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// GetEmployeeById retrieves an employee by ID
//
//	@Summary		Get employee by ID
//	@Description	Retrieve an employee by its ID
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Employee ID"
//	@Produce		json
//	@Success		200	{object}	DTO.EmployeeFull
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/employee/{id} [get]
func GetEmployeeById(c *fiber.Ctx) error {
	var employee model.Employee
	if err := GetOneBy("id", c, &employee, &[]string{"WorkSchedule.Services"}, &[]string{"Appointments"}); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &employee, &DTO.EmployeeFull{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// GetEmployeeByEmail retrieves an employee by email
//
//	@Summary		Get employee by email
//	@Description	Retrieve an employee by its email
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Failure		401				{object}	nil
//	@Param			email			path		string	true	"Employee Email"
//	@Produce		json
//	@Success		200	{object}	DTO.EmployeeFull
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/employee/email/{email} [get]
func GetEmployeeByEmail(c *fiber.Ctx) error {
	var employee model.Employee
	if err := GetOneBy("email", c, &employee, &[]string{"WorkSchedule.Services"}, &[]string{"Appointments"}); err != nil {
		return err
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &employee, &DTO.EmployeeFull{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// UpdateEmployeeById updates an employee by ID
//
//	@Summary		Update employee
//	@Description	Update an employee
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string						true	"Employee ID"
//	@Param			employee	body		DTO.UpdateEmployeeSwagger	true	"Employee"
//	@Success		200			{object}	DTO.EmployeeFull
//	@Failure		400			{object}	DTO.ErrorResponse
//	@Router			/employee/{id} [patch]
func UpdateEmployeeById(c *fiber.Ctx) error {
	var employee model.Employee
	if err := UpdateOneById(c, &employee, nil); err != nil {
		return err
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &employee, &DTO.EmployeeFull{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// DeleteEmployeeById deletes an employee by ID
//
//	@Summary		Delete employee by ID
//	@Description	Delete an employee by its ID
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Employee ID"
//	@Produce		json
//	@Success		200	{object}	DTO.EmployeeFull
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/employee/{id} [delete]
func DeleteEmployeeById(c *fiber.Ctx) error {
	return DeleteOneById(c, &model.Employee{})
}

// LoginEmployeeByPassword logs an employee in
//
//	@Summary		Login
//	@Description	Log in an employee using password
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			X-Company-ID	header	string	true	"X-Company-ID"
//	@Accept			json
//	@Produce		json
//	@Param			client	body	DTO.LoginEmployee	true	"Employee"
//	@Success		200
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Failure		401	{object}	nil
//	@Router			/employee/login [post]
func LoginEmployeeByPassword(c *fiber.Ctx) error {
	token, err := LoginByPassword(namespace.EmployeeKey.Name, &model.Employee{}, c)
	if err != nil {
		return err
	}
	c.Response().Header.Set(namespace.HeadersKey.Auth, token)
	return nil
}

// LoginEmployeeByEmailCode logs in an employee using email and validation code
//
//	@Summary		Login employee by email code
//	@Description	Login employee using email and validation code
//	@Tags			Employee
//	@Param			X-Company-ID	header	string	true	"X-Company-ID"
//	@Accept			json
//	@Produce		json
//	@Param			body	body	DTO.LoginByEmailCode	true	"Login credentials"
//	@Success		200
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/employee/login-with-code [post]
func LoginEmployeeByEmailCode(c *fiber.Ctx) error {
	token, err := LoginByEmailCode(namespace.EmployeeKey.Name, &model.Employee{}, c)
	if err != nil {
		return err
	}
	c.Response().Header.Set(namespace.HeadersKey.Auth, token)
	return nil
}

// SendEmployeeLoginValidationCodeByEmail sends a login validation code to an employee's email
//
//	@Summary		Send login validation code to employee email
//	@Description	Sends a 6-digit login validation code to the employee's email
//	@Tags			Employee
//	@Accept			json
//	@Produce		json
//	@Param			X-Company-ID	header	string	true	"X-Company-ID"
//	@Param			email			path	string	true	"Employee Email"
//	@Query			language															query	string	false	"Language code (default: en)"
//	@Success		200
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/employee/send-login-code/email/{email} [post]
func SendEmployeeLoginValidationCodeByEmail(c *fiber.Ctx) error {
	if err := SendLoginValidationCodeByEmail(c, &model.Employee{}); err != nil {
		return err
	}
	return nil
}

// ResetEmployeePasswordByEmail sets a random password of an employee using its email
//
//	@Summary		Reset employee password to a random value
//	@Description	Sets a random password of an employee using its email
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			X-Company-ID	header	string	true	"X-Company-ID"
//	@Accept			json
//	@Produce		json
//	@Param			email		path		string	true	"Employee Email"
//	@Query			language	query										string	false	"Language code (default: en)"
//	@Success		200			{object}	DTO.PasswordReseted
//	@Failure		400			{object}	DTO.ErrorResponse
//	@Router			/employee/reset-password/{email} [post]
func ResetEmployeePasswordByEmail(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("missing 'email' at params route"))
	}
	if err := SendNewPasswordByEmail(c, email, &model.Employee{}); err != nil {
		return err
	}
	return lib.ResponseFactory(c).Http200(nil)
}

// UpdateEmployeeImages updates the images of an employee
//
//	@Summary		Update employee images
//	@Description	Update the images of an employee
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Employee ID"
//	@Accept			json
//	@Produce		json
//	@Param			profile	formData	file	false	"Profile image"
//	@Success		200		{object}	dJSON.Images
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/employee/{id}/design/images [patch]
func UpdateEmployeeImages(c *fiber.Ctx) error {
	img_types_allowed := map[string]bool{"profile": true, "logo": true, "banner": true, "background": true}

	var employee model.Employee
	Design, err := UpdateImagesById(c, employee.TableName(), &employee, img_types_allowed)
	if err != nil {
		return err
	}

	return lib.ResponseFactory(c).SendDTO(200, &Design.Images, &dJSON.Images{})
}

// DeleteEmployeeImage deletes an image of an employee
//
//	@Summary		Delete employee image
//	@Description	Delete an image of an employee
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Employee ID"
//	@Param			image_type		path		string	true	"Image Type (logo, banner, favicon, background)"
//	@Produce		json
//	@Success		200	{object}	dJSON.Images
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/employee/{id}/design/images/{image_type} [delete]
func DeleteEmployeeImage(c *fiber.Ctx) error {
	img_types_allowed := map[string]bool{"profile": true, "logo": true, "banner": true, "background": true}
	var employee model.Employee
	Design, err := DeleteImageById(c, employee.TableName(), &employee, img_types_allowed)
	if err != nil {
		return err
	}

	return lib.ResponseFactory(c).SendDTO(200, &Design.Images, &dJSON.Images{})
}

// CreateEmployeeWorkSchedule creates a work schedule for an employee
//
//	@Summary		Create work schedule
//	@Description	Create a work schedule for an employee
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil		"Unauthorized"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Accept			json
//	@Param			work_schedule	body		DTO.CreateEmployeeWorkSchedule	true	"Work Schedule"
//	@Param			employee_id		path		string							true	"Employee ID"
//	@Success		200				{object}	DTO.EmployeeWorkSchedule
//	@Failure		400				{object}	DTO.ErrorResponse
//	@Router			/employee/{employee_id}/work_schedule [post]
func CreateEmployeeWorkSchedule(c *fiber.Ctx) error {
	var input DTO.CreateEmployeeWorkSchedule
	if err := c.BodyParser(&input); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	var EmployeeWorkSchedule model.EmployeeWorkSchedule

	employee_id := c.Params("employee_id")

	for i, ewr := range input.WorkRanges {
		if ewr.EmployeeID.String() != employee_id {
			return lib.Error.General.CreatedError.WithError(fmt.Errorf("work range [%d] employee ID (%s) does not match employee ID (%s) from path", i+1, ewr.EmployeeID.String(), employee_id))
		}

		start, err := lib.Parse_HHMM_To_Time(ewr.StartTime, ewr.TimeZone)
		if err != nil {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid start_time: %w", err))
		}
		end, err := lib.Parse_HHMM_To_Time(ewr.EndTime, ewr.TimeZone)
		if err != nil {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid end_time: %w", err))
		}

		services := make([]*model.Service, 0, len(ewr.Services))
		for _, serviceID := range ewr.Services {
			services = append(services, &model.Service{BaseModel: model.BaseModel{ID: serviceID.ID}})
		}

		EmployeeWorkSchedule.WorkRanges = append(EmployeeWorkSchedule.WorkRanges, model.EmployeeWorkRange{
			WorkRangeBase: model.WorkRangeBase{
				Weekday:   time.Weekday(ewr.Weekday),
				StartTime: start,
				EndTime:   end,
				TimeZone:  ewr.TimeZone,
				BranchID:  ewr.BranchID,
			},
			EmployeeID: ewr.EmployeeID,
			Services:   services,
		})
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	for _, ewr := range EmployeeWorkSchedule.WorkRanges {
		if err := tx.Create(&ewr).Error; err != nil {
			return lib.Error.General.CreatedError.WithError(err)
		}
		debug.Output("controller_CreateEmployeeWorkSchedule", ewr)
	}

	var ewr []model.EmployeeWorkRange
	if err := tx.
		Preload(clause.Associations).
		Find(&ewr, "employee_id = ?", employee_id).Error; err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}

	ews := model.EmployeeWorkSchedule{
		WorkRanges: ewr,
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &ews, &DTO.EmployeeWorkSchedule{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// GetEmployeeWorkSchedule
//
//	@Summary		Get all employee's work ranges
//	@Description	Retrieve all work ranges for an employee
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Failure		401				{object}	nil		"Unauthorized"
//	@Param			employee_id		path		string	true	"Employee ID"
//	@Produce		json
//	@Success		200	{object}	DTO.EmployeeWorkSchedule
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/employee/{employee_id}/work_schedule [get]
func GetEmployeeWorkSchedule(c *fiber.Ctx) error {
	employeeID := c.Params("employee_id")

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var ewr []model.EmployeeWorkRange
	if err := tx.
		Preload(clause.Associations).
		Find(&ewr, "employee_id = ?", employeeID).Error; err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}

	ews := model.EmployeeWorkSchedule{
		WorkRanges: ewr,
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &ews, &DTO.EmployeeWorkSchedule{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// GetEmployeeWorkRangeById retrieves a work range for an employee
//
//	@Summary		Get work range by ID
//	@Description	Retrieve a work range for an employee
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Failure		401				{object}	nil		"Unauthorized"
//	@Param			employee_id		path		string	true	"Employee ID"
//	@Param			work_range_id	path		string	true	"Work Range ID"
//	@Produce		json
//	@Success		200	{object}	DTO.EmployeeWorkRange
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/employee/{employee_id}/work_range/{work_range_id} [get]
func GetEmployeeWorkRangeById(c *fiber.Ctx) error {
	employeeID := c.Params("employee_id")
	workRangeID := c.Params("work_range_id")

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var workRange model.EmployeeWorkRange
	if err := tx.Preload(clause.Associations).First(&workRange, "id = ? AND employee_id = ?", workRangeID, employeeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	if workRange.EmployeeID.String() != employeeID {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range [%s] does not belong to employee (%s)", workRangeID, employeeID))
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &workRange, &DTO.EmployeeWorkRange{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// DeleteEmployeeWorkRange deletes a work schedule for an employee
//
//	@Summary		Delete work schedule
//	@Description	Delete a work schedule for an employee
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil		"Unauthorized"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			employee_id		path		string	true	"Employee ID"
//	@Param			work_range_id	path		string	true	"Work Range ID"
//	@Produce		json
//	@Success		200	{object}	DTO.EmployeeWorkSchedule
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/employee/{employee_id}/work_range/{work_range_id} [delete]
func DeleteEmployeeWorkRange(c *fiber.Ctx) error {
	var err error
	employee_id := c.Params("employee_id")
	work_range_id := c.Params("work_range_id")

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var employee model.Employee
	if err := tx.First(&employee, "id = ?", employee_id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.Employee.NotFound
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	var work_schedule model.EmployeeWorkRange
	if err := tx.First(&work_schedule, "id = ?", work_range_id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	if err := employee.RemoveWorkRange(tx, &work_schedule); err != nil {
		return err
	}

	var ewr []model.EmployeeWorkRange
	if err := tx.
		Preload(clause.Associations).
		Find(&ewr, "employee_id = ?", employee_id).Error; err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}

	ews := model.EmployeeWorkSchedule{
		WorkRanges: ewr,
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &ews, &DTO.EmployeeWorkSchedule{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// UpdateEmployeeWorkRange updates a work range for an employee
//
//	@Summary		Update work range
//	@Description	Update a work range for an employee
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Employee ID"
//	@Param			work_range_id	path		string	true	"Work Range ID"
//	@Accept			json
//	@Produce		json
//	@Param			work_range	body		DTO.UpdateWorkRange	true	"Work Range"
//	@Success		200			{object}	DTO.EmployeeWorkSchedule
//	@Failure		400			{object}	DTO.ErrorResponse
//	@Router			/employee/{id}/work_range/{work_range_id} [put]
func UpdateEmployeeWorkRange(c *fiber.Ctx) error {
	employee_id := c.Params("id")
	work_range_id := c.Params("work_range_id")

	var mapInput map[string]any
	if err := json.Unmarshal(c.Request().Body(), &mapInput); err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("failed to parse request body: %w", err))
	}

	if mapInput["weekday"] == nil || mapInput["start_time"] == nil || mapInput["end_time"] == nil || mapInput["time_zone"] == nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("missing required fields"))
	}

	var work_range model.EmployeeWorkRange

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	if err := tx.First(&work_range, "id = ?", work_range_id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	if work_range.EmployeeID.String() != employee_id {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range employee ID does not match employee ID from path"))
	}

	var input DTO.UpdateWorkRange
	if err := c.BodyParser(&input); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	start, err := lib.Parse_HHMM_To_Time(input.StartTime, input.TimeZone)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid start_time: %w", err))
	}
	end, err := lib.Parse_HHMM_To_Time(input.EndTime, input.TimeZone)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid end_time: %w", err))
	}

	// Atualiza o model manualmente com os dados parseados
	work_range.Weekday = time.Weekday(input.Weekday)
	work_range.StartTime = start
	work_range.EndTime = end
	work_range.TimeZone = input.TimeZone

	if err := tx.Save(&work_range).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range (%s) not found", work_range.ID))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	var ewr []model.EmployeeWorkRange
	if err := tx.
		Preload(clause.Associations).
		Find(&ewr, "employee_id = ?", employee_id).Error; err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}

	ews := model.EmployeeWorkSchedule{
		WorkRanges: ewr,
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &ews, &DTO.EmployeeWorkSchedule{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// AddEmployeeWorkRangeServices adds services to employee's work range
//
//	@Summary		Add services to employee's work range
//	@Description	Add services to an employee's work range
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Failure		401				{object}	nil
//	@Param			employee_id		path		string	true	"Employee ID"
//	@Param			work_range_id	path		string	true	"Work Range ID"
//	@Accept			json
//	@Produce		json
//	@Param			services	body		DTO.EmployeeWorkRangeServices	true	"Services"
//	@Success		200			{object}	DTO.EmployeeWorkSchedule
//	@Failure		400			{object}	DTO.ErrorResponse
//	@Router			/employee/{employee_id}/work_range/{work_range_id}/services [post]
func AddEmployeeWorkRangeServices(c *fiber.Ctx) error {
	employee_id := c.Params("employee_id")
	workRangeID := c.Params("work_range_id")

	var body DTO.EmployeeWorkRangeServices
	if err := c.BodyParser(&body); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	var employee model.Employee
	employee.ID = uuid.MustParse(employee_id)

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var wr model.EmployeeWorkRange
	if err := tx.First(&wr, "id = ?", workRangeID).Error; err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range not found"))
	}

	if wr.EmployeeID.String() != employee.ID.String() {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee ID mismatch"))
	}

	wrServices := make([]*model.Service, 0, len(body.Services))
	for _, s := range body.Services {
		wrServices = append(wrServices, &model.Service{BaseModel: model.BaseModel{ID: s.ID}})
	}

	if err := tx.Model(&wr).Association("Services").Append(wrServices); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	var ewr []model.EmployeeWorkRange
	if err := tx.
		Preload(clause.Associations).
		Find(&ewr, "employee_id = ?", employee_id).Error; err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}

	ews := model.EmployeeWorkSchedule{
		WorkRanges: ewr,
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &ews, &DTO.EmployeeWorkSchedule{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// DeleteEmployeeWorkRangeService removes a service from employee's work range
//
//	@Summary		Remove service from employee's work range
//	@Description	Remove a service from an employee's work range
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Failure		401				{object}	nil
//	@Param			employee_id		path		string	true	"Employee ID"
//	@Param			work_range_id	path		string	true	"Work Range ID"
//	@Param			service_id		path		string	true	"Service ID"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	DTO.EmployeeWorkSchedule
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/employee/{employee_id}/work_range/{work_range_id}/service/{service_id} [delete]
func DeleteEmployeeWorkRangeService(c *fiber.Ctx) error {
	employee_id := c.Params("employee_id")
	workRangeID := c.Params("work_range_id")
	serviceID := c.Params("service_id")

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var wr model.EmployeeWorkRange
	if err := tx.First(&wr, "id = ?", workRangeID).Error; err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range not found"))
	}

	if wr.EmployeeID.String() != employee_id {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee ID mismatch"))
	}

	serviceUUID, err := uuid.Parse(serviceID)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid service ID: %w", err))
	}

	service := &model.Service{BaseModel: model.BaseModel{ID: serviceUUID}}

	if err := tx.Model(&wr).Association("Services").Delete(service); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	var ewr []model.EmployeeWorkRange
	if err := tx.
		Preload(clause.Associations).
		Find(&ewr, "employee_id = ?", employee_id).Error; err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}

	ews := model.EmployeeWorkSchedule{
		WorkRanges: ewr,
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &ews, &DTO.EmployeeWorkSchedule{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// AddEmployeeService adds a service to an employee
//
//	@Summary		Add service to employee
//	@Description	Add a service to an employee
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Failure		401				{object}	nil
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			employee_id	path		string	true	"Employee ID"
//	@Param			service_id	path		string	true	"Service ID"
//	@Success		200			{object}	DTO.EmployeeFull
//	@Failure		404			{object}	DTO.ErrorResponse
//	@Router			/employee/{employee_id}/service/{service_id} [post]
func AddServiceToEmployee(c *fiber.Ctx) error {
	employee_id := c.Params("employee_id")
	service_id := c.Params("service_id")
	var employee model.Employee
	var service model.Service

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	if err := tx.Where("id = ?", service_id).Preload(clause.Associations).First(&service).Error; err != nil {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("service not found"))
	}

	if err := database.LockForUpdate(tx, &employee, "id", employee_id); err != nil {
		return err
	}

	if err := employee.AddService(tx, &service); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &employee, &DTO.EmployeeFull{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// RemoveServiceFromEmployee removes a service from an employee
//
//	@Summary		Remove service from employee
//	@Description	Remove a service from an employee
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Failure		401				{object}	nil
//	@Param			employee_id		path		string	true	"Employee ID"
//	@Param			service_id		path		string	true	"Service ID"
//	@Produce		json
//	@Success		200	{object}	DTO.EmployeeFull
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/employee/{employee_id}/service/{service_id} [delete]
func RemoveServiceFromEmployee(c *fiber.Ctx) error {
	employee_id := c.Params("employee_id")
	service_id := c.Params("service_id")
	var employee model.Employee
	var service model.Service

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	if err := tx.Where("id = ?", service_id).Preload(clause.Associations).First(&service).Error; err != nil {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("service not found"))
	}

	if err := database.LockForUpdate(tx, &employee, "id", employee_id); err != nil {
		return err
	}

	if err := employee.RemoveService(tx, &service); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &employee, &DTO.EmployeeFull{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// AddBranchToEmployee adds an employee to a branch
//
//	@Summary		Add employee to branch
//	@Description	Add an employee to a branch
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Failure		401				{object}	nil
//	@Param			branch_id		path		string	true	"Branch ID"
//	@Param			employee_id		path		string	true	"Employee ID"
//	@Produce		json
//	@Success		200	{object}	DTO.EmployeeFull
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/employee/{employee_id}/branch/{branch_id} [post]
func AddBranchToEmployee(c *fiber.Ctx) error {
	var branch model.Branch
	var employee model.Employee
	branch_id := c.Params("branch_id")
	employee_id := c.Params("employee_id")

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	if err := tx.Where("id = ?", branch_id).Preload(clause.Associations).First(&branch).Error; err != nil {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("branch not found"))
	}

	if err := database.LockForUpdate(tx, &employee, "id", employee_id); err != nil {
		return err
	}

	if err := employee.AddBranch(tx, &branch); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &employee, &DTO.EmployeeFull{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// RemoveBranchFromEmployee removes an employee from a branch
//
//	@Summary		Remove employee from branch
//	@Description	Remove an employee from a branch
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Failure		401				{object}	nil
//	@Param			branch_id		path		string	true	"Branch ID"
//	@Param			employee_id		path		string	true	"Employee ID"
//	@Produce		json
//	@Success		200	{object}	DTO.EmployeeFull
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/employee/{employee_id}/branch/{branch_id} [delete]
func RemoveBranchFromEmployee(c *fiber.Ctx) error {
	var branch model.Branch
	var employee model.Employee
	branch_id := c.Params("branch_id")
	employee_id := c.Params("employee_id")

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	if err := tx.Where("id = ?", branch_id).Preload(clause.Associations).First(&branch).Error; err != nil {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("branch not found"))
	}

	if err := database.LockForUpdate(tx, &employee, "id", employee_id); err != nil {
		return err
	}

	if err := employee.RemoveBranch(tx, &branch); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &employee, &DTO.EmployeeFull{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

func AddRoleToEmployee(c *fiber.Ctx) error {
	employee_id := c.Params("employee_id")
	role_id := c.Params("role_id")
	var employee model.Employee
	var role model.Role

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	if err := tx.Where("id = ?", role_id).Preload(clause.Associations).First(&role).Error; err != nil {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("role not found"))
	}

	if err := database.LockForUpdate(tx, &employee, "id", employee_id); err != nil {
		return err
	}

	if err := employee.AddRole(tx, &role); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &employee, &DTO.EmployeeFull{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

func RemoveRoleFromEmployee(c *fiber.Ctx) error {
	employee_id := c.Params("employee_id")
	role_id := c.Params("role_id")
	var employee model.Employee
	var role model.Role

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	if err := tx.Where("id = ?", role_id).Preload(clause.Associations).First(&role).Error; err != nil {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("role not found"))
	}

	if err := database.LockForUpdate(tx, &employee, "id", employee_id); err != nil {
		return err
	}

	if err := employee.RemoveRole(tx, &role); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &employee, &DTO.EmployeeFull{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// GetEmployeeAppointments retrieves appointments for a specific employee with pagination
//
//	@Summary		Get employee appointments
//	@Description	Retrieve appointments for a specific employee with pagination and filtering
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Failure		401				{object}	nil		"Unauthorized"
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			employee_id				path		string	true	"Employee ID"
//	@Param			page		query		int		false	"Page number"						default(1)
//	@Param			page_size	query		int		false	"Number of items per page"			default(10)
//	@Param			start_date	query		string	false	"Start date filter (DD/MM/YYYY)"	example(21/04/2025)
//	@Param			end_date	query		string	false	"End date filter (DD/MM/YYYY)"		example(31/05/2025)
//	@Param			cancelled	query		string	false	"Filter by cancelled status (true/false)"
//	@Success		200			{object}	DTO.AppointmentList
//	@Failure		400			{object}	DTO.ErrorResponse
//	@Router			/employee/{employee_id}/appointments [get]
func GetEmployeeAppointmentsById(c *fiber.Ctx) error {
	employee_id := c.Params("employee_id")

	// Validate employee_id is not empty and is a valid UUID
	if employee_id == "" {
		return lib.Error.General.ResourceNotFoundError.WithError(fmt.Errorf("employee_id parameter is required"))
	}

	// Validate UUID format
	if _, err := uuid.Parse(employee_id); err != nil {
		return lib.Error.General.ResourceNotFoundError.WithError(fmt.Errorf("invalid employee_id format"))
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	// Verify employee exists
	var count int64
	if err := tx.Model(&model.Employee{}).Where("id = ?", employee_id).Count(&count).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	if count == 0 {
		return lib.Error.General.ResourceNotFoundError.WithError(fmt.Errorf("employee not found"))
	}

	// Parse pagination parameters
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 10)
	offset := (page - 1) * pageSize

	// Parse required timezone parameter
	timezone := c.Query("timezone")
	if timezone == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("timezone parameter is required"))
	}

	// Build query with employee filter
	query := tx.Model(&model.Appointment{}).Where("employee_id = ?", employee_id)

	// Parse and validate date range filters
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr != "" || endDateStr != "" {
		var startDate, endDate time.Time

		if startDateStr != "" {
			startDate, err = time.Parse("02/01/2006", startDateStr)
			if err != nil {
				return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid start_date format, expected DD/MM/YYYY"))
			}
		}

		if endDateStr != "" {
			endDate, err = time.Parse("02/01/2006", endDateStr)
			if err != nil {
				return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid end_date format, expected DD/MM/YYYY"))
			}
			// Set end date to end of day
			endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}

		// Validate date range
		if startDateStr != "" && endDateStr != "" {
			if endDate.Before(startDate) {
				return lib.Error.General.BadRequest.WithError(fmt.Errorf("end_date must be after start_date"))
			}

			daysDiff := endDate.Sub(startDate).Hours() / 24
			if daysDiff > 90 {
				return lib.Error.General.BadRequest.WithError(fmt.Errorf("date range cannot exceed 90 days"))
			}
		}

		// Apply date filters
		if startDateStr != "" {
			query = query.Where("start_time >= ?", startDate)
		}
		if endDateStr != "" {
			query = query.Where("start_time <= ?", endDate)
		}
	}

	// Parse cancelled filter
	cancelledStr := c.Query("cancelled")
	if cancelledStr != "" {
		if cancelledStr == "true" {
			query = query.Where("is_cancelled = ?", true)
		} else if cancelledStr == "false" {
			query = query.Where("is_cancelled = ?", false)
		} else {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("cancelled parameter must be 'true' or 'false'"))
		}
	}

	// Get total count for pagination
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Get appointments with pagination
	var appointments []model.Appointment
	if err := query.
		Offset(offset).
		Limit(pageSize).
		Find(&appointments).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Convert to DTO
	bytes, err := json.Marshal(appointments)
	if err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	var appointmentsDTO []DTO.Appointment
	if err := json.Unmarshal(bytes, &appointmentsDTO); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Collect unique client IDs
	clientIDMap := make(map[uuid.UUID]bool)
	for _, apt := range appointments {
		clientIDMap[apt.ClientID] = true
	}

	// Fetch client information
	var clientInfo []DTO.ClientBasicInfo
	if len(clientIDMap) > 0 {
		clientIDs := make([]uuid.UUID, 0, len(clientIDMap))
		for id := range clientIDMap {
			clientIDs = append(clientIDs, id)
		}

		var clients []model.Client
		if err := tx.Select("id", "name", "surname", "email", "phone").
			Where("id IN ?", clientIDs).
			Find(&clients).Error; err != nil {
			return lib.Error.General.InternalError.WithError(err)
		}

		for _, client := range clients {
			clientInfo = append(clientInfo, DTO.ClientBasicInfo{
				ID:      client.ID,
				Name:    client.Name,
				Surname: client.Surname,
				Email:   client.Email,
				Phone:   client.Phone,
			})
		}
	}

	AppointmentList := DTO.AppointmentList{
		Appointments: appointmentsDTO,
		ClientInfo:   clientInfo,
		Page:         page,
		PageSize:     pageSize,
		TotalCount:   int(totalCount),
	}

	if err := lib.ResponseFactory(c).Send(200, &AppointmentList); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// SendEmployeeVerificationEmail sends a verification email to an employee
//
//	@Summary		Send employee verification email
//	@Description	Send a verification email to an employee
//	@Tags			Employee
//	@Param			company_id	path	string	true	"Company ID"
//	@Param			email		path	string	true	"Employee Email"
//	@Query			language	query	string	false	"Language for the email content"
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/employee/send-verification-code/email/{email}/{company_id} [post]
func SendEmployeeVerificationEmail(c *fiber.Ctx) error {
	return SendVerificationCodeByEmail(c, &model.Employee{})
}

// VerifyEmployeeEmail verifies an employee's email
//
//	@Summary		Verify employee email
//	@Description	Verify an employee's email
//	@Tags			Employee
//	@Param			verification_code	query	string	true	"Verification Code"
//	@Param			company_id			path	string	true	"Company ID"
//	@Param			email				path	string	true	"Employee Email"
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/client/verify-email/{email}/{code}/{company_id} [get]
func VerifyEmployeeEmail(c *fiber.Ctx) error {
	return VerifyEmail(c, &model.Employee{})
}

func Employee(Gorm *handler.Gorm) {
	endpoint := &middleware.Endpoint{DB: Gorm}
	endpoint.BulkRegisterHandler([]fiber.Handler{
		CreateEmployee,
		GetEmployeeById,
		GetEmployeeByEmail,
		UpdateEmployeeById,
		DeleteEmployeeById,
		AddServiceToEmployee,
		RemoveServiceFromEmployee,
		AddBranchToEmployee,
		RemoveBranchFromEmployee,
		LoginEmployeeByPassword,
		LoginEmployeeByEmailCode,
		SendEmployeeLoginValidationCodeByEmail,
		ResetEmployeePasswordByEmail,
		UpdateEmployeeImages,
		DeleteEmployeeImage,
		CreateEmployeeWorkSchedule,
		GetEmployeeWorkSchedule,
		GetEmployeeWorkRangeById,
		DeleteEmployeeWorkRange,
		UpdateEmployeeWorkRange,
		AddEmployeeWorkRangeServices,
		DeleteEmployeeWorkRangeService,
		GetEmployeeAppointmentsById,
		SendEmployeeVerificationEmail,
		VerifyEmployeeEmail,
	})
}
