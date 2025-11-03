package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordValidation(t *testing.T) {
	t.Run("should accept valid passwords", func(t *testing.T) {
		validPasswords := []string{
			"MyP@ss1",
			"Test123!@#",
			"Abc123!xyz",
			"P@ssw0rd123",
		}

		for _, password := range validPasswords {
			err := ValidatorV10.Var(password, "myPasswordValidation")
			assert.NoError(t, err, "Password '%s' should be valid", password)
		}
	})

	t.Run("should reject invalid passwords", func(t *testing.T) {
		invalidPasswords := []string{
			"short",
			"nouppercaseorspecial1",
			"NOLOWERCASE1!",
			"NoSpecial123",
			"NoNumber!@#",
		}

		for _, password := range invalidPasswords {
			err := ValidatorV10.Var(password, "myPasswordValidation")
			assert.Error(t, err, "Password '%s' should be invalid", password)
		}
	})
}

func TestSubdomainValidation(t *testing.T) {
	t.Run("should accept valid subdomains", func(t *testing.T) {
		validSubdomains := []string{
			"mycompany",
			"test-company",
			"company123",
			"abcdefgh",
		}

		for _, subdomain := range validSubdomains {
			err := ValidatorV10.Var(subdomain, "mySubdomainValidation")
			assert.NoError(t, err, "Subdomain '%s' should be valid", subdomain)
		}
	})

	t.Run("should reject invalid subdomains", func(t *testing.T) {
		invalidSubdomains := []string{
			"MyCompany",
			"my_company",
			"short",
			"-mycompany",
			"my--company",
		}

		for _, subdomain := range invalidSubdomains {
			err := ValidatorV10.Var(subdomain, "mySubdomainValidation")
			assert.Error(t, err, "Subdomain '%s' should be invalid", subdomain)
		}
	})
}
