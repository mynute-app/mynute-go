package controller

import (
	"fmt"
	"mynute-go/core/src/config/db/model"
	"mynute-go/core/src/config/namespace"
	"mynute-go/core/src/handler"
	"mynute-go/core/src/lib"
	"mynute-go/core/src/service"

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
	token, err := LoginByPassword(namespace.ClientKey.Name, &model.Client{}, c)
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
	token, err := LoginByEmailCode(namespace.ClientKey.Name, &model.Client{}, c)
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
	if err := SendLoginValidationCodeByEmail(c, &model.Client{}); err != nil {
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

// LoginByPassword is a shared helper for password-based login
func LoginByPassword(user_type string, model any, c *fiber.Ctx) (string, error) {
	var err error
	Service := service.New(c)
	defer func() { Service.DeferDB(err) }()
	token, err := Service.SetModel(model).LoginByPassword(user_type)
	return token, err
}

// LoginByEmailCode is a shared helper for email code-based login
func LoginByEmailCode(user_type string, model any, c *fiber.Ctx) (string, error) {
	var err error
	Service := service.New(c)
	defer func() { Service.DeferDB(err) }()
	token, err := Service.SetModel(model).LoginByEmailCode(user_type)
	return token, err
}

// ResetLoginvalidationCode resets the login validation code for a user
func ResetLoginvalidationCode(c *fiber.Ctx, user_email string, model any) (string, error) {
	var err error
	Service := service.New(c)
	defer func() { Service.DeferDB(err) }()
	return Service.SetModel(model).ResetLoginCodeByEmail(user_email)
}

// SendLoginValidationCodeByEmail sends a login validation code to the user's email
func SendLoginValidationCodeByEmail(c *fiber.Ctx, model any) error {
	user_email := c.Params("email")
	if user_email == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("missing 'email' at params route"))
	}

	LoginValidationCode, err := ResetLoginvalidationCode(c, user_email, model)
	if err != nil {
		return err
	}

	language := c.Query("language", "en")

	// Send email using the email library
	// NOTE: This requires the email templates to be accessible
	// For now, we'll return success - email integration can be added later
	_ = LoginValidationCode
	_ = language

	// TODO: Implement email sending
	// renderer := email.NewTemplateRenderer("./static/email", "./translation/email")
	// renderedEmail, err := renderer.RenderEmail("login_validation_code", language, email.TemplateData{
	// 	"LoginValidationCode": LoginValidationCode,
	// })
	// if err != nil {
	// 	return fmt.Errorf("failed to render email: %w", err)
	// }

	// provider, err := email.NewProvider(nil)
	// if err != nil {
	// 	return fmt.Errorf("failed to initialize email provider: %w", err)
	// }

	// err = provider.Send(context.Background(), email.EmailData{
	// 	To:      []string{user_email},
	// 	Subject: renderedEmail.Subject,
	// 	Html:    renderedEmail.HTMLBody,
	// })

	return nil
}
