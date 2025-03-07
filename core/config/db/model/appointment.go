package model

import (
	"time"

	"gorm.io/gorm"
)

// Fifth step: Scheduling with service and employee availability.
type Appointment struct {
	gorm.Model
	ServiceID  uint      `gorm:"not null;index;foreignKey:ServiceID;references:ID;constraint:OnDelete:CASCADE;" json:"service_id"`
	EmployeeID uint      `gorm:"not null;index;foreignKey:EmployeeID;references:ID;constraint:OnDelete:CASCADE;" json:"employee_id"`
	UserID     uint      `gorm:"not null;index;foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;" json:"user_id"`
	BranchID   uint      `gorm:"not null;index;foreignKey:BranchID;references:ID;constraint:OnDelete:CASCADE;" json:"branch_id"`
	StartTime  time.Time `gorm:"not null" json:"start_time"`
	EndTime    time.Time `gorm:"not null" json:"end_time"`
}
