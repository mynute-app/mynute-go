package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	database "agenda-kaki-go/core/config/db"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/middleware"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreateEmployee creates an employee
//
//	@Summary		Create employee
//	@Description	Create an employee
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
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

	if err := lib.ResponseFactory(c).SendDTO(200, &employee, &DTO.Employee{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// LoginEmployee logs an employee in
//
//	@Summary		Login
//	@Description	Log in an client
//	@Tags			Employee
//	@Accept			json
//	@Produce		json
//	@Param			client	body	DTO.LoginEmployee	true	"Employee"
//	@Success		200
//	@Failure		404	{object}	DTO.ErrorResponse
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
	employee.Email = email
	if err := lib.ValidatorV10.Var(employee.Email, "email"); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
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
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Employee ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Employee
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/employee/{id} [get]
func GetEmployeeById(c *fiber.Ctx) error {
	var employee model.Employee
	if err := GetOneBy("id", c, &employee); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &employee, &DTO.Employee{}); err != nil {
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
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			email			path		string	true	"Employee Email"
//	@Produce		json
//	@Success		200	{object}	DTO.Employee
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/employee/email/{email} [get]
func GetEmployeeByEmail(c *fiber.Ctx) error {
	var employee model.Employee
	if err := GetOneBy("email", c, &employee); err != nil {
		return err
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &employee, &DTO.Employee{}); err != nil {
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
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
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
	if err := lib.ResponseFactory(c).SendDTO(200, &employee, &DTO.Employee{}); err != nil {
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
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Employee ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Employee
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/employee/{id} [delete]
func DeleteEmployeeById(c *fiber.Ctx) error {
	return DeleteOneById(c, &model.Employee{})
}

// AddEmployeeService adds a service to an employee
//
//	@Summary		Add service to employee
//	@Description	Add a service to an employee
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"X-Auth-Token"
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

	if err := lib.ResponseFactory(c).SendDTO(200, &employee, &DTO.Employee{}); err != nil {
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
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			employee_id		path		string	true	"Employee ID"
//	@Param			service_id		path		string	true	"Service ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Employee
//	@Failure		404	{object}	DTO.ErrorResponse
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

	if err := lib.ResponseFactory(c).SendDTO(200, &employee, &DTO.Employee{}); err != nil {
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
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			branch_id		path		string	true	"Branch ID"
//	@Param			employee_id		path		string	true	"Employee ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Employee
//	@Failure		404	{object}	DTO.ErrorResponse
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

	if err := lib.ResponseFactory(c).SendDTO(200, &employee, &DTO.Employee{}); err != nil {
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
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			branch_id		path		string	true	"Branch ID"
//	@Param			employee_id		path		string	true	"Employee ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Employee
//	@Failure		404	{object}	DTO.ErrorResponse
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

	if err := lib.ResponseFactory(c).SendDTO(200, &employee, &DTO.Employee{}); err != nil {
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

	if err := lib.ResponseFactory(c).SendDTO(200, &employee, &DTO.Employee{}); err != nil {
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

	if err := lib.ResponseFactory(c).SendDTO(200, &employee, &DTO.Employee{}); err != nil {
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
	})
}
