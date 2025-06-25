package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	database "agenda-kaki-go/core/config/db"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/middleware"
	"agenda-kaki-go/core/service"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/go-playground/validator/v10"
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
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Accept			json
//	@Produce		json
//	@Param			employee	body		DTO.CreateEmployee	true	"Employee"
//	@Success		200			{object}	DTO.Employee
//	@Failure		400			{object}	DTO.ErrorResponse
//	@Router			/employee [post]
func CreateEmployee(c *fiber.Ctx) error {
	var employee model.Employee
	if err := Create(c, &employee); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &employee, &DTO.EmployeeFull{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// LoginEmployee logs an employee in
//
//	@Summary		Login
//	@Description	Log in an client
//	@Tags			Employee
//	@Param			X-Company-ID	header	string	true	"X-Company-ID"
//	@Accept			json
//	@Produce		json
//	@Param			client	body	DTO.LoginEmployee	true	"Employee"
//	@Success		200
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Failure		401	{object}	nil
//	@Router			/employee/login [post]
func LoginEmployee(c *fiber.Ctx) error {
	var body DTO.LoginEmployee
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var employee model.Employee
	if err := tx.Where("email = ?", body.Email).Preload(clause.Associations).First(&employee).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.Employee.NotFound
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	if !employee.Verified {
		return lib.Error.Client.NotVerified
	}

	if !handler.ComparePassword(employee.Password, body.Password) {
		return lib.Error.Auth.InvalidLogin
	}

	var dto DTO.Claims

	if employeeBytes, err := json.Marshal(&employee); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	} else {
		if err := json.Unmarshal(employeeBytes, &dto); err != nil {
			return lib.Error.General.InternalError.WithError(err)
		}
	}

	token, err := handler.JWT(c).Encode(dto)
	if err != nil {
		return err
	}

	c.Response().Header.Set(namespace.HeadersKey.Auth, token)
	return nil
}

// VerifyEmployeeEmail Does the email verification for an employee
//
//	@Summary		Verify email
//	@Description	Verify an employee's email
//	@Tags			Employee
//	@Param			X-Company-ID	header	string	true	"X-Company-ID"
//	@Accept			json
//	@Produce		json
//	@Param			email	path		string	true	"Employee Email"
//	@Param			code	path		string	true	"Verification Code"
//	@Success		200		{object}	nil
//	@Failure		404		{object}	nil
//	@Router			/employee/verify-email/{email}/{code} [post]
func VerifyEmployeeEmail(c *fiber.Ctx) error {
	email := c.Params("email")
	var employee model.Employee
	// Parse the email from the URL as it comes in the form of "john.clark%40gmail.com"
	email, err := url.QueryUnescape(email)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}
	employee.Email = email
	if err := lib.ValidatorV10.Var(employee.Email, "email"); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("email invalid"))
		} else {
			return lib.Error.General.InternalError.WithError(err)
		}
	}
	tx, end, err := database.ContextTransaction(c)
	defer end()
	if err != nil {
		return err
	}
	if err := database.LockForUpdate(tx, &employee, "email", email); err != nil {
		return err
	}
	if employee.Verified {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("email already verified"))
	}
	employee.Verified = true
	if err := tx.Save(&employee).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.Employee.NotFound
		}
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
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			id				path		string	true	"Employee ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Employee
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/employee/{id} [get]
func GetEmployeeById(c *fiber.Ctx) error {
	var employee model.Employee
	if err := GetOneBy("id", c, &employee); err != nil {
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
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			email			path		string	true	"Employee Email"
//	@Produce		json
//	@Success		200	{object}	DTO.Employee
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/employee/email/{email} [get]
func GetEmployeeByEmail(c *fiber.Ctx) error {
	var employee model.Employee
	if err := GetOneBy("email", c, &employee); err != nil {
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
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string						true	"Employee ID"
//	@Param			employee	body		DTO.UpdateEmployeeSwagger	true	"Employee"
//	@Success		200			{object}	DTO.Employee
//	@Failure		400			{object}	DTO.ErrorResponse
//	@Router			/employee/{id} [patch]
func UpdateEmployeeById(c *fiber.Ctx) error {
	var employee model.Employee
	if err := UpdateOneById(c, &employee); err != nil {
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
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			id				path		string	true	"Employee ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Employee
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/employee/{id} [delete]
func DeleteEmployeeById(c *fiber.Ctx) error {
	return DeleteOneById(c, &model.Employee{})
}

// AddEmployeeWorkSchedule creates a work schedule for an employee
//
//	@Summary		Create work schedule
//	@Description	Create a work schedule for an employee
//	@Tags			EmployeeWorkSchedule
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil	"Unauthorized"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Accept			json
//	@Param			work_schedule	body		DTO.CreateEmployeeWorkSchedule	true	"Work Schedule"
//	@Param			id				path		string	true	"Employee ID"
//	@Success		200		{object}	DTO.EmployeeFull
//	@Failure		400		{object}	lib.ErrorResponse
//	@Router			/employee/{id}/work_schedule [post]
func AddEmployeeWorkSchedule(c *fiber.Ctx) error {
	var input DTO.CreateEmployeeWorkSchedule
	if err := c.BodyParser(&input); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	var EmployeeWorkSchedule model.EmployeeWorkSchedule

	for _, wr := range input.WorkRanges {
		loc, err := time.LoadLocation(wr.TimeZone)
		if err != nil {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid timezone: %w", err))
		}

		start, err := time.ParseInLocation("15:04", wr.StartTime, loc)
		if err != nil {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid start_time: %w", err))
		}
		end, err := time.ParseInLocation("15:04", wr.EndTime, loc)
		if err != nil {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid end_time: %w", err))
		}

		services := make([]*model.Service, 0, len(wr.Services))
		for _, serviceID := range wr.Services {
			services = append(services, &model.Service{BaseModel: model.BaseModel{ID: serviceID.ID}})
		}

		EmployeeWorkSchedule.WorkRanges = append(EmployeeWorkSchedule.WorkRanges, model.EmployeeWorkRange{
			Weekday:    time.Weekday(wr.Weekday),
			StartTime:  start,
			EndTime:    end,
			TimeZone:   wr.TimeZone,
			EmployeeID: wr.EmployeeID,
			BranchID:   wr.BranchID,
			Services:   services,
		})
	}

	employee_id := c.Params("id")

	Service := service.Factory(c)
	defer Service.DeferDB()

	for i, wr := range EmployeeWorkSchedule.WorkRanges {
		if wr.EmployeeID.String() != employee_id {
			Service.MyGorm.DB.Rollback()
			return lib.Error.General.CreatedError.WithError(fmt.Errorf("work range [%d] employee ID (%s) does not match employee ID (%s) from path", i+1, wr.EmployeeID.String(), employee_id))
		}
		if err := Service.Create(&wr).Error; err != nil {
			Service.MyGorm.DB.Rollback()
			return lib.Error.General.CreatedError.WithError(err)
		}
	}

	var employee model.Employee
	if err := Service.MyGorm.DB.
		Preload(clause.Associations).
		First(&employee, "id = ?", employee_id).Error; err != nil {
		Service.MyGorm.DB.Rollback()
		return lib.Error.General.CreatedError.WithError(err)
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &employee, &DTO.EmployeeFull{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// DeleteEmployeeWorkRange deletes a work schedule for an employee
//
//	@Summary		Delete work schedule
//	@Description	Delete a work schedule for an employee
//	@Tags			EmployeeWorkSchedule
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil	"Unauthorized"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			id				path		string	true	"Employee ID"
//	@Param			work_range_id	path		string	true	"Work Range ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Employee
//	@Failure		400	{object}	lib.ErrorResponse
//	@Router			employee/{id}/work_schedule/{work_range_id} [delete]
func DeleteEmployeeWorkRange(c *fiber.Ctx) error {
	employee_id := c.Params("id")
	work_range_id := c.Params("work_range_id")

	tx, end, err := database.ContextTransaction(c)
	defer end()
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

	if err := tx.Delete(&work_schedule).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	if err := tx.Preload(clause.Associations).First(&employee, "id = ?", employee_id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &employee, &DTO.EmployeeFull{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// UpdateEmployeeWorkRange updates a work range for an employee
//
//	@Summary		Update work range
//	@Description	Update a work range for an employee
//	@Tags			EmployeeWorkRange
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			id				path		string	true	"Employee ID"
//	@Param			work_range_id	path		string	true	"Work Range ID"
//	@Accept			json
//	@Produce		json
//	@Param			work_range	body		DTO.UpdateWorkRange	true	"Work Range"
//	@Success		200		{object}	DTO.EmployeeFull
//	@Failure		400		{object}	lib.ErrorResponse
//	@Router			/employee/{id}/work_range/{work_range_id} [put]
func UpdateEmployeeWorkRange(c *fiber.Ctx) error {
	employee_id := c.Params("id")
	work_range_id := c.Params("work_range_id")

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

	loc, err := time.LoadLocation(input.TimeZone)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid timezone: %w", err))
	}

	start, err := time.ParseInLocation("15:04", input.StartTime, loc)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid start_time: %w", err))
	}

	end, err := time.ParseInLocation("15:04", input.EndTime, loc)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid end_time: %w", err))
	}

	// Atualiza o model manualmente com os dados parseados
	work_range.Weekday = time.Weekday(input.Weekday)
	work_range.StartTime = start
	work_range.EndTime = end
	work_range.TimeZone = input.TimeZone

	if err := UpdateOneById(c, &work_range); err != nil {
		return err
	}

	var employee model.Employee
	if err := tx.Preload(clause.Associations).First(&employee, "id = ?", employee_id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &employee, &DTO.EmployeeFull{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// AddEmployeeWorkRangeServices adds services to employee's work range
//
//	@Summary		Add services to employee's work range
//	@Description	Add services to an employee's work range
//	@Tags			EmployeeWorkRange
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			employee_id	path		string	true	"Employee ID"
//	@Param			work_range_id	path		string	true	"Work Range ID"
//	@Accept			json
//	@Produce		json
//	@Param			services	body		[]DTO.ServiceID	true	"Services"
//	@Success		200	{object}	DTO.EmployeeFull
//	@Failure		400	{object}	lib.ErrorResponse
//	@Router			/employee/{employee_id}/work_range/{work_range_id}/services [post]
func AddEmployeeWorkRangeServices(c *fiber.Ctx) error {
	employee_id := c.Params("employee_id")
	work_range_id := c.Params("work_range_id")

	var services []model.Service
	if err := c.BodyParser(&services); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	var employee model.Employee
	employee.ID = uuid.MustParse(employee_id)

	tx, end, err := database.ContextTransaction(c)
	defer end()
	if err != nil {
		return err
	}

	if err := employee.AddServicesToWorkRange(tx, work_range_id, services); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// DeleteEmployeeWorkRangeService removes a service from employee's work range
//
//		@Summary		Remove service from employee's work range
//		@Description	Remove a service from an employee's work range
//		@Tags			EmployeeWorkRange
//		@Security		ApiKeyAuth
//		@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//		@Failure		401				{object}	nil
//		@Param			X-Company-ID	header		string	true	"X-Company-ID"
//		@Param			employee_id	path		string	true	"Employee ID"
//		@Param			work_range_id	path		string	true	"Work Range ID"
//	 @Param 			service_id	path		string	true	"Service ID"
//		@Accept			json
//		@Produce		json
//		@Success		200	{object}	DTO.EmployeeFull
//		@Failure		400	{object}	lib.ErrorResponse
//		@Router			/employee/{employee_id}/work_range/{work_range_id}/service/{service_id} [delete]
func DeleteEmployeeWorkRangeService(c *fiber.Ctx) error {
	employee_id := c.Params("employee_id")
	work_range_id := c.Params("work_range_id")
	service_id := c.Params("service_id")

	tx, end, err := database.ContextTransaction(c)
	defer end()
	if err != nil {
		return err
	}

	var employee model.Employee
	employee.ID = uuid.MustParse(employee_id)
	if err := employee.RemoveServiceFromWorkRange(tx, work_range_id, service_id); err != nil {
		tx.Rollback()
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &employee, &DTO.EmployeeFull{}); err != nil {
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
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			employee_id	path		string	true	"Employee ID"
//	@Param			service_id	path		string	true	"Service ID"
//	@Success		200			{object}	DTO.Employee
//	@Failure		404			{object}	DTO.ErrorResponse
//	@Router			/employee/{employee_id}/service/{service_id} [post]
func AddServiceToEmployee(c *fiber.Ctx) error {
	employee_id := c.Params("employee_id")
	service_id := c.Params("service_id")
	var employee model.Employee
	var service model.Service

	tx, end, err := database.ContextTransaction(c)
	defer end()
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
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			employee_id		path		string	true	"Employee ID"
//	@Param			service_id		path		string	true	"Service ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Employee
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/employee/{employee_id}/service/{service_id} [delete]
func RemoveServiceFromEmployee(c *fiber.Ctx) error {
	employee_id := c.Params("employee_id")
	service_id := c.Params("service_id")
	var employee model.Employee
	var service model.Service

	tx, end, err := database.ContextTransaction(c)
	defer end()
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
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			branch_id		path		string	true	"Branch ID"
//	@Param			employee_id		path		string	true	"Employee ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Employee
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/employee/{employee_id}/branch/{branch_id} [post]
func AddBranchToEmployee(c *fiber.Ctx) error {
	var branch model.Branch
	var employee model.Employee
	branch_id := c.Params("branch_id")
	employee_id := c.Params("employee_id")

	tx, end, err := database.ContextTransaction(c)
	defer end()
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
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			branch_id		path		string	true	"Branch ID"
//	@Param			employee_id		path		string	true	"Employee ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Employee
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/employee/{employee_id}/branch/{branch_id} [delete]
func RemoveBranchFromEmployee(c *fiber.Ctx) error {
	var branch model.Branch
	var employee model.Employee
	branch_id := c.Params("branch_id")
	employee_id := c.Params("employee_id")

	tx, end, err := database.ContextTransaction(c)
	defer end()
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

	tx, end, err := database.ContextTransaction(c)
	defer end()
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

	tx, end, err := database.ContextTransaction(c)
	defer end()
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
		LoginEmployee,
		VerifyEmployeeEmail,
		AddEmployeeWorkSchedule,
		DeleteEmployeeWorkRange,
		UpdateEmployeeWorkRange,
		AddEmployeeWorkRangeServices,
		DeleteEmployeeWorkRangeService,
	})
}
