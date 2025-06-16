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

func isValidSubdomain(fl validator.FieldLevel) bool {
	subdomain := fl.Field().String()

	// Subdomain requirements:
	// - 8 to 30 characters
	// - Only lowercase letters, numbers, and hyphens
	// - Cannot contain special characters
	// - Cannot contain spaces
	// - Cannot contain dots
	// - Cannot contain underscores
	// - Cannot contain uppercase letters
	// - Cannot contain consecutive hyphens
	// - Cannot start or end with a hyphen

	if len(subdomain) < 8 || len(subdomain) > 30 {
		return false
	}

	hasLetter := regexp.MustCompile(`[a-z]`).MatchString
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+=<>?{}|~]`).MatchString
	hasSpace := regexp.MustCompile(`\s`).MatchString
	hasConsecutiveHyphens := regexp.MustCompile(`--`).MatchString
	hasDot := regexp.MustCompile(`\.`).MatchString
	hasUnderscore := regexp.MustCompile(`_`).MatchString
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString

	return hasLetter(subdomain) &&
		!hasSpecial(subdomain) &&
		!hasSpace(subdomain) &&
		!hasDot(subdomain) &&
		!hasUnderscore(subdomain) &&
		!hasUpper(subdomain) &&
		!hasConsecutiveHyphens(subdomain) &&
		subdomain[0] != '-' &&
		subdomain[len(subdomain)-1] != '-'
}

// init function to initialize and register the validator
func init() {
	ValidatorV10 = validator.New()
	if err := ValidatorV10.RegisterValidation("myPasswordValidation", isValidPassword); err != nil {
		panic(err)
	}
	if err := ValidatorV10.RegisterValidation("mySubdomainValidation", isValidSubdomain); err != nil {
		panic(err)
	}
}
