# Email Microservice

A dedicated microservice for handling all email operations in the mynute-go application.

## Features

- **Multiple Email Providers**: Support for Resend, MailHog, and SMTP
- **Template Rendering**: Multi-language HTML email templates with translations
- **Bulk Email Sending**: Send emails to multiple recipients
- **RESTful API**: Simple HTTP/REST interface
- **Health Checks**: Built-in health check endpoint
- **Attachments Support**: Send emails with file attachments

## Architecture

The email service is built with:
- **Fiber v2**: Fast HTTP web framework
- **Provider Pattern**: Pluggable email providers
- **Template Engine**: HTML template rendering with i18n support

## Directory Structure

```
email/
├── lib/                     # Email library
│   ├── sender.go           # Provider interface and factory
│   ├── resend.go           # Resend API implementation
│   ├── mailhog.go          # MailHog SMTP implementation
│   └── template.go         # Template renderer
├── static/                  # Static assets
│   └── email/              # Email templates
│       ├── email_verification_code.html
│       ├── login_validation_code.html
│       └── new_password.html
├── translation/             # i18n translations
│   └── email/              # Email translations
│       ├── email_verification_code.json
│       ├── login_validation_code.json
│       └── new_password.json
├── routes.go               # API route handlers
├── server.go               # Server initialization
├── Dockerfile              # Docker configuration
├── docker-compose.dev.yml  # Development compose
├── docker-compose.prod.yml # Production compose
├── .env                    # Environment variables (gitignored)
└── .env.example            # Environment template
```

## API Endpoints

### Health Check
```http
GET /health
```

Returns service health status.

### Send Email
```http
POST /api/v1/emails/send
```

Send a single email.

**Request Body:**
```json
{
  "to": "user@example.com",
  "subject": "Welcome",
  "body": "Hello, World!",
  "cc": ["cc@example.com"],
  "bcc": ["bcc@example.com"],
  "is_html": true
}
```

### Send Template Email
```http
POST /api/v1/emails/send-template
```

Send an email using a template with translations.

**Request Body:**
```json
{
  "to": "user@example.com",
  "template_name": "email_verification_code",
  "language": "en",
  "data": {
    "code": "123456",
    "username": "John Doe"
  },
  "cc": [],
  "bcc": []
}
```

### Send Bulk Email
```http
POST /api/v1/emails/send-bulk
```

Send emails to multiple recipients.

**Request Body:**
```json
{
  "recipients": [
    "user1@example.com",
    "user2@example.com"
  ],
  "subject": "Newsletter",
  "body": "Welcome to our newsletter!",
  "is_html": true
}
```

## Email Providers

### Resend (Production)
Recommended for production environments.

```env
EMAIL_PROVIDER=resend
RESEND_API_KEY=your_api_key
RESEND_DEFAULT_FROM=noreply@yourdomain.com
```

### MailHog (Development)
Local SMTP testing server.

```env
EMAIL_PROVIDER=mailhog
MAILHOG_HOST=localhost
MAILHOG_PORT=1025
MAILHOG_DEFAULT_FROM=noreply@test.local
```

### SMTP (Generic)
Use any SMTP server.

```env
EMAIL_PROVIDER=smtp
SMTP_HOST=smtp.gmail.com
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password
SMTP_PORT=587
EMAIL_FROM=noreply@yourdomain.com
```

## Running the Service

### Standalone
```bash
# Using Go
go run cmd/email-service/main.go

# Using Docker
docker-compose -f email/docker-compose.dev.yml up
```

### With All Services
```bash
# From root directory
go run main.go
```

## Environment Variables

Copy `.env.example` to `.env` and configure:

```env
APP_ENV=dev
EMAIL_SERVICE_PORT=4002
EMAIL_PROVIDER=mailhog
MAILHOG_HOST=localhost
MAILHOG_PORT=1025
```

See `.env.example` for all configuration options.

## Email Templates

Templates are located in `static/email/` and use Go's `html/template` syntax.

### Creating a New Template

1. Create HTML template file: `static/email/my_template.html`
```html
<!DOCTYPE html>
<html>
<head><title>{{.subject}}</title></head>
<body>
    <h1>{{.greeting}}</h1>
    <p>{{.message}}</p>
</body>
</html>
```

2. Create translation file: `translation/email/my_template.json`
```json
{
  "subject": "My Email Subject",
  "greeting": "Hello!",
  "message": "This is a test message"
}
```

3. Use the template:
```http
POST /api/v1/emails/send-template
{
  "to": "user@example.com",
  "template_name": "my_template",
  "language": "en",
  "data": {
    "custom_field": "Custom Value"
  }
}
```

## Development

### Running Tests
```bash
cd email/lib
go test -v
```

### Adding a New Provider

1. Create provider file: `lib/my_provider.go`
2. Implement the `Sender` interface
3. Add to factory in `lib/sender.go`

Example:
```go
func MyProvider() (Sender, error) {
    return &myProviderClient{}, nil
}
```

## Docker

### Development
```bash
docker-compose -f email/docker-compose.dev.yml up
```

### Production
```bash
docker-compose -f email/docker-compose.prod.yml up -d
```

## Security

- Use API keys for authentication (planned)
- Never commit `.env` files
- Use secure SMTP connections (TLS/SSL)
- Validate all email addresses
- Rate limit email sending (planned)

## Monitoring

The service exposes metrics at:
- Health: `http://localhost:4002/health`
- Logs: JSON-formatted structured logs

## License

See root LICENSE file.
