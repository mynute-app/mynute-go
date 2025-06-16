package lib

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateRandomIntFromRange(min, max int) int {
	if min >= max {
		panic("min must be less than max")
	}
	return rnd.Intn(max-min) + min
}

// Generates a random name with the format: "Test <str> <random_number>"
func GenerateRandomName(str string) string {
	return fmt.Sprintf("Test %v %d", str, rand.Intn(100000))
}

// Generates a random phone number in the format: +XXXXXXXXXXX
// where X is a digit (0-9) and the total length is 12 characters.
func GenerateRandomPhoneNumber() string {
	return fmt.Sprintf("+%v", GenerateRandomStrNumber(11))
}

// Creates a random string of a specified length.
// The string will consist of lowercase and uppercase letters, as well as digits.
func GenerateRandomString(length int) string {
	// Define the character set to be used for generating the random string
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Initialize the random string
	randomString := make([]byte, length)

	// Generate a random string of the desired length
	for i := range randomString {
		randomString[i] = charset[rnd.Intn(len(charset))]
	}

	// Return the generated random string
	return string(randomString)
}

// Creates a random email address based on the provided name.
// The email will be in the format: test_<name>_email_<random_number>@gmail.com
func GenerateRandomEmail(name string) string {
	provider := "@gmail.com"
	nick := fmt.Sprintf("test_%s_email_%v", name, GenerateRandomInt(5))
	return fmt.Sprintf("%v%v", nick, provider)
}

// Generates a random integer of a specified length.
func GenerateRandomInt(length int) int {
	// Define the lower and upper bounds based on the desired length
	lowerBound := int(math.Pow10(length - 1)) // 10^(n-1)
	upperBound := int(math.Pow10(length)) - 1 // 10^n - 1

	// Generate a random number in the range [lowerBound, upperBound]
	return rnd.Intn(upperBound-lowerBound+1) + lowerBound
}

func GenerateRandomStrNumber(length int) string {
	// Define the character set to be used for generating the random string
	charset := "0123456789"

	// Initialize the random string
	randomString := make([]byte, length)

	// Generate a random string of the desired length
	for i := range randomString {
		randomString[i] = charset[rnd.Intn(len(charset))]
	}

	// Return the generated random string
	return string(randomString)
}

// GenerateDate creates a date with optional parameters.
// Expected order of arguments: year, month, day, hour, minute.
func GenerateDate(params ...int) time.Time {
	now := time.Now()

	// Default values
	year, month, day, hour, minute := now.Year(), rnd.Intn(12)+1, rnd.Intn(28)+1, rnd.Intn(10)+8, 0

	// Override values if provided
	if len(params) > 0 && params[0] != 0 {
		year = params[0]
	}
	if len(params) > 1 && params[1] != 0 {
		month = params[1]
	}
	if len(params) > 2 && params[2] != 0 {
		day = params[2]
	}
	if len(params) > 3 && params[3] != 0 {
		hour = params[3]
	}
	if len(params) > 4 && params[4] != 0 {
		minute = params[4]
	}

	// Construct the date
	myTime := time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)

	// Ensure the generated time is in the future
	if myTime.Before(now) {
		myTime = now.AddDate(0, 0, rnd.Intn(30)+1).Truncate(24 * time.Hour).Add(time.Duration(rnd.Intn(10)+8) * time.Hour)
	}

	return myTime
}

// GenerateDateRFC3339 creates a date in RFC3339 format 
// eg. "2021-01-01T09:00:00Z" and also accepts optional parameters.
// Expected order of arguments: year, month, day, hour, minute.
func GenerateDateRFC3339(params ...int) string {
	return GenerateDate(params...).Format(time.RFC3339)
}

