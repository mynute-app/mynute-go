package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/service"

	"github.com/gofiber/fiber/v2"
)

type employee_controller struct {
	service.Base[model.Employee, DTO.Employee]
}

// CreateEmployee creates an employee
//
//	@Summary		Create employee
//	@Description	Create an employee
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Authorization"
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			employee	body		DTO.CreateEmployee	true	"Employee"
//	@Success		200			{object}	DTO.Employee
//	@Failure		400			{object}	DTO.ErrorResponse
//	@Router			/employee [post]
func (ec *employee_controller) CreateEmployee(c *fiber.Ctx) error {
	return ec.CreateOne(c)
}

// LoginEmployee logs an employee in
//
//	@Summary		Login
//	@Description	Log in an user
//	@Tags			Employee
//	@Accept			json
//	@Produce		json
//	@Param			user	body	DTO.LoginEmployee	true	"Employee"
//	@Success		200
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/employee/login [post]
func (ec *employee_controller) LoginEmployee(c *fiber.Ctx) error {
	var body model.Employee
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	var employee model.Employee
	if err := ec.Request.Gorm.GetOneBy("email", body.Email, &employee, ec.Associations); err != nil {
		return err
	}
	if !employee.Verified {
		return lib.Error.User.NotVerified.SendToClient(c)
	}
	if !handler.ComparePassword(employee.Password, body.Password) {
		return lib.Error.Auth.InvalidLogin.SendToClient(c)
	}
	token, err := handler.JWT(c).Encode(employee)
	if err != nil {
		return err
	}
	c.Response().Header.Set("Authorization", token)
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
func (ec *employee_controller) VerifyEmployeeEmail(c *fiber.Ctx) error {
	res := &lib.SendResponse{Ctx: c}
	email := c.Params("email")
	var employee model.Employee
	employee.Email = email
	if err := lib.ValidatorV10.Var(employee.Email, "email"); err != nil {
		return res.Send(400, err)
	}
	if err := ec.Request.Gorm.GetOneBy("email", email, &employee, []string{}); err != nil {
		return err
	}
	// code := c.Params("code")
	// }
	// if employee.VerificationCode != code {
	// 	return lib.Error.Auth.EmailCodeInvalid.SendToClient(c)
	// }
	employee.Verified = true
	if err := ec.Request.Gorm.DB.Save(&employee).Error; err != nil {
		return err
	}
	return nil
}

// GetEmployeeById retrieves an employee by ID
//
//	@Summary		Get employee by ID
//	@Description	Retrieve an employee by its ID
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Authorization"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Employee ID"
//	@Produce		json
//	@Success		200	{object}	DTO.CreateEmployee
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/employee/{id} [get]
func (ec *employee_controller) GetEmployeeById(c *fiber.Ctx) error {
	return ec.GetBy("id", c)
}

// UpdateEmployeeById updates an employee by ID
//
//	@Summary		Update employee
//	@Description	Update an employee
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Authorization"
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string						true	"Employee ID"
//	@Param			employee	body		DTO.UpdateEmployeeSwagger	true	"Employee"
//	@Success		200			{object}	DTO.CreateEmployee
//	@Failure		400			{object}	DTO.ErrorResponse
//	@Router			/employee/{id} [patch]
func (ec *employee_controller) UpdateEmployeeById(c *fiber.Ctx) error {
	return ec.UpdateOneById(c)
}

// DeleteEmployeeById deletes an employee by ID
//
//	@Summary		Delete employee by ID
//	@Description	Delete an employee by its ID
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Authorization"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Employee ID"
//	@Produce		json
//	@Success		200	{object}	DTO.CreateEmployee
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/employee/{id} [delete]
func (ec *employee_controller) DeleteEmployeeById(c *fiber.Ctx) error {
	return ec.DeleteOneById(c)
}

// AddEmployeeService adds a service to an employee
//
//	@Summary		Add service to employee
//	@Description	Add a service to an employee
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Authorization
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			employee_id	path	string	true	"Employee ID"
//	@Param			service_id	path	string	true	"Service ID"
//	@Success		200	{object}	DTO.Service
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/employee/{employee_id}/service/{service_id} [post]
func (ec *employee_controller) AddEmployeeService(c *fiber.Ctx) error {
	employee_id := c.Params("employee_id")
	service_id := c.Params("service_id")
	var employee model.Employee
	var service model.Service
	if err := ec.Request.Gorm.GetOneBy("id", employee_id, &employee, ec.Associations); err != nil {
		return err
	}
	if err := ec.Request.Gorm.GetOneBy("id", service_id, &service, []string{}); err != nil {
		return err
	}
	if employee.CompanyID != service.CompanyID {
		return lib.Error.Company.NotSame.SendToClient(c)
	}
	if err := ec.Request.Gorm.DB.Model(&employee).Association("Services").Append(&service); err != nil {
		return err
	}
	res := &lib.SendResponse{Ctx: c}
	res.SendDTO(200, &employee, &DTO.Employee{})
	return nil
}

func Employee(Gorm *handler.Gorm) *employee_controller {
	return &employee_controller{
		Base: service.Base[model.Employee, DTO.Employee]{
			Name:         namespace.HolidaysKey.Name,
			Request:      handler.Request(Gorm),
			Associations: []string{"Branches", "Company", "Services"},
		},
	}
}
