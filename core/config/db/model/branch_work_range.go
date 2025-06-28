package model

import (
	"agenda-kaki-go/core/lib"
	"fmt"
	"gorm.io/gorm"
)

type BranchWorkRange struct {
	WorkRangeBase
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

func (bwr *BranchWorkRange) BeforeUpdate(tx *gorm.DB) error {
	if err := bwr.WorkRangeBase.BeforeUpdate(tx); err != nil {
		return err
	}

	branch := &Branch{BaseModel: BaseModel{ID: bwr.BranchID}}
	if err := branch.ValidateBranchWorkRangeTime(tx, bwr); err != nil {
		return err
	}

	return nil
}

func (bwr *BranchWorkRange) Overlaps(other *BranchWorkRange) (bool, error) {
	if bwr.Weekday != other.Weekday || bwr.BranchID != other.BranchID {
		return false, nil
	}
	loc, err := bwr.GetTimeZone()
	if err != nil {
		return false, err
	}
	loc2, err := other.GetTimeZone()
	if err != nil {
		return false, err
	}

	return lib.TimeRangeOverlaps(bwr.StartTime, bwr.EndTime, loc, other.StartTime, other.EndTime, loc2), nil
}