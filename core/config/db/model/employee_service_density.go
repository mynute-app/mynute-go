package model

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmployeeServiceDensity struct {
	BaseModel
	EmployeeID uuid.UUID `json:"employee_id" gorm:"primaryKey"`
	ServiceID  uuid.UUID `json:"service_id" gorm:"primaryKey"`
	Density    uint32     `json:"density" gorm:"not null;default:1"` // Use int32 to allow negative values for unbounded
}

const EmployeeServiceDensityTableName = "employee_service_densities"

func (EmployeeServiceDensity) TableName() string  { return EmployeeServiceDensityTableName }
func (EmployeeServiceDensity) SchemaType() string { return "tenant" }
func (EmployeeServiceDensity) Indexes() map[string]string {
	return EmployeeServiceDensityIndexes(EmployeeServiceDensityTableName)
}

func EmployeeServiceDensityIndexes(table string) map[string]string {
	return map[string]string{
		"idx_employee_service": fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_employee_service ON %s (employee_id, service_id)", table),
	}
}

func (esd *EmployeeServiceDensity) BeforeCreate(tx *gorm.DB) error {
	if esd.EmployeeID == uuid.Nil || esd.ServiceID == uuid.Nil {
		return fmt.Errorf("employee_id and service_id must be set before creating an employee service density")
	}

	var count int64
	if err := tx.Model(&EmployeeServiceDensity{}).Where("employee_id = ? AND service_id = ?", esd.EmployeeID.String(), esd.ServiceID.String()).Count(&count).Error; err != nil {
		return fmt.Errorf("error checking existing employee service density: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("employee service density for employee %d and service %d already exists", esd.EmployeeID, esd.ServiceID)
	}

	var employeeTotalServiceDensity uint32
	if err := tx.Model(&Employee{}).Where("id = ?", esd.EmployeeID).Pluck("total_service_density", &employeeTotalServiceDensity).Error; err != nil {
		return fmt.Errorf("error loading employee density: %w", err)
	}

	if esd.Density > employeeTotalServiceDensity {
		return fmt.Errorf("employee service density %d exceeds employee maximum density %d", esd.Density, employeeTotalServiceDensity)
	}

	return nil
}

func (esd *EmployeeServiceDensity) BeforeUpdate(tx *gorm.DB) error {
	if esd.EmployeeID != uuid.Nil || esd.ServiceID != uuid.Nil {
		return fmt.Errorf("can not update employee_id or service_id for an existing employee service density")
	}

	if esd.Density > 0 {
		var employeeTotalServiceDensity uint32
		if err := tx.Model(&Employee{}).Where("id = ?", esd.EmployeeID).Pluck("total_service_density", &employeeTotalServiceDensity).Error; err != nil {
			return fmt.Errorf("error loading employee density: %w", err)
		}
		if esd.Density > employeeTotalServiceDensity {
			return fmt.Errorf("employee service density %d exceeds employee maximum density %d", esd.Density, employeeTotalServiceDensity)
		}
	}

	return nil
}
