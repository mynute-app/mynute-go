# UserMeta Migration Fix - Summary

This document describes the fixes applied to handle the migration of login validation code fields and Design fields from direct model fields to nested fields within the `UserMeta` JSON structure.

## Changes Made

### 1. Updated Service Layer (`core/src/service/factory.go`)

**Fixed**: `LoginByEmailCode()` method
- Added import for `mJSON "mynute-go/core/src/config/db/model/json"`
- Updated field access from direct fields to Meta field:
  - `LoginValidationCode` → `Meta.LoginValidationCode`
  - `LoginValidationExpiry` → `Meta.LoginValidationExpiry`
- Changed validation logic to work with pointer fields in UserMeta struct
- Updated database clearing logic to work with the Meta field

### 2. Updated Controller Layer (`core/src/controller/index.go`)

**Fixed**: `SendLoginValidationCodeByEmail()` function
- Updated to store validation code in `Meta.LoginValidationCode` and `Meta.LoginValidationExpiry`
- Changed from direct field access to Meta field manipulation

**Fixed**: `UpdateImagesById()` and `DeleteImageById()` functions
- Added dual support for models with direct `Design` field (Company, Branch, Service) and models with `Meta.Design` field (Client, Employee)
- Implemented conditional logic to handle both cases:
  - **Direct Design**: Company, Branch, Service models
  - **Meta Design**: Client, Employee models
- Updated database read/write operations to work with both patterns

### 3. Updated Test Files

**Fixed**: Client test files
- `test/e2e/client_test.go`: Updated `.Design` → `.Meta.Design`
- `test/src/model/client.go`: Updated image handling to use `.Meta.Design.Images`

**Fixed**: Employee test files
- `test/e2e/employee_test.go`: Updated `.Design` → `.Meta.Design`
- `test/src/model/employee.go`: Updated image handling to use `.Meta.Design.Images`

## Model Structure Changes

### Before (Direct Fields)
```go
type Client struct {
    // ... other fields ...
    LoginValidationCode   *string            `gorm:"type:varchar(6)"`
    LoginValidationExpiry *time.Time         `gorm:"type:timestamp"`
    Design                mJSON.DesignConfig `gorm:"type:jsonb"`
}
```

### After (UserMeta Structure)
```go
type Client struct {
    // ... other fields ...
    Meta mJSON.UserMeta `gorm:"type:jsonb"`
}

type UserMeta struct {
    LoginValidationCode   *string      `json:"login_validation_code,omitempty"`
    LoginValidationExpiry *time.Time   `json:"login_validation_expiry,omitempty"`
    Design                DesignConfig `json:"design"`
}
```

## Backward Compatibility

The implementation maintains backward compatibility by:

1. **Dual Model Support**: The image upload/delete functions (`UpdateImagesById`, `DeleteImageById`) now support both:
   - Models with direct `Design` field (Company, Branch, Service)
   - Models with `Meta.Design` field (Client, Employee)

2. **Runtime Detection**: Uses reflection to detect which pattern the model uses and handles accordingly.

3. **Database Operations**: Automatically chooses the correct database field (`design` vs `meta`) based on the model structure.

## Testing

All test files have been updated to use the new `.Meta.Design` pattern for Client and Employee models. The tests should now pass without issues.

## Notes

- Only Client and Employee models were migrated to use the UserMeta structure
- Company, Branch, and Service models still use direct Design fields
- The login by email code functionality works correctly with the new Meta structure
- All image upload/delete operations work with both old and new model patterns

## Files Modified

### Core Files
- `core/src/service/factory.go` - Updated LoginByEmailCode method
- `core/src/controller/index.go` - Updated SendLoginValidationCodeByEmail, UpdateImagesById, DeleteImageById

### Test Files
- `test/e2e/client_test.go` - Updated Design field access
- `test/e2e/employee_test.go` - Updated Design field access
- `test/src/model/client.go` - Updated Design field access in image operations
- `test/src/model/employee.go` - Updated Design field access in image operations

All compilation errors have been resolved and the code should now work correctly with the new UserMeta structure.