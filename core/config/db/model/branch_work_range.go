package model

import (
	"agenda-kaki-go/core/lib"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BranchWorkRange struct {
	BaseModel
	Weekday   time.Weekday `json:"weekday" gorm:"not null"`
	StartTime time.Time    `json:"start_time" gorm:"type:time;not null"`
	EndTime   time.Time    `json:"end_time" gorm:"type:time;not null"`
	TimeZone  string       `json:"timezone" gorm:"type:varchar(100);not null"`

	BranchID uuid.UUID `json:"branch_id" gorm:"type:uuid;not null;index"`
	Branch   Branch    `json:"branch" gorm:"foreignKey:BranchID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	Services []*Service `json:"services" gorm:"many2many:branch_work_range_services;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
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
	if corrected_time, err := lib.UTC_with_Zero_YMD_Date(bwr.TimeZone, bwr.StartTime, bwr.EndTime); err != nil {
		return err
	} else {
		bwr.StartTime = corrected_time.StartTime
		bwr.EndTime = corrected_time.EndTime
		bwr.TimeZone = corrected_time.TimeZone
	}
	if bwr.StartTime.Equal(bwr.EndTime) {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("start time cannot be equal to end time"))
	}
	if bwr.StartTime.After(bwr.EndTime) {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("start time cannot be after end time"))
	}
	if bwr.Weekday < 0 || bwr.Weekday > 6 {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid weekday %d", bwr.Weekday))
	}
	return nil
}

func (bwr *BranchWorkRange) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("BranchID") {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch ID cannot be changed after creation"))
	}

	if tx.Statement.Changed("StartTime") || tx.Statement.Changed("EndTime") || tx.Statement.Changed("TimeZone") {
		if corrected_time, err := lib.UTC_with_Zero_YMD_Date(bwr.TimeZone, bwr.StartTime, bwr.EndTime); err != nil {
			return err
		} else {
			bwr.StartTime = corrected_time.StartTime
			bwr.EndTime = corrected_time.EndTime
			bwr.TimeZone = corrected_time.TimeZone
		}

		if bwr.StartTime.Equal(bwr.EndTime) {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("start time cannot be equal to end time"))
		} else if bwr.StartTime.After(bwr.EndTime) {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("start time cannot be after end time"))
		}
	}

	var branch *Branch
	if err := tx.First(&branch, "id = ?", bwr.BranchID.String()).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch with ID %s not found", bwr.BranchID))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	if tx.Statement.Changed("Weekday") || tx.Statement.Changed("StartTime") || tx.Statement.Changed("EndTime") || tx.Statement.Changed("TimeZone") {
		if err := branch.ValidateBranchWorkRangeTime(tx, bwr); err != nil {
			return err
		}
	}

	return nil
}

func (bwr *BranchWorkRange) Overlaps(other *BranchWorkRange) (bool, error) {
	if bwr.Weekday != other.Weekday || bwr.BranchID != other.BranchID {
		return false, nil
	}

	loc, err := time.LoadLocation(bwr.TimeZone)
	if err != nil {
		return false, err
	}
	loc2, err := time.LoadLocation(other.TimeZone)
	if err != nil {
		return false, err
	}

	return lib.TimeRangeOverlaps(bwr.StartTime, bwr.EndTime, loc, other.StartTime, other.EndTime, loc2), nil
}