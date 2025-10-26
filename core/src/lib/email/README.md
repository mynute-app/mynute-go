# Email Package

This package provides a unified interface for sending emails with support for multiple providers and multi-language templates.

## Features

- **Multiple Email Providers**: Support for Resend (production) and MailHog (testing)
- **Template Rendering**: Multi-language HTML email templates
- **Type-Safe**: Strongly typed interfaces and data structures
- **Testable**: Easy to mock and test with comprehensive unit tests

## Providers

### Resend (Production)
Resend is used for sending real emails in production environments.

**Environment Variables:**
- `RESEND_API_KEY` - Your Resend API key (required)
- `RESEND_DEFAULT_FROM` - Default sender email address (required)

```go
provider, err := email.NewProvider("resend")
```

### MailHog (Testing/Development)
MailHog captures emails without sending them, perfect for development and e2e tests.

**Environment Variables:**
- `MAILHOG_HOST` - MailHog server host (default: "localhost")
- `MAILHOG_PORT` - MailHog SMTP port (default: "1025")
- `MAILHOG_DEFAULT_FROM` - Default sender email (default: "noreply@test.local")

```go
provider, err := email.NewProvider("mailhog")
```

**Accessing MailHog Web UI:**
- Start MailHog using Docker Compose: `docker-compose -f docker-compose.dev.yml up -d mailhog`
- Open browser: http://localhost:8025
- All emails sent via MailHog will appear in the web UI

## Template Rendering

The package includes a template renderer for creating multi-language emails.

### Supported Languages
- English (en) - default
- Portuguese (pt)
- Spanish (es)

### Usage

```go
// Initialize renderer
renderer := email.NewTemplateRenderer("./static/email", "./translation/email")

// Render email in Portuguese
htmlBody, err := renderer.RenderEmail("login_validation", "pt", email.TemplateData{
    "ValidationCode": "123456",
})
```

### Creating New Templates

1. **Create HTML Template**: Add file to `static/` directory (e.g., `static/my_email.html`)
```html
<!DOCTYPE html>
<html>
<head>
    <title>{{.title}}</title>
</head>
<body>
    <h1>{{.heading}}</h1>
    <p>{{.message}}</p>
    <p>Code: {{.ValidationCode}}</p>
</body>
</html>
```

2. **Create Translation File**: Add file to `translation/` directory (e.g., `translation/my_email.json`)
```json
{
  "en": {
    "title": "Email Title",
    "heading": "Welcome",
    "message": "Your message here"
  },
  "pt": {
    "title": "Título do Email",
    "heading": "Bem-vindo",
    "message": "Sua mensagem aqui"
  }
}
```

3. **Render the Template**:
```go
html, err := renderer.RenderEmail("my_email", "pt", email.TemplateData{
    "ValidationCode": "ABC123",
})
```

## Complete Example

```go
package main

import (
    "context"
    "fmt"
    "mynute-go/core/src/config/email"
)

func SendLoginCode(userEmail, code, language string) error {
    // Render email template
    renderer := email.NewTemplateRenderer("./static/email", "./translation/email")
    htmlBody, err := renderer.RenderEmail("login_validation", language, email.TemplateData{
        "ValidationCode": code,
    })
    if err != nil {
        return fmt.Errorf("failed to render email: %w", err)
    }

    // Choose provider based on environment
    providerName := "mailhog" // Use "resend" in production
    provider, err := email.NewProvider(providerName)
    if err != nil {
        return fmt.Errorf("failed to initialize provider: %w", err)
    }

    // Send email
    err = provider.Send(context.Background(), email.EmailData{
        To:      []string{userEmail},
        Subject: "Your Login Code",
        Html:    htmlBody,
    })
    if err != nil {
        return fmt.Errorf("failed to send email: %w", err)
    }

    return nil
}
```

## Testing

Run all tests:
```bash
go test ./core/src/config/email/... -v
```

Run specific provider tests:
```bash
go test ./core/src/config/email/... -v -run TestMailHog
go test ./core/src/config/email/... -v -run TestResend
```

## E2E Testing with MailHog

1. **Start MailHog**:
```bash
docker-compose -f docker-compose.dev.yml up -d mailhog
```

2. **Configure your tests to use MailHog**:
```go
os.Setenv("EMAIL_PROVIDER", "mailhog")
```

3. **Send test emails and verify in Web UI**:
- Open http://localhost:8025
- Send your test emails
- Verify content, recipients, and formatting

4. **Stop MailHog**:
```bash
docker-compose -f docker-compose.dev.yml stop mailhog
```

### Programmatic Email Verification

MailHog provides an API client for programmatic email verification in E2E tests:

```go
// Initialize MailHog adapter
provider, _ := email.NewProvider("mailhog")
mailhog := provider.(*email.MailHogAdapter)

// Send test email
err := mailhog.Send(ctx, email.EmailData{
    To:      []string{"test@example.com"},
    Subject: "Validation Code",
    Html:    "<h1>Your code is 123456</h1>",
})

// Retrieve the email
msg, err := mailhog.GetLatestMessageTo("test@example.com")
if err != nil {
    t.Fatal("Email not received")
}

// Extract validation code
code, err := msg.ExtractValidationCode()
assert.Equal(t, "123456", code)

// Verify subject
assert.Equal(t, "Validation Code", msg.GetSubject())

// Clean up after test
mailhog.DeleteAllMessages()
```

### MailHog API Methods

#### Retrieve Messages
```go
// Get all messages
messages, err := mailhog.GetMessages()

// Get latest message to specific recipient
msg, err := mailhog.GetLatestMessageTo("user@example.com")
```

#### Extract Information
```go
// Get email body (HTML or plain text)
body := msg.GetMessageBody()

// Get email subject
subject := msg.GetSubject()

// Extract validation code (tries multiple patterns)
code, err := msg.ExtractValidationCode()

// Extract code with custom regex pattern
code, err := msg.ExtractCode(`\b[A-Z]{3}\d{3}\b`)
```

#### Clean Up
```go
// Delete specific message
err := mailhog.DeleteMessage(messageID)

// Delete all messages
err := mailhog.DeleteAllMessages()
```

### Complete E2E Test Example

```go
func TestUserRegistration(t *testing.T) {
    // Setup
    provider, _ := email.NewProvider("mailhog")
    mailhog := provider.(*email.MailHogAdapter)
    
    // Clean mailbox before test
    mailhog.DeleteAllMessages()
    
    // Trigger registration (sends email)
    userEmail := "newuser@example.com"
    RegisterUser(userEmail)
    
    // Wait briefly for email to arrive
    time.Sleep(100 * time.Millisecond)
    
    // Retrieve the registration email
    msg, err := mailhog.GetLatestMessageTo(userEmail)
    require.NoError(t, err)
    
    // Verify email content
    assert.Contains(t, msg.GetSubject(), "Welcome")
    assert.Contains(t, msg.GetMessageBody(), "registration")
    
    // Extract and use validation code
    code, err := msg.ExtractValidationCode()
    require.NoError(t, err)
    
    // Verify registration with code
    err = VerifyRegistration(userEmail, code)
    assert.NoError(t, err)
    
    // Cleanup
    mailhog.DeleteAllMessages()
}
```

### Code Extraction Patterns

The `ExtractValidationCode()` method tries multiple patterns automatically:
- 6-digit codes: `123456`
- 4-8 digit codes: `1234`, `12345678`
- 6-character alphanumeric: `ABC123`
- 3 letters + 3 digits: `XYZ789`

For custom patterns, use `ExtractCode(pattern)`:
```go
// Extract UUID
code, _ := msg.ExtractCode(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)

// Extract token
token, _ := msg.ExtractCode(`token=[A-Za-z0-9]{32}`)
```


## Environment Setup

### Development (.env)
```env
# Use MailHog for development
MAILHOG_HOST=localhost
MAILHOG_PORT=1025
MAILHOG_DEFAULT_FROM=noreply@test.local
```

### Production (.env)
```env
# Use Resend for production
RESEND_API_KEY=your_actual_api_key
RESEND_DEFAULT_FROM=noreply@yourdomain.com
```

## API Reference

### EmailData
```go
type EmailData struct {
    From        string            // Sender email (optional, uses default if empty)
    To          []string          // Recipient emails (required)
    Subject     string            // Email subject
    Html        string            // HTML email body
    Text        string            // Plain text body (fallback)
    Cc          []string          // CC recipients
    Bcc         []string          // BCC recipients
    ReplyTo     string            // Reply-To address
    Headers     map[string]string // Custom headers
    Tags        []Tag             // Email tags (provider-specific)
    Attachments []*Attachment     // File attachments
    ScheduledAt string            // Schedule send time (provider-specific)
}
```

### Sender Interface
```go
type Sender interface {
    Send(ctx context.Context, data EmailData) error
}
```

### TemplateData
```go
type TemplateData map[string]interface{}
```
