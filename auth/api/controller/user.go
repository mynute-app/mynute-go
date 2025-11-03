package controller

import (
	"fmt"
	DTO "mynute-go/auth/dto"
	"mynute-go/auth/lib"
	authModel "mynute-go/auth/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	var client authModel.Client
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
	var client authModel.Client
	if err := GetOneBy("email", c, &client); err != nil {
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
	var client authModel.Client
	if err := GetOneBy("id", c, &client); err != nil {
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
	var client authModel.Client
	if err := UpdateOneById(c, &client); err != nil {
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
	return DeleteOneById(c, &authModel.Client{})
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
	var employee authModel.Employee
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
	var employee authModel.Employee
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
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Failure		401				{object}	nil
//	@Param			email			path		string	true	"Employee Email"
//	@Produce		json
//	@Success		200	{object}	DTO.EmployeeFull
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/users/employee/email/{email} [get]
func GetEmployeeByEmail(c *fiber.Ctx) error {
	var employee authModel.Employee
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
	var employee authModel.Employee
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
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Employee ID"
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/users/employee/{id} [delete]
func DeleteEmployeeById(c *fiber.Ctx) error {
	return DeleteOneById(c, &authModel.Employee{})
}

// =====================
// SHARED HELPERS
// =====================

func CreateUser(c *fiber.Ctx, modelInstance any) error {
	// Get database session
	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	// Parse the request body into the model
	if err := c.BodyParser(modelInstance); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	// Create the user record
	if err := tx.Create(modelInstance).Error; err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}

	return nil
}

func GetOneBy(param string, c *fiber.Ctx, modelInstance any) error {
	// Get database session
	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	// Get the parameter value from context
	paramValue := c.Params(param)
	if paramValue == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("missing parameter: %s", param))
	}

	// Validate UUID if param is "id"
	if param == "id" {
		if _, err := uuid.Parse(paramValue); err != nil {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid UUID format"))
		}
	}

	// Query the database
	query := tx.Model(modelInstance).Where(param+" = ?", paramValue)

	// Execute the query
	if err := query.First(modelInstance).Error; err != nil {
		if err.Error() == "record not found" {
			return lib.Error.General.RecordNotFound
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

func UpdateOneById(c *fiber.Ctx, modelInstance any) error {
	// Get database session
	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	// Get ID from params
	id := c.Params("id")
	if id == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("missing id parameter"))
	}

	// Validate UUID
	if _, err := uuid.Parse(id); err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid UUID format"))
	}

	// Parse updates from body
	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	// Update the record
	if err := tx.Model(modelInstance).Where("id = ?", id).Updates(updates).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Fetch the updated record
	if err := tx.Where("id = ?", id).First(modelInstance).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

func DeleteOneById(c *fiber.Ctx, modelInstance any) error {
	// Get database session
	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	// Get ID from params
	id := c.Params("id")
	if id == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("missing id parameter"))
	}

	// Validate UUID
	if _, err := uuid.Parse(id); err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid UUID format"))
	}

	// Delete the record
	if err := tx.Where("id = ?", id).Delete(modelInstance).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Return success
	return lib.ResponseFactory(c).Send(200, map[string]string{"message": "Deleted successfully"})
}
