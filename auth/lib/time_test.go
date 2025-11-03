package lib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetTimestampVersion(t *testing.T) {
	t.Run("should return timestamp in correct format", func(t *testing.T) {
		result := GetTimestampVersion()

		assert.Equal(t, 14, len(result), "Timestamp should be 14 characters (YYYYMMDDHHMMSS)")
	})

	t.Run("should return valid timestamp that can be parsed", func(t *testing.T) {
		result := GetTimestampVersion()

		// Try to parse the timestamp
		_, err := time.Parse("20060102150405", result)

		assert.NoError(t, err)
	})

	t.Run("should return different timestamps when called sequentially", func(t *testing.T) {
		first := GetTimestampVersion()
		time.Sleep(1 * time.Second)
		second := GetTimestampVersion()

		assert.NotEqual(t, first, second)
	})
}

func TestGetTimeZone(t *testing.T) {
	t.Run("should load valid timezone", func(t *testing.T) {
		validTimezones := []string{
			"America/New_York",
			"Europe/London",
			"Asia/Tokyo",
			"UTC",
			"America/Sao_Paulo",
		}

		for _, tz := range validTimezones {
			loc, err := GetTimeZone(tz)

			assert.NoError(t, err, "Timezone '%s' should be valid", tz)
			assert.NotNil(t, loc)
		}
	})

	t.Run("should return error for empty timezone", func(t *testing.T) {
		loc, err := GetTimeZone("")

		assert.Error(t, err)
		assert.Nil(t, loc)
		assert.Contains(t, err.Error(), "time_zone cannot be empty")
	})

	t.Run("should return error for invalid timezone", func(t *testing.T) {
		loc, err := GetTimeZone("Invalid/Timezone")

		assert.Error(t, err)
		assert.Nil(t, loc)
		assert.Contains(t, err.Error(), "invalid time_zone")
	})

	t.Run("should return correct location object", func(t *testing.T) {
		loc, err := GetTimeZone("America/New_York")

		assert.NoError(t, err)
		assert.NotNil(t, loc)
		assert.Equal(t, "America/New_York", loc.String())
	})
}

func TestLocalTime2UTC(t *testing.T) {
	t.Run("should convert local time to UTC", func(t *testing.T) {
		// Create a local time (e.g., 14:30 in zero year)
		localTime := time.Date(1, 1, 1, 14, 30, 0, 0, time.UTC)

		utcTime, err := LocalTime2UTC("America/New_York", localTime)

		assert.NoError(t, err)
		assert.NotNil(t, utcTime)
	})

	t.Run("should return error for empty timezone", func(t *testing.T) {
		localTime := time.Date(1, 1, 1, 14, 30, 0, 0, time.UTC)

		utcTime, err := LocalTime2UTC("", localTime)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "time_zone cannot be empty")
		assert.True(t, utcTime.IsZero())
	})

	t.Run("should return error for time with non-zero seconds", func(t *testing.T) {
		localTime := time.Date(1, 1, 1, 14, 30, 15, 0, time.UTC)

		utcTime, err := LocalTime2UTC("America/New_York", localTime)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "time must have zero seconds and nanoseconds")
		assert.True(t, utcTime.IsZero())
	})

	t.Run("should return error for time with non-zero nanoseconds", func(t *testing.T) {
		localTime := time.Date(1, 1, 1, 14, 30, 0, 123456, time.UTC)

		utcTime, err := LocalTime2UTC("America/New_York", localTime)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "time must have zero seconds and nanoseconds")
		assert.True(t, utcTime.IsZero())
	})

	t.Run("should accept time with zero seconds and nanoseconds", func(t *testing.T) {
		localTime := time.Date(1, 1, 1, 9, 0, 0, 0, time.UTC)

		utcTime, err := LocalTime2UTC("America/Sao_Paulo", localTime)

		assert.NoError(t, err)
		assert.NotNil(t, utcTime)
		assert.Equal(t, 0, utcTime.Second())
		assert.Equal(t, 0, utcTime.Nanosecond())
	})
}

func TestGenerateRandomEmail(t *testing.T) {
	t.Run("should generate email with correct format", func(t *testing.T) {
		name := "user"

		result := GenerateRandomEmail(name)

		assert.Contains(t, result, "test_")
		assert.Contains(t, result, "_email_")
		assert.Contains(t, result, "@gmail.com")
		assert.Contains(t, result, name)
	})

	t.Run("should generate different emails", func(t *testing.T) {
		name := "testuser"
		emails := make(map[string]bool)

		for i := 0; i < 50; i++ {
			email := GenerateRandomEmail(name)
			emails[email] = true
		}

		assert.Greater(t, len(emails), 1)
	})

	t.Run("should handle empty name", func(t *testing.T) {
		result := GenerateRandomEmail("")

		assert.Contains(t, result, "@gmail.com")
		assert.Contains(t, result, "test_")
		assert.Contains(t, result, "_email_")
	})
}

func TestGenerateRandomInt(t *testing.T) {
	t.Run("should generate number with correct number of digits", func(t *testing.T) {
		testCases := []struct {
			length int
			min    int
			max    int
		}{
			{1, 1, 9},
			{2, 10, 99},
			{3, 100, 999},
			{4, 1000, 9999},
			{5, 10000, 99999},
		}

		for _, tc := range testCases {
			for i := 0; i < 20; i++ {
				result := GenerateRandomInt(tc.length)
				assert.GreaterOrEqual(t, result, tc.min, "Should be >= %d for length %d", tc.min, tc.length)
				assert.LessOrEqual(t, result, tc.max, "Should be <= %d for length %d", tc.max, tc.length)
			}
		}
	})

	t.Run("should generate different numbers", func(t *testing.T) {
		numbers := make(map[int]bool)

		for i := 0; i < 50; i++ {
			num := GenerateRandomInt(3)
			numbers[num] = true
		}

		assert.Greater(t, len(numbers), 1)
	})
}
