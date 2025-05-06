package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	database "agenda-kaki-go/core/config/db"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/middleware"
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
//	@Param			Authorization	header		string	true	"Authorization"
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
	var body model.Employee
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	tx, err := database.Session(c)
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
func VerifyEmployeeEmail(c *fiber.Ctx) error {
	email := c.Params("email")
	var employee model.Employee
	employee.Email = email
	if err := lib.ValidatorV10.Var(employee.Email, "email"); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}
	tx, end, err := database.Transaction(c)
	defer end()
	if err != nil {
		return err
	}
	if err := database.LockForUpdate(tx, &employee, email); err != nil {
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
//	@Param			Authorization	header		string	true	"Authorization"
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
//	@Param			Authorization	header		string	true	"Authorization"
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
//	@Param			Authorization	header		string	true	"Authorization"
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
//	@Param			Authorization	header		string	true	"Authorization"
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
//	@Param			Authorization	header		string	true	"Authorization"
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
	
	tx, end, err := database.Transaction(c)
	defer end()
	if err != nil {
		return err
	}

	if err := tx.Where("id = ?", service_id).Preload(clause.Associations).First(&service).Error; err != nil {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("service not found"))
	}

	if err := database.LockForUpdate(tx, &employee, employee_id); err != nil {
		return err
	}

	if employee.CompanyID != service.CompanyID {
		return lib.Error.Company.NotSame
	}

	// TODO: Continue from here...

	if err := ec.Request.Gorm.DB.Model(&employee).Association("Services").Append(&service); err != nil {
		return err
	}
	if err := ec.Request.Gorm.GetOneBy("id", employee_id, &employee); err != nil {
		return err
	}
	res := &lib.SendResponseStruct{Ctx: c}
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
func RemoveServiceFromEmployee(c *fiber.Ctx) error {
	employee_id := c.Params("employee_id")
	service_id := c.Params("service_id")
	var employee model.Employee
	var service model.Service
	if err := ec.Request.Gorm.GetOneBy("id", employee_id, &employee); err != nil {
		return err
	}
	if err := ec.Request.Gorm.GetOneBy("id", service_id, &service); err != nil {
		return err
	}
	if employee.CompanyID != service.CompanyID {
		return lib.Error.Company.NotSame
	}
	if err := ec.Request.Gorm.DB.Model(&employee).Association("Services").Delete(&service); err != nil {
		return err
	}
	if err := ec.Request.Gorm.GetOneBy("id", employee_id, &employee); err != nil {
		return err
	}
	res := &lib.SendResponseStruct{Ctx: c}
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
func AddBranchToEmployee(c *fiber.Ctx) error {
	var branch model.Branch
	var employee model.Employee
	branch_id := c.Params("branch_id")
	employee_id := c.Params("employee_id")
	if err := ec.Request.Gorm.GetOneBy("id", employee_id, &employee); err != nil {
		return err
	}
	if err := ec.Request.Gorm.GetOneBy("id", branch_id, &branch); err != nil {
		return err
	}
	if employee.CompanyID != branch.CompanyID {
		return lib.Error.Company.NotSame
	}
	if err := ec.Request.Gorm.DB.Model(&employee).Association("Branches").Append(&branch); err != nil {
		return err
	}
	if err := ec.Request.Gorm.GetOneBy("id", employee_id, &employee); err != nil {
		return err
	}
	res := &lib.SendResponseStruct{Ctx: c}
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
func RemoveBranchFromEmployee(c *fiber.Ctx) error {
	var branch model.Branch
	var employee model.Employee
	branch_id := c.Params("branch_id")
	employee_id := c.Params("employee_id")
	if err := ec.Request.Gorm.GetOneBy("id", employee_id, &employee); err != nil {
		return err
	}
	if err := ec.Request.Gorm.GetOneBy("id", branch_id, &branch); err != nil {
		return err
	}
	if employee.CompanyID != branch.CompanyID {
		return lib.Error.Company.NotSame
	}
	if err := ec.Request.Gorm.DB.Model(&employee).Association("Branches").Delete(&branch); err != nil {
		return err
	}
	if err := ec.Request.Gorm.GetOneBy("id", employee_id, &employee); err != nil {
		return err
	}
	res := &lib.SendResponseStruct{Ctx: c}
	if err := res.SendDTO(200, &employee, &DTO.Employee{}); err != nil {
		return err
	}
	return nil
}

func AddRoleToEmployee(c *fiber.Ctx) error {
	employee_id := c.Params("employee_id")
	role_id := c.Params("role_id")
	var employee model.Employee
	var role model.Role
	if err := ec.Request.Gorm.GetOneBy("id", employee_id, &employee); err != nil {
		return err
	}
	if err := ec.Request.Gorm.GetOneBy("id", role_id, &role); err != nil {
		return err
	}
	if role.CompanyID != nil && employee.CompanyID != *role.CompanyID {
		return lib.Error.Company.NotSame
	}
	if err := ec.Request.Gorm.DB.Model(&employee).Association("Roles").Append(&role); err != nil {
		return err
	}
	if err := ec.Request.Gorm.GetOneBy("id", employee_id, &employee); err != nil {
		return err
	}
	res := &lib.SendResponseStruct{Ctx: c}
	if err := res.SendDTO(200, &employee, &DTO.Employee{}); err != nil {
		return err
	}
	return nil
}

func RemoveRoleFromEmployee(c *fiber.Ctx) error {
	employee_id := c.Params("employee_id")
	role_id := c.Params("role_id")
	var employee model.Employee
	var role model.Role
	if err := ec.Request.Gorm.GetOneBy("id", employee_id, &employee); err != nil {
		return err
	}
	if err := ec.Request.Gorm.GetOneBy("id", role_id, &role); err != nil {
		return err
	}
	if role.CompanyID != nil && employee.CompanyID != *role.CompanyID {
		return lib.Error.Company.NotSame
	}
	if err := ec.Request.Gorm.DB.Model(&employee).Association("Roles").Delete(&role); err != nil {
		return err
	}
	if err := ec.Request.Gorm.GetOneBy("id", employee_id, &employee); err != nil {
		return err
	}
	res := &lib.SendResponseStruct{Ctx: c}
	if err := res.SendDTO(200, &employee, &DTO.Employee{}); err != nil {
		return err
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
