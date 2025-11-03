package controller

import (
	DTO "mynute-go/core/src/config/api/dto"
	"mynute-go/core/src/config/db/model"
	"mynute-go/core/src/lib"
	"mynute-go/core/src/service"

	"github.com/gofiber/fiber/v2"
)

// =====================
// CLIENT MANAGEMENT
// =====================

// CreateClient creates a client
//
//	@Summary		Create client
//	@Description	Create a client
//	@Tags			Client
//	@Accept			json
//	@Produce		json
//	@Param			client	body		DTO.CreateClient	true	"Client"
//	@Success		200		{object}	DTO.Client
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/users/client [post]
func CreateClient(c *fiber.Ctx) error {
	var client model.Client
	if err := CreateUser(c, &client); err != nil {
		return err
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &client, &DTO.Client{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// GetClientByEmail retrieves a client by email
//
//	@Summary		Get client by email
//	@Description	Retrieve a client by its email
//	@Tags			Client
//	@Param			email	path		string	true	"Client Email"
//	@Produce		json
//	@Success		200	{object}	DTO.Client
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/users/client/email/{email} [get]
func GetClientByEmail(c *fiber.Ctx) error {
	var client model.Client
	if err := GetOneBy("email", c, &client, nil, &[]string{"Appointments"}); err != nil {
		return err
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &client, &DTO.Client{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// GetClientById retrieves a client by ID
//
//	@Summary		Get client by ID
//	@Description	Retrieve a client by its ID
//	@Tags			Client
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Client ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Client
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/users/client/{id} [get]
func GetClientById(c *fiber.Ctx) error {
	var client model.Client
	if err := GetOneBy("id", c, &client, nil, &[]string{"Appointments"}); err != nil {
		return err
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &client, &DTO.Client{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// UpdateClientById updates a client by ID
//
//	@Summary		Update client
//	@Description	Update a client
//	@Tags			Client
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"Client ID"
//	@Param			client	body		DTO.UpdateClient	true	"Client"
//	@Success		200		{object}	DTO.Client
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/users/client/{id} [patch]
func UpdateClientById(c *fiber.Ctx) error {
	var client model.Client
	if err := UpdateOneById(c, &client, nil); err != nil {
		return err
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &client, &DTO.Client{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// DeleteClientById deletes a client by ID
//
//	@Summary		Delete client by ID
//	@Description	Delete a client by its ID
//	@Tags			Client
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Client ID"
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/users/client/{id} [delete]
func DeleteClientById(c *fiber.Ctx) error {
	return DeleteOneById(c, &model.Client{})
}

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
//	@Success		200			{object}	DTO.EmployeeFull
//	@Failure		400			{object}	DTO.ErrorResponse
//	@Router			/users/employee [post]
func CreateEmployee(c *fiber.Ctx) error {
	var employee model.Employee
	if err := CreateUser(c, &employee); err != nil {
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
//	@Router			/users/employee/{id} [get]
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
//	@Router			/users/employee/email/{email} [get]
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
//	@Router			/users/employee/{id} [patch]
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
//	@Success		200	{object}	nil
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/users/employee/{id} [delete]
func DeleteEmployeeById(c *fiber.Ctx) error {
	return DeleteOneById(c, &model.Employee{})
}

// =====================
// SHARED HELPERS
// =====================

func CreateUser(c *fiber.Ctx, model any) error {
	var err error
	Service := service.New(c)
	defer func() { Service.DeferDB(err) }()
	if err := Service.SetModel(model).Create().Error; err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}
	return nil
}

func GetOneBy(param string, c *fiber.Ctx, model any, nested_preload *[]string, do_not_load *[]string) error {
	var err error
	Service := service.New(c)
	defer func() { Service.DeferDB(err) }()
	if err = Service.
		SetModel(model).
		SetNestedPreload(nested_preload).
		SetDoNotLoad(do_not_load).
		GetBy(param).Error; err != nil {
		return err
	}
	return nil
}

func UpdateOneById(c *fiber.Ctx, model any, nested_preload *[]string) error {
	var err error
	Service := service.New(c)
	defer func() { Service.DeferDB(err) }()
	if err = Service.SetModel(model).SetNestedPreload(nested_preload).UpdateOneById().Error; err != nil {
		return err
	}
	return nil
}

func DeleteOneById(c *fiber.Ctx, model any) error {
	var err error
	Service := service.New(c)
	defer func() { Service.DeferDB(err) }()
	if err = Service.SetModel(model).DeleteOneById().Error; err != nil {
		return err
	}
	return nil
}
