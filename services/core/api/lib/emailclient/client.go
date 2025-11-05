package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Client handles communication with the email service
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// New creates a new email service client
func New() *Client {
	baseURL := os.Getenv("EMAIL_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:4002"
	}

	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendEmailRequest represents a simple email request
type SendEmailRequest struct {
	To      string   `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	CC      []string `json:"cc,omitempty"`
	BCC     []string `json:"bcc,omitempty"`
	IsHTML  bool     `json:"is_html"`
}

// SendTemplateRequest represents a template email request
type SendTemplateRequest struct {
	To           string                 `json:"to"`
	TemplateName string                 `json:"template_name"`
	Language     string                 `json:"language"`
	Data         map[string]interface{} `json:"data"`
	CC           []string               `json:"cc,omitempty"`
	BCC          []string               `json:"bcc,omitempty"`
}

// SendBulkRequest represents a bulk email request
type SendBulkRequest struct {
	Recipients []string `json:"recipients"`
	Subject    string   `json:"subject"`
	Body       string   `json:"body"`
	IsHTML     bool     `json:"is_html"`
}

// EmailResponse represents the API response
type EmailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
	Details string `json:"details,omitempty"`
}

// Send sends a simple email
func (c *Client) Send(req *SendEmailRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.HTTPClient.Post(
		c.BaseURL+"/api/v1/emails/send",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to send request to email service: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var emailResp EmailResponse
	if err := json.Unmarshal(body, &emailResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if emailResp.Details != "" {
			return fmt.Errorf("email service error: %s - %s", emailResp.Error, emailResp.Details)
		}
		return fmt.Errorf("email service error (status %d): %s", resp.StatusCode, emailResp.Error)
	}

	return nil
}

// SendTemplate sends an email using a template
func (c *Client) SendTemplate(req *SendTemplateRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.HTTPClient.Post(
		c.BaseURL+"/api/v1/emails/send-template",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to send request to email service: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var emailResp EmailResponse
	if err := json.Unmarshal(body, &emailResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if emailResp.Details != "" {
			return fmt.Errorf("email service error: %s - %s", emailResp.Error, emailResp.Details)
		}
		return fmt.Errorf("email service error (status %d): %s", resp.StatusCode, emailResp.Error)
	}

	return nil
}

// SendBulk sends emails to multiple recipients
func (c *Client) SendBulk(req *SendBulkRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.HTTPClient.Post(
		c.BaseURL+"/api/v1/emails/send-bulk",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to send request to email service: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var emailResp EmailResponse
	if err := json.Unmarshal(body, &emailResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		if emailResp.Details != "" {
			return fmt.Errorf("email service error: %s - %s", emailResp.Error, emailResp.Details)
		}
		return fmt.Errorf("email service error (status %d): %s", resp.StatusCode, emailResp.Error)
	}

	return nil
}

// HealthCheck checks if the email service is available
func (c *Client) HealthCheck() error {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/health")
	if err != nil {
		return fmt.Errorf("email service is not reachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("email service is unhealthy (status %d)", resp.StatusCode)
	}

	return nil
}
