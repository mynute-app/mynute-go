package lib

import (
	"fmt"
	"net/url"

	"github.com/go-playground/validator/v10"
)

// PrepareEmail validates and decodes an email parameter
func PrepareEmail(user_email string) (string, error) {
	if user_email == "" {
		return "", Error.General.BadRequest.WithError(fmt.Errorf("email parameter is empty"))
	}

	cleanedEmail, err := url.QueryUnescape(user_email)
	if err != nil {
		return "", Error.General.BadRequest.WithError(err)
	}

	if err := ValidatorV10.Var(cleanedEmail, "email"); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return "", Error.General.BadRequest.WithError(fmt.Errorf("invalid email: %w", err))
		} else {
			return "", Error.General.InternalError.WithError(err)
		}
	}

	return cleanedEmail, nil
}

