package models

import "gorm.io/gorm"

// Fifth step: Scheduling with service and employee availability.
type Schedule struct {
	gorm.Model
	ServiceID  uint      `json:"service_id"`
	EmployeeID uint      `json:"employee_id"`
	BranchID   uint      `json:"branch_id"`
	Duration   TimeRange `json:"duration"`
}