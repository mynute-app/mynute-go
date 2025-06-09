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

func (a *Appointment) CreateRandomly(s int, cy *Company, ct *Client, e *Employee, token, company_id string) error {
	if a.Created != nil {
		return fmt.Errorf("appointment already created, cannot create again")
	}
	preferredLocation := time.UTC // Choose your timezone (e.g., UTC)
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
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	preferredLocation := time.UTC // Choose your timezone (e.g., UTC)
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

func (a *Appointment) Create(status int, x_auth_token string, x_company_id *string, startTime *string, b *Branch, e *Employee, s *Service, cy *Company, ct *Client) error {
	companyIDStr := cy.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
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
	}
	http.Send(A)
	http.ParseResponse(&a.Created)
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
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
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
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
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
	a.Created = nil // Clear the created appointment
	// Clear references to Employee, Branch, and Client
	a.Employee = nil
	a.Branch = nil
	a.Client = nil
	a.Service = nil // Clear the service reference as well
	a.Company = nil // Clear the company reference as well
	return nil
}

type FoundAppointmentSlot struct {
	StartTimeRFC3339 string
	BranchID         string
	ServiceID        string
}

const (
	slotSearchHorizonDays = 14               // Example: search up to 2 weeks ahead
	slotSearchTimeStep    = 15 * time.Minute // Example: check every 15 minutes
)

func (a *Appointment) FindValidAppointmentSlot(employee *Employee, preferredLocation *time.Location) (*FoundAppointmentSlot, bool, error) {
	if preferredLocation == nil {
		return nil, false, fmt.Errorf("preferredLocation is nil; timezone must be explicitly passed")
	}

	// fmt.Printf("---- Starting findValidAppointmentSlot for Employee ID: %s ----\n", employee.Created.ID.String())

	workSchedule := employee.Created.WorkSchedule
	weekdaySchedules := map[time.Weekday][]mJSON.WorkRange{
		time.Sunday:    workSchedule.Sunday,
		time.Monday:    workSchedule.Monday,
		time.Tuesday:   workSchedule.Tuesday,
		time.Wednesday: workSchedule.Wednesday,
		time.Thursday:  workSchedule.Thursday,
		time.Friday:    workSchedule.Friday,
		time.Saturday:  workSchedule.Saturday,
	}

	now := time.Now().In(preferredLocation)
	searchStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, preferredLocation)

	branchCache := make(map[string]*model.Branch)
	serviceCache := make(map[string]*model.Service)

	for dayOffset := range slotSearchHorizonDays {
		currentDate := searchStart.AddDate(0, 0, dayOffset)
		currentWeekday := currentDate.Weekday()
		workRanges := weekdaySchedules[currentWeekday]

		for iWr, wr := range workRanges {
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
			for _, wrSrvcID := range wr.Services {
				sID := wrSrvcID.String()
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
				startTime, err := parseTimeWithLocation(currentDate, wr.Start, preferredLocation)
				if err != nil {
					return nil, false, fmt.Errorf("failed to parse start time for work range #%d: %w", iWr, err)
				}
				endTime, err := parseTimeWithLocation(currentDate, wr.End, preferredLocation)
				if err != nil || !startTime.Before(endTime) {
					return nil, false, fmt.Errorf("invalid time range for work range #%d: %w", iWr, err)
				}
				for t := startTime; t.Add(duration).Before(endTime) || t.Add(duration).Equal(endTime); t = t.Add(slotSearchTimeStep) {
					if t.Before(now) {
						continue
					}
					tEnd := t.Add(duration)
					overlap := false
					for _, appt := range employee.Created.Appointments {
						start := appt.StartTime.In(preferredLocation)
						end := appt.EndTime
						if end.IsZero() && appt.Service != nil {
							end = start.Add(time.Duration(appt.Service.Duration) * time.Minute)
						}
						if start.Before(tEnd) && end.After(t) {
							overlap = true
							break
						}
					}
					if overlap {
						continue
					}
					return &FoundAppointmentSlot{
						StartTimeRFC3339: t.Format(time.RFC3339),
						BranchID:         branchID,
						ServiceID:        sID,
					}, true, nil
				}
			}
		}
	}
	return nil, false, fmt.Errorf("no valid appointment slot found for employee %s", employee.Created.ID.String())
}

// Helper to parse HH:MM or HH:MM:SS time string into a full time.Time on a specific date/location
func parseTimeWithLocation(targetDate time.Time, timeStr string, loc *time.Location) (time.Time, error) {
	layout := "15:04" // Default HH:MM
	colonCount := 0
	for _, r := range timeStr {
		if r == ':' {
			colonCount++
		}
	}
	if colonCount == 2 { // Detect HH:MM:SS
		layout = "15:04:05"
	}

	parsedTime, err := time.ParseInLocation(layout, timeStr, loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time string '%s' with layout '%s': %w", timeStr, layout, err)
	}
	// Combine the date part from targetDate with the time parts from parsedTime
	return time.Date(
		targetDate.Year(), targetDate.Month(), targetDate.Day(),
		parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), 0, // Nanoseconds set to 0
		loc,
	), nil
}
