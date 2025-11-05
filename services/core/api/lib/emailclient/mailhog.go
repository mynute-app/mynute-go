package email

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"
)

// MailHogMessage represents an email message from MailHog
type MailHogMessage struct {
	ID   string `json:"ID"`
	From struct {
		Relayed string `json:"Relayed"`
		Mailbox string `json:"Mailbox"`
		Domain  string `json:"Domain"`
		Params  string `json:"Params"`
	} `json:"From"`
	To []struct {
		Relayed string `json:"Relayed"`
		Mailbox string `json:"Mailbox"`
		Domain  string `json:"Domain"`
		Params  string `json:"Params"`
	} `json:"To"`
	Content struct {
		Headers map[string][]string `json:"Headers"`
		Body    string              `json:"Body"`
		Size    int                 `json:"Size"`
		MIME    interface{}         `json:"MIME"`
	} `json:"Content"`
	Created time.Time   `json:"Created"`
	MIME    interface{} `json:"MIME"`
	Raw     struct {
		From string   `json:"From"`
		To   []string `json:"To"`
		Data string   `json:"Data"`
		Helo string   `json:"Helo"`
	} `json:"Raw"`
}

// MailHogClient provides access to MailHog for testing
type MailHogClient struct {
	BaseURL string
	Client  *http.Client
}

// NewMailHogClient creates a new MailHog client for testing
func NewMailHogClient() *MailHogClient {
	host := os.Getenv("MAILHOG_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("MAILHOG_API_PORT")
	if port == "" {
		port = "8025"
	}

	return &MailHogClient{
		BaseURL: fmt.Sprintf("http://%s:%s", host, port),
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetMessages retrieves messages from MailHog
func (m *MailHogClient) GetMessages() ([]MailHogMessage, error) {
	resp, err := m.Client.Get(m.BaseURL + "/api/v2/messages")
	if err != nil {
		return nil, fmt.Errorf("failed to get messages from MailHog: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read MailHog response: %w", err)
	}

	var result struct {
		Total    int              `json:"total"`
		Count    int              `json:"count"`
		Start    int              `json:"start"`
		Messages []MailHogMessage `json:"items"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse MailHog response: %w", err)
	}

	return result.Messages, nil
}

// DeleteAll deletes all messages from MailHog
func (m *MailHogClient) DeleteAll() error {
	req, err := http.NewRequest("DELETE", m.BaseURL+"/api/v1/messages", nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	resp, err := m.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete messages from MailHog: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("MailHog delete failed with status: %d", resp.StatusCode)
	}

	return nil
}

// FindMessageByRecipient finds the first message sent to a specific email
func (m *MailHogClient) FindMessageByRecipient(email string) (*MailHogMessage, error) {
	messages, err := m.GetMessages()
	if err != nil {
		return nil, err
	}

	for _, msg := range messages {
		for _, to := range msg.To {
			if to.Mailbox+"@"+to.Domain == email {
				return &msg, nil
			}
		}
	}

	return nil, fmt.Errorf("no message found for recipient: %s", email)
}

// GetSubject returns the email subject
func (msg *MailHogMessage) GetSubject() string {
	if len(msg.Content.Headers["Subject"]) == 0 {
		return ""
	}
	return msg.Content.Headers["Subject"][0]
}

// GetBody returns the email body
func (msg *MailHogMessage) GetBody() string {
	return msg.Content.Body
}

// ExtractPassword extracts password from email body (simple implementation)
// This looks for common password patterns in the email
func (msg *MailHogMessage) ExtractPassword() (string, error) {
	body := msg.Content.Body

	// This is a simplified extraction - adjust based on your email template
	// Looking for patterns like "Password: XXXXX" or "NewPassword: XXXXX"
	patterns := []string{
		`(?i)password:\s*([^\s<]+)`,
		`(?i)newpassword:\s*([^\s<]+)`,
		`(?i)your password is:\s*([^\s<]+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(body)
		if len(matches) > 1 {
			return matches[1], nil
		}
	}

	return "", fmt.Errorf("could not extract password from email body")
}

// ExtractVerificationCode extracts verification code from email body
func (msg *MailHogMessage) ExtractVerificationCode() (string, error) {
	body := msg.Content.Body

	patterns := []string{
		`(?i)verification code:\s*([A-Z0-9]+)`,
		`(?i)code:\s*([A-Z0-9]{6,})`,
		`(?i)verificationcode:\s*([A-Z0-9]+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(body)
		if len(matches) > 1 {
			return matches[1], nil
		}
	}

	return "", fmt.Errorf("could not extract verification code from email body")
}

// ExtractValidationCode is an alias for ExtractVerificationCode
func (msg *MailHogMessage) ExtractValidationCode() (string, error) {
	return msg.ExtractVerificationCode()
}

// GetLatestMessageTo returns the most recent message sent to a specific email
func (m *MailHogClient) GetLatestMessageTo(email string) (*MailHogMessage, error) {
	return m.FindMessageByRecipient(email)
}
