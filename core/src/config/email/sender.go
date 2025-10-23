package email

import (
	"context"
	"fmt"
)

// EmailData holds all the necessary information for an email.
type EmailData struct {
	From          string
	To            []string
	Subject       string
	HTMLBody      string
	PlainTextBody *string
}

// Sender defines the contract for any email sending service.
// This allows for easy mocking and swapping of implementations (e.g., for testing).
type Sender interface {
	Send(ctx context.Context, data EmailData) error
}

func NewProvider(provider string) (Sender, error) {
	switch provider {
	case "resend":
		return Resend()
	default:
		return nil, fmt.Errorf("email provider (%s) not implemented", provider)
	}
}
