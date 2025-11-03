package lib

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomIntFromRange(t *testing.T) {
	t.Run("should generate number within range", func(t *testing.T) {
		min, max := 1, 10

		for i := 0; i < 100; i++ {
			result := GenerateRandomIntFromRange(min, max)
			assert.GreaterOrEqual(t, result, min)
			assert.LessOrEqual(t, result, max)
		}
	})

	t.Run("should generate exact value when min equals max", func(t *testing.T) {
		value := 5

		result := GenerateRandomIntFromRange(value, value)

		assert.Equal(t, value, result)
	})

	t.Run("should panic when min greater than max", func(t *testing.T) {
		assert.Panics(t, func() {
			GenerateRandomIntFromRange(10, 5)
		})
	})

	t.Run("should handle zero and negative numbers", func(t *testing.T) {
		min, max := -10, 10

		for i := 0; i < 50; i++ {
			result := GenerateRandomIntFromRange(min, max)
			assert.GreaterOrEqual(t, result, min)
			assert.LessOrEqual(t, result, max)
		}
	})

	t.Run("should include both min and max in possible results", func(t *testing.T) {
		min, max := 0, 1
		foundMin, foundMax := false, false

		// Run enough times to likely get both values
		for i := 0; i < 100 && !(foundMin && foundMax); i++ {
			result := GenerateRandomIntFromRange(min, max)
			if result == min {
				foundMin = true
			}
			if result == max {
				foundMax = true
			}
		}

		assert.True(t, foundMin, "Should generate min value")
		assert.True(t, foundMax, "Should generate max value")
	})
}

func TestGenerateRandomName(t *testing.T) {
	t.Run("should generate name with correct format", func(t *testing.T) {
		str := "User"

		result := GenerateRandomName(str)

		assert.True(t, strings.HasPrefix(result, "Test User "))
		assert.Contains(t, result, "Test")
		assert.Contains(t, result, str)
	})

	t.Run("should generate different names", func(t *testing.T) {
		str := "Company"
		names := make(map[string]bool)

		for i := 0; i < 50; i++ {
			name := GenerateRandomName(str)
			names[name] = true
		}

		// Should have generated at least some different names
		assert.Greater(t, len(names), 1)
	})

	t.Run("should handle empty string", func(t *testing.T) {
		result := GenerateRandomName("")

		assert.True(t, strings.HasPrefix(result, "Test  "))
	})

	t.Run("should handle special characters in string", func(t *testing.T) {
		str := "Test@#$"

		result := GenerateRandomName(str)

		assert.Contains(t, result, str)
		assert.True(t, strings.HasPrefix(result, "Test"))
	})
}

func TestGenerateRandomPhoneNumber(t *testing.T) {
	t.Run("should generate phone number with correct format", func(t *testing.T) {
		result := GenerateRandomPhoneNumber()

		assert.True(t, strings.HasPrefix(result, "+"))
		assert.Equal(t, 12, len(result), "Should be 12 characters total (+XXXXXXXXXXX)")
	})

	t.Run("should generate only numeric characters after plus", func(t *testing.T) {
		result := GenerateRandomPhoneNumber()

		// Remove the + prefix and check if remaining are digits
		phoneNumber := strings.TrimPrefix(result, "+")
		for _, char := range phoneNumber {
			assert.True(t, char >= '0' && char <= '9', "Character should be a digit")
		}
	})

	t.Run("should generate different phone numbers", func(t *testing.T) {
		phoneNumbers := make(map[string]bool)

		for i := 0; i < 50; i++ {
			phone := GenerateRandomPhoneNumber()
			phoneNumbers[phone] = true
		}

		// Should have generated at least some different numbers
		assert.Greater(t, len(phoneNumbers), 1)
	})
}

func TestGenerateRandomString(t *testing.T) {
	t.Run("should generate string of correct length", func(t *testing.T) {
		lengths := []int{5, 10, 20, 50}

		for _, length := range lengths {
			result := GenerateRandomString(length)
			assert.Equal(t, length, len(result), "String should be %d characters", length)
		}
	})

	t.Run("should generate different strings", func(t *testing.T) {
		length := 10
		strings := make(map[string]bool)

		for i := 0; i < 50; i++ {
			str := GenerateRandomString(length)
			strings[str] = true
		}

		// Should have generated many different strings
		assert.Greater(t, len(strings), 30)
	})

	t.Run("should only contain alphanumeric characters", func(t *testing.T) {
		result := GenerateRandomString(100)

		for _, char := range result {
			isValid := (char >= 'a' && char <= 'z') ||
				(char >= 'A' && char <= 'Z') ||
				(char >= '0' && char <= '9')
			assert.True(t, isValid, "Character should be alphanumeric")
		}
	})

	t.Run("should handle zero length", func(t *testing.T) {
		result := GenerateRandomString(0)

		assert.Equal(t, "", result)
	})

	t.Run("should generate string with mix of upper, lower, and digits", func(t *testing.T) {
		// Generate a reasonably long string
		result := GenerateRandomString(100)

		hasLower := false
		hasUpper := false
		hasDigit := false

		for _, char := range result {
			if char >= 'a' && char <= 'z' {
				hasLower = true
			}
			if char >= 'A' && char <= 'Z' {
				hasUpper = true
			}
			if char >= '0' && char <= '9' {
				hasDigit = true
			}
		}

		// With 100 characters, we should have all types (statistically very likely)
		assert.True(t, hasLower, "Should contain lowercase letters")
		assert.True(t, hasUpper, "Should contain uppercase letters")
		assert.True(t, hasDigit, "Should contain digits")
	})
}

func TestGenerateRandomStrNumber(t *testing.T) {
	t.Run("should generate numeric string of correct length", func(t *testing.T) {
		length := 11

		result := GenerateRandomStrNumber(length)

		assert.Equal(t, length, len(result))
	})

	t.Run("should generate only numeric characters", func(t *testing.T) {
		result := GenerateRandomStrNumber(20)

		for _, char := range result {
			assert.True(t, char >= '0' && char <= '9', "Character should be a digit")
		}
	})

	t.Run("should generate different numbers", func(t *testing.T) {
		numbers := make(map[string]bool)

		for i := 0; i < 50; i++ {
			num := GenerateRandomStrNumber(10)
			numbers[num] = true
		}

		assert.Greater(t, len(numbers), 1)
	})
}
