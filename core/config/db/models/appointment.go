package models

import (
	"time"

	"gorm.io/gorm"
)

// Fifth step: Scheduling with service and employee availability.
type Appointment struct {
	gorm.Model
	ServiceID  uint      `json:"service_id"`
	Service    Service   `gorm:"foreignKey:ServiceID"`
	EmployeeID uint      `json:"employee_id"`
	Employee   Employee  `gorm:"foreignKey:EmployeeID"`
	BranchID   uint      `json:"branch_id"`
	Branch     Branch    `gorm:"foreignKey:BranchID"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
}
