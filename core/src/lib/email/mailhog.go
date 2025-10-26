package email

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strings"
	"time"
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

// --- MailHog API Client for E2E Testing ---

// MailHogMessage represents an email message from MailHog API
type MailHogMessage struct {
	ID      string                 `json:"ID"`
	From    MailHogPath            `json:"From"`
	To      []MailHogPath          `json:"To"`
	Content MailHogContent         `json:"Content"`
	Created time.Time              `json:"Created"`
	MIME    *MailHogMIME           `json:"MIME"`
	Raw     map[string]interface{} `json:"Raw"`
}

// MailHogPath represents an email address
type MailHogPath struct {
	Relays  interface{} `json:"Relays"`
	Mailbox string      `json:"Mailbox"`
	Domain  string      `json:"Domain"`
	Params  string      `json:"Params"`
}

// MailHogContent represents email content
type MailHogContent struct {
	Headers map[string][]string `json:"Headers"`
	Body    string              `json:"Body"`
	Size    int                 `json:"Size"`
	MIME    interface{}         `json:"MIME"`
}

// MailHogMIME represents MIME content
type MailHogMIME struct {
	Parts []MailHogMIMEPart `json:"Parts"`
}

// MailHogMIMEPart represents a MIME part
type MailHogMIMEPart struct {
	Headers map[string][]string `json:"Headers"`
	Body    string              `json:"Body"`
	Size    int                 `json:"Size"`
	MIME    interface{}         `json:"MIME"`
}

// MailHogMessagesResponse represents the API response
type MailHogMessagesResponse struct {
	Total    int              `json:"total"`
	Count    int              `json:"count"`
	Start    int              `json:"start"`
	Messages []MailHogMessage `json:"items"`
}

// GetAPIURL returns the MailHog API URL
func (m *MailHogAdapter) GetAPIURL() string {
	apiPort := os.Getenv("MAILHOG_API_PORT")
	if apiPort == "" {
		apiPort = "8025" // Default MailHog API port
	}
	return fmt.Sprintf("http://%s:%s/api", m.host, apiPort)
}

// GetMessages retrieves all messages from MailHog
func (m *MailHogAdapter) GetMessages() ([]MailHogMessage, error) {
	url := fmt.Sprintf("%s/v2/messages", m.GetAPIURL())

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages from mailhog: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("mailhog API returned status %d", resp.StatusCode)
	}

	var result MailHogMessagesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode mailhog response: %w", err)
	}

	return result.Messages, nil
}

// GetLatestMessageTo retrieves the latest message sent to a specific email address
func (m *MailHogAdapter) GetLatestMessageTo(email string) (*MailHogMessage, error) {
	messages, err := m.GetMessages()
	if err != nil {
		return nil, err
	}

	// Search from most recent to oldest
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		for _, to := range msg.To {
			recipientEmail := fmt.Sprintf("%s@%s", to.Mailbox, to.Domain)
			if recipientEmail == email {
				return &msg, nil
			}
		}
	}

	return nil, fmt.Errorf("no message found for recipient: %s", email)
}

// GetMessageBody returns the email body (HTML or plain text)
func (msg *MailHogMessage) GetMessageBody() string {
	// Try to get HTML body first
	if msg.MIME != nil && len(msg.MIME.Parts) > 0 {
		for _, part := range msg.MIME.Parts {
			contentType := part.Headers["Content-Type"]
			if len(contentType) > 0 && strings.Contains(contentType[0], "text/html") {
				return part.Body
			}
		}
		// Fallback to first part
		if len(msg.MIME.Parts) > 0 {
			return msg.MIME.Parts[0].Body
		}
	}

	// Fallback to raw body
	return msg.Content.Body
}

// GetSubject returns the email subject
func (msg *MailHogMessage) GetSubject() string {
	if subject, ok := msg.Content.Headers["Subject"]; ok && len(subject) > 0 {
		return subject[0]
	}
	return ""
}

// ExtractCode extracts a numeric code from the email body using a regex pattern
// Default pattern matches 4-6 digit codes
func (msg *MailHogMessage) ExtractCode(pattern ...string) (string, error) {
	body := msg.GetMessageBody()

	// Default pattern: 4-6 digit code
	regexPattern := `\b\d{4,6}\b`
	if len(pattern) > 0 {
		regexPattern = pattern[0]
	}

	re := regexp.MustCompile(regexPattern)
	matches := re.FindStringSubmatch(body)

	if len(matches) == 0 {
		return "", fmt.Errorf("no code found matching pattern: %s", regexPattern)
	}

	return matches[0], nil
}

// ExtractValidationCode is a convenience method for extracting validation codes
// It looks for common patterns like "123456" or "ABC123"
func (msg *MailHogMessage) ExtractValidationCode() (string, error) {
	body := msg.GetMessageBody()

	// Try to find 6-digit codes that are not hex colors
	// Pattern: find all 6-digit numbers
	re := regexp.MustCompile(`\b\d{6}\b`)
	matches := re.FindAllStringIndex(body, -1)

	for _, match := range matches {
		start := match[0]
		end := match[1]
		code := body[start:end]

		// Check if this is preceded by a # (hex color)
		if start > 0 && body[start-1] == '#' {
			continue // Skip hex colors
		}

		// Found a valid code!
		return code, nil
	}

	// Try other patterns if 6-digit numeric didn't work
	patterns := []string{
		`\b\d{4,5}\b`,       // 4-5 digit code
		`\b\d{7,8}\b`,       // 7-8 digit code
		`\b[A-Z0-9]{6}\b`,   // 6-character alphanumeric
		`\b[A-Z]{3}\d{3}\b`, // 3 letters + 3 digits
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(body)
		if len(matches) > 0 {
			return matches[0], nil
		}
	}

	return "", fmt.Errorf("no validation code found in email body")
}

// DeleteMessage deletes a specific message from MailHog
func (m *MailHogAdapter) DeleteMessage(messageID string) error {
	url := fmt.Sprintf("%s/v1/messages/%s", m.GetAPIURL(), messageID)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("mailhog API returned status %d", resp.StatusCode)
	}

	return nil
}

// DeleteAllMessages deletes all messages from MailHog
func (m *MailHogAdapter) DeleteAllMessages() error {
	url := fmt.Sprintf("%s/v1/messages", m.GetAPIURL())

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete messages: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("mailhog API returned status %d", resp.StatusCode)
	}

	return nil
}
