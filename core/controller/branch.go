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
)

// CreateBranch creates a branch
//
//	@Summary		Create branch
//	@Description	Create a branch
//	@Tags			Branch
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Accept			json
//	@Produce		json
//	@Param			branch	body		DTO.CreateBranch	true	"Branch"
//	@Success		201		{object}	DTO.Branch
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/branch [post]
func CreateBranch(c *fiber.Ctx) error {
	var branch model.Branch
	if err := Create(c, &branch); err != nil {
		return err
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &branch, &DTO.Branch{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// GetBranchById retrieves a branch by ID
//
//	@Summary		Get branch by ID
//	@Description	Retrieve a branch by its ID
//	@Tags			Branch
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			id				path		string	true	"Branch ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Branch
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/branch/{id} [get]
func GetBranchById(c *fiber.Ctx) error {
	var branch model.Branch
	if err := GetOneBy("id", c, &branch); err != nil {
		return err
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &branch, &DTO.Branch{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// GetBranchByName retrieves a branch by name
//
//	@Summary		Get branch by name
//	@Description	Retrieve a branch by its name
//	@Tags			Branch
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			name			path		string	true	"Branch Name"
//	@Produce		json
//	@Success		200	{object}	DTO.Branch
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/branch/name/{name} [get]
func GetBranchByName(c *fiber.Ctx) error {
	var branch model.Branch
	if err := GetOneBy("name", c, &branch); err != nil {
		return err
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &branch, &DTO.Branch{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// UpdateBranch updates a branch by ID
//
//	@Summary		Update branch
//	@Description	Update a branch
//	@Tags			Branch
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"Branch ID"
//	@Param			branch	body		DTO.UpdateBranch	true	"Branch"
//	@Success		200		{object}	DTO.Branch
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/branch/{id} [patch]
func UpdateBranchById(c *fiber.Ctx) error {
	var branch model.Branch

	if err := UpdateOneById(c, &branch); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &branch, &DTO.Branch{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// DeleteBranchById deletes a branch by ID
//
//	@Summary		Delete branch by ID
//	@Description	Delete a branch by its ID
//	@Tags			Branch
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Branch
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/branch/{id} [delete]
func DeleteBranchById(c *fiber.Ctx) error {
	return DeleteOneById(c, &model.Branch{})
}

// GetEmployeeServicesByBranchId retrieves all services of an employee included in the branch ID
//
//	@Summary		Get employee services included in the branch ID
//	@Description	Retrieve all services of an employee included in the branch ID
//	@Tags			Branch
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			branch_id		path		string	true	"Branch ID"
//	@Param			employee_id		path		string	true	"Employee ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Service
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/branch/{branch_id}/employee/{employee_id}/services [get]
func GetEmployeeServicesByBranchId(c *fiber.Ctx) error {
	var employee model.Employee
	branchID := c.Params("branch_id")
	employeeID := c.Params("employee_id")

	tx, err := lib.Session(c)
	if err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Verifica se o employee existe
	if err := tx.
		Preload("Services", "branch_id = ?", branchID).
		Where("id = ?", employeeID).
		First(&employee).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.Employee.NotFound
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	res := &lib.SendResponseStruct{Ctx: c}
	if err := res.SendDTO(200, &employee.Services, &[]DTO.Service{}); err != nil {
		return err
	}
	return nil
}

// AddServiceToBranch adds a service to a branch
//
//	@Summary		Add service to branch
//	@Description	Add a service to a branch
//	@Tags			Branch
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			branch_id		path		string	true	"Branch ID"
//	@Param			service_id		path		string	true	"Service ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Branch
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/branch/{branch_id}/service/{service_id} [post]
func AddServiceToBranch(c *fiber.Ctx) error {
	var branch model.Branch
	var service model.Service
	branch_id := c.Params("branch_id")
	service_id := c.Params("service_id")
	if branch_id == "" {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("missing branch_id in the url"))
	} else if service_id == "" {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("missing service_id in the url"))
	}
	tx, end, err := database.ContextTransaction(c)
	defer end()
	if err != nil {
		return err
	}
	if err := tx.Where("id = ?", service_id).First(&service).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.UpdatedError.WithError(fmt.Errorf("service not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}
	if err := database.LockForUpdate(tx, &branch, "id", branch_id); err != nil {
		return err
	}
	if err := branch.AddService(tx, &service); err != nil {
		return err
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &branch, &DTO.Branch{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// RemoveServiceFromBranch removes a service from a branch
//
//	@Summary		Remove service from branch
//	@Description	Remove a service from a branch
//	@Tags			Branch
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			branch_id		path		string	true	"Branch ID"
//	@Param			service_id		path		string	true	"Service ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Branch
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/branch/{branch_id}/service/{service_id} [delete]
func RemoveServiceFromBranch(c *fiber.Ctx) error {
	var branch model.Branch
	var service model.Service
	branch_id := c.Params("branch_id")
	service_id := c.Params("service_id")
	if branch_id == "" {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("missing branch_id in the url"))
	} else if service_id == "" {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("missing service_id in the url"))
	}
	tx, end, err := database.ContextTransaction(c)
	defer end()
	if err != nil {
		return err
	}
	if err := tx.Where("id = ?", service_id).First(&service).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.UpdatedError.WithError(fmt.Errorf("service not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}
	if err := database.LockForUpdate(tx, &branch, "id", branch_id); err != nil {
		return err
	}
	if err := branch.RemoveService(tx, &service); err != nil {
		return err
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &branch, &DTO.Branch{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// CreateBranch creates a branch
func Branch(Gorm *handler.Gorm) {
	endpoint := &middleware.Endpoint{DB: Gorm}
	endpoint.BulkRegisterHandler([]fiber.Handler{
		CreateBranch,
		GetBranchById,
		GetBranchByName,
		UpdateBranchById,
		DeleteBranchById,
		GetEmployeeServicesByBranchId,
		AddServiceToBranch,
		RemoveServiceFromBranch,
	})
}
