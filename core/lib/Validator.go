package lib

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// Global validator instance
var ValidatorV10 *validator.Validate

// Custom password validation function
func isValidPassword(fl validator.FieldLevel) bool {
	pswd := fl.Field().String()

	// Password requirements:
	// - 6 to 16 characters
	// - At least one uppercase letter
	// - At least one lowercase letter
	// - At least one number
	// - At least one special character (!@#$%^&*)

	if len(pswd) < 6 || len(pswd) > 16 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString
	hasLower := regexp.MustCompile(`[a-z]`).MatchString
	hasDigit := regexp.MustCompile(`\d`).MatchString
	hasSpecial := regexp.MustCompile(`[!@#$%^&*]`).MatchString

	return hasUpper(pswd) && hasLower(pswd) && hasDigit(pswd) && hasSpecial(pswd)
}

// init function to initialize and register the validator
func init() {
	ValidatorV10 = validator.New()
	if err := ValidatorV10.RegisterValidation("myPasswordValidation", isValidPassword); err != nil {
		panic(err)
	}
}