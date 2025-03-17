package lib

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

// Helper function to generate a random name
func GenerateRandomName(str string) string {
	return fmt.Sprintf("Test %v %d", str, rand.Intn(100000))
}

func GenerateRandomPhoneNumber() string {
	return fmt.Sprintf("+%v", GenerateRandomStrNumber(11))
}

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

func GenerateRandomEmail(name string) string {
	provider := "@gmail.com"
	nick := fmt.Sprintf("test_%s_email_%v", name, GenerateRandomInt(5))
	return fmt.Sprintf("%v%v", nick, provider)
}

func GenerateRandomInt(length int) int {
	// Create a new random source

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
