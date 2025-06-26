package lib

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TimeRangeResult encapsula o resultado corrigido com horário UTC
type TimeRangeResult struct {
	StartTime time.Time
	EndTime   time.Time
	TimeZone  string
}

// Add this new function alongside your existing DatabaseUtc2LocalTime.

// LocalTime2UTC takes a local time-of-day (represented by a time.Time object)
// and converts it into its canonical UTC representation on the zero date (0001-01-01)
// while avoiding the LTM (Local Mean Time) issues that can occur with historical timezones
// which is the case for (0001-01-01) date.
func LocalTime2UTC(tz string, localTime time.Time) (time.Time, error) {
	if localTime.Second() != 0 || localTime.Nanosecond() != 0 {
		return time.Time{}, fmt.Errorf("time must have zero seconds and nanoseconds, got %d seconds and %d nanoseconds", localTime.Second(), localTime.Nanosecond())
	}

	if tz == "" {
		return time.Time{}, fmt.Errorf("timezone cannot be empty")
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timezone '%s': %w", tz, err)
	}

	// ========================================================================
	// --- Start of Timezone Conversion Logic ---

	// This logic correctly prepares a local time for database storage as UTC.

	// 1. Use a MODERN base date to build a clean local timestamp.
	//    This defensively uses ONLY the Hour and Minute from the input `localTime`,
	//    cleaning any unwanted seconds. It also ensures the correct modern timezone
	//    offset is used, avoiding historical LMT issues.
	modernTimestampInLocal := time.Date(2000, 1, 1,
		localTime.Hour(), localTime.Minute(), 0, 0, // Explicitly set seconds and nanoseconds to 0
		loc)

	// 2. Convert this clean, modern local timestamp to its UTC equivalent.
	modernTimestampInUTC := modernTimestampInLocal.In(time.UTC)

	// 3. Finally, transplant the clean UTC time components back onto our conceptual
	//    "zero date" (year 1). This is the canonical format for database storage.
	finalUTCTime := time.Date(1, 1, 1,
		modernTimestampInUTC.Hour(), modernTimestampInUTC.Minute(), modernTimestampInUTC.Second(), modernTimestampInUTC.Nanosecond(),
		time.UTC)

	// --- End of Timezone Conversion Logic ---
	// ========================================================================

	return finalUTCTime, nil
}

// DUtc2LocalTime takes a time-of-day (represented by a time.Time object) and
// correctly converts it to the local wall-clock time
// in the specified timezone on the conceptual "zero date" (0001-01-01) while
// avoiding the LTM (Local Mean Time) issues that can occur with historical timezones
// which is the case for (0001-01-01) date.
func Utc2LocalTime(tz string, dbTime time.Time) (time.Time, error) {
	if dbTime.Second() != 0 || dbTime.Nanosecond() != 0 {
		return time.Time{}, fmt.Errorf("time must have zero seconds and nanoseconds, got %d seconds and %d nanoseconds", dbTime.Second(), dbTime.Nanosecond())
	}
	if tz == "" {
		return time.Time{}, fmt.Errorf("timezone cannot be empty")
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timezone '%s': %w", tz, err)
	}

	// ========================================================================
	// --- Start of Timezone Conversion Logic with Documentation ---

	// The following logic correctly converts a UTC time-of-day from the database
	// into its proper local time-of-day representation.

	// 1. Get a pure UTC instant from the input time to establish a reliable baseline.
	//    This removes any location info (like the server's time.Local) that the
	//    database driver might have attached.
	utcInstant := dbTime.In(time.UTC)

	// 2. Use a MODERN base date (e.g., year 2000) to perform the timezone conversion.
	//    This is CRITICAL to get the correct, modern timezone offset.
	//    WHY: Before standardized timezones (circa 1914 in Brazil), cities used their
	//    own Local Mean Time (LMT) based on longitude. Go's time library is
	//    historically accurate. If we ask it to convert a time in the year 1 for
	//    "America/Sao_Paulo", it will correctly use the LMT offset (-03:06),
	//    which is not what we want for a recurring modern schedule. Using a modern
	//    date forces Go to use the standard offset (e.g., -03:00).
	modernTimestampInUTC := time.Date(2000, 1, 1,
		utcInstant.Hour(), utcInstant.Minute(), utcInstant.Second(), utcInstant.Nanosecond(),
		time.UTC)

	// 3. Convert this modern UTC timestamp to the target location. The result
	//    will now have the correct modern offset.
	modernTimestampInLocal := modernTimestampInUTC.In(loc)

	// 4. Finally, transplant the correct local time components (hour, minute, etc.)
	//    back onto our conceptual "zero date" (year 1). This gives us the final,
	//    clean time-of-day object for application use (e.g., 09:00:00-03:00).
	finalLocalTime := time.Date(1, 1, 1,
		modernTimestampInLocal.Hour(), modernTimestampInLocal.Minute(), modernTimestampInLocal.Second(), modernTimestampInLocal.Nanosecond(),
		loc)

	// --- End of Timezone Conversion Logic ---
	// ========================================================================

	return finalLocalTime, nil
}

func UTC_with_Zero_YMD_Date(tz string, start time.Time, end time.Time) (TimeRangeResult, error) {
	if tz == "" {
		return TimeRangeResult{}, fmt.Errorf("timezone cannot be empty")
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return TimeRangeResult{}, fmt.Errorf("invalid timezone: %w", err)
	}

	CorrectTimeRange := TimeRangeResult{}
	CorrectTimeRange.TimeZone = tz

	// 1. Take the input time (which represents a UTC instant, regardless of its initial Location)
	//    and create a canonical UTC time.Time object on the zero date.
	//    This strips away any arbitrary date (like today's date) or loc (like time.Local)
	//    that the driver might have added.
	//    Also, the canonical UTC time.Time created avoid the
	utcTime := start.In(time.UTC)
	canonicalUTCStart := time.Date(1, 1, 1, utcTime.Hour(), utcTime.Minute(), utcTime.Second(), utcTime.Nanosecond(), time.UTC)

	// 2. Convert this canonical UTC instant to the target local timezone.
	//    This is the key conversion. `In()` correctly calculates the local wall-clock time.
	CorrectTimeRange.StartTime = canonicalUTCStart.In(loc)

	// 3. Repeat for EndTime
	utcTime = end.In(time.UTC)
	canonicalUTCEnd := time.Date(1, 1, 1, utcTime.Hour(), utcTime.Minute(), utcTime.Second(), utcTime.Nanosecond(), time.UTC)
	CorrectTimeRange.EndTime = canonicalUTCEnd.In(loc)

	return CorrectTimeRange, nil
}

// aplica corretamente o timezone informado sobre os horários e converte para UTC com data base 0001-01-01
// func UTC_with_Zero_YMD_Date(tz string, start time.Time, end time.Time) (TimeRangeResult, error) {
// 	if tz == "" {
// 		return TimeRangeResult{}, fmt.Errorf("timezone cannot be empty")
// 	}
// 	location, err := time.LoadLocation(tz)
// 	if err != nil {
// 		return TimeRangeResult{}, fmt.Errorf("invalid timezone: %w", err)
// 	}

// 	startInTZ := time.Date(1, 1, 1, start.Hour(), start.Minute(), start.Second(), 0, location)
// 	endInTZ := time.Date(1, 1, 1, end.Hour(), end.Minute(), end.Second(), 0, location)

// 	return TimeRangeResult{
// 		StartTime: startInTZ.UTC(),
// 		EndTime:   endInTZ.UTC(),
// 		TimeZone:  tz,
// 	}, nil
// }

// OnlyTime extrai apenas a hora do dia em forma de Duration
func OnlyTime(t time.Time) time.Duration {
	return time.Duration(t.Hour())*time.Hour +
		time.Duration(t.Minute())*time.Minute +
		time.Duration(t.Second())*time.Second
}

// TimeRangeOverlaps correctly checks if two time ranges, A and B, have any inclusive overlap.
// "Inclusive" means that intervals touching at the edges (e.g., [8-10] and [10-12])
// are considered to be overlapping.
// It handles both full date-times and time-of-day comparisons, including overnight shifts.
func TimeRangeOverlaps(aStart, aEnd time.Time, aTZ *time.Location, bStart, bEnd time.Time, bTZ *time.Location) bool {
	// ... (The defensive .In(aTZ) logic remains the same) ...
	if aTZ != nil {
		aStart = aStart.In(aTZ).In(time.UTC)
		aEnd = aEnd.In(aTZ).In(time.UTC)
	}
	if bTZ != nil {
		bStart = bStart.In(bTZ).In(time.UTC)
		bEnd = bEnd.In(bTZ).In(time.UTC)
	}

	isOnlyTime := aStart.Year() == 1 && bStart.Year() == 1

	if isOnlyTime {
		// --- Logic for Time-of-Day Only, handles overnight shifts ---
		aStartDur := OnlyTime(aStart)
		aEndDur := OnlyTime(aEnd)
		bStartDur := OnlyTime(bStart)
		bEndDur := OnlyTime(bEnd)

		day := 24 * time.Hour
		if aEndDur <= aStartDur {
			aEndDur += day
		}
		if bEndDur <= bStartDur {
			bEndDur += day
		}

		// This helper function now performs the correct inclusive check
		doTheyOverlap := func(as, ae, bs, be time.Duration) bool {
			// They do NOT overlap if A is entirely before B, OR B is entirely before A.
			// Otherwise, they must overlap.
			a_is_before_b := ae <= bs
			b_is_before_a := be <= as
			return !(a_is_before_b || b_is_before_a)
		}

		overlaps := doTheyOverlap(aStartDur, aEndDur, bStartDur, bEndDur)
		overlaps = overlaps || doTheyOverlap(aStartDur, aEndDur, bStartDur+day, bEndDur+day)
		overlaps = overlaps || doTheyOverlap(aStartDur+day, aEndDur+day, bStartDur, bEndDur)

		return overlaps
	}

	// --- Standard Logic for Full Date-Times (Inclusive) ---
	// They do NOT overlap if A ends on or before B starts, OR if B ends on or before A starts.
	// We return the negation of "do not overlap".
	isABeforeB := aEnd.Before(bStart) || aEnd.Equal(bStart)
	isBBeforeA := bEnd.Before(aStart) || bEnd.Equal(aStart)

	return !(isABeforeB || isBBeforeA)
}

// ParseTimeHHMMWithDateBase parses a "HH:MM" string using the given timezone location
// and returns a time.Time with date 0001-01-01.
func ParseTimeHHMMWithDateBase(input string, loc *time.Location) (time.Time, error) {
	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid time format: %s", input)
	}
	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid hour in time: %s", input)
	}
	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid minute in time: %s", input)
	}
	return time.Date(1, 1, 1, hour, minute, 0, 0, loc), nil
}
