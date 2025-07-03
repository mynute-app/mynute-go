package modelT

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	mJSON "agenda-kaki-go/core/config/db/model/json"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/test/handlers"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Appointment struct {
	Created  *model.Appointment
	Employee *Employee
	Company  *Company
	Client   *Client
	Branch   *Branch
	Service  *Service
}

func (a *Appointment) Reset() *Appointment {
	a.Created = nil
	a.Employee = nil
	a.Company = nil
	a.Client = nil
	a.Branch = nil
	a.Service = nil
	return a
}

func (a *Appointment) CreateRandomly(s int, cy *Company, ct *Client, e *Employee, token, company_id string) error {
	if a.Created != nil {
		return fmt.Errorf("appointment already created, cannot create again")
	}
	preferredLocation := time.UTC // Choose your time_zone (e.g., UTC)
	appointmentSlot, found, err := a.FindValidAppointmentSlot(e, preferredLocation)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("setup failed: could not find a valid appointment slot for initial booking")
	}
	var branch *Branch
	for _, b := range e.Company.Branches {
		if b.Created.ID.String() == appointmentSlot.BranchID {
			branch = b
			break
		}
	}
	var service *Service
	for _, s := range e.Company.Services {
		if s.Created.ID.String() == appointmentSlot.ServiceID {
			service = s
			break
		}
	}
	if branch == nil {
		return fmt.Errorf("branch with ID %s not loaded at company %s modelT structure", appointmentSlot.BranchID, cy.Created.ID.String())
	}
	if service == nil {
		return fmt.Errorf("service with ID %s not loaded at company %s modelT structure", appointmentSlot.ServiceID, cy.Created.ID.String())
	}
	http := handler.NewHttpClient()
	if err := http.
		Method("POST").
		URL("/appointment").
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, token).
		Header(namespace.HeadersKey.Company, company_id).
		Send(map[string]any{
			"branch_id":   appointmentSlot.BranchID,
			"service_id":  appointmentSlot.ServiceID,
			"employee_id": e.Created.ID.String(),
			"company_id":  cy.Created.ID.String(),
			"client_id":   ct.Created.ID.String(),
			"start_time":  appointmentSlot.StartTimeRFC3339, // Use found start time
		}).
		ParseResponse(&a.Created).Error; err != nil {
		return fmt.Errorf("failed to create appointment: %w", err)
	}
	if a.Created.ID == uuid.Nil && s != 200 && s != 201 {
		a.Reset()
		return nil
	}
	if err := ct.GetByEmail(200); err != nil {
		return err
	}
	if err := e.GetById(200, nil, nil); err != nil {
		return err
	}
	if err := branch.GetById(200, cy.Owner.X_Auth_Token, &company_id); err != nil {
		return err
	}
	// Assign the loaded entities to the Appointment struct
	a.Branch = branch
	a.Service = service
	a.Employee = e
	a.Company = cy
	a.Client = ct
	// Assign the appointment to the Employee, Branch and Client
	e.Appointments = append(e.Appointments, a)
	branch.Appointments = append(branch.Appointments, a)
	ct.Appointments = append(ct.Appointments, a)
	return nil
}

func (a *Appointment) RescheduleRandomly(s int, x_auth_token string, x_company_id *string) error {
	companyIDStr := a.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	preferredLocation := time.UTC // Choose your time_zone (e.g., UTC)
	appointmentSlot, found, err := a.FindValidAppointmentSlot(a.Employee, preferredLocation)
	if err != nil {
		return fmt.Errorf("failed to find valid appointment slot: %w", err)
	}
	if !found {
		return fmt.Errorf("no valid appointment slot found for employee %s in company %s", a.Employee.Created.ID.String(), a.Company.Created.ID.String())
	}
	var branch *Branch
	for _, b := range a.Company.Branches {
		if b.Created.ID.String() == appointmentSlot.BranchID {
			branch = b
			break
		}
	}
	var service *Service
	for _, s := range a.Company.Services {
		if s.Created.ID.String() == appointmentSlot.ServiceID {
			service = s
			break
		}
	}
	if branch == nil {
		return fmt.Errorf("branch with ID %s not loaded at company %s modelT structure", appointmentSlot.BranchID, a.Company.Created.ID.String())
	}
	if service == nil {
		return fmt.Errorf("service with ID %s not loaded at company %s modelT structure", appointmentSlot.ServiceID, a.Company.Created.ID.String())
	}
	if err := handler.NewHttpClient().
		Method("PATCH").
		URL("/appointment/"+a.Created.ID.String()).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Header(namespace.HeadersKey.Company, cID).
		Send(map[string]any{
			"branch_id":  appointmentSlot.BranchID,
			"service_id": appointmentSlot.ServiceID,
			"start_time": appointmentSlot.StartTimeRFC3339,
		}).Error; err != nil {
		return fmt.Errorf("failed to reschedule appointment: %w", err)
	}
	a.Branch = branch
	a.Service = service
	if err := a.Employee.GetById(200, nil, nil); err != nil {
		return err
	}
	if err := a.Branch.GetById(200, a.Company.Owner.X_Auth_Token, nil); err != nil {
		return err
	}
	if err := a.Client.GetByEmail(200); err != nil {
		return err
	}
	a.Created.BranchID = uuid.MustParse(appointmentSlot.BranchID)
	a.Created.ServiceID = uuid.MustParse(appointmentSlot.ServiceID)
	a.Created.StartTime, err = time.Parse(time.RFC3339, appointmentSlot.StartTimeRFC3339)
	if err != nil {
		return fmt.Errorf("failed to parse new start time: %w", err)
	}
	return nil
}

func (a *Appointment) Create(status int, x_auth_token string, x_company_id *string, startTime *string, tz string, b *Branch, e *Employee, s *Service, cy *Company, ct *Client) error {
	companyIDStr := cy.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	http := handler.NewHttpClient()
	http.Method("POST")
	http.URL("/appointment")
	http.ExpectedStatus(status)
	http.Header(namespace.HeadersKey.Company, cID)
	http.Header(namespace.HeadersKey.Auth, x_auth_token)
	if startTime == nil {
		tempStartTime := lib.GenerateDateRFC3339(2027, 10, 29)
		startTime = &tempStartTime
	}
	A := DTO.CreateAppointment{
		BranchID:   b.Created.ID,
		EmployeeID: e.Created.ID,
		ServiceID:  s.Created.ID,
		ClientID:   ct.Created.ID,
		CompanyID:  cy.Created.ID,
		StartTime:  *startTime,
		TimeZone:   tz,
	}
	if err := http.Send(A).Error; err != nil {
		return fmt.Errorf("failed to create appointment: %w", err)
	}
	if err := http.ParseResponse(&a.Created).Error; err != nil {
		return fmt.Errorf("failed to parse appointment response: %w", err)
	}
	if err := b.GetById(200, cy.Owner.X_Auth_Token, x_company_id); err != nil {
		return err
	}
	if err := e.GetById(200, nil, x_company_id); err != nil {
		return err
	}
	if err := s.GetById(200, cy.Owner.X_Auth_Token, x_company_id); err != nil {
		return err
	}
	if err := cy.GetById(200, cy.Owner.X_Auth_Token, x_company_id); err != nil {
		return err
	}
	if err := ct.GetByEmail(200); err != nil {
		return err
	}
	var ClientAppointment mJSON.ClientAppointment
	aCreatedByte, err := json.Marshal(a.Created)
	if err != nil {
		return fmt.Errorf("failed to marshal appointment: %w", err)
	}
	err = json.Unmarshal(aCreatedByte, &ClientAppointment)
	if err != nil {
		return fmt.Errorf("failed to unmarshal appointment: %w", err)
	}
	return nil
}

func (a *Appointment) GetById(s int, x_auth_token string, x_company_id *string) error {
	companyIDStr := a.Created.CompanyID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("GET").
		URL("/appointment/"+a.Created.ID.String()).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Header(namespace.HeadersKey.Company, cID).
		Send(nil).
		ParseResponse(&a.Created).Error; err != nil {
		return fmt.Errorf("failed to get appointment %s: %w", a.Created.ID.String(), err)
	}
	return nil
}

func (a *Appointment) Cancel(s int, x_auth_token string, x_company_id *string) error {
	companyIDStr := a.Created.CompanyID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("DELETE").
		URL("/appointment/"+a.Created.ID.String()).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Header(namespace.HeadersKey.Company, cID).
		Send(nil).Error; err != nil {
		return fmt.Errorf("failed to cancel appointment %s: %w", a.Created.ID.String(), err)
	}
	if s >= 400 {
		return nil
	}
	// Delete appointment from Employee, Branch and Client
	for i, appt := range a.Employee.Appointments {
		if appt.Created.ID == a.Created.ID {
			a.Employee.Appointments = slices.Delete(a.Employee.Appointments, i, i+1)
			break
		}
	}
	for i, appt := range a.Branch.Appointments {
		if appt.Created.ID == a.Created.ID {
			a.Branch.Appointments = slices.Delete(a.Branch.Appointments, i, i+1)
			break
		}
	}
	for i, appt := range a.Client.Appointments {
		if appt.Created.ID == a.Created.ID {
			a.Client.Appointments = slices.Delete(a.Client.Appointments, i, i+1)
			break
		}
	}
	a.Reset() // Reset the Appointment struct after cancellation
	return nil
}

type FoundAppointmentSlot struct {
	StartTimeRFC3339 string
	BranchID         string
	ServiceID        string
	TimeZone         string // Timezone in IANA format, e.g., "America/New_York"
}

const (
	slotSearchHorizonDays = 14               // Example: search up to 2 weeks ahead
	slotSearchTimeStep    = 15 * time.Minute // Example: check every 15 minutes
)

func (a *Appointment) FindValidAppointmentSlot(employee *Employee, preferredLocation *time.Location) (*FoundAppointmentSlot, bool, error) {
	if preferredLocation == nil {
		return nil, false, fmt.Errorf("preferredLocation is nil; time_zone must be explicitly passed")
	}

	// fmt.Printf("---- Starting findValidAppointmentSlot for Employee ID: %s ----\n", employee.Created.ID.String())

	weekdaySchedules := map[time.Weekday][]model.EmployeeWorkRange{
		time.Sunday:    employee.Created.GetWorkRangeForDay(time.Sunday),
		time.Monday:    employee.Created.GetWorkRangeForDay(time.Monday),
		time.Tuesday:   employee.Created.GetWorkRangeForDay(time.Tuesday),
		time.Wednesday: employee.Created.GetWorkRangeForDay(time.Wednesday),
		time.Thursday:  employee.Created.GetWorkRangeForDay(time.Thursday),
		time.Friday:    employee.Created.GetWorkRangeForDay(time.Friday),
		time.Saturday:  employee.Created.GetWorkRangeForDay(time.Saturday),
	}

	now := time.Now().In(preferredLocation)
	searchStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, preferredLocation)

	branchCache := make(map[string]*model.Branch)
	serviceCache := make(map[string]*model.Service)

	for dayOffset := range slotSearchHorizonDays {
		currentDate := searchStart.AddDate(0, 0, dayOffset)
		currentWeekday := currentDate.Weekday()
		workRanges := weekdaySchedules[currentWeekday]

		for _, wr := range workRanges {
			branch, ok := branchCache[wr.BranchID.String()]
			if !ok {
				var branchModel model.Branch
				if err := handler.NewHttpClient().
					Header(namespace.HeadersKey.Company, employee.Company.Created.ID.String()).
					Header(namespace.HeadersKey.Auth, employee.Company.Owner.X_Auth_Token).
					Method("GET").
					URL("/branch/" + wr.BranchID.String()).
					ExpectedStatus(200).
					Send(nil).
					ParseResponse(&branchModel).Error; err != nil {
					return nil, false, fmt.Errorf("failed to get branch by ID %s: %w", wr.BranchID.String(), err)
				}
				branch = &branchModel
				branchCache[wr.BranchID.String()] = branch
			}
			branchID := branch.ID.String()
			for _, wrSrvc := range wr.Services {
				sID := wrSrvc.ID.String()
				service, ok := serviceCache[sID]
				if !ok {
					var serviceModel model.Service
					if err := handler.NewHttpClient().
						Header(namespace.HeadersKey.Company, employee.Company.Created.ID.String()).
						Header(namespace.HeadersKey.Auth, employee.Company.Owner.X_Auth_Token).
						Method("GET").
						URL("/service/" + sID).
						ExpectedStatus(200).
						Send(nil).
						ParseResponse(&serviceModel).Error; err != nil {
						return nil, false, fmt.Errorf("failed to get service by ID %s: %w", sID, err)
					}
					service = &serviceModel
					serviceCache[sID] = service
				}
				duration := time.Duration(service.Duration) * time.Minute
				startTime := wr.StartTime.In(preferredLocation)
				endTime := wr.EndTime.In(preferredLocation)
				// startTime, err := parseTimeWithLocation(currentDate, wr.StartTime.String(), wr.TimeZone, preferredLocation)
				// if err != nil {
				// 	return nil, false, fmt.Errorf("failed to parse start time for work range #%d: %w", iWr, err)
				// }
				// endTime, err := parseTimeWithLocation(currentDate, wr.EndTime.String(), wr.TimeZone, preferredLocation)
				// if err != nil || !startTime.Before(endTime) {
				// 	return nil, false, fmt.Errorf("invalid time range for work range #%d: %w", iWr, err)
				// }
				for tStart := startTime; tStart.Add(duration).Before(endTime) || tStart.Add(duration).Equal(endTime); tStart = tStart.Add(slotSearchTimeStep) {
					if tStart.Before(now) {
						continue
					}
					tEnd := tStart.Add(duration)
					overlap := false
					for _, appt := range employee.Created.Appointments {
						start := appt.StartTime.In(preferredLocation)
						end := appt.EndTime.In(preferredLocation)
						if end.IsZero() && appt.Service != nil {
							end = start.Add(time.Duration(appt.Service.Duration) * time.Minute)
						}
						if start.Before(tEnd) && end.After(tStart) {
							overlap = true
							break
						}
					}
					if overlap {
						continue
					}
					return &FoundAppointmentSlot{
						StartTimeRFC3339: tStart.Format(time.RFC3339),
						BranchID:         branchID,
						ServiceID:        sID,
						TimeZone:         preferredLocation.String(),
					}, true, nil
				}
			}
		}
	}
	return nil, false, fmt.Errorf("no valid appointment slot found for employee %s", employee.Created.ID.String())
}

// parseTimeWithLocation parses a time string in "HH:MM" (15:04) or "HH:MM:SS" (15:04:05) format and combines it with a given date.
// It also supports full time.Time string formats like "YYYY-MM-DD HH:MM:SS -LLLL -LL" (0001-01-01 07:00:00 -0300 -03), extracting only the time part.
//
// Parameters:
// - targetDate: the date to apply the parsed time to.
// - timeStr: a string representing the time, either in "HH:MM", "HH:MM:SS", or full time.Time.String() format.
// - loc: the time.Location to apply.
//
// Returns:
// - A time.Time with the date from targetDate and the time from timeStr.
// - An error if parsing fails.
func parseTimeWithLocation(targetDate time.Time, timeStr string, original_time_zone string, preferred_time_zone *time.Location) (time.Time, error) {
	// If the input looks like a full time.Time string, extract just the "HH:MM:SS" part
	if strings.Contains(timeStr, " ") && strings.Contains(timeStr, "-") && strings.Contains(timeStr, ":") {
		parts := strings.Fields(timeStr)
		if len(parts) >= 2 {
			timeStr = parts[1] // Extract only the time portion, e.g., "07:00:00"
		}
	}

	// Determine the expected layout based on number of colons
	layout := "15:04"
	if strings.Count(timeStr, ":") == 2 {
		layout = "15:04:05"
	}

	time_original_time_zone, err := time.LoadLocation(original_time_zone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time zone '%s': %w", time_original_time_zone, err)
	}

	timeParsed, err := time.ParseInLocation(layout, timeStr, time_original_time_zone)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time string '%s' with layout '%s': %w", timeStr, layout, err)
	}

	timeInLoc := timeParsed.In(preferred_time_zone)

	return time.Date(
		targetDate.Year(), targetDate.Month(), targetDate.Day(),
		timeInLoc.Hour(), timeInLoc.Minute(), timeInLoc.Second(), 0,
		preferred_time_zone,
	), nil
}
