package model

import (
	"agenda-kaki-go/core/lib"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkRangeBase struct {
	BaseModel
	Weekday   time.Weekday `json:"weekday" gorm:"not null"`
	StartTime time.Time    `json:"start_time" gorm:"not null;type:timestamptz"`
	EndTime   time.Time    `json:"end_time" gorm:"not null;type:timestamptz"`
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
	// if err := wr.LocalTime2UTC(); err != nil {
	// 	return err
	// }
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
	// if err := wr.LocalTime2UTC(); err != nil {
	// 	return err
	// }

	return nil
}

func (wr *WorkRangeBase) AfterFind(tx *gorm.DB) error {
	// if err := wr.Utc2LocalTime(); err != nil {
	// 	return fmt.Errorf("work range (%s) failed to convert UTC to local time (%s): %w", wr.ID, wr.TimeZone, err)
	// }

	return nil
}

func (wr *WorkRangeBase) ValidateTime() error {
	if wr.StartTime.Equal(wr.EndTime) {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("start time cannot be equal to end time"))
	}
	if wr.StartTime.After(wr.EndTime) {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("start time cannot be after end time"))
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
	// Get Timezone from start time
	stl := wr.StartTime.Location()
	etl := wr.EndTime.Location()
	if stl != etl {
		return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("start time and end time must have the same timezone"))
	}
	return stl, nil
}

func (wr *WorkRangeBase) GetTimeZoneString() (string, error) {
	loc, err := wr.GetTimeZone()
	if err != nil {
		return "", err
	}
	return loc.String(), nil
}

// LocalTime2UTC converts the start and end times of the WorkRange from local time to UTC.
// DO NOT EVER USE before the ConvertToBranchTimeZone function.
// func (wr *WorkRangeBase) LocalTime2UTC() error {
// 	var err error
// 	startTimeUTC, err := lib.LocalTime2UTC(wr.TimeZone, wr.StartTime)
// 	if err != nil {
// 		return fmt.Errorf("invalid start time %s: %w", wr.StartTime, err)
// 	}
// 	endTimeUTC, err := lib.LocalTime2UTC(wr.TimeZone, wr.EndTime)
// 	if err != nil {
// 		return fmt.Errorf("invalid end time %s: %w", wr.EndTime, err)
// 	}
// 	wr.StartTime = startTimeUTC
// 	wr.EndTime = endTimeUTC
// 	return nil
// }

// func (wr *WorkRangeBase) Utc2LocalTime() error {
// 	var err error
// 	wr.StartTime, err = lib.Utc2LocalTime(wr.TimeZone, wr.StartTime)
// 	if err != nil {
// 		return fmt.Errorf("failed to convert start time to local: %w", err)
// 	}
// 	wr.EndTime, err = lib.Utc2LocalTime(wr.TimeZone, wr.EndTime)
// 	if err != nil {
// 		return fmt.Errorf("failed to convert end time to local: %w", err)
// 	}
// 	return nil
// }

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
	var bTZ_str string
	if err := tx.
		Model(&Branch{}).
		Where("id = ?", wr.BranchID.String()).
		Select("time_zone").
		Row().Scan(&bTZ_str); err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch with ID %s not found", wr.BranchID))
		}
		return lib.Error.General.InternalError.WithError(err)
	}
	bTZ, err := lib.GetTimeZone(bTZ_str)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch with ID %s has invalid time zone %s: %w", wr.BranchID, bTZ_str, err))
	}
	startTimeInBranchTZ := wr.StartTime.In(bTZ)
	endTimeInBranchTZ := wr.EndTime.In(bTZ)
	wr.StartTime = startTimeInBranchTZ
	wr.EndTime = endTimeInBranchTZ
	return nil
}
