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

func GetTimeZone(tz string) (*time.Location, error) {
	if tz == "" {
		return nil, fmt.Errorf("time_zone cannot be empty")
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, fmt.Errorf("invalid time_zone '%s': %w", tz, err)
	}

	return loc, nil
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
		return time.Time{}, fmt.Errorf("time_zone cannot be empty")
	}

	// Ignore if already in UTC, as this function assumes localTime is in the specified time_zone.
	if localTime.Location() == time.UTC {
		return localTime, nil
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time_zone '%s': %w", tz, err)
	}

	// ========================================================================
	// --- Start of Timezone Conversion Logic ---

	// This logic correctly prepares a local time for database storage as UTC.

	// 1. Use a MODERN base date to build a clean local timestamp.
	//    This defensively uses ONLY the Hour and Minute from the input `localTime`,
	//    cleaning any unwanted seconds. It also ensures the correct modern time_zone
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
// in the specified time_zone on the conceptual "zero date" (0001-01-01) while
// avoiding the LTM (Local Mean Time) issues that can occur with historical timezones
// which is the case for (0001-01-01) date.
func Utc2LocalTime(tz string, utcTime time.Time) (time.Time, error) {
	if utcTime.Second() != 0 || utcTime.Nanosecond() != 0 {
		return time.Time{}, fmt.Errorf("time must have zero seconds and nanoseconds, got %d seconds and %d nanoseconds", utcTime.Second(), utcTime.Nanosecond())
	}
	if tz == "" {
		return time.Time{}, fmt.Errorf("time_zone cannot be empty")
	}
	// Ignore if not in UTC, as this function assumes utcTime is in UTC.
	if utcTime.Location() != time.UTC {
		return utcTime, nil
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time_zone '%s': %w", tz, err)
	}

	// ========================================================================
	// --- Start of Timezone Conversion Logic with Documentation ---

	// The following logic correctly converts a UTC time-of-day from the database
	// into its proper local time-of-day representation.

	// 1. Get a pure UTC instant from the input time to establish a reliable baseline.
	//    This removes any location info (like the server's time.Local) that the
	//    database driver might have attached.
	utcInstant := utcTime.In(time.UTC)

	// 2. Use a MODERN base date (e.g., year 2000) to perform the time_zone conversion.
	//    This is CRITICAL to get the correct, modern time_zone offset.
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

		// This helper function performs inclusive check
		TimeDurationOverlap := func(as, ae, bs, be time.Duration) bool {
			// They do NOT overlap if A ends before when B starts, OR if B ends
			// before when A starts. Otherwise, they must overlap.
			// Therefore, we return the negation of "do not overlap".
			a_is_before_b := ae <= bs
			b_is_before_a := be <= as
			return !(a_is_before_b || b_is_before_a)
		}

		overlaps := TimeDurationOverlap(aStartDur, aEndDur, bStartDur, bEndDur)

		return overlaps
	}

	// This helper function performs inclusive check
	TimeRangeOverlap := func(as, ae, bs, be time.Time) bool {

		// --- Standard Logic for Full Date-Times (Inclusive) ---
		// They do NOT overlap if A ends before when B starts, OR if B ends
		// before when A starts. Otherwise, they must overlap.
		// Therefore, we return the negation of "do not overlap".
		a_is_before_b := ae.Before(bs)
		b_is_before_a := be.Before(as)

		return !(a_is_before_b || b_is_before_a)
	}

	overlaps := TimeRangeOverlap(aStart, aEnd, bStart, bEnd)
	return overlaps
}

// TimeRangeFullyContained checks if the time range [bStart, bEnd] is fully contained
// within the time range [aStart, aEnd], supporting both full date-time and time-of-day-only comparisons.
//
// The comparison is inclusive, meaning edge-aligned ranges are considered contained.
// For example, [08:00, 12:00] fully contains [08:00, 12:00].
//
// If the times are "time-of-day only" (i.e., Year == 1), it handles overnight shifts correctly.
// Time zones can be passed to convert inputs before comparing.
//
// Parameters:
//   - aStart, aEnd: Start and end of the containing range
//   - aTZ: Optional time zone for aStart/aEnd (nil = use time as-is)
//   - bStart, bEnd: Start and end of the inner (possibly contained) range
//   - bTZ: Optional time zone for bStart/bEnd (nil = use time as-is)
//
// Returns:
//   - true if [bStart, bEnd] is entirely within [aStart, aEnd] (inclusive), false otherwise.
func TimeRangeFullyContained(aStart, aEnd time.Time, aTZ *time.Location, bStart, bEnd time.Time, bTZ *time.Location) bool {
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
		aStartDur := OnlyTime(aStart)
		aEndDur := OnlyTime(aEnd)
		bStartDur := OnlyTime(bStart)
		bEndDur := OnlyTime(bEnd)

		// Normalize overnight shifts
		normalize := func(start, end time.Duration) (time.Duration, time.Duration) {
			if end <= start {
				end += 24 * time.Hour
			}
			return start, end
		}

		aStartDur, aEndDur = normalize(aStartDur, aEndDur)
		bStartDur, bEndDur = normalize(bStartDur, bEndDur)

		return bStartDur >= aStartDur && bEndDur <= aEndDur
	}

	// Full date-time containment (inclusive)
	return !bStart.Before(aStart) && !bEnd.After(aEnd)
}

// Parse_HHMM_To_Time parses either a "HH:MM" string using the given time_zone location
// or a full ISO datetime string with offset like "2020-01-01T08:00:00-03:00".
// It always returns a time.Time with date 2020-01-01 and the correct time zone location.
func Parse_HHMM_To_Time(input string, tz string) (time.Time, error) {
	loc, err := GetTimeZone(tz)
	if err != nil {
		return time.Time{}, err
	}

	// Case 1: full ISO datetime with offset
	if t, err := time.Parse(time.RFC3339, input); err == nil {
		// Convert to desired location
		t = t.In(loc)
		return time.Date(2020, 1, 1, t.Hour(), t.Minute(), 0, 0, loc), nil
	}

	// Case 2: HH:MM
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

	return time.Date(2020, 1, 1, hour, minute, 0, 0, loc), nil
}
