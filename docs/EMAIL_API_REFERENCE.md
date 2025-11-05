# Email Service Quick Reference

Quick guide for using the email microservice in your code.

## Service Information

- **URL**: `http://localhost:4002`
- **Port**: 4002
- **Health Check**: `GET /health`

## API Endpoints

### 1. Send Simple Email

**Endpoint**: `POST /api/v1/emails/send`

**Request:**
```json
{
  "to": "user@example.com",
  "subject": "Welcome",
  "body": "Hello World!",
  "cc": ["copy@example.com"],
  "bcc": ["blind@example.com"],
  "is_html": true
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "Email sent successfully",
  "to": "user@example.com"
}
```

**Go Example:**
```go
import (
    "bytes"
    "encoding/json"
    "net/http"
)

func sendEmail(to, subject, body string) error {
    payload := map[string]interface{}{
        "to":      to,
        "subject": subject,
        "body":    body,
        "is_html": true,
    }
    
    jsonData, _ := json.Marshal(payload)
    resp, err := http.Post(
        "http://localhost:4002/api/v1/emails/send",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    return nil
}
```

---

### 2. Send Template Email

**Endpoint**: `POST /api/v1/emails/send-template`

**Request:**
```json
{
  "to": "user@example.com",
  "template_name": "email_verification_code",
  "language": "en",
  "data": {
    "code": "123456",
    "username": "John Doe"
  }
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "Template email sent successfully",
  "to": "user@example.com"
}
```

**Go Example:**
```go
func sendVerificationEmail(email, code, username string) error {
    payload := map[string]interface{}{
        "to":            email,
        "template_name": "email_verification_code",
        "language":      "en",
        "data": map[string]interface{}{
            "code":     code,
            "username": username,
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
    
    return nil
}
```

---

### 3. Send Bulk Emails

**Endpoint**: `POST /api/v1/emails/send-bulk`

**Request:**
```json
{
  "recipients": [
    "user1@example.com",
    "user2@example.com",
    "user3@example.com"
  ],
  "subject": "Newsletter",
  "body": "Check out our latest updates!",
  "is_html": false
}
```

**Response (200):**
```json
{
  "success": true,
  "total": 3,
  "sent": 3,
  "failed": 0
}
```

**Response (206 - Partial Success):**
```json
{
  "success": true,
  "total": 3,
  "sent": 2,
  "failed": 1,
  "failed_recipients": ["user3@example.com"]
}
```

**Go Example:**
```go
func sendNewsletter(recipients []string, subject, body string) error {
    payload := map[string]interface{}{
        "recipients": recipients,
        "subject":    subject,
        "body":       body,
        "is_html":    true,
    }
    
    jsonData, _ := json.Marshal(payload)
    resp, err := http.Post(
        "http://localhost:4002/api/v1/emails/send-bulk",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    return nil
}
```

---

## Available Templates

| Template Name | Languages | Variables |
|--------------|-----------|-----------|
| `email_verification_code` | en | code, username |
| `login_validation_code` | en | code, username |
| `new_password` | en | password, username |

## Adding New Templates

1. Create HTML file: `email/static/email/my_template.html`
```html
<!DOCTYPE html>
<html>
<head><title>{{.subject}}</title></head>
<body>
    <h1>{{.title}}</h1>
    <p>Hello {{.username}},</p>
    <p>{{.message}}</p>
</body>
</html>
```

2. Create translation file: `email/translation/email/my_template.json`
```json
{
  "subject": "My Email Subject",
  "title": "Welcome!",
  "message": "This is a template message"
}
```

3. Use in API:
```json
{
  "to": "user@example.com",
  "template_name": "my_template",
  "language": "en",
  "data": {
    "username": "John"
  }
}
```

## Email Providers

### Development (MailHog)
```env
APP_ENV=dev
EMAIL_PROVIDER=mailhog
MAILHOG_HOST=localhost
MAILHOG_PORT=1025
```

View emails at: http://localhost:8025

### Production (Resend)
```env
APP_ENV=prod
EMAIL_PROVIDER=resend
RESEND_API_KEY=re_xxxxxxxxxxxxx
RESEND_DEFAULT_FROM=noreply@yourdomain.com
```

### Generic SMTP
```env
EMAIL_PROVIDER=smtp
SMTP_HOST=smtp.gmail.com
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password
SMTP_PORT=587
EMAIL_FROM=noreply@yourdomain.com
```

## Error Handling

**Common Errors:**

| Status | Error | Cause |
|--------|-------|-------|
| 400 | Invalid request body | Malformed JSON |
| 500 | Failed to send email | Provider error |
| 500 | Failed to render email template | Template not found or data error |

**Example Error Response:**
```json
{
  "error": "Failed to send email",
  "details": "authentication failed"
}
```

## Testing

### Using curl
```bash
# Test simple email
curl -X POST http://localhost:4002/api/v1/emails/send \
  -H "Content-Type: application/json" \
  -d '{"to":"test@test.com","subject":"Test","body":"Hello","is_html":false}'

# Test template email
curl -X POST http://localhost:4002/api/v1/emails/send-template \
  -H "Content-Type: application/json" \
  -d '{"to":"test@test.com","template_name":"email_verification_code","language":"en","data":{"code":"123456"}}'

# Check health
curl http://localhost:4002/health
```

### Using Go HTTP Client

**With Context and Timeout:**
```go
import (
    "context"
    "time"
)

func sendEmailWithTimeout(to, subject, body string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    payload := map[string]interface{}{
        "to":      to,
        "subject": subject,
        "body":    body,
        "is_html": true,
    }
    
    jsonData, _ := json.Marshal(payload)
    req, _ := http.NewRequestWithContext(
        ctx,
        "POST",
        "http://localhost:4002/api/v1/emails/send",
        bytes.NewBuffer(jsonData),
    )
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("email sending failed: status %d", resp.StatusCode)
    }
    
    return nil
}
```

## Troubleshooting

### Email Not Sending

1. Check service is running:
   ```bash
   curl http://localhost:4002/health
   ```

2. Check logs:
   ```bash
   # If running via main.go
   # Logs appear in console
   
   # If running via Docker
   docker logs mynute-email
   ```

3. Verify environment variables in `email/.env`

### MailHog Not Receiving Emails

1. Check MailHog is running:
   ```bash
   curl http://localhost:8025
   ```

2. Verify MailHog settings in `email/.env`:
   ```env
   MAILHOG_HOST=localhost
   MAILHOG_PORT=1025
   ```

### Template Errors

1. Ensure template file exists in `email/static/email/`
2. Ensure translation file exists in `email/translation/email/`
3. Check template syntax in HTML file
4. Verify all variables in translation JSON match template

## Performance Tips

1. **Bulk Sending**: Use `/send-bulk` instead of multiple `/send` calls
2. **Timeouts**: Set appropriate timeouts for HTTP requests
3. **Async**: Consider sending emails asynchronously in goroutines
4. **Retry**: Implement retry logic for failed sends

## Support

- Email Service Docs: `email/README.md`
- Migration Guide: `docs/EMAIL_MIGRATION.md`
- Architecture: `docs/MICROSERVICES_ARCHITECTURE.md`
- Logs: Check service console output or Docker logs
