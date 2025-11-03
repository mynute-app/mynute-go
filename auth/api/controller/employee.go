package controller

import (
	DTO "mynute-go/auth/config/dto"
	"mynute-go/auth/lib"
	authModel "mynute-go/auth/config/db/model"

	"github.com/gofiber/fiber/v2"
)

// =====================
// EMPLOYEE MANAGEMENT
// =====================

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
//	@Success		200			{object}	DTO.EmployeeBase
//	@Failure		400			{object}	DTO.ErrorResponse
//	@Router			/users/employee [post]
func CreateEmployee(c *fiber.Ctx) error {
	var user authModel.User
	user.Type = "employee"
	if err := CreateUser(c, &user); err != nil {
		return err
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &user, &DTO.EmployeeBase{}); err != nil {
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
//	@Success		200	{object}	DTO.EmployeeBase
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/users/employee/{id} [get]
func GetEmployeeById(c *fiber.Ctx) error {
	var user authModel.User
	if err := GetOneBy("id", c, &user); err != nil {
		return err
	}
	if user.Type != "employee" {
		return lib.Error.General.RecordNotFound
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &user, &DTO.EmployeeBase{}); err != nil {
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
//	@Success		200	{object}	DTO.EmployeeBase
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/users/employee/email/{email} [get]
func GetEmployeeByEmail(c *fiber.Ctx) error {
	var user authModel.User
	if err := GetOneBy("email", c, &user); err != nil {
		return err
	}
	if user.Type != "employee" {
		return lib.Error.General.RecordNotFound
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &user, &DTO.EmployeeBase{}); err != nil {
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
//	@Success		200			{object}	DTO.EmployeeBase
//	@Failure		400			{object}	DTO.ErrorResponse
//	@Router			/users/employee/{id} [patch]
func UpdateEmployeeById(c *fiber.Ctx) error {
	var user authModel.User
	if err := UpdateOneById(c, &user); err != nil {
		return err
	}
	if user.Type != "employee" {
		return lib.Error.General.RecordNotFound
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &user, &DTO.EmployeeBase{}); err != nil {
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
//	@Success		200	{object}	nil
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/users/employee/{id} [delete]
func DeleteEmployeeById(c *fiber.Ctx) error {
	return DeleteOneById(c, &authModel.User{})
}

