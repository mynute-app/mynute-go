package utilsT

import (
	"agenda-kaki-go/core/config/db/model"
	mJSON "agenda-kaki-go/core/config/db/model/json"
	"agenda-kaki-go/core/config/namespace"
	handlerT "agenda-kaki-go/core/test/handlers"
	modelT "agenda-kaki-go/core/test/models"
	"fmt"
	"time"
)

type FoundAppointmentSlot struct {
	StartTimeRFC3339 string
	BranchID         string
	ServiceID        string
}

const (
	slotSearchHorizonDays = 14               // Example: search up to 2 weeks ahead
	slotSearchTimeStep    = 15 * time.Minute // Example: check every 15 minutes
)

func FindValidAppointmentSlotV2(employee *modelT.Employee, preferredLocation *time.Location) (*FoundAppointmentSlot, bool, error) {
	if preferredLocation == nil {
		return nil, false, fmt.Errorf("preferredLocation is nil; timezone must be explicitly passed")
	}

	fmt.Printf("---- Starting findValidAppointmentSlot for Employee ID: %s ----\n", employee.Created.ID.String())

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
				if err := handlerT.NewHttpClient().
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
					if err := handlerT.NewHttpClient().
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

func RescheduleAppointmentRandomly(s int, employee *modelT.Employee, company *modelT.Company, appointment_id, token string) error {
	preferredLocation := time.UTC // Choose your timezone (e.g., UTC)
	appointmentSlot, found, err := FindValidAppointmentSlotV2(employee, preferredLocation)
	if err != nil {
		return fmt.Errorf("failed to find valid appointment slot: %w", err)
	}
	if !found {
		return fmt.Errorf("no valid appointment slot found for employee %s in company %s", employee.Created.ID.String(), company.Created.ID.String())
	}
	if err := handlerT.NewHttpClient().
		Method("PATCH").
		URL("/appointment/"+appointment_id).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, token).
		Header(namespace.HeadersKey.Company, company.Created.ID.String()).
		Send(map[string]any{
			"branch_id":  appointmentSlot.BranchID,
			"service_id": appointmentSlot.ServiceID,
			"start_time": appointmentSlot.StartTimeRFC3339,
		}).Error; err != nil {
		return fmt.Errorf("failed to reschedule appointment: %w", err)
	}
	if err := employee.GetById(200, nil, nil); err != nil {
		return err
	}
	if err := company.GetById(200, company.Owner.X_Auth_Token, nil); err != nil {
		return err
	}
	return nil
}

func CreateAppointmentRandomly(s int, company *modelT.Company, client *modelT.Client, employee *modelT.Employee, token, company_id string, a *modelT.Appointment) error {
	preferredLocation := time.UTC // Choose your timezone (e.g., UTC)
	appointmentSlot, found, err := FindValidAppointmentSlotV2(employee, preferredLocation)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("setup failed: could not find a valid appointment slot for initial booking")
	}
	http := handlerT.NewHttpClient()
	if err := http.
		Method("POST").
		URL("/appointment").
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, token).
		Header(namespace.HeadersKey.Company, company_id).
		Send(map[string]any{
			"branch_id":   appointmentSlot.BranchID,
			"service_id":  appointmentSlot.ServiceID,
			"employee_id": employee.Created.ID.String(),
			"company_id":  company.Created.ID.String(),
			"client_id":   client.Created.ID.String(),
			"start_time":  appointmentSlot.StartTimeRFC3339, // Use found start time
		}).Error; err != nil {
		return fmt.Errorf("failed to create appointment: %w", err)
	}
	if a != nil {
		http.ParseResponse(&a.Created)
	}
	if err := company.GetById(200, company.Owner.X_Auth_Token, nil); err != nil {
		return err
	}
	if err := client.GetByEmail(200); err != nil {
		return err
	}
	if err := employee.GetById(200, nil, nil); err != nil {
		return err
	}
	return nil
}

func GetAppointment(s int, appointment_id string, company_id, token string, a *modelT.Appointment) error {
	http := handlerT.NewHttpClient()
	if err := http.
		Method("GET").
		URL("/appointment/"+appointment_id).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, token).
		Header(namespace.HeadersKey.Company, company_id).
		Send(nil).Error; err != nil {
		return fmt.Errorf("failed to get appointment by ID: %w", err)
	}
	if a != nil {
		if err := http.ParseResponse(&a.Created).Error; err != nil {
			return fmt.Errorf("failed to parse appointment response: %w", err)
		}
	}
	return nil
}

func CancelAppointment(s int, appointment_id, company_id, token string) error {
	if err := handlerT.NewHttpClient().
		Method("DELETE").
		URL("/appointment/"+appointment_id).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, token).
		Header(namespace.HeadersKey.Company, company_id).
		Send(nil).Error; err != nil {
		return fmt.Errorf("failed to cancel appointment: %w", err)
	}
	return nil
}
