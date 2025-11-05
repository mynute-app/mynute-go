package auth_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

// ==================== MODELS ====================

// EndPoint represents an endpoint from the auth service
type EndPoint struct {
	ID               uuid.UUID  `json:"id"`
	ControllerName   string     `json:"controller_name"`
	Description      string     `json:"description"`
	Method           string     `json:"method"`
	Path             string     `json:"path"`
	DenyUnauthorized bool       `json:"deny_unauthorized"`
	NeedsCompanyId   bool       `json:"needs_company_id"`
	ResourceID       *uuid.UUID `json:"resource_id,omitempty"`
}

// Policy represents an access control policy
type Policy struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Effect      string          `json:"effect"` // "Allow" or "Deny"
	EndPointID  uuid.UUID       `json:"end_point_id"`
	Conditions  json.RawMessage `json:"conditions,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// User represents a user from the auth service (Client, Employee, or Admin)
type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TokenValidationResponse represents the response from token validation
type TokenValidationResponse struct {
	Valid     bool      `json:"valid"`
	UserID    uuid.UUID `json:"user_id,omitempty"`
	Email     string    `json:"email,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

// AccessCheckRequest represents a request to check access
type AccessCheckRequest struct {
	UserID    uuid.UUID `json:"user_id"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	CompanyID *uuid.UUID `json:"company_id,omitempty"`
}

// AccessCheckResponse represents the response from access check
type AccessCheckResponse struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason,omitempty"`
}

// CreateEndpointRequest represents a request to create an endpoint
type CreateEndpointRequest struct {
	ControllerName   string     `json:"controller_name"`
	Description      string     `json:"description"`
	Method           string     `json:"method"`
	Path             string     `json:"path"`
	DenyUnauthorized bool       `json:"deny_unauthorized"`
	NeedsCompanyId   bool       `json:"needs_company_id"`
	ResourceID       *uuid.UUID `json:"resource_id,omitempty"`
}

// UpdateEndpointRequest represents a request to update an endpoint
type UpdateEndpointRequest struct {
	ControllerName   *string    `json:"controller_name,omitempty"`
	Description      *string    `json:"description,omitempty"`
	Method           *string    `json:"method,omitempty"`
	Path             *string    `json:"path,omitempty"`
	DenyUnauthorized *bool      `json:"deny_unauthorized,omitempty"`
	NeedsCompanyId   *bool      `json:"needs_company_id,omitempty"`
	ResourceID       *uuid.UUID `json:"resource_id,omitempty"`
}

// CreatePolicyRequest represents a request to create a policy
type CreatePolicyRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Effect      string          `json:"effect"`
	EndPointID  uuid.UUID       `json:"end_point_id"`
	Conditions  json.RawMessage `json:"conditions,omitempty"`
}

// UpdatePolicyRequest represents a request to update a policy
type UpdatePolicyRequest struct {
	Name        *string         `json:"name,omitempty"`
	Description *string         `json:"description,omitempty"`
	Effect      *string         `json:"effect,omitempty"`
	EndPointID  *uuid.UUID      `json:"end_point_id,omitempty"`
	Conditions  json.RawMessage `json:"conditions,omitempty"`
}

// ==================== CLIENT ====================

// AuthClient handles communication with the auth service
type AuthClient struct {
	BaseURL string
	Client  *http.Client
}

// NewAuthClient creates a new auth service client
func NewAuthClient() *AuthClient {
	// Get auth service URL from environment or use default
	authURL := os.Getenv("AUTH_SERVICE_URL")
	if authURL == "" {
		authURL = "http://localhost:4001"
	}

	return &AuthClient{
		BaseURL: authURL,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FetchEndpoints retrieves all endpoints from the auth service
func (ac *AuthClient) FetchEndpoints() ([]*EndPoint, error) {
	url := fmt.Sprintf("%s/endpoints", ac.BaseURL)

	resp, err := ac.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch endpoints from auth service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("auth service returned status %d: %s", resp.StatusCode, string(body))
	}

	var endpoints []*EndPoint
	if err := json.NewDecoder(resp.Body).Decode(&endpoints); err != nil {
		return nil, fmt.Errorf("failed to decode endpoints response: %w", err)
	}

	return endpoints, nil
}

// IsAvailable checks if the auth service is reachable
func (ac *AuthClient) IsAvailable() bool {
	url := fmt.Sprintf("%s/health", ac.BaseURL)

	resp, err := ac.Client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// ==================== ENDPOINT OPERATIONS ====================

// CreateEndpoint creates a new endpoint in the auth service
func (ac *AuthClient) CreateEndpoint(req CreateEndpointRequest) (*EndPoint, error) {
	return doRequest[EndPoint](ac, "POST", "/endpoints", req)
}

// GetEndpoint retrieves a specific endpoint by ID
func (ac *AuthClient) GetEndpoint(id uuid.UUID) (*EndPoint, error) {
	return doRequest[EndPoint](ac, "GET", fmt.Sprintf("/endpoints/%s", id), nil)
}

// UpdateEndpoint updates an existing endpoint
func (ac *AuthClient) UpdateEndpoint(id uuid.UUID, req UpdateEndpointRequest) (*EndPoint, error) {
	return doRequest[EndPoint](ac, "PATCH", fmt.Sprintf("/endpoints/%s", id), req)
}

// DeleteEndpoint deletes an endpoint by ID
func (ac *AuthClient) DeleteEndpoint(id uuid.UUID) error {
	_, err := doRequest[any](ac, "DELETE", fmt.Sprintf("/endpoints/%s", id), nil)
	return err
}

// ==================== POLICY OPERATIONS ====================

// FetchPolicies retrieves all policies from the auth service
func (ac *AuthClient) FetchPolicies() ([]*Policy, error) {
	return doRequestList[Policy](ac, "GET", "/policies", nil)
}

// CreatePolicy creates a new policy
func (ac *AuthClient) CreatePolicy(req CreatePolicyRequest) (*Policy, error) {
	return doRequest[Policy](ac, "POST", "/policies", req)
}

// GetPolicy retrieves a specific policy by ID
func (ac *AuthClient) GetPolicy(id uuid.UUID) (*Policy, error) {
	return doRequest[Policy](ac, "GET", fmt.Sprintf("/policies/%s", id), nil)
}

// UpdatePolicy updates an existing policy
func (ac *AuthClient) UpdatePolicy(id uuid.UUID, req UpdatePolicyRequest) (*Policy, error) {
	return doRequest[Policy](ac, "PATCH", fmt.Sprintf("/policies/%s", id), req)
}

// DeletePolicy deletes a policy by ID
func (ac *AuthClient) DeletePolicy(id uuid.UUID) error {
	_, err := doRequest[any](ac, "DELETE", fmt.Sprintf("/policies/%s", id), nil)
	return err
}

// ==================== USER MANAGEMENT ====================

// GetClientByEmail retrieves a client by email
func (ac *AuthClient) GetClientByEmail(email string) (*User, error) {
	return doRequest[User](ac, "GET", fmt.Sprintf("/users/client/email/%s", email), nil)
}

// GetClientByID retrieves a client by ID
func (ac *AuthClient) GetClientByID(id uuid.UUID) (*User, error) {
	return doRequest[User](ac, "GET", fmt.Sprintf("/users/client/%s", id), nil)
}

// GetEmployeeByEmail retrieves an employee by email
func (ac *AuthClient) GetEmployeeByEmail(email string) (*User, error) {
	return doRequest[User](ac, "GET", fmt.Sprintf("/users/employee/email/%s", email), nil)
}

// GetEmployeeByID retrieves an employee by ID
func (ac *AuthClient) GetEmployeeByID(id uuid.UUID) (*User, error) {
	return doRequest[User](ac, "GET", fmt.Sprintf("/users/employee/%s", id), nil)
}

// GetAdminByID retrieves an admin by ID
func (ac *AuthClient) GetAdminByID(id uuid.UUID) (*User, error) {
	return doRequest[User](ac, "GET", fmt.Sprintf("/users/admin/%s", id), nil)
}

// ListAdmins retrieves all admins
func (ac *AuthClient) ListAdmins() ([]*User, error) {
	return doRequestList[User](ac, "GET", "/users/admin", nil)
}

// ==================== AUTHENTICATION ====================

// ValidateToken validates a user token
func (ac *AuthClient) ValidateToken(token string) (*TokenValidationResponse, error) {
	req := map[string]string{"token": token}
	return doRequest[TokenValidationResponse](ac, "POST", "/auth/validate", req)
}

// ValidateAdminToken validates an admin token
func (ac *AuthClient) ValidateAdminToken(token string) (*TokenValidationResponse, error) {
	req := map[string]string{"token": token}
	return doRequest[TokenValidationResponse](ac, "POST", "/auth/validate-admin", req)
}

// ==================== AUTHORIZATION ====================

// CheckAccess checks if a user has access to a specific endpoint
func (ac *AuthClient) CheckAccess(req AccessCheckRequest) (*AccessCheckResponse, error) {
	return doRequest[AccessCheckResponse](ac, "POST", "/authorize/by-method-and-path", req)
}

// ==================== HELPER FUNCTIONS ====================

// doRequest is a generic helper for making HTTP requests
func doRequest[T any](ac *AuthClient, method, path string, body interface{}) (*T, error) {
	url := fmt.Sprintf("%s%s", ac.BaseURL, path)

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := ac.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("auth service returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// doRequestList is a helper for making requests that return a list
func doRequestList[T any](ac *AuthClient, method, path string, body interface{}) ([]*T, error) {
	url := fmt.Sprintf("%s%s", ac.BaseURL, path)

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := ac.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("auth service returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result []*T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}
