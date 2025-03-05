package model

import (
	"time"

	"gorm.io/gorm"
)

// Fifth step: Scheduling with service and employee availability.
type Appointment struct {
	gorm.Model
	ServiceID  uint      `gorm:"not null" json:"service_id"`
	Service    Service   `gorm:"foreignKey:ServiceID"`
	EmployeeID uint      `gorm:"not null" json:"employee_id"`
	Employee   User      `gorm:"foreignKey:EmployeeID"`
	UserID     uint      `gorm:"not null" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID"`
	BranchID   uint      `gorm:"not null" json:"branch_id"`
	Branch     Branch    `gorm:"foreignKey:BranchID"`
	StartTime  time.Time `gorm:"not null" json:"start_time"`
	EndTime    time.Time `gorm:"not null" json:"end_time"`
}
