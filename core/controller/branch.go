package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	database "agenda-kaki-go/core/config/db"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/middleware"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreateBranch creates a branch
//
//	@Summary		Create branch
//	@Description	Create a branch
//	@Tags			Branch
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Accept			json
//	@Produce		json
//	@Param			branch	body		DTO.CreateBranch	true	"Branch"
//	@Success		200		{object}	DTO.Branch
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/branch [post]
func CreateBranch(c *fiber.Ctx) error {
	var branch model.Branch
	if err := Create(c, &branch); err != nil {
		return err
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &branch, &DTO.BranchFull{}); err != nil {
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
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
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
	if err := lib.ResponseFactory(c).SendDTO(200, &branch, &DTO.BranchFull{}); err != nil {
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
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
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
	if err := lib.ResponseFactory(c).SendDTO(200, &branch, &DTO.BranchFull{}); err != nil {
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
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
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

	if err := lib.ResponseFactory(c).SendDTO(200, &branch, &DTO.BranchFull{}); err != nil {
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
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Branch
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/branch/{id} [delete]
func DeleteBranchById(c *fiber.Ctx) error {
	return DeleteOneById(c, &model.Branch{})
}

// CreateBranchWorkSchedule creates a work schedule for a branch
//
//	@Summary		Create work schedule for a branch
//	@Description	Create a work schedule for a branch
//	@Tags			BranchWorkSchedule
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header	string	true	"X-Company-ID"
//	@Param			id				path	string	true	"Branch ID"
//	@Param			work_schedule	body	DTO.CreateBranchWorkSchedule	true	"Branch Work Schedule"
//	@Success		200	{object}	DTO.BranchFull
//	@Failure		400	{object}	lib.ErrorResponse
//	@Router			/branch/{id}/work_schedule [post]
func CreateBranchWorkSchedule(c *fiber.Ctx) error {
	var input DTO.CreateBranchWorkSchedule
	if err := c.BodyParser(&input); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	var schedule []model.BranchWorkRange
	branchID := c.Params("id")

	for i, bwr := range input.WorkRanges {
		loc, err := time.LoadLocation(bwr.TimeZone)
		if err != nil {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid timezone at index %d: %w", i, err))
		}
		start, err := lib.ParseTimeHHMMWithDateBase(bwr.StartTime, loc)
		if err != nil {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid start_time at index %d: %w", i, err))
		}
		end, err := lib.ParseTimeHHMMWithDateBase(bwr.EndTime, loc)
		if err != nil {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid end_time at index %d: %w", i, err))
		}

		services := make([]*model.Service, 0, len(bwr.Services))
		for _, s := range bwr.Services {
			services = append(services, &model.Service{BaseModel: model.BaseModel{ID: s.ID}})
		}

		schedule = append(schedule, model.BranchWorkRange{
			Weekday:   time.Weekday(bwr.Weekday),
			StartTime: start,
			EndTime:   end,
			TimeZone:  bwr.TimeZone,
			BranchID:  uuid.MustParse(branchID),
			Services:  services,
		})
	}

	tx, end, err := database.ContextTransaction(c)
	defer end()
	if err != nil {
		return err
	}

	for _, wr := range schedule {
		if err := tx.Create(&wr).Error; err != nil {
			tx.Rollback()
			return lib.Error.General.CreatedError.WithError(err)
		}
	}

	var branch model.Branch
	if err := tx.Preload(clause.Associations).First(&branch, "id = ?", branchID).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return lib.ResponseFactory(c).SendDTO(200, &branch, &DTO.BranchFull{})
}

// GetBranchWorkRange
//
//	@Summary		Get branch work range By ID
//	@Description	Retrieve a branch's work range by its ID
//	@Tags			BranchWorkSchedule
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header	string	true	"X-Company-ID"
//	@Param			id				path	string	true	"Branch ID"
//	@Param			work_range_id	path	string	true	"Work Range ID"
//	@Produce		json
//	@Success		200	{object}	DTO.BranchWorkRange
//	@Failure		400	{object}	lib.ErrorResponse
func GetBranchWorkRange(c *fiber.Ctx) error {
	branchID := c.Params("id")
	workRangeID := c.Params("work_range_id")

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var wr model.BranchWorkRange
	if err := tx.First(&wr, "id = ?", workRangeID).Error; err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range not found"))
	}

	if wr.BranchID.String() != branchID {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch ID mismatch"))
	}

	return lib.ResponseFactory(c).SendDTO(200, &wr, &DTO.BranchWorkRange{})
}

// UpdateBranchWorkRange updates a work range for a branch
//
//	@Summary		Update branch work range
//	@Description	Update a branch's work range
//	@Tags			BranchWorkSchedule
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header	string	true	"X-Company-ID"
//	@Param			id				path	string	true	"Branch ID"
//	@Param			work_range_id	path	string	true	"Work Range ID"
//	@Param			work_range		body	DTO.UpdateWorkRange	true	"Work Range"
//	@Success		200	{object}	DTO.BranchFull
//	@Failure		400	{object}	lib.ErrorResponse
//	@Router			/branch/{id}/work_range/{work_range_id} [put]
func UpdateBranchWorkRange(c *fiber.Ctx) error {
	branchID := c.Params("id")
	workRangeID := c.Params("work_range_id")

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var wr model.BranchWorkRange
	if err := tx.First(&wr, "id = ?", workRangeID).Error; err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range not found"))
	}

	if wr.BranchID.String() != branchID {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch ID mismatch"))
	}

	var input DTO.UpdateWorkRange
	if err := c.BodyParser(&input); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	loc, err := time.LoadLocation(input.TimeZone)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid timezone: %w", err))
	}
	start, err := lib.ParseTimeHHMMWithDateBase(input.StartTime, loc)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid start time: %w", err))
	}
	end, err := lib.ParseTimeHHMMWithDateBase(input.EndTime, loc)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid end time: %w", err))
	}

	wr.Weekday = time.Weekday(input.Weekday)
	wr.StartTime = start
	wr.EndTime = end
	wr.TimeZone = input.TimeZone

	if err := UpdateOneById(c, &wr); err != nil {
		return err
	}

	var branch model.Branch
	if err := tx.Preload(clause.Associations).First(&branch, "id = ?", branchID).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return lib.ResponseFactory(c).SendDTO(200, &branch, &DTO.BranchFull{})
}

// DeleteBranchWorkRange deletes a work range for a branch
//
//	@Summary		Delete branch work range
//	@Description	Delete a branch's work range
//	@Tags			BranchWorkSchedule
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header	string	true	"X-Company-ID"
//	@Param			id				path	string	true	"Branch ID"
//	@Param			work_range_id	path	string	true	"Work Range ID"
//	@Success		200	{object}	DTO.Branch
//	@Failure		400	{object}	lib.ErrorResponse
//	@Router			/branch/{id}/work_range/{work_range_id} [delete]
func DeleteBranchWorkRange(c *fiber.Ctx) error {
	branchID := c.Params("id")
	workRangeID := c.Params("work_range_id")

	tx, end, err := database.ContextTransaction(c)
	defer end()
	if err != nil {
		return err
	}

	var wr model.BranchWorkRange
	if err := tx.First(&wr, "id = ?", workRangeID).Error; err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range not found"))
	}

	if wr.BranchID.String() != branchID {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch ID mismatch"))
	}

	if err := tx.Delete(&wr).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	var branch model.Branch
	if err := tx.Preload(clause.Associations).First(&branch, "id = ?", branchID).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return lib.ResponseFactory(c).SendDTO(200, &branch, &DTO.BranchFull{})
}

// AddBranchWorkRangeServices adds services to a branch's work range
//
//	@Summary		Add services to branch work range
//	@Description	Add services to a branch's work range
//	@Tags			BranchWorkSchedule
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header	string	true	"X-Company-ID"
//	@Param			id				path	string	true	"Branch ID"
//	@Param			work_range_id	path	string	true	"Work Range ID"
//	@Param			services		body	[]DTO.ServiceID	true	"Services"
//	@Success		200	{object}	DTO.BranchFull
//	@Failure		400	{object}	lib.ErrorResponse
//	@Router			/branch/{id}/work_range/{work_range_id}/services [post]
func AddBranchWorkRangeServices(c *fiber.Ctx) error {
	branchID := c.Params("id")
	workRangeID := c.Params("work_range_id")

	var serviceIDs []DTO.ServiceID
	if err := c.BodyParser(&serviceIDs); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	tx, end, err := database.ContextTransaction(c)
	defer end()
	if err != nil {
		return err
	}

	var wr model.BranchWorkRange
	if err := tx.First(&wr, "id = ?", workRangeID).Error; err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range not found"))
	}

	if wr.BranchID.String() != branchID {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch ID mismatch"))
	}

	services := make([]*model.Service, 0, len(serviceIDs))
	for _, s := range serviceIDs {
		services = append(services, &model.Service{BaseModel: model.BaseModel{ID: s.ID}})
	}

	if err := tx.Model(&wr).Association("Services").Append(services); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// DeleteBranchWorkRangeService removes a service from a branch's work range
//
//	@Summary		Remove service from branch work range
//	@Description	Remove a service from a branch's work range
//	@Tags			BranchWorkSchedule
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header	string	true	"X-Company-ID"
//	@Param			id				path	string	true	"Branch ID"
//	@Param			work_range_id	path	string	true	"Work Range ID"
//	@Param			service_id		path	string	true	"Service ID"
//	@Success		200	{object}	DTO.BranchFull
//	@Failure		400	{object}	lib.ErrorResponse
//	@Router			/branch/{id}/work_range/{work_range_id}/service/{service_id} [delete]
func DeleteBranchWorkRangeService(c *fiber.Ctx) error {
	branchID := c.Params("id")
	workRangeID := c.Params("work_range_id")
	serviceID := c.Params("service_id")

	tx, end, err := database.ContextTransaction(c)
	defer end()
	if err != nil {
		return err
	}

	var wr model.BranchWorkRange
	if err := tx.First(&wr, "id = ?", workRangeID).Error; err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range not found"))
	}

	if wr.BranchID.String() != branchID {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch ID mismatch"))
	}

	serviceUUID := uuid.MustParse(serviceID)
	if err := tx.Model(&wr).Association("Services").Delete(&model.Service{BaseModel: model.BaseModel{ID: serviceUUID}}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// GetEmployeeServicesByBranchId retrieves all services of an employee included in the branch ID
//
//	@Summary		Get employee services included in the branch ID
//	@Description	Retrieve all services of an employee included in the branch ID
//	@Tags			Branch
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
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
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
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
	if err := lib.ResponseFactory(c).SendDTO(200, &branch, &DTO.BranchFull{}); err != nil {
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
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
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
	if err := lib.ResponseFactory(c).SendDTO(200, &branch, &DTO.BranchFull{}); err != nil {
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
		CreateBranchWorkSchedule,
		GetBranchWorkRange,
		UpdateBranchWorkRange,
		DeleteBranchWorkRange,
		AddBranchWorkRangeServices,
		DeleteBranchWorkRangeService,
	})
}
