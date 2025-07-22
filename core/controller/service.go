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
//	@Param			X-Company-ID		header		string	true	"X-Company-ID"
//	@Param			id			path		string	true	"Service ID"
//	@Param			timezone			query		string	true	"Client Time Zone (IANA format, e.g., America/New_York)"
//	@Param			date_forward_start	query		number	true	"The start date for the forward search in number format"
//	@Param			date_forward_end	query		number	true	"The end date for the forward search in number format"
//	@Produce		json
//	@Success		200	{object}	DTO.ServiceAvailability
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/service/{id}/availability [get]
func GetServiceAvailability(c *fiber.Ctx) error {
	// Step 1: Validate and parse query parameters
	serviceIDStr := c.Params("id")
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

	startDate := time.Now().In(loc).AddDate(0, 0, dfs)
	endDate := time.Now().In(loc).AddDate(0, 0, dfe).Add(24 * time.Hour) // Para incluir o dia final inteiro

	// =========================================================================
	// Step 2: Fetch all necessary data in fewer, more efficient queries
	// =========================================================================

	// --- 2a. Fetch all employee work ranges for the service, with Employee and Branch info
	var empRanges []model.EmployeeWorkRange
	err = tx.
		Joins("JOIN employee_work_range_services es ON es.employee_work_range_id = employee_work_ranges.id").
		Where("es.service_id = ?", serviceID).
		Preload("Employee"). // Preload Employee data to get SlotTimeDiff and density
		Preload("Branch").   // Preload Branch data for the response
		Find(&empRanges).Error
	if err != nil {
		return err
	}

	// --- 2b. Fetch all appointments for the relevant employees and date range in ONE query.
	// This is the single most important optimization.
	type appointmentCountResult struct {
		EmployeeID uuid.UUID
		StartTime  time.Time
		Count      int64
	}
	var appointmentCounts []appointmentCountResult

	// Extrai os IDs dos funcionários para filtrar a query de agendamentos
	// --- Extrai os IDs ÚNICOS dos funcionários para filtrar a query de agendamentos ---
	// Usamos um map como um "set" para garantir que cada ID apareça apenas uma vez.
	employeeIDSet := make(map[uuid.UUID]struct{})
	for _, er := range empRanges {
		employeeIDSet[er.EmployeeID] = struct{}{}
	}

	// Agora, convertemos o set de volta para um slice para usar na query do GORM.
	employeeIDs := make([]uuid.UUID, 0, len(employeeIDSet))
	for id := range employeeIDSet {
		employeeIDs = append(employeeIDs, id)
	}

	if len(employeeIDs) > 0 {
		err = tx.Model(&model.Appointment{}).
			Select("employee_id, start_time, count(*) as count").
			Where("company_id = ? AND service_id = ? AND employee_id IN ? AND is_cancelled = false", companyID, serviceID, employeeIDs).
			Where("start_time >= ? AND start_time < ?", startDate, endDate). // Date range filter
			Group("employee_id, start_time").
			Find(&appointmentCounts).Error
		if err != nil {
			return err
		}
	}

	// --- 2c. Fetch service-specific densities for employees.
	var densities []model.EmployeeServiceDensity
	err = tx.Where("service_id = ? AND employee_id IN ?", serviceID, employeeIDs).Find(&densities).Error
	if err != nil {
		return err
	}

	// =========================================================================
	// Step 3: Index all fetched data into maps for fast lookups (O(1) access)
	// =========================================================================

	// --- 3a. Index appointment counts by employee and start time.
	// Key: "employeeID-2023-10-27T10:30:00Z", Value: count
	appointmentSlotMap := make(map[string]int64)
	for _, ac := range appointmentCounts {
		key := fmt.Sprintf("%s-%s", ac.EmployeeID.String(), ac.StartTime.In(loc).Format(time.RFC3339))
		appointmentSlotMap[key] = ac.Count
	}

	// --- 3b. Index service-specific densities by employee ID.
	densityMap := make(map[uuid.UUID]uint32)
	for _, d := range densities {
		densityMap[d.EmployeeID] = d.Density
	}

	// --- 3c. Index employee work ranges by Weekday and then by BranchID.
	// This avoids nested loops later.
	rangesByDayAndBranch := make(map[time.Weekday]map[uuid.UUID][]model.EmployeeWorkRange)
	for _, r := range empRanges {
		if _, ok := rangesByDayAndBranch[r.Weekday]; !ok {
			rangesByDayAndBranch[r.Weekday] = make(map[uuid.UUID][]model.EmployeeWorkRange)
		}
		rangesByDayAndBranch[r.Weekday][r.BranchID] = append(rangesByDayAndBranch[r.Weekday][r.BranchID], r)
	}

	// =========================================================================
	// Step 4: Process availability using the in-memory maps (NO DB QUERIES HERE)
	// =========================================================================
	employeeInfoMap := map[uuid.UUID]DTO.EmployeeBase{}
	branchInfoMap := map[uuid.UUID]DTO.BranchBase{}
	availabilityMap := map[string]map[uuid.UUID]map[string][]uuid.UUID{} // date → branch → time → []employeeID

	// Reset endDate to not have the extra 24h for the loop
	endDate = endDate.Add(-24 * time.Hour)
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		weekday := d.Weekday()

		// Directly get the branches and employees working on this weekday
		if branchesForDay, ok := rangesByDayAndBranch[weekday]; ok {
			for branchID, empRangesInBranch := range branchesForDay {
				for _, empRange := range empRangesInBranch {
					emp := empRange.Employee // Already preloaded

					// Generate slots for this employee's work range
					slot := time.Date(d.Year(), d.Month(), d.Day(), empRange.StartTime.Hour(), empRange.StartTime.Minute(), 0, 0, loc)
					endOfDay := time.Date(d.Year(), d.Month(), d.Day(), empRange.EndTime.Hour(), empRange.EndTime.Minute(), 0, 0, loc)

					for slot.Before(endOfDay) {

						// Check availability using the map - THIS IS THE KEY
						lookupKey := fmt.Sprintf("%s-%s", emp.ID.String(), slot.Format(time.RFC3339))
						currentBookings := appointmentSlotMap[lookupKey]

						// Determine max capacity for this employee and service
						maxCapacity := emp.TotalServiceDensity
						if specificDensity, hasSpecific := densityMap[emp.ID]; hasSpecific {
							maxCapacity = specificDensity
						}

						if uint32(currentBookings) < maxCapacity {
							// This slot is available, add it to the results
							dateStr := d.Format("2006-01-02")
							timeStr := slot.Format("15:04")

							if _, ok := availabilityMap[dateStr]; !ok {
								availabilityMap[dateStr] = map[uuid.UUID]map[string][]uuid.UUID{}
							}
							if _, ok := availabilityMap[dateStr][branchID]; !ok {
								availabilityMap[dateStr][branchID] = map[string][]uuid.UUID{}
							}
							availabilityMap[dateStr][branchID][timeStr] = append(availabilityMap[dateStr][branchID][timeStr], emp.ID)

							// Populate info maps if not already present
							if _, ok := branchInfoMap[branchID]; !ok {
								branchInfoMap[branchID] = DTO.BranchBase{ /* ... map from empRange.Branch ... */ }
							}
							if _, ok := employeeInfoMap[emp.ID]; !ok {
								employeeInfoMap[emp.ID] = DTO.EmployeeBase{ /* ... map from emp ... */ }
							}
						}

						slot = slot.Add(time.Minute * time.Duration(emp.SlotTimeDiff))
					}
				}
			}
		}
	}

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
