package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	DTO "mynute-go/core/src/api/dto"
	dJSON "mynute-go/core/src/api/dto/json"
	"mynute-go/core/src/config/db/model"
	"mynute-go/core/src/handler"
	"mynute-go/core/src/lib"
	"mynute-go/core/src/api/middleware"
	"mynute-go/debug"
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
	if err := debug.Output("controller_CreateService", service); err != nil {
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
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			id				path		string	true	"Service ID"
//	@Param			id				path		string	true	"Service ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Service
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/service/{id} [get]
func GetServiceById(c *fiber.Ctx) error {
	var service model.Service

	if err := GetOneBy("id", c, &service, nil, nil); err != nil {
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
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			id				path		string	true	"Service ID"
//	@Param			name			path		string	true	"Service Name"
//	@Produce		json
//	@Success		200	{object}	DTO.Service
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/service/name/{name} [get]
func GetServiceByName(c *fiber.Ctx) error {
	var service model.Service

	if err := GetOneBy("name", c, &service, nil, nil); err != nil {
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
//	@Param			profile	formData	file	false	"Profile image"
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
//	@Param			X-Company-ID		header	string	true	"X-Company-ID"
//	@Param			id					path	string	true	"Service ID"
//	@Param			timezone			query	string	false	"Client Time Zone (IANA format, e.g., America/New_York)"
//	@Param			date_forward_start	query	number	true	"The start date for the forward search in number format"
//	@Param			date_forward_end	query	number	true	"The end date for the forward search in number format"
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
	if dfe-dfs > 31 {
		return lib.Error.General.BadRequest.WithError(errors.New("date search range (date_forward_end - date_forward_start) must not exceed 31 days"))
	}
	if dfe > 100 {
		return lib.Error.General.BadRequest.WithError(errors.New("date_forward_end must not exceed 100 days in the future"))
	}
	if dfs < 0 {
		return lib.Error.General.BadRequest.WithError(errors.New("date_forward_start must not be negative"))
	}

	companyIDStr := c.Get("X-Company-ID")
	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(errors.New("invalid X-Company-ID"))
	}

	timezone := c.Query("timezone")
	if timezone == "" {
		timezone = "UTC"
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid timezone: %s", timezone))
	}

	tx, err := lib.Session(c)
	if err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Truncate to midnight in the target timezone to avoid partial day inclusion
	now := time.Now().In(loc)
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	startDate := midnight.AddDate(0, 0, dfs)
	endDate := midnight.AddDate(0, 0, dfe).Add(24 * time.Hour) // Para incluir o dia final inteiro

	// =========================================================================
	// Step 2: Fetch all necessary data in fewer, more efficient queries
	// =========================================================================

	// --- 2a. Fetch all employee work ranges for the service, with Employee and Branch info
	var empRanges []model.EmployeeWorkRange
	err = tx.
		Joins("JOIN employee_work_range_services es ON es.employee_work_range_id = employee_work_ranges.id").
		Joins("JOIN employee_branches eb ON eb.employee_id = employee_work_ranges.employee_id AND eb.branch_id = employee_work_ranges.branch_id").
		Where("es.service_id = ?", serviceID).
		Preload("Employee"). // Preload Employee data to get SlotTimeDiff and density
		Preload("Branch").   // Preload Branch data for the response
		Find(&empRanges).Error
	if err != nil {
		return err
	}

	// --- 2b. Fetch all appointments for the relevant employees and date range in ONE query.

	// Extrai os IDs dos funcionários para filtrar a query de agendamentos.
	employeeIDSet := make(map[uuid.UUID]struct{})
	for _, er := range empRanges {
		employeeIDSet[er.EmployeeID] = struct{}{}
	}

	// Agora, convertemos o set de volta para um slice para usar na query do GORM.
	employeeIDs := make([]uuid.UUID, 0, len(employeeIDSet))
	for id := range employeeIDSet {
		employeeIDs = append(employeeIDs, id)
	}

	// TODO: Add pagination
	// TODO: Add caching

	type appointmentCountResult struct {
		EmployeeID uuid.UUID
		StartTime  time.Time
		Count      int64
	}
	
	var appointmentCounts []appointmentCountResult

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
					if emp.SlotTimeDiff <= 0 {
						continue // Use continue to proceed with the next range
					}

					// --- Timezone-Correct Shift Calculation ---
					// 1. Load the branch's specific timezone. Default to UTC if invalid or missing.
					branchLoc, err := time.LoadLocation(empRange.Branch.TimeZone)
					if err != nil {
						branchLoc = time.UTC // Fallback to UTC
					}

					// 2. Construct the shift start and end times in the branch's actual timezone.
					shiftStartInBranchTZ := time.Date(d.Year(), d.Month(), d.Day(), empRange.StartTime.Hour(), empRange.StartTime.Minute(), 0, 0, branchLoc)
					shiftEndInBranchTZ := time.Date(d.Year(), d.Month(), d.Day(), empRange.EndTime.Hour(), empRange.EndTime.Minute(), 0, 0, branchLoc)

					// 3. Convert the branch's shift times to the user's requested timezone for the loop.
					slot := shiftStartInBranchTZ.In(loc)
					endOfDay := shiftEndInBranchTZ.In(loc)

					// 4. For today's date, ensure we don't show slots from the past.
					if dfs == 0 && d.Year() == now.Year() && d.Month() == now.Month() && d.Day() == now.Day() {
						if slot.Before(now) {
							slot = now
						}
					}

					// 5. Align the final start time to the next valid slot boundary.
					if emp.SlotTimeDiff > 0 {
						// Calculate alignment based on the user's timezone view of the slot.
						minutesSinceMidnight := slot.Hour()*60 + slot.Minute()
						slotDiff := int(emp.SlotTimeDiff)
						if remainder := minutesSinceMidnight % slotDiff; remainder != 0 {
							minutesToAdd := slotDiff - remainder
							slot = slot.Add(time.Duration(minutesToAdd) * time.Minute)
						}
					}

				for slot.Before(endOfDay) {

					// Check availability using the map - THIS IS THE KEY
					lookupKey := fmt.Sprintf("%s-%s", emp.UserID.String(), slot.Format(time.RFC3339))
					currentBookings := appointmentSlotMap[lookupKey]

					// Determine max capacity for this employee and service
					maxCapacity := emp.TotalServiceDensity
					if specificDensity, hasSpecific := densityMap[emp.UserID]; hasSpecific {
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
						availabilityMap[dateStr][branchID][timeStr] = append(availabilityMap[dateStr][branchID][timeStr], emp.UserID)

						// Populate info maps if not already present
						if _, ok := branchInfoMap[branchID]; !ok {
							empRangeBranchBytes, err := json.Marshal(empRange.Branch)
							if err != nil {
								return fmt.Errorf("failed to marshal branch info: %w", err)
							}
							var dtoBranchBase DTO.BranchBase
							if err := json.Unmarshal(empRangeBranchBytes, &dtoBranchBase); err != nil {
								return fmt.Errorf("failed to unmarshal branch info: %w", err)
							}
							branchInfoMap[branchID] = dtoBranchBase
						}
						if _, ok := employeeInfoMap[emp.UserID]; !ok {
							empBytes, err := json.Marshal(emp)
							if err != nil {
								return fmt.Errorf("failed to marshal employee info: %w", err)
						}
						var dtoEmployeeBase DTO.EmployeeBase
						if err := json.Unmarshal(empBytes, &dtoEmployeeBase); err != nil {
							return fmt.Errorf("failed to unmarshal employee info: %w", err)
						}
						employeeInfoMap[emp.UserID] = dtoEmployeeBase
					}
				}

				slot = slot.Add(time.Minute * time.Duration(emp.SlotTimeDiff))
			}
		}
	}
}	availableDateMap := map[string]map[uuid.UUID]*DTO.AvailableDate{}

	client_public_id := c.Query("client_public_id")
	var clientAppointments []model.ClientAppointment
	if client_public_id != "" {
		if err := lib.ChangeToPublicSchemaByContext(c); err != nil {
			return err
		}
		if err := tx.Model(&model.ClientAppointment{}).
			Where("client_id = ? AND is_cancelled = ?", client_public_id, false).
			Where("start_time >= ? AND start_time < ?", startDate, endDate).
			Find(&clientAppointments).Error; err != nil {
			return lib.Error.General.InternalError.WithError(err)
		}
		if err := lib.ChangeToCompanySchemaByContext(c); err != nil {
			return err
		}
	}

	// Get service duration to help with conflict checking
	var serviceDuration uint16
	if client_public_id != "" {
		if err := tx.
		Model(&model.Service{}).
		Where("id = ?", serviceID).
		Pluck("duration", &serviceDuration).Error; err != nil {
			return lib.Error.General.InternalError.WithError(err)
		}
	}

	for date, branches := range availabilityMap {
		if _, ok := availableDateMap[date]; !ok {
			availableDateMap[date] = map[uuid.UUID]*DTO.AvailableDate{}
		}

		for branchID, slots := range branches {
			if _, ok := availableDateMap[date][branchID]; !ok {
				availableDateMap[date][branchID] = &DTO.AvailableDate{
					Date:           date,
					BranchID:       branchID,
					AvailableTimes: []DTO.AvailableTime{},
				}
			}

			slotLoop:
			for timeStr, empIDs := range slots {
				// Filter out times where the client already has an appointment
				if client_public_id != "" {
					// Create a time object for the current slot
					slotTime_RFC3339 := fmt.Sprintf("%sT%s:00", date, timeStr)
					slotTime, err := time.Parse("2006-01-02T15:04:05", slotTime_RFC3339)
					if err != nil {
						return lib.Error.General.InternalError.WithError(fmt.Errorf("failed to parse slot time: %w", err))
					}

					for _, appt := range clientAppointments {
						// Same start time means conflict
						if appt.StartTime.Equal(slotTime) {
							continue slotLoop
						}
						// If slotTime is between appt.StartTime and appt.EndTime then there is a conflict
						if appt.StartTime.Before(slotTime) && appt.EndTime.After(slotTime) {
							continue slotLoop
						}
						// If appt.StartTime is After slotTime 
						// AND 
						// appt.StartTime is Before slotTime + service duration then there is a conflict
						slotTime_end := slotTime.Add(time.Minute * time.Duration(serviceDuration))
						if appt.StartTime.After(slotTime) && appt.StartTime.Before(slotTime_end) {
							continue slotLoop
						}
					}
				}
				availableDateMap[date][branchID].AvailableTimes = append(
					availableDateMap[date][branchID].AvailableTimes,
					DTO.AvailableTime{
						Time:        timeStr,
						EmployeesID: empIDs,
					},
				)
			}
		}
	}

	// Flatten into slice
	var availableDates []DTO.AvailableDate
	for _, branches := range availableDateMap {
		for _, ad := range branches {
			availableDates = append(availableDates, *ad)
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

	Availability := DTO.ServiceAvailability{
		ServiceID:      serviceID,
		AvailableDates: availableDates,
		EmployeeInfo:   employeeInfo,
		BranchInfo:     branchInfo,
	}

	debug.Output("controller_GetServiceAvailability", Availability)

	return lib.ResponseFactory(c).Send(200, Availability)
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


