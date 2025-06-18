// filter/employee_query_builder.go
package schedule_filter

import (
	"fmt"
	"strings"

	"agenda-kaki-go/core/config/db/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// employeeQueryBuilder constructs a GORM query to find employees based on various criteria.
type employeeQueryBuilder struct {
	query *gorm.DB
}

// newEmployeeQueryBuilder initializes a new builder.
func newEmployeeQueryBuilder(tx *gorm.DB) *employeeQueryBuilder {
	return &employeeQueryBuilder{
		query: tx.Model(&model.Employee{}),
	}
}

func (b *employeeQueryBuilder) byID(id uuid.UUID) *employeeQueryBuilder {
	b.query = b.query.Where("employees.id = ?", id)
	return b
}

func (b *employeeQueryBuilder) inBranch(id uuid.UUID) *employeeQueryBuilder {
	b.query = b.query.Joins("JOIN employee_branches ON employee_branches.employee_id = employees.id").
		Where("employee_branches.branch_id = ?", id)
	return b
}

func (b *employeeQueryBuilder) providesService(id uuid.UUID) *employeeQueryBuilder {
	b.query = b.query.Joins("JOIN employee_services ON employee_services.employee_id = employees.id").
		Where("employee_services.service_id = ?", id)
	return b
}

// availableOn adds conditions for employee work schedules stored in JSONB.
func (b *employeeQueryBuilder) availableOn(params *ScheduleQueryParams) *employeeQueryBuilder {
	// If both are nil, do nothing.
	if params.Weekday == nil && params.StartTime == nil {
		return b
	}

	var jsonClauses []string
	var args []any

	if params.Weekday != nil && params.StartTime != nil {
		// Filter by a specific day and time.
		dayKey := strings.ToLower(*params.Weekday)
		timeStr := params.StartTime.Format("15:04")

		// This complex clause checks if an element exists in the JSON array for the given day
		// where the provided time is within the start/end range of the work shift.
		jsonClauses = append(jsonClauses, `jsonb_path_exists(work_schedule, '$.??[*] ? (@.start <= ? && @.end > ?)')`)
		args = append(args, dayKey, timeStr, timeStr)
	} else if params.Weekday != nil {
		// Filter by any availability on a given day.
		dayKey := strings.ToLower(*params.Weekday)
		// Checks if the JSON array for the given day exists and is not empty.
		jsonClauses = append(jsonClauses, `jsonb_array_length(work_schedule -> ?) > 0`)
		args = append(args, dayKey)
	} else if params.StartTime != nil {
		// Filter by a specific time on ANY day of the week.
		timeStr := params.StartTime.Format("15:04")
		// This is the most complex query. It checks all weekday arrays.
		jsonClauses = append(jsonClauses, `(
            jsonb_path_exists(work_schedule, '$.monday[*] ? (@.start <= ? && @.end > ?)') OR
            jsonb_path_exists(work_schedule, '$.tuesday[*] ? (@.start <= ? && @.end > ?)') OR
            jsonb_path_exists(work_schedule, '$.wednesday[*] ? (@.start <= ? && @.end > ?)') OR
            jsonb_path_exists(work_schedule, '$.thursday[*] ? (@.start <= ? && @.end > ?)') OR
            jsonb_path_exists(work_schedule, '$.friday[*] ? (@.start <= ? && @.end > ?)') OR
            jsonb_path_exists(work_schedule, '$.saturday[*] ? (@.start <= ? && @.end > ?)') OR
            jsonb_path_exists(work_schedule, '$.sunday[*] ? (@.start <= ? && @.end > ?)')
        )`)
		// Add the time argument 7 times
		for i := 0; i < 7; i++ {
			args = append(args, timeStr, timeStr)
		}
	}

	if len(jsonClauses) > 0 {
		b.query = b.query.Where(strings.Join(jsonClauses, " AND "), args...)
	}

	return b
}

// getEmployees executes the query and returns a slice of employees.
func (b *employeeQueryBuilder) getEmployees() ([]*model.Employee, error) {
	var employees []*model.Employee
	if err := b.query.Distinct().Find(&employees).Error; err != nil {
		return nil, fmt.Errorf("failed to get employees: %w", err)
	}
	return employees, nil
}

// getEmployeeIDs executes the query and returns only the employee IDs.
func (b *employeeQueryBuilder) getEmployeeIDs() ([]uuid.UUID, error) {
	var ids []uuid.UUID
	err := b.query.Distinct().Pluck("employees.id", &ids).Error
	if err != nil {
		return nil, fmt.Errorf("failed to pluck employee ids: %w", err)
	}
	return ids, nil
}