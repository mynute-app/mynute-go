package lib

import (
	"errors"
	"fmt"
	"regexp"
)

func ValidateName(name string, structStr string) error {
	isValid := len(name) >= 3
	if !isValid {
		errStr := fmt.Sprintf("%s.name must be at least 3 characters long", structStr)
		return errors.New(errStr)
	}
	return nil
}

func ValidateTaxID(taxID string) bool {
	return regexp.MustCompile(`^\d{14}$`).MatchString(taxID)
}