// filter/schedule_strategy.go
package schedule_filter

import (
	"fmt"
	"time"

	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/lib" // Assuming your custom errors are in this package

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ScheduleStrategy defines the interface for different schedule fetching strategies.
type ScheduleStrategy interface {
	Fetch(params *ScheduleQueryParams) (any, error)
}

// --- Services Strategy ---
type ServicesStrategy struct {
	tx *gorm.DB
}

func (s *ServicesStrategy) Fetch(params *ScheduleQueryParams) (any, error) {
	if params.allExceptGetAreNil() {
		return nil, params.toLibError(fmt.Errorf("at least one filter is required when getting services"))
	}

	// The core of the logic is finding the right employees first.
	builder := newEmployeeQueryBuilder(s.tx).availableOn(params)
	if params.BranchID != nil {
		builder.inBranch(*params.BranchID)
	}
	if params.EmployeeID != nil {
		builder.byID(*params.EmployeeID)
	}
	// Note: ServiceID is not used to filter employees here, but to filter the final services.

	employeeIDs, err := builder.getEmployeeIDs()
	if err != nil {
		return nil, err
	}
	if len(employeeIDs) == 0 {
		return []*model.Service{}, nil // Return empty slice, not an error
	}

	// Now find the services provided by these employees.
	query := s.tx.Model(&model.Service{}).
		Joins("JOIN employee_services es ON es.service_id = services.id").
		Where("es.employee_id IN ?", employeeIDs)

	// If a branch is specified, the service must also be available at that branch.
	if params.BranchID != nil {
		query = query.Joins("JOIN branch_services bs ON bs.service_id = services.id").
			Where("bs.branch_id = ?", *params.BranchID)
	}

	var services []*model.Service
	err = query.Distinct().Find(&services).Error
	if err != nil {
		return nil, lib.Error.General.DatabaseError.WithError(err)
	}
	return services, nil
}

// --- Branches Strategy ---
type BranchesStrategy struct {
	tx *gorm.DB
}

func (s *BranchesStrategy) Fetch(params *ScheduleQueryParams) (any, error) {
	if params.allExceptGetAreNil() {
		var branches []*model.Branch
		err := s.tx.Find(&branches).Error
		if err != nil {
			return nil, lib.Error.General.DatabaseError.WithError(err)
		}
		return branches, nil
	}

	// Base query to find branches that exist. We will join and filter based on employees/services.
	query := s.tx.Model(&model.Branch{}).
		Joins("JOIN employee_branches eb ON eb.branch_id = branches.id")

	// Find employees matching the criteria first to get their IDs.
	empBuilder := newEmployeeQueryBuilder(s.tx).availableOn(params)
	if params.EmployeeID != nil {
		empBuilder.byID(*params.EmployeeID)
	}
	if params.ServiceID != nil {
		empBuilder.providesService(*params.ServiceID)
	}

	employeeIDs, err := empBuilder.getEmployeeIDs()
	if err != nil {
		return nil, err
	}
	// If specific filters were applied and no employees were found, then no branches can match.
	if len(employeeIDs) == 0 && !params.allExceptGetAreNil() {
		return []*model.Branch{}, nil
	}

	// Filter the branches to those where the found employees work.
	if len(employeeIDs) > 0 {
		query = query.Where("eb.employee_id IN ?", employeeIDs)
	}

	// If a service is specified, the branch must also directly offer it.
	if params.ServiceID != nil {
		query = query.Joins("JOIN branch_services bs ON bs.branch_id = branches.id").
			Where("bs.service_id = ?", *params.ServiceID)
	}

	var branches []*model.Branch
	err = query.Distinct().Find(&branches).Error
	if err != nil {
		return nil, lib.Error.General.DatabaseError.WithError(err)
	}
	return branches, nil
}

// --- Employees Strategy ---
type EmployeesStrategy struct {
	tx *gorm.DB
}

func (s *EmployeesStrategy) Fetch(params *ScheduleQueryParams) (any, error) {
	if params.allExceptGetAreNil() {
		var employees []*model.Employee
		err := s.tx.Preload("Branches").Preload("Services").Find(&employees).Error
		if err != nil {
			return nil, lib.Error.General.DatabaseError.WithError(err)
		}
		return employees, nil
	}

	builder := newEmployeeQueryBuilder(s.tx).availableOn(params)
	if params.BranchID != nil {
		builder.inBranch(*params.BranchID)
	}
	if params.ServiceID != nil {
		builder.providesService(*params.ServiceID)
	}

	employees, err := builder.getEmployees()
	if err != nil {
		return nil, lib.Error.General.DatabaseError.WithError(err)
	}
	return employees, nil
}

// --- TimeSlots Strategy ---

// TimeSlot represents a single available appointment slot.
type TimeSlot struct {
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	BranchID   uuid.UUID `json:"branch_id"`
	EmployeeID uuid.UUID `json:"employee_id"`
}

type TimeSlotsStrategy struct {
	tx *gorm.DB
}

func (s *TimeSlotsStrategy) Fetch(params *ScheduleQueryParams) (any, error) {
	if params.allExceptGetAreNil() {
		return nil, params.toLibError(fmt.Errorf("at least one filter is required when getting time_slots"))
	}

	// 1. Find all potentially available employees
	builder := newEmployeeQueryBuilder(s.tx).availableOn(params)
	if params.BranchID != nil {
		builder.inBranch(*params.BranchID)
	}
	if params.EmployeeID != nil {
		builder.byID(*params.EmployeeID)
	}
	if params.ServiceID != nil {
		builder.providesService(*params.ServiceID)
	}

	employees, err := builder.getEmployees()
	if err != nil {
		return nil, lib.Error.General.DatabaseError.WithError(err)
	}
	if len(employees) == 0 {
		return []TimeSlot{}, nil
	}

	var targetDate time.Time
	if params.StartTime != nil {
		targetDate = *params.StartTime
	} else {
		return nil, params.toLibError(fmt.Errorf("a specific start_time (including date) is required to generate time slots"))
	}

	var existingAppointments []model.Appointment
	loc := targetDate.Location()

	startOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, loc)
	endOfDay := startOfDay.Add(24 * time.Hour)

	employeeIDs := make([]uuid.UUID, len(employees))
	for i, e := range employees {
		employeeIDs[i] = e.ID
	}

	err = s.tx.Model(&model.Appointment{}).
		Where("employee_id IN ?", employeeIDs).
		Where("start_time >= ? AND start_time < ?", startOfDay, endOfDay).
		Find(&existingAppointments).Error

	if err != nil {
		return nil, lib.Error.General.DatabaseError.WithError(err)
	}

	// Create a quick lookup map of blocked start times
	blockedSlots := make(map[time.Time]bool)
	for _, appt := range existingAppointments {
		blockedSlots[appt.StartTime.In(loc)] = true
	}

	var allSlots []TimeSlot

	// 2. Generate slots for each employee
	for _, emp := range employees {
		// ***** FIX IS HERE: The unused 'dayOfWeek' variable was removed. *****
		workRanges := emp.GetWorkRangeForDay(targetDate.Weekday())

		for _, wr := range workRanges {
			// If a branch is filtered, only consider work ranges for that branch
			if params.BranchID != nil && wr.BranchID != *params.BranchID {
				continue
			}

			start, _ := time.ParseInLocation("15:04", wr.StartTime.String(), loc)
			end, _ := time.ParseInLocation("15:04", wr.EndTime.String(), loc)

			slotStart := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), start.Hour(), start.Minute(), 0, 0, loc)
			slotEndBoundary := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), end.Hour(), end.Minute(), 0, 0, loc)

			// Use the specific employee's slot time difference, default to 30 if not set.
			slotDurationMinutes := emp.SlotTimeDiff
			if slotDurationMinutes == 0 {
				slotDurationMinutes = 30
			}
			slotDuration := time.Duration(slotDurationMinutes) * time.Minute

			for currentSlotStart := slotStart; currentSlotStart.Add(slotDuration).Before(slotEndBoundary) || currentSlotStart.Add(slotDuration).Equal(slotEndBoundary); currentSlotStart = currentSlotStart.Add(slotDuration) {
				if _, isBlocked := blockedSlots[currentSlotStart]; isBlocked {
					continue // Skip this slot
				}
				// If a specific start time was requested, only return that single matching slot
				if params.StartTime != nil {
					if currentSlotStart.Hour() == params.StartTime.Hour() && currentSlotStart.Minute() == params.StartTime.Minute() {
						allSlots = append(allSlots, TimeSlot{
							StartTime:  currentSlotStart,
							EndTime:    currentSlotStart.Add(slotDuration),
							BranchID:   wr.BranchID,
							EmployeeID: emp.ID,
						})
						break // Found the specific slot, no need to check further in this range
					}
				} else {
					allSlots = append(allSlots, TimeSlot{
						StartTime:  currentSlotStart,
						EndTime:    currentSlotStart.Add(slotDuration),
						BranchID:   wr.BranchID,
						EmployeeID: emp.ID,
					})
				}
			}
		}
	}

	return allSlots, nil
}
