package controller

import (
	"encoding/json"
	"fmt"
	"mynute-go/auth/config/namespace"
	DTO "mynute-go/auth/dto"
	"mynute-go/auth/handler"
	"mynute-go/auth/lib"
	"mynute-go/auth/model"
	mJSON "mynute-go/auth/model/json"
	"reflect"
	"time"

	"github.com/gofiber/fiber/v2"
)

// =====================
// CLIENT AUTH
// =====================

// LoginClientByPassword logs a client in
//
//	@Summary		Login client
//	@Description	Log in a client using password
//	@Tags			Client/Auth
//	@Accept			json
//	@Produce		json
//	@Param			client	body	DTO.LoginClient	true	"Client credentials"
//	@Success		200		"Token returned in X-Auth-Token header"
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/auth/client/login [post]
func LoginClientByPassword(c *fiber.Ctx) error {
	token, err := LoginByPassword(namespace.ClientKey.Name, c)
	if err != nil {
		return err
	}
	c.Response().Header.Set(namespace.HeadersKey.Auth, token)
	return nil
}

// LoginClientByEmailCode logs in a client using email and validation code
//
//	@Summary		Login client by email code
//	@Description	Login client using email and validation code
//	@Tags			Client/Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body	DTO.LoginByEmailCode	true	"Login credentials"
//	@Success		200		"Token returned in X-Auth-Token header"
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/auth/client/login-with-code [post]
func LoginClientByEmailCode(c *fiber.Ctx) error {
	token, err := LoginByEmailCode(namespace.ClientKey.Name, c)
	if err != nil {
		return err
	}
	c.Response().Header.Set(namespace.HeadersKey.Auth, token)
	return nil
}

// SendClientLoginValidationCodeByEmail sends a login validation code to a client's email
//
//	@Summary		Send client login validation code by email
//	@Description	Send a login validation code to a client's email
//	@Tags			Client/Auth
//	@Param			email	path	string	true	"Client Email"
//	@Query			language	query	string	false	"Language code (default: en)"
//	@Produce		json
//	@Success		200		{object}	nil
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/auth/client/send-login-code/email/{email} [post]
func SendClientLoginValidationCodeByEmail(c *fiber.Ctx) error {
	if err := SendLoginValidationCodeByEmail(namespace.ClientKey.Name, c); err != nil {
		return err
	}
	return nil
}

// =====================
// EMPLOYEE AUTH
// =====================

// LoginEmployeeByPassword logs an employee in
//
//	@Summary		Login employee
//	@Description	Log in an employee using password
//	@Tags			Employee/Auth
//	@Security		ApiKeyAuth
//	@Param			X-Company-ID	header	string	true	"X-Company-ID"
//	@Accept			json
//	@Produce		json
//	@Param			employee	body	DTO.LoginEmployee	true	"Employee credentials"
//	@Success		200			"Token returned in X-Auth-Token header"
//	@Failure		400			{object}	DTO.ErrorResponse
//	@Failure		401			{object}	nil
//	@Router			/auth/employee/login [post]
func LoginEmployeeByPassword(c *fiber.Ctx) error {
	token, err := LoginByPassword(namespace.EmployeeKey.Name, &model.Employee{}, c)
	if err != nil {
		return err
	}
	c.Response().Header.Set(namespace.HeadersKey.Auth, token)
	return nil
}

// LoginEmployeeByEmailCode logs in an employee using email and validation code
//
//	@Summary		Login employee by email code
//	@Description	Login employee using email and validation code
//	@Tags			Employee/Auth
//	@Param			X-Company-ID	header	string	true	"X-Company-ID"
//	@Accept			json
//	@Produce		json
//	@Param			body	body	DTO.LoginByEmailCode	true	"Login credentials"
//	@Success		200		"Token returned in X-Auth-Token header"
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/auth/employee/login-with-code [post]
func LoginEmployeeByEmailCode(c *fiber.Ctx) error {
	token, err := LoginByEmailCode(namespace.EmployeeKey.Name, &model.Employee{}, c)
	if err != nil {
		return err
	}
	c.Response().Header.Set(namespace.HeadersKey.Auth, token)
	return nil
}

// SendEmployeeLoginValidationCodeByEmail sends a login validation code to an employee's email
//
//	@Summary		Send employee login validation code by email
//	@Description	Sends a 6-digit login validation code to the employee's email
//	@Tags			Employee/Auth
//	@Accept			json
//	@Produce		json
//	@Param			X-Company-ID	header	string	true	"X-Company-ID"
//	@Param			email			path	string	true	"Employee Email"
//	@Query			language		query	string	false	"Language code (default: en)"
//	@Success		200
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/auth/employee/send-login-code/email/{email} [post]
func SendEmployeeLoginValidationCodeByEmail(c *fiber.Ctx) error {
	if err := SendLoginValidationCodeByEmail(c, &model.Employee{}); err != nil {
		return err
	}
	return nil
}

// =====================
// ADMIN AUTH
// =====================

// AdminLoginByPassword handles admin authentication
//
//	@Summary		Admin login
//	@Description	Authenticate admin user and return JWT token in X-Auth-Token header
//	@Tags			Admin/Auth
//	@Accept			json
//	@Produce		json
//	@Param			login	body	DTO.AdminLoginRequest	true	"Admin login credentials"
//	@Success		200		"Token returned in X-Auth-Token header"
//	@Failure		401		{object}	DTO.ErrorResponse
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/auth/admin/login [post]
func AdminLoginByPassword(c *fiber.Ctx) error {
	token, err := LoginByPassword(namespace.AdminKey.Name, &model.Admin{}, c)
	if err != nil {
		return err
	}
	c.Response().Header.Set(namespace.HeadersKey.Auth, token)
	return nil
}

// =====================
// TOKEN VALIDATION
// =====================

// ValidateToken validates a JWT token and returns the user claims
// This endpoint is used by the business service to validate tokens
//
//	@Summary		Validate token
//	@Description	Validate a JWT token and return user claims
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			X-Auth-Token	header	string	true	"JWT token to validate"
//	@Success		200				{object}	DTO.Claims
//	@Failure		401				{object}	DTO.ErrorResponse
//	@Router			/auth/validate [post]
func ValidateToken(c *fiber.Ctx) error {
	// Parse the token from header
	claims, err := handler.JWT(c).WhoAreYou()
	if err != nil {
		return lib.Error.Auth.InvalidToken.WithError(err)
	}
	if claims == nil {
		return lib.Error.Auth.InvalidToken.WithError(fmt.Errorf("no token provided"))
	}

	// Return the claims as JSON
	return lib.ResponseFactory(c).Send(200, claims)
}

// ValidateAdminToken validates an admin JWT token and returns the admin claims
// This endpoint is used by the business service to validate admin tokens
//
//	@Summary		Validate admin token
//	@Description	Validate an admin JWT token and return admin claims
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			X-Auth-Token	header	string	true	"JWT admin token to validate"
//	@Success		200				{object}	DTO.AdminClaims
//	@Failure		401				{object}	DTO.ErrorResponse
//	@Router			/auth/validate-admin [post]
func ValidateAdminToken(c *fiber.Ctx) error {
	// Parse the admin token from header
	claims, err := handler.JWT(c).WhoAreYouAdmin()
	if err != nil {
		return lib.Error.Auth.InvalidToken.WithError(err)
	}
	if claims == nil {
		return lib.Error.Auth.InvalidToken.WithError(fmt.Errorf("no admin token provided"))
	}

	// Return the claims as JSON
	return lib.ResponseFactory(c).Send(200, claims)
}

// =====================
// SHARED LOGIN HELPERS
// =====================

// LoginByPassword is a shared helper for password-based login using unified User model
func LoginByPassword(user_type string, c *fiber.Ctx) (string, error) {
	// Parse request body
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return "", lib.Error.General.BadRequest.WithError(err)
	}

	// Get database session
	tx, err := lib.Session(c)
	if err != nil {
		return "", err
	}

	// Find user by email and type
	var user model.User
	if err := tx.Where("email = ? AND type = ?", body.Email, user_type).First(&user).Error; err != nil {
		if err.Error() == "record not found" {
			return "", lib.Error.Client.NotFound
		}
		return "", lib.Error.General.InternalError.WithError(err)
	}

	// Validate user status and password
	if !user.Verified {
		return "", lib.Error.Client.NotVerified
	}
	if !handler.ComparePassword(user.Password, body.Password) {
		return "", lib.Error.Auth.InvalidLogin
	}

	// Create JWT claims from user data
	userBytes, err := json.Marshal(&user)
	if err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}

	var claims DTO.Claims
	if err := json.Unmarshal(userBytes, &claims); err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}
	claims.Type = user_type

	// Generate JWT token
	token, err := handler.JWT(c).Encode(&claims)
	if err != nil {
		return "", err
	}

	return token, nil
}

// LoginByEmailCode is a shared helper for email code-based login using unified User model
func LoginByEmailCode(user_type string, c *fiber.Ctx) (string, error) {
	// Parse request body
	var body DTO.LoginByEmailCode
	if err := c.BodyParser(&body); err != nil {
		return "", lib.Error.General.BadRequest.WithError(err)
	}

	// Get database session
	tx, err := lib.Session(c)
	if err != nil {
		return "", err
	}

	// Find user by email and type
	var user model.User
	if err := tx.Where("email = ? AND type = ?", body.Email, user_type).First(&user).Error; err != nil {
		if err.Error() == "record not found" {
			return "", lib.Error.Client.NotFound
		}
		return "", lib.Error.General.InternalError.WithError(err)
	}

	// Check if code exists
	if user.Meta.Login.ValidationCode == nil {
		return "", lib.Error.Auth.InvalidLogin.WithError(fmt.Errorf("no validation code found"))
	}

	// Check if expiry exists
	if user.Meta.Login.ValidationExpiry == nil {
		return "", lib.Error.Auth.InvalidLogin.WithError(fmt.Errorf("validation code has expired"))
	}

	// Check if code has expired
	if time.Now().After(*user.Meta.Login.ValidationExpiry) {
		return "", lib.Error.Auth.InvalidLogin.WithError(fmt.Errorf("validation code has expired"))
	}

	// Verify the code
	if *user.Meta.Login.ValidationCode != body.Code {
		return "", lib.Error.Auth.InvalidLogin.WithError(fmt.Errorf("invalid validation code"))
	}

	// Clear the validation code after successful use
	if err := tx.Model(&user).Updates(map[string]interface{}{
		"meta": map[string]interface{}{
			"login": map[string]interface{}{
				"validation_code":   nil,
				"validation_expiry": nil,
			},
		},
		"verified": true,
	}).Error; err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}

	// Generate JWT token
	userBytes, err := json.Marshal(&user)
	if err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}
	
	var claims DTO.Claims
	if err := json.Unmarshal(userBytes, &claims); err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}
	claims.Type = user_type

	token, err := handler.JWT(c).Encode(&claims)
	if err != nil {
		return "", err
	}
	
	return token, nil
}

// GenerateLoginValidationCode generates and stores a validation code, returning it for external use
func GenerateLoginValidationCode(c *fiber.Ctx, user_email string, user_type string) (string, error) {
	// Get database session
	tx, err := lib.Session(c)
	if err != nil {
		return "", err
	}

	// URL decode the email
	user_email, err = lib.PrepareEmail(user_email)
	if err != nil {
		return "", lib.Error.General.BadRequest.WithError(err)
	}

	// Find user by email
	if err := tx.Model(modelInstance).Where("email = ?", user_email).First(modelInstance).Error; err != nil {
		if err.Error() == "record not found" {
			return "", lib.Error.General.RecordNotFound
		}
		return "", lib.Error.General.InternalError.WithError(err)
	}

	// Use reflection to access Meta field
	val := reflect.ValueOf(modelInstance)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	metaField := val.FieldByName("Meta")
	if !metaField.IsValid() {
		return "", lib.Error.General.InternalError.WithError(fmt.Errorf("model does not have a Meta field"))
	}

	// Generate 6-digit validation code
	code := lib.GenerateRandomInt(6)
	codeString := fmt.Sprintf("%06d", code)

	// Set expiration to 15 minutes
	expiryTime := time.Now().Add(15 * time.Minute)
	now := time.Now()

	// Update the database with new validation code
	updates := map[string]interface{}{
		"meta": map[string]interface{}{
			"login": map[string]interface{}{
				"validation_code":         codeString,
				"validation_expiry":       expiryTime,
				"validation_requested_at": now,
			},
		},
	}

	if err := tx.Model(modelInstance).Updates(updates).Error; err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}

	return codeString, nil
}

// SendLoginValidationCodeByEmail generates a validation code and returns it in the response
// The caller (business service) can then send the code via email
func SendLoginValidationCodeByEmail(c *fiber.Ctx, model any) error {
	user_email := c.Params("email")
	if user_email == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("missing 'email' at params route"))
	}

	validationCode, err := GenerateLoginValidationCode(c, user_email, model)
	if err != nil {
		return err
	}

	// Return the validation code in the response for the business service to send via email
	return lib.ResponseFactory(c).Send(200, map[string]interface{}{
		"validation_code": validationCode,
		"email":           user_email,
		"message":         "Validation code generated. Use this code to login or have your service send it via email.",
	})
}
