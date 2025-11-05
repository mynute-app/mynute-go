# Email Service Migration Guide

This document explains how email functionality was migrated from the core service to a dedicated email microservice.

## Overview

The email functionality has been extracted from the core service into a dedicated microservice to:
- Improve separation of concerns
- Enable independent scaling of email operations
- Isolate email provider dependencies
- Simplify testing and maintenance

## What Was Moved

### 1. Email Library (`core/src/lib/email/` → `email/lib/`)

**Files Moved:**
- `sender.go` - Provider interface and factory
- `resend.go` - Resend API implementation
- `mailhog.go` - MailHog SMTP implementation
- `template.go` - Template renderer with i18n support
- All test files (`*_test.go`)

### 2. Email Templates (`core/static/email/` → `email/static/email/`)

**Templates Moved:**
- `email_verification_code.html`
- `login_validation_code.html`
- `new_password.html`

### 3. Email Translations (`core/translation/email/` → `email/translation/email/`)

**Translation Files Moved:**
- `email_verification_code.json`
- `login_validation_code.json`
- `new_password.json`

### 4. Environment Configuration

**Email-specific variables moved to `email/.env`:**
- `EMAIL_SERVICE_PORT`
- `EMAIL_PROVIDER`
- `RESEND_API_KEY`
- `MAILHOG_HOST`, `MAILHOG_PORT`
- `SMTP_*` settings

## Architecture Changes

### Before Migration

```
core/
├── src/
│   └── lib/
│       └── email/          # Email library
├── static/
│   └── email/              # Email templates
└── translation/
    └── email/              # Email translations
```

All services imported from `mynute-go/services/core/src/lib/email`.

### After Migration

```
email/                      # Dedicated microservice
├── lib/                    # Email library (moved)
├── static/email/           # Templates (moved)
├── translation/email/      # Translations (moved)
├── routes.go               # REST API handlers
└── server.go               # Service initialization
```

Email functionality is now accessed via HTTP API calls to the email service.

## API Changes

### Old Approach (Direct Library Usage)
```go
// In any service
import "mynute-go/services/core/src/lib/email"

provider, _ := email.NewProvider(nil)
emailData := email.EmailData{
    To:      []string{"user@example.com"},
    Subject: "Hello",
    Html:    "<p>World</p>",
}
provider.Send(ctx, emailData)
```

### New Approach (HTTP API)
```go
// In any service
import "net/http"

payload := map[string]interface{}{
    "to":      "user@example.com",
    "subject": "Hello",
    "body":    "<p>World</p>",
    "is_html": true,
}

resp, err := http.Post(
    "http://localhost:4002/api/v1/emails/send",
    "application/json",
    bytes.NewBuffer(jsonPayload),
)
```

Or use template emails:
```go
payload := map[string]interface{}{
    "to":            "user@example.com",
    "template_name": "email_verification_code",
    "language":      "en",
    "data": map[string]interface{}{
        "code":     "123456",
        "username": "John Doe",
    },
}

resp, err := http.Post(
    "http://localhost:4002/api/v1/emails/send-template",
    "application/json",
    bytes.NewBuffer(jsonPayload),
)
```

## Migration Steps Performed

1. ✅ Created email microservice structure
   - Created `email/server.go`
   - Created `email/routes.go`
   - Created `email/Dockerfile`
   - Created `email/docker-compose.{dev,prod}.yml`

2. ✅ Moved email library files
   - Copied `core/src/lib/email/*.go` → `email/lib/`
   - Package declarations remain as `package email`

3. ✅ Moved email templates
   - Copied `core/static/email/*` → `email/static/email/`

4. ✅ Moved email translations
   - Copied `core/translation/email/*` → `email/translation/email/`

5. ✅ Created environment configuration
   - Created `email/.env.example`
   - Created `email/.env` (gitignored)

6. ✅ Implemented REST API
   - `POST /api/v1/emails/send` - Send simple email
   - `POST /api/v1/emails/send-template` - Send template email
   - `POST /api/v1/emails/send-bulk` - Send bulk emails

7. ✅ Integrated email library with API
   - Initialize email provider on server startup
   - Initialize template renderer
   - Wire up handlers to use email library

8. ✅ Updated main.go
   - Added email service to concurrent startup
   - Service runs on port 4002

## Updating Services to Use Email API

### For Auth Service

Replace direct email library calls with HTTP API calls:

```go
// Old
import "mynute-go/services/core/src/lib/email"

// New
import (
    "bytes"
    "encoding/json"
    "net/http"
)

func sendVerificationEmail(userEmail, code string) error {
    payload := map[string]interface{}{
        "to":            userEmail,
        "template_name": "email_verification_code",
        "language":      "en",
        "data": map[string]interface{}{
            "code": code,
        },
    }
    
    jsonData, _ := json.Marshal(payload)
    resp, err := http.Post(
        "http://localhost:4002/api/v1/emails/send-template",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to send email: status %d", resp.StatusCode)
    }
    
    return nil
}
```

### For Core Service

Same approach - replace library imports with HTTP calls.

## Environment Variables

### Email Service (.env)
```env
APP_ENV=dev
EMAIL_SERVICE_PORT=4002
EMAIL_PROVIDER=mailhog
MAILHOG_HOST=localhost
MAILHOG_PORT=1025
MAILHOG_DEFAULT_FROM=noreply@test.local
```

### Other Services
Remove email-specific variables. Only need:
```env
EMAIL_SERVICE_URL=http://localhost:4002
```

## Running the Services

### Start All Services
```bash
go run main.go
```

### Start Email Service Only
```bash
go run cmd/email-service/main.go
```

### Docker Compose
```bash
docker-compose -f email/docker-compose.dev.yml up
```

## Testing

### Test Email Sending
```bash
curl -X POST http://localhost:4002/api/v1/emails/send \
  -H "Content-Type: application/json" \
  -d '{
    "to": "test@example.com",
    "subject": "Test Email",
    "body": "<h1>Hello World</h1>",
    "is_html": true
  }'
```

### Test Template Email
```bash
curl -X POST http://localhost:4002/api/v1/emails/send-template \
  -H "Content-Type: application/json" \
  -d '{
    "to": "test@example.com",
    "template_name": "email_verification_code",
    "language": "en",
    "data": {
      "code": "123456"
    }
  }'
```

### Check Health
```bash
curl http://localhost:4002/health
```

## Rollback Plan

If issues arise:

1. Revert main.go to not start email service
2. Update services to use direct library imports again
3. Point imports back to `mynute-go/services/core/src/lib/email`
4. Restart services

## Benefits

✅ **Separation of Concerns**: Email logic isolated from business logic
✅ **Independent Scaling**: Scale email service independently
✅ **Better Testing**: Test email functionality in isolation
✅ **Provider Flexibility**: Switch email providers without affecting other services
✅ **Reduced Coupling**: Services only depend on HTTP API contract
✅ **Easier Maintenance**: Email-related changes confined to one service

## Next Steps

1. Add API authentication (API keys)
2. Implement rate limiting
3. Add email queuing for bulk operations
4. Add retry logic with exponential backoff
5. Implement email delivery tracking
6. Add webhooks for delivery status
7. Create email service client library for Go services

## Cleanup (Optional)

After confirming all services work with the email API:

1. Remove email library from core service:
   ```bash
   rm -rf core/src/lib/email
   ```

2. Remove email templates from core:
   ```bash
   rm -rf core/static/email
   ```

3. Remove email translations from core:
   ```bash
   rm -rf core/translation/email
   ```

## Support

For issues or questions:
- Check email service logs: `docker logs mynute-email`
- Verify email service is running: `curl http://localhost:4002/health`
- Check MailHog UI (dev): `http://localhost:8025`
- Review email service README: `email/README.md`
