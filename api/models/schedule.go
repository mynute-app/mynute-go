package models

import (
	"time"

	"gorm.io/gorm"
)

// Fifth step: Scheduling with service and employee availability.
type Schedule struct {
	gorm.Model
	ServiceID  uint      `json:"service_id"`
	EmployeeID uint      `json:"employee_id"`
	BranchID   uint      `json:"branch_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
}