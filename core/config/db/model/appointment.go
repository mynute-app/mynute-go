package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Fifth step: Scheduling with service and employee availability.
type Appointment struct {
	gorm.Model

	ServiceID uint     `gorm:"not null;index" json:"service_id"`
	Service   *Service `gorm:"foreignKey:ServiceID;references:ID;constraint:OnDelete:CASCADE;"`

	EmployeeID uint      `gorm:"not null;index" json:"employee_id"`
	Employee   *Employee `gorm:"foreignKey:EmployeeID;references:ID;constraint:OnDelete:CASCADE;"`

	UserID uint  `gorm:"not null;index" json:"user_id"`
	User   *User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;"`

	BranchID uint    `gorm:"not null;index" json:"branch_id"`
	Branch   *Branch `gorm:"foreignKey:BranchID;references:ID;constraint:OnDelete:CASCADE;"`

	CompanyID uint     `gorm:"not null;index" json:"company_id"`
	Company   *Company `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE;" json:"company"`

	StartTime time.Time `gorm:"not null" json:"start_time"`
	EndTime   time.Time `gorm:"not null" json:"end_time"`
}

func (a *Appointment) BeforeCreate(tx *gorm.DB) error {
	if a.StartTime.Before(time.Now()) {
		return errors.New("start time must be in the future")
	}
	if err := tx.First(&a.Service, a.ServiceID).Error; err != nil {
		return err
	}
	a.EndTime = a.StartTime.Add(time.Duration(a.Service.Duration) * time.Minute)
	if a.EndTime.Before(a.StartTime) {
		return errors.New("end time must be after start time")
	}
	return nil
}
