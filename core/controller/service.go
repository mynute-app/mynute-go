package controller

import (
	"errors"
	"fmt"
	DTO "mynute-go/core/config/api/dto"
	dJSON "mynute-go/core/config/api/dto/json"
	"mynute-go/core/config/db/model"
	"mynute-go/core/handler"
	"mynute-go/core/lib"
	"mynute-go/core/middleware"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateService creates a service
//
//	@Summary		Create service
//	@Description	Create a service
//	@Tags			Service
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Accept			json
//	@Produce		json
//	@Param			service	body		DTO.CreateService	true	"Service"
//	@Success		200		{object}	DTO.Service
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/service [post]
func CreateService(c *fiber.Ctx) error {
	var service model.Service
	if err := Create(c, &service); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &service, &DTO.Service{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// GetServiceById retrieves a service by ID
//
//	@Summary		Get service by ID
//	@Description	Retrieve a service by its ID
//	@Tags			Service
//	@Param			id	path	string	true	"Service ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Service
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/service/{id} [get]
func GetServiceById(c *fiber.Ctx) error {
	var service model.Service

	if err := GetOneBy("id", c, &service, nil); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &service, &DTO.Service{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// GetServiceByName retrieves a service by name
//
//	@Summary		Get service by name
//	@Description	Retrieve a service by its name
//	@Tags			Service
//	@Param			name	path	string	true	"Service Name"
//	@Produce		json
//	@Success		200	{object}	DTO.Service
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/service/name/{name} [get]
func GetServiceByName(c *fiber.Ctx) error {
	var service model.Service

	if err := GetOneBy("name", c, &service, nil); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &service, &DTO.Service{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// UpdateServiceById updates a service by ID
//
//	@Summary		Update service by ID
//	@Description	Update a service by its ID
//	@Tags			Service
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			id				path		string	true	"Service ID"
//	@Accept			json
//	@Produce		json
//	@Param			service	body		DTO.Service	true	"Service"
//	@Success		200		{object}	DTO.Service
//	@Failure		404		{object}	nil
//	@Router			/service/{id} [patch]
func UpdateServiceById(c *fiber.Ctx) error {
	var service model.Service

	if err := UpdateOneById(c, &service, nil); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &service, &DTO.Service{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// DeleteServiceById deletes a service by ID
//
//	@Summary		Delete service by ID
//	@Description	Delete a service by its ID
//	@Tags			Service
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			id				path		string	true	"Service ID"
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		404	{object}	nil
//	@Router			/service/{id} [delete]
func DeleteServiceById(c *fiber.Ctx) error {
	return DeleteOneById(c, &model.Service{})
}

// UpdateServiceImages updates images of a service
//
//	@Summary		Update service images
//	@Description	Update images of a service
//	@Tags			Service
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			id				path		string	true	"Service ID"
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			images	formData	dJSON.Images	true	"Images"
//	@Success		200		{object}	dJSON.Images
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/service/{id}/design/images [patch]
func UpdateServiceImages(c *fiber.Ctx) error {
	img_types_allowed := map[string]bool{"profile": true}

	var service model.Service
	Design, err := UpdateImagesById(c, service.TableName(), &service, img_types_allowed)
	if err != nil {
		return err
	}

	return lib.ResponseFactory(c).SendDTO(200, &Design.Images, &dJSON.Images{})
}

// DeleteServiceImage deletes images of a service
//
//	@Summary		Delete service images
//	@Description	Delete images of a service
//	@Tags			Service
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			id				path		string	true	"Service ID"
//	@Produce		json
//	@Success		200	{object}	dJSON.Images
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/service/{id}/design/images/{image_type} [delete]
func DeleteServiceImage(c *fiber.Ctx) error {
	img_types_allowed := map[string]bool{"profile": true}

	var service model.Service
	Design, err := DeleteImageById(c, service.TableName(), &service, img_types_allowed)
	if err != nil {
		return err
	}

	return lib.ResponseFactory(c).SendDTO(200, &Design.Images, &dJSON.Images{})
}

// GetServiceAvailability retrieves the availability of a service
//
//	@Summary		Get service availability
//	@Description	Retrieve the availability of a service for the next 30 days
//	@Tags			Service
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			service_id	query		string	true	"Service ID"
//	@Param			timezone		query		string	true	"Client Time Zone (IANA format, e.g., America/New_York)"
//	@Param      date_forward_start query number true "The start date for the forward search in number format"
//	@Param      date_forward_end query number true "The end date for the forward search in number format"
//	@Produce		json
//	@Success		200	{object}	DTO.ServiceAvailability
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/service/{id}/availability [get]
func GetServiceAvailability(c *fiber.Ctx) error {
	serviceIDStr := c.Query("id")
	if serviceIDStr == "" {
		return lib.Error.General.BadRequest.WithError(errors.New("id query parameter is required"))
	}
	serviceID, err := uuid.Parse(serviceIDStr)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(errors.New("invalid id"))
	}
	date_forward_start := c.Query("date_forward_start")
	if date_forward_start == "" {
		return lib.Error.General.BadRequest.WithError(errors.New("date_forward_start is required"))
	}
	date_forward_end := c.Query("date_forward_end")
	if date_forward_end == "" {
		return lib.Error.General.BadRequest.WithError(errors.New("date_forward_end is required"))
	}
	dfs, err := strconv.Atoi(date_forward_start)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(errors.New("invalid date_forward_start"))
	}
	dfe, err := strconv.Atoi(date_forward_end)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(errors.New("invalid date_forward_end"))
	}
	if dfe <= dfs {
		return lib.Error.General.BadRequest.WithError(errors.New("date_forward_end must be greater than date_forward_start"))
	}
	if dfe-dfs > 30 {
		return lib.Error.General.BadRequest.WithError(errors.New("date search range (date_forward_end - date_forward_start) must not exceed 30 days"))
	}

	companyIDStr := c.Get("X-Company-ID")
	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(errors.New("invalid X-Company-ID"))
	}

	timezone := c.Query("timezone")
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid timezone: %s", timezone))
	}

	tx, err := lib.Session(c)
	if err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Step 2: Load all BranchWorkRanges with the given service
	var branchRanges []model.BranchWorkRange
	err = tx.
		Joins("JOIN branch_work_range_services bs ON bs.branch_work_range_id = branch_work_ranges.id").
		Where("bs.service_id = ?", serviceID).
		Preload("Branch").
		Find(&branchRanges).Error
	if err != nil {
		return err
	}

	// Step 3: Load all EmployeeWorkRanges with the given service
	var empRanges []model.EmployeeWorkRange
	err = tx.
		Joins("JOIN employee_work_range_services es ON es.employee_work_range_id = employee_work_ranges.id").
		Where("es.service_id = ?", serviceID).
		Preload("Employee").
		Find(&empRanges).Error
	if err != nil {
		return err
	}

	// Step 4: Load employee densities (map[employeeID]density)
	var densities []model.EmployeeServiceDensity
	err = tx.Where("service_id = ?", serviceID).Find(&densities).Error
	if err != nil {
		return err
	}
	densityMap := make(map[uuid.UUID]uint32)
	for _, d := range densities {
		densityMap[d.EmployeeID] = d.Density
	}

	// Step 5: Prepare maps for reuse
	employeeInfoMap := map[uuid.UUID]DTO.EmployeeBase{}
	branchInfoMap := map[uuid.UUID]DTO.BranchBase{}
	availabilityMap := map[string]map[uuid.UUID]map[string][]uuid.UUID{} // date → branch → time → []employeeID

	startDate := time.Now().In(loc).AddDate(0, 0, dfs)
	endDate := time.Now().In(loc).AddDate(0, 0, dfe)

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		weekday := d.Weekday()

		for _, br := range branchRanges {
			if br.Weekday != weekday {
				continue
			}
			branchID := br.BranchID
			branchInfoMap[branchID] = DTO.BranchBase{
				ID:                  br.Branch.ID,
				CompanyID:           br.Branch.CompanyID,
				Name:                br.Branch.Name,
				Street:              br.Branch.Street,
				Number:              br.Branch.Number,
				Complement:          br.Branch.Complement,
				Neighborhood:        br.Branch.Neighborhood,
				ZipCode:             br.Branch.ZipCode,
				City:                br.Branch.City,
				State:               br.Branch.State,
				Country:             br.Branch.Country,
				TimeZone:            br.Branch.TimeZone,
				TotalServiceDensity: br.Branch.TotalServiceDensity,
				Design:              br.Branch.Design,
			}

			for _, emp := range empRanges {
				if emp.Weekday != weekday || emp.BranchID != branchID {
					continue
				}

				empID := emp.Employee.ID
				if _, ok := employeeInfoMap[empID]; !ok {
					employeeInfoMap[empID] = DTO.EmployeeBase{
						ID:                  empID,
						CompanyID:           emp.Employee.CompanyID,
						Name:                emp.Employee.Name,
						Surname:             emp.Employee.Surname,
						TimeZone:            emp.Employee.TimeZone,
						TotalServiceDensity: emp.Employee.TotalServiceDensity,
						Design:              emp.Employee.Design,
					}
				}

				slot := emp.StartTime.In(loc)
				end := emp.EndTime.In(loc)

				for slot.Before(end) {
					slotEnd := slot.Add(time.Minute * time.Duration(emp.Employee.SlotTimeDiff))

					// appointment overlap check
					var count int64
					err := tx.Model(&model.Appointment{}).
						Where("employee_id = ? AND service_id = ? AND company_id = ? AND is_cancelled = false", empID, serviceID, companyID).
						Where("start_time <= ? AND end_time > ?", slot, slot).
						Count(&count).Error
					if err != nil {
						return err
					}

					max := emp.Employee.TotalServiceDensity
					if val, ok := densityMap[empID]; ok && val > 0 {
						max = val
					}

					if uint32(count) < max {
						dateStr := d.Format("2006-01-02")
						timeStr := slot.Format("15:04")

						if _, ok := availabilityMap[dateStr]; !ok {
							availabilityMap[dateStr] = map[uuid.UUID]map[string][]uuid.UUID{}
						}
						if _, ok := availabilityMap[dateStr][branchID]; !ok {
							availabilityMap[dateStr][branchID] = map[string][]uuid.UUID{}
						}
						availabilityMap[dateStr][branchID][timeStr] = append(availabilityMap[dateStr][branchID][timeStr], empID)
					}

					slot = slotEnd
				}
			}
		}
	}

	// Step 6: Build final response
	var availableDays []DTO.AvailableDay
	for date, branches := range availabilityMap {
		for branchID, slots := range branches {
			var timeSlots []DTO.AvailableTimeSlot
			for timeStr, empIDs := range slots {
				timeSlots = append(timeSlots, DTO.AvailableTimeSlot{
					Time:        timeStr,
					EmployeesID: empIDs,
				})
			}
			availableDays = append(availableDays, DTO.AvailableDay{
				Date:      date,
				BranchID:  branchID,
				TimeSlots: timeSlots,
			})
		}
	}

	// Convert maps to slices
	var employeeInfo []DTO.EmployeeBase
	for _, e := range employeeInfoMap {
		employeeInfo = append(employeeInfo, e)
	}
	var branchInfo []DTO.BranchBase
	for _, b := range branchInfoMap {
		branchInfo = append(branchInfo, b)
	}

	return lib.ResponseFactory(c).Send(200, &DTO.ServiceAvailability{
		ServiceID:     serviceID,
		AvailableDays: availableDays,
		EmployeeInfo:  employeeInfo,
		BranchInfo:    branchInfo,
	})
}

// Service returns a service_controller
func Service(Gorm *handler.Gorm) {
	endpoint := &middleware.Endpoint{DB: Gorm}
	endpoint.BulkRegisterHandler([]fiber.Handler{
		CreateService,
		GetServiceById,
		GetServiceByName,
		UpdateServiceById,
		DeleteServiceById,
		UpdateServiceImages,
		DeleteServiceImage,
		GetServiceAvailability,
	})
}
