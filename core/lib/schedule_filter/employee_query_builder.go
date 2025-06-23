package ScheduleFilter

import (
	"fmt"

	"agenda-kaki-go/core/config/db/model" // Assuming model.WorkRange is here

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type employeeQueryBuilder struct {
	query *gorm.DB
}

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

// availableOn filters employees based on their work schedules and the provided time parameters.
func (b *employeeQueryBuilder) availableOn(params *ScheduleQueryParams) *employeeQueryBuilder {
	if params.OriginalStartTime == nil {
		return b
	}

	workRangeSubQuery := b.query.Session(&gorm.Session{NewDB: true}).Model(&model.WorkRange{})

	if params.IsSpecificDateQuery && params.QueryTimeOfDay != nil {
		// Specific date (weekday) and specific time of day
		// params.QueryTimeOfDay is already normalized (0000-01-01 HH:MM:SS UTC)
		// We query the TIME columns WorkRange.StartTime and WorkRange.EndTime
		workRangeSubQuery = workRangeSubQuery.Where("work_ranges.weekday = ?", params.QueryWeekday).
			Where("work_ranges.start_time <= ?", *params.QueryTimeOfDay). // GORM should handle time comparison
			Where("work_ranges.end_time > ?", *params.QueryTimeOfDay)   // Assuming end_time is exclusive upper bound for a start

	} else if params.IsTimeOfDayQuery && params.QueryTimeOfDay != nil {
		// Specific time of day, any working day
		workRangeSubQuery = workRangeSubQuery.Where("work_ranges.start_time <= ?", *params.QueryTimeOfDay).
			Where("work_ranges.end_time > ?", *params.QueryTimeOfDay)

	} else if params.IsSpecificDateQuery {
        // Specific date (weekday), any time during their shift on that day
        workRangeSubQuery = workRangeSubQuery.Where("work_ranges.weekday = ?", params.QueryWeekday)
    } else {
		return b // No further availability filter from OriginalStartTime
	}

    if params.BranchID != nil { // If filtering by branch, ensure work range is for that branch
        workRangeSubQuery = workRangeSubQuery.Where("work_ranges.branch_id = ?", *params.BranchID)
    }

	b.query = b.query.Where("employees.id IN (?)", workRangeSubQuery.Select("employee_id").Distinct())
	return b
}

func (b *employeeQueryBuilder) getEmployees() ([]*model.Employee, error) {
	var employees []*model.Employee
	// Consider preloading relevant associations if needed by the caller, e.g., WorkRanges
	// For example: .Preload("WorkRanges")
	if err := b.query.Preload("WorkSchedule").Distinct().Find(&employees).Error; err != nil {
		return nil, fmt.Errorf("failed to get employees: %w", err)
	}
	return employees, nil
}

func (b *employeeQueryBuilder) getEmployeeIDs() ([]uuid.UUID, error) {
	var ids []uuid.UUID
	err := b.query.Distinct().Pluck("employees.id", &ids).Error
	if err != nil {
		return nil, fmt.Errorf("failed to pluck employee ids: %w", err)
	}
	return ids, nil
}