package model

import (
	"fmt"
	"mynute-go/core/src/lib"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkRangeBase struct {
	BaseModel
	Weekday   time.Weekday `json:"weekday" gorm:"not null"`
	StartTime time.Time    `json:"start_time" gorm:"not null;type:timestamptz"`
	EndTime   time.Time    `json:"end_time" gorm:"not null;type:timestamptz"`
	TimeZone  string       `json:"time_zone" gorm:"not null;type:varchar(255)" validate:"required,myTimezoneValidation"` // Time zone of the work range, e.g., "America/New_York"
	BranchID  uuid.UUID    `json:"branch_id" gorm:"type:uuid;not null;index:idx_branch_id"`
	Branch    Branch       `json:"branch" gorm:"foreignKey:BranchID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" validate:"-"`
}

func (wr *WorkRangeBase) BeforeCreate(tx *gorm.DB) error {
	if err := lib.MyCustomStructValidator(wr); err != nil {
		return err
	}
	if err := wr.ValidateTime(); err != nil {
		return err
	}
	if err := wr.ConvertToBranchTimeZone(tx); err != nil {
		return err
	}
	return nil
}

func (wr *WorkRangeBase) BeforeUpdate(tx *gorm.DB) error {
	if wr.StartTime.IsZero() || wr.EndTime.IsZero() {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("start time and end time must be present when updating a work range"))
	} else if wr.Weekday.String() == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("weekday must be present when updating a work range"))
	}
	if err := wr.ValidateTime(); err != nil {
		return err
	}
	if err := wr.ConvertToBranchTimeZone(tx); err != nil {
		return err
	}
	return nil
}

func (wr *WorkRangeBase) AfterFind(tx *gorm.DB) error {
	loc, err := wr.GetTimeZone()
	if err != nil {
		return err
	}
	wr.StartTime = wr.StartTime.In(loc)
	wr.EndTime = wr.EndTime.In(loc)
	return nil
}

func (wr *WorkRangeBase) ValidateTime() error {
	if wr.StartTime.Equal(wr.EndTime) {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("start time cannot be equal to end time"))
	}
	// If end time is not at midnight, start time cannot be after end time
	if wr.StartTime.After(wr.EndTime) && !(wr.EndTime.Hour() == 0 && wr.EndTime.Minute() == 0) {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("start time (%s) cannot be after end time (%s)", wr.StartTime, wr.EndTime))
	}
	if wr.Weekday < 0 || wr.Weekday > 6 {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid weekday %d", wr.Weekday))
	}
	if wr.StartTime.Second() != 0 {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("start time has seconds: %d", wr.StartTime.Second()))
	}
	if wr.EndTime.Second() != 0 {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("end time has seconds: %d", wr.EndTime.Second()))
	}
	if wr.StartTime.Second()%15 != 0 {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("start time seconds must be zero or divisible by 15, got %d", wr.StartTime.Second()))
	}
	if wr.EndTime.Second()%15 != 0 {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("end time seconds must be zero or divisible by 15, got %d", wr.EndTime.Second()))
	}
	return nil
}

func (wr *WorkRangeBase) GetTimeZone() (*time.Location, error) {
	loc, err := lib.GetTimeZone(wr.TimeZone)
	if err != nil {
		return nil, err
	}
	return loc, nil
}

func (wr *WorkRangeBase) GetTimeZoneString() (string, error) {
	loc, err := wr.GetTimeZone()
	if err != nil {
		return "", err
	}
	return loc.String(), nil
}

func (wr *WorkRangeBase) Overlaps(other *WorkRangeBase) (bool, error) {
	if wr.Weekday != other.Weekday {
		return false, nil
	}
	loc1, err := wr.GetTimeZone()
	if err != nil {
		return false, err
	}
	loc2, err := other.GetTimeZone()
	if err != nil {
		return false, err
	}
	return lib.TimeRangeOverlaps(wr.StartTime, wr.EndTime, loc1, other.StartTime, other.EndTime, loc2), nil
}

// ConvertToBranchTimeZone converts the start and end times of the BranchWorkRange to the branch's time zone.
// It retrieves the branch's time zone from the database and applies it to the start and end times.
// Necessary to maintain consistency with Branch time zone as it wouldn't make sense to have a BranchWorkRange in a different time zone than the Branch itself.
func (wr *WorkRangeBase) ConvertToBranchTimeZone(tx *gorm.DB) error {
	var bTZ string
	if err := tx.
		Model(&Branch{}).
		Where("id = ?", wr.BranchID.String()).
		Select("time_zone").
		Row().Scan(&bTZ); err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch with ID %s not found", wr.BranchID))
		}
		return lib.Error.General.InternalError.WithError(err)
	}
	if bTZ == wr.TimeZone {
		return nil
	}
	bLoc, err := lib.GetTimeZone(bTZ)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch with ID %s has invalid time zone %s: %w", wr.BranchID, bTZ, err))
	}
	wr.StartTime = wr.StartTime.In(bLoc)
	wr.EndTime = wr.EndTime.In(bLoc)
	wr.TimeZone = bTZ
	return nil
}

