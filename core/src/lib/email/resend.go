package email

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/resend/resend-go/v2"
)

// --- Resend Implementation ---

// ResendAdapter is an implementation of the Sender interface that uses the Resend API.
type ResendAdapter struct {
	client      *resend.Client
	defaultFrom string
}

// NewResendAdapter initializes and returns a new ResendAdapter.
// It should be created once when the application starts.
func Resend() (*ResendAdapter, error) {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return nil, errors.New("RESEND_API_KEY environment variable is not set")
	}

	defaultFrom := os.Getenv("RESEND_DEFAULT_FROM")
	if defaultFrom == "" {
		return nil, errors.New("RESEND_DEFAULT_FROM environment variable is not set")
	}

	client := resend.NewClient(apiKey)

	return &ResendAdapter{
		client:      client,
		defaultFrom: defaultFrom,
	}, nil
}

// Send sends an email using the Resend service.
func (r *ResendAdapter) Send(ctx context.Context, data EmailData) error {
	if len(data.To) == 0 {
		return errors.New("email must have at least one recipient")
	}

	from := data.From
	if from == "" {
		from = r.defaultFrom
	}

	// Validate and clean email addresses
	cleanedTo := make([]string, len(data.To))
	for i, email := range data.To {
		// Trim whitespace
		trimmed := strings.TrimSpace(email)
		// URL decode the email (fixes %40 encoding)
		unescaped, _ := url.QueryUnescape(trimmed)
		cleanedTo[i] = unescaped
		if cleanedTo[i] == "" {
			return fmt.Errorf("recipient email at index %d is empty", i)
		}
	}

	params := &resend.SendEmailRequest{
		From:    from,
		To:      cleanedTo,
		Subject: data.Subject,
		Html:    data.Html,
	}

	_, err := r.client.Emails.SendWithContext(ctx, params)
	if err != nil {
		for email := range cleanedTo {
			log.Printf("Sending email to: %s\n", cleanedTo[email])
		}
		log.Printf(">>> Resend email error: %v\n", err)
		return fmt.Errorf("failed to send email via resend (from: %s, to: %v): %w", from, cleanedTo, err)
	}

	return nil
}
