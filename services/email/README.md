# Email Microservice

A dedicated microservice for handling all email operations in the mynute-go application.

## Overview

The Email Service is responsible for **email delivery only**. It does not store templates or business logic - that responsibility belongs to the Core Service. 

**Email Service provides two delivery methods:**
1. **Send**: Direct email to one or more recipients (each with their own CC/BCC)
2. **Template Merge Send**: Receives template HTML and translations from calling service (e.g., Core), merges data, and sends

## Architecture Principle

**Separation of Concerns:**
- **Email Service**: Handles email **delivery** (sending via providers like MailHog/Resend)
- **Core Service**: Manages email **templates and translations** stored at `services/core/api/view/`

## Features

- **Multiple Email Providers**: Support for Resend, MailHog, and SMTP
- **Template Rendering**: Multi-language HTML email templates with translations
- **Template Merge Endpoint**: Accepts template HTML + translations from other services
- **Bulk Email Sending**: Send emails to multiple recipients
- **RESTful API**: Simple HTTP/REST interface
- **Health Checks**: Built-in health check endpoint
- **Attachments Support**: Send emails with file attachments
- **Fully Independent**: No dependencies on Core or Auth services

## Service Details

- **Port**: `4002`
- **Swagger**: `http://localhost:4002/swagger/index.html`
- **Health Check**: `http://localhost:4002/health`

## Directory Structure

```
email/
├── api/
│   ├── controller/          # HTTP controllers
│   │   └── email.go        # Email endpoints
│   └── lib/                # Email library
│       ├── sender.go       # Provider interface and factory
│       ├── resend.go       # Resend API implementation
│       ├── mailhog.go      # MailHog SMTP implementation
│       └── template.go     # Template renderer
├── config/
│   └── dto/                # Data transfer objects
│       └── email.go        # Request/response DTOs
├── lib/
│   └── env.go              # Environment loading
├── docs/                   # Swagger documentation
├── routes.go               # API route definitions
├── server.go               # Server initialization
├── Dockerfile              # Docker configuration
├── docker-compose.dev.yml  # Development compose
├── docker-compose.prod.yml # Production compose
├── loki-config.yaml        # Logging configuration
├── prometheus.yml          # Metrics configuration
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

Send emails to one or more recipients with plain text or HTML content. Each recipient can have their own CC and BCC lists.

**Single Recipient Example:**
```json
{
  "subject": "Welcome",
  "body": "Hello, World!",
  "is_html": true,
  "recipients": [
    {
      "to": ["user@example.com"],
      "cc": ["cc@example.com"],
      "bcc": ["bcc@example.com"]
    }
  ]
}
```

**Multiple Recipients Example:**
```json
{
  "subject": "Newsletter",
  "body": "<h1>Monthly Update</h1><p>Here's what's new...</p>",
  "is_html": true,
  "recipients": [
    {
      "to": ["user1@example.com"],
      "cc": ["manager1@example.com"]
    },
    {
      "to": ["user2@example.com", "user3@example.com"],
      "bcc": ["archive@example.com"]
    }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "message": "All emails sent successfully",
  "total": 2,
  "sent": 2,
  "failed": 0,
  "results": [
    {
      "to": ["user1@example.com"],
      "success": true
    },
    {
      "to": ["user2@example.com", "user3@example.com"],
      "success": true
    }
  ]
}
```

### Send Template Merge Email
```http
POST /api/v1/emails/send-template-merge
```

**Purpose**: Receives template HTML and translations from calling services (like Core) for rendering and delivery.

**Use Case**: Core Service stores templates at `services/core/api/view/html/email/*.html` and translations at `services/core/api/view/translation/email/*.json`. Core reads these files and sends them to Email Service for merging and delivery.

**Request Body:**
```json
{
  "to": ["user@example.com"],
  "template_html": "<html><body><h1>{{.greeting}}</h1><p>{{.message}}</p></body></html>",
  "translations": {
    "subject": "Email Verification",
    "greeting": "Hello",
    "message": "Please verify your email"
  },
  "data": {
    "username": "John Doe"
  },
  "cc": ["cc@example.com"],
  "bcc": ["bcc@example.com"]
}
```

**How it works:**
1. Calling service (Core) loads template HTML from its filesystem
2. Calling service loads translations for requested language
3. Calling service sends both to `/send-template-merge`
4. Email Service merges `translations` + `data` into `template_html`
5. Email Service sends the rendered email via configured provider

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

### API Documentation

Interactive Swagger documentation is available at:
```
http://localhost:4002/swagger/index.html
```

### Regenerating Swagger Docs
```bash
swag init -g cmd/email-service/main.go -o email/docs --parseDependency --parseInternal
```

Or use the Makefile:
```bash
make swagger-email
```

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
