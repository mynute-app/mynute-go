# Login by Email Code Implementation

This document describes the implementation of the login by email code functionality for both clients and employees.

## Overview

The login by email code feature allows users (clients and employees) to receive a 6-digit validation code via email and use it to authenticate instead of their regular password. This provides an alternative authentication method that's useful for temporary access or enhanced security.

## Implementation Details

### 1. Database Schema Changes

Added new fields to both `Client` and `Employee` models:
- `LoginValidationCode *string` - Stores the 6-digit validation code (nullable)
- `LoginValidationExpiry *time.Time` - Stores when the code expires (nullable)

These fields are:
- Nullable to allow clearing them after use
- Temporary (automatically cleared after successful login or expiration)
- Set to expire 15 minutes after generation

### 2. DTOs (Data Transfer Objects)

Created a generic DTO for login by email code:
```go
// In core/src/config/api/dto/auth.go
type LoginByEmailCode struct {
    Email string `json:"email" example:"john.doe@example.com"`
    Code  string `json:"code" example:"123456"`
}
```

### 3. Controller Functions

#### Generic Controllers (core/src/controller/index.go)
- `LoginByEmailCode(user_type string, model any, c *fiber.Ctx) (string, error)` - Generic login by code function
- Updated `SendLoginValidationCodeByEmail()` - Now stores the generated code in the database with expiration

#### Client-Specific Controllers (core/src/controller/client.go)
- `LoginClientByEmailCode(c *fiber.Ctx) error` - Client login by email code endpoint
- `SendClientLoginValidationCodeByEmail(c *fiber.Ctx) error` - Already existed, now works with updated backend

#### Employee-Specific Controllers (core/src/controller/employee.go)
- `LoginEmployeeByEmailCode(c *fiber.Ctx) error` - Employee login by email code endpoint
- `SendEmployeeLoginValidationCodeByEmail(c *fiber.Ctx) error` - New function to send codes to employees

### 4. Service Layer

#### New Service Method (core/src/service/factory.go)
```go
func (s *service) LoginByEmailCode(user_type string) (string, error)
```

This method:
1. Parses the email and code from the request body
2. Finds the user by email
3. Validates the user is verified
4. Checks if a validation code exists and hasn't expired
5. Compares the provided code with the stored code
6. Clears the validation code after successful authentication
7. Generates and returns a JWT token

### 5. API Endpoints

#### Client Endpoints
- `POST /client/send-login-code/email/{email}?lang=en` - Send validation code to client
- `POST /client/login/code` - Login client with email and code

#### Employee Endpoints  
- `POST /employee/send-login-code/email/{email}?lang=en` - Send validation code to employee
- `POST /employee/login/code` - Login employee with email and code

### 6. Security Features

- **Code Expiration**: Codes expire after 15 minutes
- **One-Time Use**: Codes are automatically cleared after successful login
- **User Verification**: Only verified users can use this login method
- **Code Validation**: Strict comparison of provided vs stored codes

### 7. Error Handling

The implementation handles various error cases:
- User not found
- User not verified
- No validation code found
- Validation code expired
- Invalid validation code
- Database errors during code storage/retrieval

## Usage Flow

1. **Request Code**: User requests a login code via email
   ```bash
   POST /client/send-login-code/email/user@example.com?lang=en
   ```

2. **Receive Code**: User receives a 6-digit code via email (expires in 15 minutes)

3. **Login with Code**: User submits email and code to login
   ```bash
   POST /client/login/code
   Content-Type: application/json
   
   {
     "email": "user@example.com",
     "code": "123456"
   }
   ```

4. **Success**: User receives JWT token in response header and code is cleared from database

## Files Modified

- `core/src/config/api/dto/auth.go` - Added `LoginByEmailCode` DTO
- `core/src/config/db/model/client.go` - Added validation code fields to `ClientMeta`
- `core/src/config/db/model/employee.go` - Added validation code fields to `Employee`
- `core/src/controller/index.go` - Added `LoginByEmailCode()` and updated `SendLoginValidationCodeByEmail()`
- `core/src/controller/client.go` - Implemented `LoginClientByEmailCode()` and registered it
- `core/src/controller/employee.go` - Added `LoginEmployeeByEmailCode()` and `SendEmployeeLoginValidationCodeByEmail()`
- `core/src/service/factory.go` - Added `LoginByEmailCode()` service method

## Database Migration Required

The new fields in the Client and Employee models will require a database migration to add:
- `login_validation_code VARCHAR(6) NULL`
- `login_validation_expiry TIMESTAMP NULL`

## Testing

To test the implementation:

1. Ensure database has the new fields
2. Start the application
3. Send a validation code request
4. Check email for the 6-digit code
5. Use the code to login within 15 minutes
6. Verify JWT token is returned and code is cleared