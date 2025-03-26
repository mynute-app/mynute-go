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
//	@Description	Log in an client
//	@Tags			Employee
//	@Accept			json
//	@Produce		json
//	@Param			client	body	DTO.LoginEmployee	true	"Employee"
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
		return lib.Error.Client.NotVerified.SendToClient(c)
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
	if err := ec.Request.Gorm.UpdateOneById(fmt.Sprintf("%v", employee.ID), &model.Employee{}, map[string]interface{}{"verified": true}, []string{}); err != nil {
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
//	@Success		200	{object}	DTO.Employee
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/employee/{id} [get]
func (ec *employee_controller) GetEmployeeById(c *fiber.Ctx) error {
	return ec.GetBy("id", c)
}

// GetEmployeeByEmail retrieves an employee by email
//
//	@Summary		Get employee by email
//	@Description	Retrieve an employee by its email
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Authorization"
//	@Failure		401				{object}	nil
//	@Param			email			path		string	true	"Employee Email"
//	@Produce		json
//	@Success		200	{object}	DTO.Employee
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/employee/email/{email} [get]
func (ec *employee_controller) GetEmployeeByEmail(c *fiber.Ctx) error {
	return ec.GetBy("email", c)
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
//	@Success		200			{object}	DTO.Employee
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
//	@Success		200	{object}	DTO.Employee
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
//	@Param			Authorization	header		string	true	"Authorization"
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			employee_id	path		string	true	"Employee ID"
//	@Param			service_id	path		string	true	"Service ID"
//	@Success		200			{object}	DTO.Employee
//	@Failure		404			{object}	DTO.ErrorResponse
//	@Router			/employee/{employee_id}/service/{service_id} [post]
func (ec *employee_controller) AddServiceToEmployee(c *fiber.Ctx) error {
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
	if err := ec.Request.Gorm.GetOneBy("id", employee_id, &employee, ec.Associations); err != nil {
		return err
	}
	res := &lib.SendResponse{Ctx: c}
	if err := res.SendDTO(200, &employee, &DTO.Employee{}); err != nil {
		return err
	}
	return nil
}

// RemoveServiceFromEmployee removes a service from an employee
//
//	@Summary		Remove service from employee
//	@Description	Remove a service from an employee
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Authorization"
//	@Failure		401				{object}	nil
//	@Param			employee_id		path		string	true	"Employee ID"
//	@Param			service_id		path		string	true	"Service ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Employee
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/employee/{employee_id}/service/{service_id} [delete]
func (ec *employee_controller) RemoveServiceFromEmployee(c *fiber.Ctx) error {
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
	if err := ec.Request.Gorm.DB.Model(&employee).Association("Services").Delete(&service); err != nil {
		return err
	}
	if err := ec.Request.Gorm.GetOneBy("id", employee_id, &employee, ec.Associations); err != nil {
		return err
	}
	res := &lib.SendResponse{Ctx: c}
	if err := res.SendDTO(200, &employee, &DTO.Employee{}); err != nil {
		return err
	}
	return nil
}

// AddBranchToEmployee adds an employee to a branch
//
//	@Summary		Add employee to branch
//	@Description	Add an employee to a branch
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Authorization"
//	@Failure		401				{object}	nil
//	@Param			branch_id		path		string	true	"Branch ID"
//	@Param			employee_id		path		string	true	"Employee ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Employee
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/employee/{employee_id}/branch/{branch_id} [post]
func (ec *employee_controller) AddBranchToEmployee(c *fiber.Ctx) error {
	var branch model.Branch
	var employee model.Employee
	branch_id := c.Params("branch_id")
	employee_id := c.Params("employee_id")
	if err := ec.Request.Gorm.GetOneBy("id", employee_id, &employee, ec.Associations); err != nil {
		return err
	}
	if err := ec.Request.Gorm.GetOneBy("id", branch_id, &branch, []string{}); err != nil {
		return err
	}
	if employee.CompanyID != branch.CompanyID {
		return lib.Error.Company.NotSame.SendToClient(c)
	}
	if err := ec.Request.Gorm.DB.Model(&employee).Association("Branches").Append(&branch); err != nil {
		return err
	}
	if err := ec.Request.Gorm.GetOneBy("id", employee_id, &employee, ec.Associations); err != nil {
		return err
	}
	res := &lib.SendResponse{Ctx: c}
	if err := res.SendDTO(200, &employee, &DTO.Employee{}); err != nil {
		return err
	}
	return nil
}

// RemoveBranchFromEmployee removes an employee from a branch
//
//	@Summary		Remove employee from branch
//	@Description	Remove an employee from a branch
//	@Tags			Employee
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Authorization"
//	@Failure		401				{object}	nil
//	@Param			branch_id		path		string	true	"Branch ID"
//	@Param			employee_id		path		string	true	"Employee ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Employee
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/employee/{employee_id}/branch/{branch_id} [delete]
func (ec *employee_controller) RemoveBranchFromEmployee(c *fiber.Ctx) error {
	var branch model.Branch
	var employee model.Employee
	branch_id := c.Params("branch_id")
	employee_id := c.Params("employee_id")
	if err := ec.Request.Gorm.GetOneBy("id", employee_id, &employee, ec.Associations); err != nil {
		return err
	}
	if err := ec.Request.Gorm.GetOneBy("id", branch_id, &branch, []string{}); err != nil {
		return err
	}
	if employee.CompanyID != branch.CompanyID {
		return lib.Error.Company.NotSame.SendToClient(c)
	}
	if err := ec.Request.Gorm.DB.Model(&employee).Association("Branches").Delete(&branch); err != nil {
		return err
	}
	if err := ec.Request.Gorm.GetOneBy("id", employee_id, &employee, ec.Associations); err != nil {
		return err
	}
	res := &lib.SendResponse{Ctx: c}
	if err := res.SendDTO(200, &employee, &DTO.Employee{}); err != nil {
		return err
	}
	return nil
}

func Employee(Gorm *handler.Gorm) *employee_controller {
	ec := &employee_controller{
		Base: service.Base[model.Employee, DTO.Employee]{
			Name:         namespace.HolidaysKey.Name,
			Request:      handler.Request(Gorm),
			Associations: []string{"Branches", "Company", "Services", "Appointments"},
		},
	}
	route := &handler.Route{DB: Gorm.DB}
	route.Register("/employee", "post", "private", ec.CreateEmployee, "Create employee").Save()
	route.Register("/employee/login", "post", "public", ec.LoginEmployee, "Login employee").Save()
	route.Register("/employee/verify-email/:email/:code", "post", "public", ec.VerifyEmployeeEmail, "Verify employee email").Save()
	route.Register("/employee/:id", "get", "private", ec.GetEmployeeById, "Get employee by ID").Save()
	route.Register("/employee/email/:email", "get", "private", ec.GetEmployeeByEmail, "Get employee by email").Save()
	route.Register("/employee/:id", "patch", "private", ec.UpdateEmployeeById, "Update employee by ID").Save()
	route.Register("/employee/:id", "delete", "private", ec.DeleteEmployeeById, "Delete employee by ID").Save()
	route.Register("/employee/:employee_id/service/:service_id", "post", "private", ec.AddServiceToEmployee, "Add service to employee").Save()
	route.Register("/employee/:employee_id/service/:service_id", "delete", "private", ec.RemoveServiceFromEmployee, "Remove service from employee").Save()
	route.Register("/employee/:employee_id/branch/:branch_id", "post", "private", ec.AddBranchToEmployee, "Add employee to branch").Save()
	route.Register("/employee/:employee_id/branch/:branch_id", "delete", "private", ec.RemoveBranchFromEmployee, "Remove employee from branch").Save()
	return ec
}
