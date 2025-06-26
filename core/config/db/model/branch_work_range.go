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
	BranchID  uuid.UUID    `json:"branch_id" gorm:"type:uuid;not null;index:idx_branch_id"`
	Branch    Branch       `json:"branch" gorm:"foreignKey:BranchID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Services  []*Service   `json:"services" gorm:"many2many:branch_work_range_services;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
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

func (bwr *BranchWorkRange) AfterFind(tx *gorm.DB) error {
	var err error

	bwr.StartTime, err = lib.Utc2LocalTime(bwr.TimeZone, bwr.StartTime)
	if err != nil {
		return fmt.Errorf("branch work range (%s) failed to convert start time to local time: %w", bwr.ID, err)
	}

	bwr.EndTime, err = lib.Utc2LocalTime(bwr.TimeZone, bwr.EndTime)
	if err != nil {
		return fmt.Errorf("branch work range (%s) failed to convert end time to local time: %w", bwr.ID, err)
	}

	return nil
}

func (bwr *BranchWorkRange) BeforeCreate(tx *gorm.DB) error {
	var err error
	bwr.StartTime, err = lib.LocalTime2UTC(bwr.TimeZone, bwr.StartTime)

	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid start time %s: %w", bwr.StartTime, err))
	}

	bwr.EndTime, err = lib.LocalTime2UTC(bwr.TimeZone, bwr.EndTime)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid end time %s: %w", bwr.EndTime, err))
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

	if bwr.StartTime.Second() != 0 {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("parsing of start time generated a time with seconds, which is not allowed: %d", bwr.StartTime.Second()))
	} else if bwr.EndTime.Second() != 0 {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("parsing of end time generated a time with seconds, which is not allowed: %d", bwr.EndTime.Second()))
	}
	return nil
}

func (bwr *BranchWorkRange) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("BranchID") {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch ID cannot be changed after creation"))
	}

	if tx.Statement.Changed("StartTime") || tx.Statement.Changed("EndTime") || tx.Statement.Changed("TimeZone") {
		var err error
		bwr.StartTime, err = lib.LocalTime2UTC(bwr.TimeZone, bwr.StartTime)

		if err != nil {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid start time %s: %w", bwr.StartTime, err))
		}

		bwr.EndTime, err = lib.LocalTime2UTC(bwr.TimeZone, bwr.EndTime)
		if err != nil {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid end time %s: %w", bwr.EndTime, err))
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
