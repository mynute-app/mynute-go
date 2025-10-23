package email

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

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

	params := &resend.SendEmailRequest{
		From:    from,
		To:      data.To,
		Subject: data.Subject,
		Html:    data.HTMLBody,
	}

	sent, err := r.client.Emails.SendWithContext(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to send email via resend: %w", err)
	}

	log.Printf("Email sent successfully to %v with ID: %s\n", data.To, sent.Id)
	return nil
}