package email

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// --- MailHog API Tests ---

func TestMailHogAdapter_GetAPIURL(t *testing.T) {
	t.Run("should return default API URL", func(t *testing.T) {
		adapter := &MailHogAdapter{
			host: "localhost",
			port: "1025",
		}

		url := adapter.GetAPIURL()
		assert.Equal(t, "http://localhost:8025/api", url)
	})

	t.Run("should use custom API port from environment", func(t *testing.T) {
		os.Setenv("MAILHOG_API_PORT", "9025")
		defer os.Unsetenv("MAILHOG_API_PORT")

		adapter := &MailHogAdapter{
			host: "mailhog-server",
			port: "1025",
		}

		url := adapter.GetAPIURL()
		assert.Equal(t, "http://mailhog-server:9025/api", url)
	})
}

func TestMailHogAdapter_GetMessages(t *testing.T) {
	t.Run("should retrieve messages successfully", func(t *testing.T) {
		// Create mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v2/messages", r.URL.Path)

			response := MailHogMessagesResponse{
				Total: 2,
				Count: 2,
				Start: 0,
				Messages: []MailHogMessage{
					{
						ID:   "msg1",
						From: MailHogPath{Mailbox: "sender", Domain: "example.com"},
						To: []MailHogPath{
							{Mailbox: "recipient", Domain: "example.com"},
						},
						Created: time.Now(),
						Content: MailHogContent{
							Headers: map[string][]string{
								"Subject": {"Test Subject"},
							},
							Body: "Test Body",
						},
					},
					{
						ID:   "msg2",
						From: MailHogPath{Mailbox: "sender2", Domain: "example.com"},
						To: []MailHogPath{
							{Mailbox: "recipient2", Domain: "example.com"},
						},
						Created: time.Now(),
						Content: MailHogContent{
							Headers: map[string][]string{
								"Subject": {"Test Subject 2"},
							},
							Body: "Test Body 2",
						},
					},
				},
			}

			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		// Set API URL to mock server
		os.Setenv("MAILHOG_API_PORT", strings.Split(server.URL, ":")[2])
		defer os.Unsetenv("MAILHOG_API_PORT")

		adapter := &MailHogAdapter{
			host: "localhost",
			port: "1025",
		}

		messages, err := adapter.GetMessages()

		require.NoError(t, err)
		assert.Len(t, messages, 2)
		assert.Equal(t, "msg1", messages[0].ID)
		assert.Equal(t, "msg2", messages[1].ID)
	})

	t.Run("should return error on HTTP failure", func(t *testing.T) {
		adapter := &MailHogAdapter{
			host: "invalid-host",
			port: "1025",
		}

		_, err := adapter.GetMessages()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get messages")
	})
}

func TestMailHogAdapter_GetLatestMessageTo(t *testing.T) {
	now := time.Now()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := MailHogMessagesResponse{
			Total: 3,
			Count: 3,
			Start: 0,
			Messages: []MailHogMessage{
				{
					ID:      "msg1",
					To:      []MailHogPath{{Mailbox: "user1", Domain: "example.com"}},
					Content: MailHogContent{Body: "First email"},
					Created: now.Add(-2 * time.Minute),
				},
				{
					ID:      "msg2",
					To:      []MailHogPath{{Mailbox: "user2", Domain: "example.com"}},
					Content: MailHogContent{Body: "Second email"},
					Created: now.Add(-1 * time.Minute),
				},
				{
					ID:      "msg3",
					To:      []MailHogPath{{Mailbox: "user1", Domain: "example.com"}},
					Content: MailHogContent{Body: "Third email"},
					Created: now,
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	os.Setenv("MAILHOG_API_PORT", strings.Split(server.URL, ":")[2])
	defer os.Unsetenv("MAILHOG_API_PORT")

	adapter := &MailHogAdapter{host: "localhost", port: "1025"}

	t.Run("should get latest message for recipient", func(t *testing.T) {
		msg, err := adapter.GetLatestMessageTo("user1@example.com")

		require.NoError(t, err)
		assert.Equal(t, "msg3", msg.ID)
		assert.Equal(t, "Third email", msg.Content.Body)
	})

	t.Run("should return error for non-existent recipient", func(t *testing.T) {
		msg, err := adapter.GetLatestMessageTo("nonexistent@example.com")

		assert.Error(t, err)
		assert.Nil(t, msg)
		assert.Contains(t, err.Error(), "no message found")
	})
}

func TestMailHogMessage_GetMessageBody(t *testing.T) {
	t.Run("should return HTML body from MIME parts", func(t *testing.T) {
		msg := &MailHogMessage{
			MIME: &MailHogMIME{
				Parts: []MailHogMIMEPart{
					{
						Headers: map[string][]string{
							"Content-Type": {"text/html; charset=UTF-8"},
						},
						Body: "<h1>HTML Content</h1>",
					},
				},
			},
			Content: MailHogContent{
				Body: "Fallback body",
			},
		}

		body := msg.GetMessageBody()
		assert.Equal(t, "<h1>HTML Content</h1>", body)
	})

	t.Run("should return first MIME part when no HTML", func(t *testing.T) {
		msg := &MailHogMessage{
			MIME: &MailHogMIME{
				Parts: []MailHogMIMEPart{
					{
						Headers: map[string][]string{
							"Content-Type": {"text/plain; charset=UTF-8"},
						},
						Body: "Plain text content",
					},
				},
			},
			Content: MailHogContent{
				Body: "Fallback body",
			},
		}

		body := msg.GetMessageBody()
		assert.Equal(t, "Plain text content", body)
	})

	t.Run("should fallback to raw body", func(t *testing.T) {
		msg := &MailHogMessage{
			Content: MailHogContent{
				Body: "Raw body content",
			},
		}

		body := msg.GetMessageBody()
		assert.Equal(t, "Raw body content", body)
	})
}

func TestMailHogMessage_GetSubject(t *testing.T) {
	t.Run("should return subject from headers", func(t *testing.T) {
		msg := &MailHogMessage{
			Content: MailHogContent{
				Headers: map[string][]string{
					"Subject": {"Test Email Subject"},
				},
			},
		}

		subject := msg.GetSubject()
		assert.Equal(t, "Test Email Subject", subject)
	})

	t.Run("should return empty string when no subject", func(t *testing.T) {
		msg := &MailHogMessage{
			Content: MailHogContent{
				Headers: map[string][]string{},
			},
		}

		subject := msg.GetSubject()
		assert.Equal(t, "", subject)
	})
}

func TestMailHogMessage_ExtractCode(t *testing.T) {
	t.Run("should extract 6-digit code with default pattern", func(t *testing.T) {
		msg := &MailHogMessage{
			Content: MailHogContent{
				Body: "Your validation code is: 123456. Please use it to login.",
			},
		}

		code, err := msg.ExtractCode()
		require.NoError(t, err)
		assert.Equal(t, "123456", code)
	})

	t.Run("should extract code with custom pattern", func(t *testing.T) {
		msg := &MailHogMessage{
			Content: MailHogContent{
				Body: "Your code is ABC123 for verification.",
			},
		}

		code, err := msg.ExtractCode(`[A-Z]{3}\d{3}`)
		require.NoError(t, err)
		assert.Equal(t, "ABC123", code)
	})

	t.Run("should return error when no code found", func(t *testing.T) {
		msg := &MailHogMessage{
			Content: MailHogContent{
				Body: "No code in this email",
			},
		}

		code, err := msg.ExtractCode()
		assert.Error(t, err)
		assert.Equal(t, "", code)
		assert.Contains(t, err.Error(), "no code found")
	})
}

func TestMailHogMessage_ExtractValidationCode(t *testing.T) {
	t.Run("should extract 6-digit validation code", func(t *testing.T) {
		msg := &MailHogMessage{
			Content: MailHogContent{
				Body: "<h1>Your code is 987654</h1>",
			},
		}

		code, err := msg.ExtractValidationCode()
		require.NoError(t, err)
		assert.Equal(t, "987654", code)
	})

	t.Run("should extract alphanumeric code", func(t *testing.T) {
		msg := &MailHogMessage{
			Content: MailHogContent{
				Body: "Verification code: ABC123",
			},
		}

		code, err := msg.ExtractValidationCode()
		require.NoError(t, err)
		assert.Equal(t, "ABC123", code)
	})

	t.Run("should extract from HTML body", func(t *testing.T) {
		msg := &MailHogMessage{
			MIME: &MailHogMIME{
				Parts: []MailHogMIMEPart{
					{
						Headers: map[string][]string{
							"Content-Type": {"text/html"},
						},
						Body: "<div>Code: <strong>456789</strong></div>",
					},
				},
			},
		}

		code, err := msg.ExtractValidationCode()
		require.NoError(t, err)
		assert.Equal(t, "456789", code)
	})

	t.Run("should return error when no validation code found", func(t *testing.T) {
		msg := &MailHogMessage{
			Content: MailHogContent{
				Body: "Just plain text with no code",
			},
		}

		code, err := msg.ExtractValidationCode()
		assert.Error(t, err)
		assert.Equal(t, "", code)
		assert.Contains(t, err.Error(), "no validation code found")
	})
}

func TestMailHogAdapter_DeleteMessage(t *testing.T) {
	t.Run("should delete message successfully", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodDelete, r.Method)
			assert.Equal(t, "/api/v1/messages/msg123", r.URL.Path)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		os.Setenv("MAILHOG_API_PORT", strings.Split(server.URL, ":")[2])
		defer os.Unsetenv("MAILHOG_API_PORT")

		adapter := &MailHogAdapter{host: "localhost", port: "1025"}
		err := adapter.DeleteMessage("msg123")

		assert.NoError(t, err)
	})
}

func TestMailHogAdapter_DeleteAllMessages(t *testing.T) {
	t.Run("should delete all messages successfully", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodDelete, r.Method)
			assert.Equal(t, "/api/v1/messages", r.URL.Path)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		os.Setenv("MAILHOG_API_PORT", strings.Split(server.URL, ":")[2])
		defer os.Unsetenv("MAILHOG_API_PORT")

		adapter := &MailHogAdapter{host: "localhost", port: "1025"}
		err := adapter.DeleteAllMessages()

		assert.NoError(t, err)
	})
}

