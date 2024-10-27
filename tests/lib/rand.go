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

func GenerateRandomIntOfExactly(length int) int {
	// Create a new random source
	
	// Define the lower and upper bounds based on the desired length
	lowerBound := int(math.Pow10(length - 1)) // 10^(n-1)
	upperBound := int(math.Pow10(length)) - 1 // 10^n - 1

	// Generate a random number in the range [lowerBound, upperBound]
	return rnd.Intn(upperBound-lowerBound+1) + lowerBound
}
