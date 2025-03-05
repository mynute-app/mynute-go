package controller

import (
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/middleware"
	"agenda-kaki-go/core/service"

	"github.com/gofiber/fiber/v2"
)

type employee_controller struct {
	service.Base[model.Employee, model.Employee]
}

// CreateEmployee creates an employee
//
//	@Summary		Create employee
//	@Description	Create an employee
//	@Tags			Employee
//	@Accept			json
//	@Produce		json
//	@Param			employee	body		DTO.CreateEmployee	true	"Employee"
//	@Success		200			{object}	DTO.Employee
//	@Failure		400			{object}	DTO.ErrorResponse
//	@Router			/employee [post]
func (cc *employee_controller) CreateEmployee(c *fiber.Ctx) error {
	return cc.CreateOne(c)
}

// GetEmployeeById retrieves an employee by ID
//
//	@Summary		Get employee by ID
//	@Description	Retrieve an employee by its ID
//	@Tags			Employee
//	@Param			id	path	string	true	"Employee ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Employee
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/employee/{id} [get]
func (cc *employee_controller) GetEmployeeById(c *fiber.Ctx) error {
	return cc.GetBy("id", c)
}

// UpdateEmployeeById updates an employee by ID
//
//	@Summary		Update employee
//	@Description	Update an employee
//	@Tags			Employee
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string						true	"Employee ID"
//	@Param			employee	body		DTO.UpdateEmployeeSwagger	true	"Employee"
//	@Success		200			{object}	DTO.Employee
//	@Failure		400			{object}	DTO.ErrorResponse
//	@Router			/employee/{id} [patch]
func (cc *employee_controller) UpdateEmployeeById(c *fiber.Ctx) error {
	return cc.UpdateOneById(c)
}

// DeleteEmployeeById deletes an employee by ID
//
//	@Summary		Delete employee by ID
//	@Description	Delete an employee by its ID
//	@Tags			Employee
//	@Param			id	path	string	true	"Employee ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Employee
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/employee/{id} [delete]
func (cc *employee_controller) DeleteEmployeeById(c *fiber.Ctx) error {
	return cc.DeleteOneById(c)
}

func Employee(Gorm *handler.Gorm) *employee_controller {
	return &employee_controller{
		Base: service.Base[model.Employee, model.Employee]{
			Name:         namespace.HolidaysKey.Name,
			Request:      handler.Request(Gorm),
			Middleware:   middleware.Holidays(Gorm),
			Associations: []string{},
		},
	}
}
