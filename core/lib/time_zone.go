package lib

import (
	"errors"
	"fmt"
	"time"
)

type corrected_time struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	TimeZone  string    `json:"timezone"`
}

func UTC_with_Zero_YMD_Date(TimeZone string, StartTime time.Time, EndTime time.Time) (*corrected_time,error) {
	if TimeZone == "" {
		return nil, Error.General.BadRequest.WithError(errors.New("time zone cannot be empty"))
	}
	loc, err := time.LoadLocation(TimeZone)
	if err != nil {
		return nil, Error.General.BadRequest.WithError(fmt.Errorf("invalid time zone %s: %w", TimeZone, err))
	}
	start := time.Date(1, 1, 1, StartTime.Hour(), StartTime.Minute(), StartTime.Second(), 0, loc)
	end := time.Date(1, 1, 1, EndTime.Hour(), EndTime.Minute(), EndTime.Second(), 0, loc)

	correct_time := &corrected_time{
		StartTime: start.UTC(),
		EndTime:   end.UTC(),
		TimeZone:  TimeZone,
	}
	return correct_time, nil
}
func TimeRangeOverlaps(aStart, aEnd time.Time, aTZ *time.Location, bStart, bEnd time.Time, bTZ *time.Location) bool {
	aStart = aStart.In(aTZ)
	aEnd = aEnd.In(aTZ)
	bStart = bStart.In(bTZ)
	bEnd = bEnd.In(bTZ)

	aStartBeforeOrEqualBStart := aStart.Before(bStart) || aStart.Equal(bStart)
	aEndAfterOrEqualBEnd := aEnd.After(bEnd) || aEnd.Equal(bEnd)
	bStartBeforeOrEqualAStart := bStart.Before(aStart) || bStart.Equal(aStart)
	bEndAfterOrEqualAEnd := bEnd.After(aEnd) || bEnd.Equal(aEnd)

	aEqualsB := aStart.Equal(bStart) && aEnd.Equal(bEnd)
	aContainsB := aStartBeforeOrEqualBStart && aEndAfterOrEqualBEnd
	bContainsA := bStartBeforeOrEqualAStart && bEndAfterOrEqualAEnd
	aContainsBStart := aStartBeforeOrEqualBStart && bEndAfterOrEqualAEnd
	bContainsAStart := bStartBeforeOrEqualAStart && aEndAfterOrEqualBEnd

	return aEqualsB || aContainsB || bContainsA || aContainsBStart || bContainsAStart
}

