package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/service"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

type branch_controller struct {
	service.Base[model.Branch, DTO.Branch]
}

// CreateBranch creates a branch
//
//	@Summary		Create branch
//	@Description	Create a branch
//	@Tags			Branch
//	@Accept			json
//	@Produce		json
//	@Param			branch	body		DTO.CreateBranch	true	"Branch"
//	@Success		201		{object}	DTO.Branch
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/branch [post]
func (cc *branch_controller) CreateBranch(c *fiber.Ctx) error {
	return cc.CreateOne(c)
}

// GetBranchById retrieves a branch by ID
//
//	@Summary		Get branch by ID
//	@Description	Retrieve a branch by its ID
//	@Tags			Branch
//	@Param			id	path	string	true	"Branch ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Branch
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/branch/{id} [get]
func (cc *branch_controller) GetBranchById(c *fiber.Ctx) error {
	return cc.GetBy("id", c)
}

// GetBranchByName retrieves a branch by name
//
//	@Summary		Get branch by name
//	@Description	Retrieve a branch by its name
//	@Tags			Branch
//	@Param			name	path	string	true	"Branch Name"
//	@Produce		json
//	@Success		200	{object}	DTO.Branch
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/branch/name/{name} [get]
func (cc *branch_controller) GetBranchByName(c *fiber.Ctx) error {
	return cc.GetBy("name", c)
}

// UpdateBranch updates a branch by ID
//
//	@Summary		Update branch
//	@Description	Update a branch
//	@Tags			Branch
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Authorization"
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"Branch ID"
//	@Param			branch	body		DTO.UpdateBranch	true	"Branch"
//	@Success		200		{object}	DTO.Branch
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/branch/{id} [patch]
func (cc *branch_controller) UpdateBranchById(c *fiber.Ctx) error {
	return cc.UpdateOneById(c)
}

// DeleteBranchById deletes a branch by ID
//
//	@Summary		Delete branch by ID
//	@Description	Delete a branch by its ID
//	@Tags			Branch
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Authorization"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Branch ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Branch
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/branch/{id} [delete]
func (cc *branch_controller) DeleteBranchById(c *fiber.Ctx) error {
	return cc.DeleteOneById(c)
}

// GetEmployeeServicesByBranchId retrieves all services of an employee included in the branch ID
//
//	@Summary		Get employee services included in the branch ID
//	@Description	Retrieve all services of an employee included in the branch ID
//	@Tags			Branch
//	@Param			branch_id	path	string	true	"Branch ID"
//	@Param			employee_id	path	string	true	"Employee ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Service
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/branch/{branch_id}/employee/{employee_id}/services [get]
func (cc *branch_controller) GetEmployeeServicesByBranchId(c *fiber.Ctx) error {
	var branch model.Branch
	var employee model.Employee
	branch_id := c.Params("branch_id")
	employee_id := c.Params("employee_id")
	if err := cc.Request.Gorm.GetOneBy("id", branch_id, &branch, cc.Associations); err != nil {
		return err
	}
	if err := cc.Request.Gorm.GetOneBy("id", employee_id, &employee, []string{}); err != nil {
		return err
	}
	if employee.CompanyID != branch.CompanyID {
		return lib.Error.Company.NotSame.SendToClient(c)
	}
	if err := cc.Request.Gorm.DB.Model(&branch).Association("Services").Find(&employee.Services); err != nil {
		return err
	}
	res := &lib.SendResponse{Ctx: c}
	res.SendDTO(200, &employee.Services, &DTO.Service{})
	return nil
}

// AddEmployeeToBranch adds an employee to a branch
//
//	@Summary		Add employee to branch
//	@Description	Add an employee to a branch
//	@Tags			Branch
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Authorization
//	@Failure		401				{object}	nil
//	@Param			branch_id	path		string	true	"Branch ID"
//	@Param			employee_id	path		string	true	"Employee ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Branch
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/branch/{branch_id}/employee/{employee_id} [post]
func (cc *branch_controller) AddEmployeeToBranch(c *fiber.Ctx) error {
	var branch model.Branch
	var employee model.Employee
	branch_id := c.Params("branch_id")
	employee_id := c.Params("employee_id")
	if err := cc.Request.Gorm.GetOneBy("id", branch_id, &branch, cc.Associations); err != nil {
		return err
	}
	if err := cc.Request.Gorm.GetOneBy("id", employee_id, &employee, []string{}); err != nil {
		return err
	}
	if employee.CompanyID != branch.CompanyID {
		return lib.Error.Company.NotSame.SendToClient(c)
	}
	if err := cc.Request.Gorm.DB.Model(&branch).Association("Employees").Append(&employee); err != nil {
		return err
	}
	res := &lib.SendResponse{Ctx: c}
	res.SendDTO(200, &branch, &DTO.Branch{})
	return nil
}

// RemoveEmployeeFromBranch removes an employee from a branch
//
//	@Summary		Remove employee from branch
//	@Description	Remove an employee from a branch
//	@Tags			Branch
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Authorization
//	@Failure		401				{object}	nil
//	@Param			branch_id	path		string	true	"Branch ID"
//	@Param			employee_id	path		string	true	"Employee ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Branch
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/branch/{branch_id}/employee/{employee_id} [delete]
func (cc *branch_controller) RemoveEmployeeFromBranch(c *fiber.Ctx) error {
	var branch model.Branch
	var employee model.Employee
	branch_id := c.Params("branch_id")
	employee_id := c.Params("employee_id")
	if err := cc.Request.Gorm.GetOneBy("id", branch_id, &branch, cc.Associations); err != nil {
		return err
	}
	if err := cc.Request.Gorm.GetOneBy("id", employee_id, &employee, []string{}); err != nil {
		return err
	}
	if employee.CompanyID != branch.CompanyID {
		return lib.Error.Company.NotSame.SendToClient(c)
	}
	if err := cc.Request.Gorm.DB.Model(&branch).Association("Employees").Delete(&employee); err != nil {
		return err
	}
	res := &lib.SendResponse{Ctx: c}
	res.SendDTO(200, &branch, &DTO.Branch{})
	return nil
}

// AddServiceToBranch adds a service to a branch
//
//	@Summary		Add service to branch
//	@Description	Add a service to a branch
//	@Tags			Branch
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Authorization
//	@Failure		401				{object}	nil
//	@Param			branch_id	path		string	true	"Branch ID"
//	@Param			service_id	path		string	true	"Service ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Branch
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/branch/{branch_id}/service/{service_id} [post]
func (cc *branch_controller) AddServiceToBranch(c *fiber.Ctx) error {
	var branch model.Branch
	var service model.Service
	branch_id := c.Params("branch_id")
	service_id := c.Params("service_id")
	if err := cc.Request.Gorm.GetOneBy("id", branch_id, &branch, cc.Associations); err != nil {
		return err
	}
	if err := cc.Request.Gorm.GetOneBy("id", service_id, &service, []string{}); err != nil {
		return err
	}
	if service.CompanyID != branch.CompanyID {
		return lib.Error.Company.NotSame.SendToClient(c)
	}
	if err := cc.Request.Gorm.DB.Model(&branch).Association("Services").Append(&service); err != nil {
		return err
	}
	res := &lib.SendResponse{Ctx: c}
	// branch_marchal, err := json.Marshal(&branch)
	// if err != nil {
	// 	return err
	// }
	// var DTO DTO.Branch
	// if err := json.Unmarshal(branch_marchal, &DTO); err != nil {
	// 	return err
	// }
	return res.SendDTO(200, &branch, &DTO.Branch{})
	// res.Http200(&DTO)
	// return nil
}

// RemoveServiceFromBranch removes a service from a branch
//
//	@Summary		Remove service from branch
//	@Description	Remove a service from a branch
//	@Tags			Branch
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Authorization
//	@Failure		401				{object}	nil
//	@Param			branch_id	path		string	true	"Branch ID"
//	@Param			service_id	path		string	true	"Service ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Branch
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/branch/{branch_id}/service/{service_id} [delete]
func (cc *branch_controller) RemoveServiceFromBranch(c *fiber.Ctx) error {
	var branch model.Branch
	var service model.Service
	branch_id := c.Params("branch_id")
	service_id := c.Params("service_id")
	if err := cc.Request.Gorm.GetOneBy("id", branch_id, &branch, cc.Associations); err != nil {
		return err
	}
	if err := cc.Request.Gorm.GetOneBy("id", service_id, &service, []string{}); err != nil {
		return err
	}
	if service.CompanyID != branch.CompanyID {
		return lib.Error.Company.NotSame.SendToClient(c)
	}
	if err := cc.Request.Gorm.DB.Model(&branch).Association("Services").Delete(&service); err != nil {
		return err
	}
	res := &lib.SendResponse{Ctx: c}
	branch_marchal, err := json.Marshal(&branch)
	if err != nil {
		return err
	}
	var DTO DTO.Branch
	if err := json.Unmarshal(branch_marchal, &DTO); err != nil {
		return err
	}
	res.Http200(&DTO)
	return nil
}

// CreateBranch creates a branch
func Branch(Gorm *handler.Gorm) *branch_controller {
	return &branch_controller{
		Base: service.Base[model.Branch, DTO.Branch]{
			Name:         namespace.UserKey.Name,
			Request:      handler.Request(Gorm),
			Associations: []string{"Employees", "Services", "Company"},
		},
	}
}
