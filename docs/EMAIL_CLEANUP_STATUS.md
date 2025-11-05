# Email Library Cleanup - Important Notes

## Current Situation

The email library **has been restored** to `core/src/lib/email/` because it is still actively used by production code in the core service.

## Why We Kept It

### Active Usage in Core Service

**File: `core/src/api/controller/index.go`**
- `SendLoginValidationCodeByEmail()` - Sends login validation codes
- `ResetPasswordByEmail()` - Sends password reset emails
- `SendNewPasswordByEmail()` - Sends new password emails

These functions use the email library directly:
```go
provider, err := email.NewProvider(nil)
provider.Send(context.Background(), email.EmailData{...})
```

### Active Usage in Test Files

**Files:**
- `core/test/src/model/admin.go`
- `core/test/src/model/employee.go`
- `core/test/src/model/client.go`

These test files use MailHog client to verify emails were sent:
```go
mailhog, err := email.MailHog()
messages := mailhog.GetMessages(...)
```

## Current Architecture

```
┌─────────────────────────────────────────────────┐
│                 Core Service                     │
│  ┌──────────────────────────────────────────┐  │
│  │  core/src/lib/email/                     │  │
│  │  - sender.go                             │  │
│  │  - resend.go                             │  │
│  │  - mailhog.go                            │  │
│  │  - template.go                           │  │
│  └──────────────────────────────────────────┘  │
│           ↑                                     │
│           │ (direct import)                     │
│  ┌────────┴─────────────────────────────────┐  │
│  │  core/src/api/controller/index.go        │  │
│  │  - SendLoginValidationCodeByEmail()      │  │
│  │  - ResetPasswordByEmail()                │  │
│  │  - SendNewPasswordByEmail()              │  │
│  └──────────────────────────────────────────┘  │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│              Email Service                       │
│  ┌──────────────────────────────────────────┐  │
│  │  email/lib/                              │  │
│  │  - sender.go (copy)                      │  │
│  │  - resend.go (copy)                      │  │
│  │  - mailhog.go (copy)                     │  │
│  │  - template.go (copy)                    │  │
│  └──────────────────────────────────────────┘  │
│           ↑                                     │
│  ┌────────┴─────────────────────────────────┐  │
│  │  email/routes.go                         │  │
│  │  - POST /api/v1/emails/send              │  │
│  │  - POST /api/v1/emails/send-template     │  │
│  │  - POST /api/v1/emails/send-bulk         │  │
│  └──────────────────────────────────────────┘  │
└─────────────────────────────────────────────────┘
```

## What Was Accomplished

✅ **Email Service Created**: Fully functional standalone microservice
✅ **Email Library Duplicated**: Library exists in both places
✅ **REST API Working**: Email service has complete API
✅ **Documentation Complete**: Full documentation for email service

## Migration Plan (Phase 2)

To fully migrate to the email service architecture, follow these steps:

### Step 1: Update Core Controller Functions

Replace direct email library usage with HTTP API calls.

**Before:**
```go
provider, err := email.NewProvider(nil)
err = provider.Send(context.Background(), email.EmailData{
    To:      []string{user_email},
    Subject: renderedEmail.Subject,
    Html:    renderedEmail.HTMLBody,
})
```

**After:**
```go
import (
    "bytes"
    "encoding/json"
    "net/http"
)

payload := map[string]interface{}{
    "to":            user_email,
    "template_name": "login_validation_code",
    "language":      language,
    "data": map[string]interface{}{
        "LoginValidationCode": LoginValidationCode,
    },
}

jsonData, _ := json.Marshal(payload)
resp, err := http.Post(
    "http://localhost:4002/api/v1/emails/send-template",
    "application/json",
    bytes.NewBuffer(jsonData),
)
```

### Step 2: Update Test Files

Update test files to either:
- Mock the email service API
- Use the email service's test endpoints
- Keep using MailHog directly for integration tests

### Step 3: Create Email Client Library (Optional)

Create a reusable Go client library:

**File: `core/src/lib/emailclient/client.go`**
```go
package emailclient

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type Client struct {
    BaseURL string
}

func New() *Client {
    return &Client{
        BaseURL: "http://localhost:4002",
    }
}

func (c *Client) SendTemplate(to, templateName, language string, data map[string]interface{}) error {
    payload := map[string]interface{}{
        "to":            to,
        "template_name": templateName,
        "language":      language,
        "data":          data,
    }
    
    jsonData, _ := json.Marshal(payload)
    resp, err := http.Post(
        c.BaseURL+"/api/v1/emails/send-template",
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

### Step 4: Environment Configuration

Add to `core/.env`:
```env
EMAIL_SERVICE_URL=http://localhost:4002
```

### Step 5: Remove Email Library from Core

After all code is updated:
```bash
Remove-Item -Path "c:\code\mynute-go\core\src\lib\email" -Recurse -Force
Remove-Item -Path "c:\code\mynute-go\core\static\email" -Recurse -Force
Remove-Item -Path "c:\code\mynute-go\core\translation\email" -Recurse -Force
```

## Files That Need Updates

### Production Code
- [ ] `core/src/api/controller/index.go`
  - [ ] `SendLoginValidationCodeByEmail()`
  - [ ] `ResetPasswordByEmail()`
  - [ ] `SendNewPasswordByEmail()`

### Test Code
- [ ] `core/test/src/model/admin.go`
- [ ] `core/test/src/model/employee.go`
- [ ] `core/test/src/model/client.go`

## Estimated Effort

- **Small**: ~2-4 hours
  - Update 3 controller functions
  - Update test files
  - Testing
  
- **Medium** (with client library): ~4-6 hours
  - Create email client library
  - Update all usages
  - Add error handling
  - Testing

## Benefits After Migration

✅ **True Microservice Separation**: Email logic completely isolated
✅ **Independent Scaling**: Scale email service based on email volume
✅ **No Duplicate Code**: Single source of truth for email functionality
✅ **Easier Provider Switching**: Change email provider without touching core
✅ **Better Testing**: Mock email service in tests
✅ **Centralized Email Logs**: All email activity in one service

## Current Status

**Email Library Location**: 
- ✅ `email/lib/` - Used by email service
- ✅ `core/src/lib/email/` - Still used by core service (TEMPORARY)

**Email Templates**:
- Templates were NOT restored to core (they're only in email service now)
- Core controller uses template renderer which needs templates
- This may cause errors if core tries to render templates

**Email Translations**:
- Translations were NOT restored to core
- Same issue as templates

## Immediate Actions Needed

### Option 1: Restore Templates and Translations (Quick Fix)
```bash
# Restore templates
Copy-Item -Path "c:\code\mynute-go\email\static\email\*" -Destination "c:\code\mynute-go\core\static\email\" -Recurse -Force

# Restore translations  
Copy-Item -Path "c:\code\mynute-go\email\translation\email\*" -Destination "c:\code\mynute-go\core\translation\email\" -Recurse -Force
```

### Option 2: Update Code to Use Email Service API (Proper Fix)
Follow the migration plan above.

## Recommendation

**Short-term**: Restore templates and translations to core so existing code works
**Long-term**: Update core service to use email service API and remove duplicates

## Questions?

See:
- `docs/EMAIL_MIGRATION.md` - Full migration guide
- `docs/EMAIL_API_REFERENCE.md` - API usage examples
- `email/README.md` - Email service documentation
