package model

import (
	"fmt"
	"gorm.io/gorm"
)

type BranchWorkSchedule struct {
	WorkRanges []BranchWorkRange `json:"branch_work_ranges"`
}

type BranchWorkRange struct {
	WorkRangeBase
	Services []*Service `gorm:"many2many:branch_work_range_services;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"services" validate:"-"`
}

const BranchWorkRangeTableName = "branch_work_ranges"

func (BranchWorkRange) TableName() string  { return BranchWorkRangeTableName }
func (BranchWorkRange) SchemaType() string { return "tenant" }
func (BranchWorkRange) Indexes() map[string]string {
	return BranchWorkRangeIndexes(BranchWorkRangeTableName)
}

func BranchWorkRangeIndexes(table string) map[string]string {
	return map[string]string{
		"idx_branch_weekday":            fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_branch_weekday ON %s (branch_id, weekday)", table),
		"idx_branch_start_time":         fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_branch_start_time ON %s (branch_id, start_time)", table),
		"idx_branch_weekday_start_time": fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_branch_weekday_start_time ON %s (branch_id, weekday, start_time)", table),
	}
}

func (bwr *BranchWorkRange) BeforeCreate(tx *gorm.DB) error {
	if err := bwr.WorkRangeBase.BeforeCreate(tx); err != nil {
		return err
	}

	branch := &Branch{BaseModel: BaseModel{ID: bwr.BranchID}}
	if err := branch.ValidateBranchWorkRangeTime(tx, bwr); err != nil {
		return err
	}

	if err := branch.HasServices(tx, bwr.Services); err != nil {
		return err
	}

	return nil
}

func (bwr *BranchWorkRange) BeforeUpdate(tx *gorm.DB) error {
	var old BranchWorkRange
	if err := tx.Model(&old).Where("id = ?", bwr.ID).First(&old).Error; err != nil {
		return fmt.Errorf("error fetching old work range: %w", err)
	}

	if bwr.BranchID != old.BranchID {
		return fmt.Errorf("branch ID cannot be changed after creation")
	}

	if err := bwr.WorkRangeBase.BeforeUpdate(tx); err != nil {
		return err
	}

	branch := &Branch{BaseModel: BaseModel{ID: bwr.BranchID}}
	if err := branch.ValidateBranchWorkRangeTime(tx, bwr); err != nil {
		return err
	}

	return nil
}
