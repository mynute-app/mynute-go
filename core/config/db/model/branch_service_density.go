package model

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BranchServiceDensity struct {
	BaseModel
	BranchID  uuid.UUID `json:"branch_id" gorm:"primaryKey"`
	Branch    Branch    `json:"branch" gorm:"foreignKey:BranchID;references:ID;constraint:OnDelete:CASCADE;"`
	ServiceID uuid.UUID `json:"service_id" gorm:"primaryKey"`
	Service   Service   `json:"service" gorm:"foreignKey:ServiceID;references:ID;constraint:OnDelete:CASCADE;"`
	Density   int32     `json:"density" gorm:"not null;default:1"` // Use int32 to allow negative values for unbounded
}

const BranchServiceDensityTableName = "branch_service_densities"

func (BranchServiceDensity) TableName() string  { return BranchServiceDensityTableName }
func (BranchServiceDensity) SchemaType() string { return "tenant" }
func (BranchServiceDensity) Indexes() map[string]string {
	return BranchServiceDensityIndexes(BranchServiceDensityTableName)
}

func BranchServiceDensityIndexes(table string) map[string]string {
	return map[string]string{
		"idx_branch_service": fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_branch_service ON %s (branch_id, service_id)", table),
	}
}

func (bsd *BranchServiceDensity) BeforeCreate(tx *gorm.DB) error {
	if bsd.BranchID == uuid.Nil || bsd.ServiceID == uuid.Nil {
		return fmt.Errorf("branch_id and service_id must be set before creating a branch service density")
	}

	var count int64
	if err := tx.Model(&BranchServiceDensity{}).Where("branch_id = ? AND service_id = ?", bsd.BranchID, bsd.ServiceID).Count(&count).Error; err != nil {
		return fmt.Errorf("error checking existing branch service density: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("branch service density for branch %d and service %d already exists", bsd.BranchID, bsd.ServiceID)
	}

	var branchTotalServiceDensity int32
	if err := tx.Model(&Branch{}).Where("id = ?", bsd.BranchID).Pluck("total_service_density", &branchTotalServiceDensity).Error; err != nil {
		return fmt.Errorf("error loading branch density: %w", err)
	}

	if branchTotalServiceDensity >= 0 && bsd.Density > branchTotalServiceDensity {
		return fmt.Errorf("branch service density %d exceeds branch maximum density %d", bsd.Density, branchTotalServiceDensity)
	}

	return nil
}

func (bsd *BranchServiceDensity) BeforeUpdate(tx *gorm.DB) error {
	if bsd.BranchID != uuid.Nil || bsd.ServiceID != uuid.Nil {
		return fmt.Errorf("can not update branch_id or service_id for an existing branch service density")
	}

	if bsd.Density > 0 {
		var branchTotalServiceDensity int32
		if err := tx.Model(&Branch{}).Where("id = ?", bsd.BranchID).Pluck("total_service_density", &branchTotalServiceDensity).Error; err != nil {
			return fmt.Errorf("error loading branch density: %w", err)
		}
		if branchTotalServiceDensity >= 0 && bsd.Density > branchTotalServiceDensity {
			return fmt.Errorf("branch service density %d exceeds branch maximum density %d", bsd.Density, branchTotalServiceDensity)
		}
	}

	return nil
}
