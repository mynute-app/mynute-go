package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	dJSON "agenda-kaki-go/core/config/api/dto/json"
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
//	@Success		200		{object}	DTO.BranchFull
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
//	@Success		200	{object}	DTO.BranchFull
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
//	@Success		200	{object}	DTO.BranchFull
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
//	@Success		200		{object}	DTO.BranchFull
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
//	@Success		200	{object}	nil
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/branch/{id} [delete]
func DeleteBranchById(c *fiber.Ctx) error {
	return DeleteOneById(c, &model.Branch{})
}

// UpdateBranchImages updates a branch's images
//
//	@Summary		Update branch images
//	@Description	Update a branch's images
//	@Tags			Branch
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			id				path		string	true	"Branch ID"
//	@Accept			json
//	@Produce		json
//	@Param			profile	formData	file	false	"Profile image"
//	@Success		200		{object}	dJSON.Images
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/branch/{id}/design/images [patch]
func UpdateBranchImages(c *fiber.Ctx) error {
	img_types_allowed := map[string]bool{"profile": true}

	var branch model.Branch
	Design, err := UpdateImagesById(c, branch.TableName(), &branch, img_types_allowed)
	if err != nil {
		return err
	}

	return lib.ResponseFactory(c).SendDTO(200, &Design.Images, &dJSON.Images{})
}

// DeleteBranchImage deletes a branch's image
//
//	@Summary		Delete branch image
//	@Description	Delete a branch's image
//	@Tags			Branch
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			id				path		string	true	"Branch ID"
//	@Param			image_type		path		string	true	"Image Type"
//	@Success		200				{object}	dJSON.Images
//	@Failure		400				{object}	DTO.ErrorResponse
//	@Router			/branch/{id}/design/images/{image_type} [delete]
func DeleteBranchImage(c *fiber.Ctx) error {
	img_types_allowed := map[string]bool{"profile": true}
	var branch model.Branch
	Design, err := DeleteImageById(c, branch.TableName(), &branch, img_types_allowed)
	if err != nil {
		return err
	}

	return lib.ResponseFactory(c).SendDTO(200, &Design.Images, &dJSON.Images{})
}

// CreateBranchWorkSchedule creates a work schedule for a branch
//
//	@Summary		Create work schedule for a branch
//	@Description	Create a work schedule for a branch
//	@Tags			BranchWorkSchedule
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string							true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string							true	"X-Company-ID"
//	@Param			id				path		string							true	"Branch ID"
//	@Param			work_schedule	body		DTO.CreateBranchWorkSchedule	true	"Branch Work Schedule"
//	@Success		200				{object}	DTO.BranchFull
//	@Failure		400				{object}	DTO.ErrorResponse
//	@Router			/branch/{id}/work_schedule [post]
func CreateBranchWorkSchedule(c *fiber.Ctx) error {
	var err error

	var input DTO.CreateBranchWorkSchedule
	if err := c.BodyParser(&input); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	var schedule []model.BranchWorkRange
	branch_id := c.Params("id")

	for i, bwr := range input.WorkRanges {
		start, err := lib.Parse_HHMM_To_Time(bwr.StartTime, bwr.TimeZone)
		if err != nil {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid start_time at index %d: %w", i, err))
		}
		end, err := lib.Parse_HHMM_To_Time(bwr.EndTime, bwr.TimeZone)
		if err != nil {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid end_time at index %d: %w", i, err))
		}

		services := make([]*model.Service, 0, len(bwr.Services))
		for _, s := range bwr.Services {
			services = append(services, &model.Service{BaseModel: model.BaseModel{ID: s.ID}})
		}

		schedule = append(schedule, model.BranchWorkRange{
			WorkRangeBase: model.WorkRangeBase{
				Weekday:   time.Weekday(bwr.Weekday),
				StartTime: start,
				EndTime:   end,
				BranchID:  uuid.MustParse(branch_id),
			},
			Services: services,
		})

		if bwr.Weekday > 100 {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid weekday %d at index %d", bwr.Weekday, i))
		}
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	for _, wr := range schedule {
		if err := tx.Create(&wr).Error; err != nil {
			return lib.Error.General.CreatedError.WithError(err)
		}
	}

	var bwr []model.BranchWorkRange
	if err := tx.
		Preload(clause.Associations).
		Find(&bwr, "branch_id = ?", branch_id).Error; err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}

	bws := model.BranchWorkSchedule{
		WorkRanges: bwr,
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &bws, &DTO.BranchWorkSchedule{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
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
//	@Failure		400	{object}	DTO.ErrorResponse
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
//	@Param			X-Auth-Token	header		string				true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string				true	"X-Company-ID"
//	@Param			id				path		string				true	"Branch ID"
//	@Param			work_range_id	path		string				true	"Work Range ID"
//	@Param			work_range		body		DTO.UpdateWorkRange	true	"Work Range"
//	@Success		200				{object}	DTO.BranchFull
//	@Failure		400				{object}	DTO.ErrorResponse
//	@Router			/branch/{id}/work_range/{work_range_id} [put]
func UpdateBranchWorkRange(c *fiber.Ctx) error {
	branch_id := c.Params("id")
	workRangeID := c.Params("work_range_id")

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var work_range model.BranchWorkRange
	if err := tx.First(&work_range, "id = ?", workRangeID).Error; err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range not found"))
	}

	if work_range.BranchID.String() != branch_id {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range branch ID does not match with branch ID from path"))
	}

	var input DTO.UpdateWorkRange
	if err := c.BodyParser(&input); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	if work_range.TimeZone != input.TimeZone {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range time zone (%s) does not match with input time zone (%s)", work_range.TimeZone, input.TimeZone))
	}

	start, err := lib.Parse_HHMM_To_Time(input.StartTime, input.TimeZone)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid start time: %w", err))
	}
	end, err := lib.Parse_HHMM_To_Time(input.EndTime, input.TimeZone)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid end time: %w", err))
	}

	work_range.Weekday = time.Weekday(input.Weekday)
	work_range.StartTime = start
	work_range.EndTime = end
	work_range.TimeZone = input.TimeZone

	if err := tx.Save(&work_range).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range (%s) not found", work_range.ID))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	var bwr []*model.BranchWorkRange
	if err := tx.
		Preload(clause.Associations).
		Find(&bwr, "branch_id = ?", branch_id).Error; err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}

	var dto []*DTO.BranchWorkRange

	if err := lib.ResponseFactory(c).SendDTO(200, &bwr, &dto); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// DeleteBranchWorkRange deletes a work range for a branch
//
//	@Summary		Delete branch work range
//	@Description	Delete a branch's work range
//	@Tags			BranchWorkSchedule
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			id				path		string	true	"Branch ID"
//	@Param			work_range_id	path		string	true	"Work Range ID"
//	@Success		200				{object}	DTO.BranchWorkRange
//	@Failure		400				{object}	DTO.ErrorResponse
//	@Router			/branch/{id}/work_range/{work_range_id} [delete]
func DeleteBranchWorkRange(c *fiber.Ctx) error {
	var err error
	branch_id := c.Params("id")
	workRangeID := c.Params("work_range_id")

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var wr model.BranchWorkRange
	if err := tx.First(&wr, "id = ?", workRangeID).Error; err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range not found"))
	}

	if wr.BranchID.String() != branch_id {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch ID mismatch"))
	}

	if err := tx.Delete(&wr).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	var bwr []*model.BranchWorkRange
	if err := tx.
		Preload(clause.Associations).
		Find(&bwr, "branch_id = ?", branch_id).Error; err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}

	var dto []*DTO.BranchWorkRange

	if err := lib.ResponseFactory(c).SendDTO(200, &bwr, &dto); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// AddBranchWorkRangeServices adds services to a branch's work range
//
//	@Summary		Add services to branch work range
//	@Description	Add services to a branch's work range
//	@Tags			BranchWorkSchedule
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string			true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string			true	"X-Company-ID"
//	@Param			id				path		string			true	"Branch ID"
//	@Param			work_range_id	path		string			true	"Work Range ID"
//	@Param			services		body		[]DTO.ServiceID	true	"Services"
//	@Success		200				{object}	DTO.BranchFull
//	@Failure		400				{object}	DTO.ErrorResponse
//	@Router			/branch/{id}/work_range/{work_range_id}/services [post]
func AddBranchWorkRangeServices(c *fiber.Ctx) error {
	var err error
	branch_id := c.Params("id")
	workRangeID := c.Params("work_range_id")

	var serviceIDs []DTO.ServiceID
	if err := c.BodyParser(&serviceIDs); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var wr model.BranchWorkRange
	if err := tx.First(&wr, "id = ?", workRangeID).Error; err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range not found"))
	}

	if wr.BranchID.String() != branch_id {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch ID mismatch"))
	}

	services := make([]*model.Service, 0, len(serviceIDs))
	for _, s := range serviceIDs {
		services = append(services, &model.Service{BaseModel: model.BaseModel{ID: s.ID}})
	}

	if err := tx.Model(&wr).Association("Services").Append(services); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	var bwr []*model.BranchWorkRange
	if err := tx.
		Preload(clause.Associations).
		Find(&bwr, "branch_id = ?", branch_id).Error; err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}

	var dto []*DTO.BranchWorkRange

	if err := lib.ResponseFactory(c).SendDTO(200, &bwr, &dto); err != nil {
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
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			id				path		string	true	"Branch ID"
//	@Param			work_range_id	path		string	true	"Work Range ID"
//	@Param			service_id		path		string	true	"Service ID"
//	@Success		200				{object}	DTO.BranchFull
//	@Failure		400				{object}	DTO.ErrorResponse
//	@Router			/branch/{id}/work_range/{work_range_id}/service/{service_id} [delete]
func DeleteBranchWorkRangeService(c *fiber.Ctx) error {
	var err error
	branch_id := c.Params("id")
	workRangeID := c.Params("work_range_id")
	serviceID := c.Params("service_id")

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var wr model.BranchWorkRange
	if err := tx.First(&wr, "id = ?", workRangeID).Error; err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range not found"))
	}

	if wr.BranchID.String() != branch_id {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch ID mismatch"))
	}

	serviceUUID := uuid.MustParse(serviceID)
	if err := tx.Model(&wr).Association("Services").Delete(&model.Service{BaseModel: model.BaseModel{ID: serviceUUID}}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	var bwr []*model.BranchWorkRange
	if err := tx.
		Preload(clause.Associations).
		Find(&bwr, "branch_id = ?", branch_id).Error; err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}

	var dto []*DTO.BranchWorkRange

	if err := lib.ResponseFactory(c).SendDTO(200, &bwr, &dto); err != nil {
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
	branch_id := c.Params("branch_id")
	employee_id := c.Params("employee_id")

	tx, err := lib.Session(c)
	if err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	var employee model.Employee

	if err := tx.
		Preload("Services", "branch_id = ?", branch_id).
		First(&employee, "id = ?", employee_id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.Employee.NotFound
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	res := &lib.SendResponseStruct{Ctx: c}
	if err := res.SendDTO(200, &employee.Services, &[]*DTO.Service{}); err != nil {
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
//	@Success		200	{object}	DTO.BranchFull
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/branch/{branch_id}/service/{service_id} [post]
func AddServiceToBranch(c *fiber.Ctx) error {
	var err error
	var branch model.Branch
	var service model.Service
	branch_id := c.Params("branch_id")
	service_id := c.Params("service_id")
	if branch_id == "" {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("missing branch_id in the url"))
	} else if service_id == "" {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("missing service_id in the url"))
	}
	tx, err := lib.Session(c)
	if err != nil {
		return err
	}
	if err := tx.First(&service, "id = ?", service_id).Error; err != nil {
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
//	@Success		200	{object}	DTO.BranchFull
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/branch/{branch_id}/service/{service_id} [delete]
func RemoveServiceFromBranch(c *fiber.Ctx) error {
	var err error
	var branch model.Branch
	var service model.Service
	branch_id := c.Params("branch_id")
	service_id := c.Params("service_id")
	if branch_id == "" {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("missing branch_id in the url"))
	} else if service_id == "" {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("missing service_id in the url"))
	}
	tx, err := lib.Session(c)
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
		UpdateBranchImages,
		DeleteBranchImage,
		CreateBranchWorkSchedule,
		GetBranchWorkRange,
		UpdateBranchWorkRange,
		DeleteBranchWorkRange,
		AddBranchWorkRangeServices,
		DeleteBranchWorkRangeService,
	})
}
