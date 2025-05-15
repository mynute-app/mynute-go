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
	// - 8 to 20 characters
	// - Must have at least one letter
	// - Must not contain special characters
	// - Must not contain spaces
	// - Must not contain consecutive hyphens
	// - Must not start or end with a hyphen
	// - Must not contain dots
	// - Must not contain underscores
	// - Must not contain uppercase letters

	hasEnoughChars := len(subdomain) >= 8 && len(subdomain) <= 20
	hasLetter := regexp.MustCompile(`[a-z]`).MatchString
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+=<>?{}|~]`).MatchString
	hasSpace := regexp.MustCompile(`\s`).MatchString
	hasConsecutiveHyphens := regexp.MustCompile(`--`).MatchString
	startsWithHyphen := subdomain[0] == '-'
	endsWithHyphen := subdomain[len(subdomain)-1] == '-'
	hasDot := regexp.MustCompile(`\.`).MatchString
	hasUnderscore := regexp.MustCompile(`_`).MatchString
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString
	shouldntHaveThese := !(hasLetter(subdomain) &&
		hasSpecial(subdomain) &&
		hasSpace(subdomain) &&
		hasConsecutiveHyphens(subdomain) &&
		startsWithHyphen &&
		endsWithHyphen &&
		hasDot(subdomain) &&
		hasUnderscore(subdomain) &&
		hasUpper(subdomain))
	return hasEnoughChars && shouldntHaveThese
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
