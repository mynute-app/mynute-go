package email

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMailHog(t *testing.T) {
	t.Run("should create adapter with default values", func(t *testing.T) {
		os.Unsetenv("MAILHOG_HOST")
		os.Unsetenv("MAILHOG_PORT")
		os.Unsetenv("MAILHOG_DEFAULT_FROM")

		adapter, err := MailHog()

		assert.NoError(t, err)
		assert.NotNil(t, adapter)
		assert.Equal(t, "localhost", adapter.host)
		assert.Equal(t, "1025", adapter.port)
		assert.Equal(t, "noreply@test.local", adapter.defaultFrom)
	})

	t.Run("should create adapter with custom values from environment", func(t *testing.T) {
		os.Setenv("MAILHOG_HOST", "mailhog-server")
		os.Setenv("MAILHOG_PORT", "2525")
		os.Setenv("MAILHOG_DEFAULT_FROM", "custom@example.com")
		defer os.Unsetenv("MAILHOG_HOST")
		defer os.Unsetenv("MAILHOG_PORT")
		defer os.Unsetenv("MAILHOG_DEFAULT_FROM")

		adapter, err := MailHog()

		assert.NoError(t, err)
		assert.NotNil(t, adapter)
		assert.Equal(t, "mailhog-server", adapter.host)
		assert.Equal(t, "2525", adapter.port)
		assert.Equal(t, "custom@example.com", adapter.defaultFrom)
	})
}

func TestMailHogAdapter_Send(t *testing.T) {
	t.Run("should return error if no recipients", func(t *testing.T) {
		adapter := &MailHogAdapter{
			host:        "localhost",
			port:        "1025",
			defaultFrom: "test@example.com",
		}
		data := EmailData{
			To:      []string{},
			Subject: "Test",
			Html:    "<h1>Test</h1>",
		}

		err := adapter.Send(context.Background(), data)

		assert.Error(t, err)
		assert.Equal(t, "email must have at least one recipient", err.Error())
	})

	// Note: The following tests verify the message building logic
	// Actual SMTP sending would require a running MailHog instance
	t.Run("should use default from address if not provided", func(t *testing.T) {
		adapter := &MailHogAdapter{
			host:        "localhost",
			port:        "1025",
			defaultFrom: "default@example.com",
		}
		data := EmailData{
			To:      []string{"recipient@example.com"},
			Subject: "Test Subject",
			Html:    "<h1>Test</h1>",
		}

		from := data.From
		if from == "" {
			from = adapter.defaultFrom
		}
		message := adapter.buildMessage(from, data)

		// Verify the from address is set to default
		assert.Contains(t, message, "From: default@example.com")
	})

	t.Run("should use custom from address when provided", func(t *testing.T) {
		adapter := &MailHogAdapter{
			host:        "localhost",
			port:        "1025",
			defaultFrom: "default@example.com",
		}
		data := EmailData{
			From:    "custom@example.com",
			To:      []string{"recipient@example.com"},
			Subject: "Test Subject",
			Html:    "<h1>Test</h1>",
		}

		message := adapter.buildMessage(data.From, data)

		assert.Contains(t, message, "From: custom@example.com")
	})
}

func TestMailHogAdapter_BuildMessage(t *testing.T) {
	adapter := &MailHogAdapter{
		host:        "localhost",
		port:        "1025",
		defaultFrom: "test@example.com",
	}

	t.Run("should build basic HTML email message", func(t *testing.T) {
		data := EmailData{
			From:    "sender@example.com",
			To:      []string{"recipient@example.com"},
			Subject: "Test Subject",
			Html:    "<h1>Hello World</h1>",
		}

		message := adapter.buildMessage(data.From, data)

		assert.Contains(t, message, "From: sender@example.com")
		assert.Contains(t, message, "To: recipient@example.com")
		assert.Contains(t, message, "Subject: Test Subject")
		assert.Contains(t, message, "Content-Type: text/html; charset=UTF-8")
		assert.Contains(t, message, "<h1>Hello World</h1>")
		assert.Contains(t, message, "MIME-Version: 1.0")
	})

	t.Run("should build plain text email message", func(t *testing.T) {
		data := EmailData{
			From:    "sender@example.com",
			To:      []string{"recipient@example.com"},
			Subject: "Test Subject",
			Text:    "Hello World",
		}

		message := adapter.buildMessage(data.From, data)

		assert.Contains(t, message, "From: sender@example.com")
		assert.Contains(t, message, "To: recipient@example.com")
		assert.Contains(t, message, "Subject: Test Subject")
		assert.Contains(t, message, "Content-Type: text/plain; charset=UTF-8")
		assert.Contains(t, message, "Hello World")
	})

	t.Run("should include multiple recipients", func(t *testing.T) {
		data := EmailData{
			From:    "sender@example.com",
			To:      []string{"recipient1@example.com", "recipient2@example.com"},
			Subject: "Test Subject",
			Html:    "<h1>Test</h1>",
		}

		message := adapter.buildMessage(data.From, data)

		assert.Contains(t, message, "To: recipient1@example.com, recipient2@example.com")
	})

	t.Run("should include CC recipients", func(t *testing.T) {
		data := EmailData{
			From:    "sender@example.com",
			To:      []string{"recipient@example.com"},
			Cc:      []string{"cc1@example.com", "cc2@example.com"},
			Subject: "Test Subject",
			Html:    "<h1>Test</h1>",
		}

		message := adapter.buildMessage(data.From, data)

		assert.Contains(t, message, "Cc: cc1@example.com, cc2@example.com")
	})

	t.Run("should include BCC recipients", func(t *testing.T) {
		data := EmailData{
			From:    "sender@example.com",
			To:      []string{"recipient@example.com"},
			Bcc:     []string{"bcc1@example.com", "bcc2@example.com"},
			Subject: "Test Subject",
			Html:    "<h1>Test</h1>",
		}

		message := adapter.buildMessage(data.From, data)

		assert.Contains(t, message, "Bcc: bcc1@example.com, bcc2@example.com")
	})

	t.Run("should include Reply-To header", func(t *testing.T) {
		data := EmailData{
			From:    "sender@example.com",
			To:      []string{"recipient@example.com"},
			ReplyTo: "replyto@example.com",
			Subject: "Test Subject",
			Html:    "<h1>Test</h1>",
		}

		message := adapter.buildMessage(data.From, data)

		assert.Contains(t, message, "Reply-To: replyto@example.com")
	})

	t.Run("should include custom headers", func(t *testing.T) {
		data := EmailData{
			From:    "sender@example.com",
			To:      []string{"recipient@example.com"},
			Subject: "Test Subject",
			Html:    "<h1>Test</h1>",
			Headers: map[string]string{
				"X-Custom-Header": "custom-value",
				"X-Priority":      "1",
			},
		}

		message := adapter.buildMessage(data.From, data)

		assert.Contains(t, message, "X-Custom-Header: custom-value")
		assert.Contains(t, message, "X-Priority: 1")
	})

	t.Run("should prefer HTML over text when both provided", func(t *testing.T) {
		data := EmailData{
			From:    "sender@example.com",
			To:      []string{"recipient@example.com"},
			Subject: "Test Subject",
			Html:    "<h1>HTML Content</h1>",
			Text:    "Text Content",
		}

		message := adapter.buildMessage(data.From, data)

		assert.Contains(t, message, "Content-Type: text/html; charset=UTF-8")
		assert.Contains(t, message, "<h1>HTML Content</h1>")
		// Text should not be in the message when HTML is present
		lines := strings.Split(message, "\r\n")
		found := false
		for _, line := range lines {
			if line == "Text Content" {
				found = true
				break
			}
		}
		assert.False(t, found, "Text content should not be present when HTML is provided")
	})

	t.Run("should handle empty body gracefully", func(t *testing.T) {
		data := EmailData{
			From:    "sender@example.com",
			To:      []string{"recipient@example.com"},
			Subject: "Test Subject",
		}

		message := adapter.buildMessage(data.From, data)

		assert.Contains(t, message, "From: sender@example.com")
		assert.Contains(t, message, "To: recipient@example.com")
		assert.Contains(t, message, "Subject: Test Subject")
		assert.Contains(t, message, "Content-Type: text/plain; charset=UTF-8")
	})

	t.Run("should properly format message with CRLF line endings", func(t *testing.T) {
		data := EmailData{
			From:    "sender@example.com",
			To:      []string{"recipient@example.com"},
			Subject: "Test Subject",
			Html:    "<h1>Test</h1>",
		}

		message := adapter.buildMessage(data.From, data)

		// Check that headers are separated by CRLF
		assert.Contains(t, message, "\r\n")
		// Check that there's a blank line (CRLF CRLF) between headers and body
		assert.Contains(t, message, "\r\n\r\n")
	})
}
