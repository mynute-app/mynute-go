package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mynute-go/core/src/config/namespace"
	"mynute-go/core/src/lib"
	"net/http"
	"os"
	"time"

	DTO "mynute-go/core/src/config/api/dto"

	"github.com/gofiber/fiber/v2"
)

// AuthServiceClient handles communication with the auth service
type AuthServiceClient struct {
	BaseURL string
	Client  *http.Client
}

// NewAuthServiceClient creates a new auth service client
func NewAuthServiceClient() *AuthServiceClient {
	baseURL := os.Getenv("AUTH_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:4001" // Default for local development
	}

	return &AuthServiceClient{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// ValidateToken validates a JWT token by calling the auth service
func (a *AuthServiceClient) ValidateToken(token string) (*DTO.Claims, error) {
	req, err := http.NewRequest("POST", a.BaseURL+"/auth/validate", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(namespace.HeadersKey.Auth, token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call auth service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("auth service returned status %d: %s", resp.StatusCode, string(body))
	}

	var claims DTO.Claims
	if err := json.NewDecoder(resp.Body).Decode(&claims); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &claims, nil
}

// ValidateAdminToken validates an admin JWT token by calling the auth service
func (a *AuthServiceClient) ValidateAdminToken(token string) (*DTO.AdminClaims, error) {
	req, err := http.NewRequest("POST", a.BaseURL+"/auth/validate-admin", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(namespace.HeadersKey.Auth, token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call auth service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("auth service returned status %d: %s", resp.StatusCode, string(body))
	}

	var claims DTO.AdminClaims
	if err := json.NewDecoder(resp.Body).Decode(&claims); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &claims, nil
}

// DenyUnauthorizedViaAuthService is an alternative middleware that validates tokens
// by calling the auth service instead of doing it locally
func DenyUnauthorizedViaAuthService(c *fiber.Ctx) error {
	authClient := NewAuthServiceClient()

	// Get token from header
	token := c.Get(namespace.HeadersKey.Auth)
	if token == "" {
		return lib.Error.Auth.InvalidToken.WithError(fmt.Errorf("no token provided"))
	}

	// Validate token with auth service
	claims, err := authClient.ValidateToken(token)
	if err != nil {
		return lib.Error.Auth.InvalidToken.WithError(err)
	}

	// Store claims in context for downstream handlers
	c.Locals(namespace.RequestKey.Auth_Claims, claims)

	return c.Next()
}

// ProxyAuthServiceLogin proxies login requests to the auth service
// This can be used as a temporary bridge during migration
func ProxyAuthServiceLogin(userType string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authClient := NewAuthServiceClient()

		// Read request body
		bodyBytes, err := io.ReadAll(c.Request().BodyStream())
		if err != nil {
			return lib.Error.General.BadRequest.WithError(err)
		}

		// Determine the auth service endpoint based on user type
		var endpoint string
		switch userType {
		case namespace.ClientKey.Name:
			endpoint = "/auth/client/login"
		case namespace.EmployeeKey.Name:
			endpoint = "/auth/employee/login"
		case namespace.AdminKey.Name:
			endpoint = "/auth/admin/login"
		default:
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("unknown user type: %s", userType))
		}

		// Create request to auth service
		req, err := http.NewRequest("POST", authClient.BaseURL+endpoint, bytes.NewReader(bodyBytes))
		if err != nil {
			return lib.Error.General.InternalError.WithError(err)
		}

		// Copy headers
		req.Header.Set("Content-Type", "application/json")
		if companyID := c.Get("X-Company-ID"); companyID != "" {
			req.Header.Set("X-Company-ID", companyID)
		}

		// Make request to auth service
		resp, err := authClient.Client.Do(req)
		if err != nil {
			return lib.Error.General.InternalError.WithError(fmt.Errorf("failed to call auth service: %w", err))
		}
		defer resp.Body.Close()

		// Copy response headers (especially X-Auth-Token)
		for key, values := range resp.Header {
			for _, value := range values {
				c.Response().Header.Set(key, value)
			}
		}

		// Copy response status and body
		c.Status(resp.StatusCode)
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return lib.Error.General.InternalError.WithError(err)
		}

		return c.Send(respBody)
	}
}
