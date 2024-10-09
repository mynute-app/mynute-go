package lib

import (
	"fmt"
	"math/rand"
)

// Helper function to generate a random name
func GenerateRandomName(str string) string {
	return fmt.Sprintf("Test %v %d", str, rand.Intn(1000))
}