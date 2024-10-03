package lib

import "regexp"

func ValidateName(name string) bool {
	return len(name) >= 3
}

func ValidateTaxID(taxID string) bool {
	return regexp.MustCompile(`^\d{14}$`).MatchString(taxID)
}