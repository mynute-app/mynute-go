package auth_client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

// EndPoint represents an endpoint from the auth service
type EndPoint struct {
	ID               uuid.UUID `json:"id"`
	ControllerName   string    `json:"controller_name"`
	Description      string    `json:"description"`
	Method           string    `json:"method"`
	Path             string    `json:"path"`
	DenyUnauthorized bool      `json:"deny_unauthorized"`
	NeedsCompanyId   bool      `json:"needs_company_id"`
}

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
