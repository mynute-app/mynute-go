package email

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
)

// --- MailHog Implementation ---

// MailHogAdapter is an implementation of the Sender interface that uses MailHog for testing.
// MailHog captures emails without actually sending them, making it perfect for e2e tests.
type MailHogAdapter struct {
	host        string
	port        string
	defaultFrom string
}

// NewMailHogAdapter initializes and returns a new MailHogAdapter.
// It should be created once when the application starts.
func MailHog() (*MailHogAdapter, error) {
	host := os.Getenv("MAILHOG_HOST")
	if host == "" {
		host = "localhost" // Default to localhost for local testing
	}

	port := os.Getenv("MAILHOG_PORT")
	if port == "" {
		port = "1025" // Default MailHog SMTP port
	}

	defaultFrom := os.Getenv("MAILHOG_DEFAULT_FROM")
	if defaultFrom == "" {
		defaultFrom = "noreply@test.local" // Default for testing
	}

	return &MailHogAdapter{
		host:        host,
		port:        port,
		defaultFrom: defaultFrom,
	}, nil
}

// Send sends an email using the MailHog service.
func (m *MailHogAdapter) Send(ctx context.Context, data EmailData) error {
	if len(data.To) == 0 {
		return errors.New("email must have at least one recipient")
	}

	from := data.From
	if from == "" {
		from = m.defaultFrom
	}

	// Build the email message
	message := m.buildMessage(from, data)

	// Connect to MailHog SMTP server
	addr := fmt.Sprintf("%s:%s", m.host, m.port)

	// MailHog doesn't require authentication
	err := smtp.SendMail(
		addr,
		nil, // No auth required for MailHog
		from,
		data.To,
		[]byte(message),
	)

	if err != nil {
		return fmt.Errorf("failed to send email via mailhog: %w", err)
	}

	log.Printf("Email sent successfully to %v via MailHog (%s)\n", data.To, addr)
	return nil
}

// buildMessage constructs the email message with headers and body
func (m *MailHogAdapter) buildMessage(from string, data EmailData) string {
	var builder strings.Builder

	// Required headers
	builder.WriteString(fmt.Sprintf("From: %s\r\n", from))
	builder.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(data.To, ", ")))
	builder.WriteString(fmt.Sprintf("Subject: %s\r\n", data.Subject))

	// Optional headers
	if len(data.Cc) > 0 {
		builder.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(data.Cc, ", ")))
	}

	if len(data.Bcc) > 0 {
		builder.WriteString(fmt.Sprintf("Bcc: %s\r\n", strings.Join(data.Bcc, ", ")))
	}

	if data.ReplyTo != "" {
		builder.WriteString(fmt.Sprintf("Reply-To: %s\r\n", data.ReplyTo))
	}

	// Custom headers
	for key, value := range data.Headers {
		builder.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	// MIME headers for HTML email
	builder.WriteString("MIME-Version: 1.0\r\n")

	if data.Html != "" {
		builder.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		builder.WriteString("\r\n")
		builder.WriteString(data.Html)
	} else if data.Text != "" {
		builder.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		builder.WriteString("\r\n")
		builder.WriteString(data.Text)
	} else {
		builder.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		builder.WriteString("\r\n")
	}

	return builder.String()
}
