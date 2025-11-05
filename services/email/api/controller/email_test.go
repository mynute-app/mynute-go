package controller

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"mynute-go/services/email/api/lib"
	"mynute-go/services/email/config/dto"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSender is a mock implementation of the Sender interface
type MockSender struct {
	mock.Mock
}

func (m *MockSender) Send(ctx context.Context, emailData lib.EmailData) error {
	args := m.Called(ctx, emailData)
	return args.Error(0)
}

// MockTemplateRenderer is a mock implementation for template rendering
type MockTemplateRendererStruct struct {
	mock.Mock
}

func (m *MockTemplateRendererStruct) RenderFromString(templateHTML string, data lib.TemplateData) (string, error) {
	args := m.Called(templateHTML, data)
	return args.String(0), args.Error(1)
}

func setupTestApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})
	return app
}

func TestSendEmail(t *testing.T) {
	t.Run("should send email to single recipient successfully", func(t *testing.T) {
		mockSender := new(MockSender)
		EmailProvider = mockSender

		expectedEmailData := lib.EmailData{
			To:      []string{"user@example.com"},
			Subject: "Test Subject",
			Html:    "<h1>Test</h1>",
			Cc:      []string{},
			Bcc:     []string{},
		}

		mockSender.On("Send", mock.Anything, expectedEmailData).Return(nil)

		app := setupTestApp()
		app.Post("/send", SendEmail)

		reqBody := `{
			"subject": "Test Subject",
			"body": "<h1>Test</h1>",
			"is_html": true,
			"recipients": [
				{
					"to": ["user@example.com"],
					"cc": [],
					"bcc": []
				}
			]
		}`

		req := httptest.NewRequest("POST", "/send", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var response dto.EmailSuccessResponse
		json.Unmarshal(body, &response)

		assert.True(t, response.Success)
		assert.Equal(t, "All emails sent successfully", response.Message)
		assert.Equal(t, 1, response.Total)
		assert.Equal(t, 1, response.Sent)
		assert.Equal(t, 0, response.Failed)
		assert.Len(t, response.Results, 1)
		assert.True(t, response.Results[0].Success)
		assert.Equal(t, []string{"user@example.com"}, response.Results[0].To)

		mockSender.AssertExpectations(t)
	})

	t.Run("should send email to multiple recipients successfully", func(t *testing.T) {
		mockSender := new(MockSender)
		EmailProvider = mockSender

		expectedEmailData1 := lib.EmailData{
			To:      []string{"user1@example.com"},
			Subject: "Newsletter",
			Html:    "<h1>Newsletter</h1>",
			Cc:      []string{"manager1@example.com"},
			Bcc:     []string{},
		}

		expectedEmailData2 := lib.EmailData{
			To:      []string{"user2@example.com"},
			Subject: "Newsletter",
			Html:    "<h1>Newsletter</h1>",
			Cc:      []string{},
			Bcc:     []string{"archive@example.com"},
		}

		mockSender.On("Send", mock.Anything, expectedEmailData1).Return(nil)
		mockSender.On("Send", mock.Anything, expectedEmailData2).Return(nil)

		app := setupTestApp()
		app.Post("/send", SendEmail)

		reqBody := `{
			"subject": "Newsletter",
			"body": "<h1>Newsletter</h1>",
			"is_html": true,
			"recipients": [
				{
					"to": ["user1@example.com"],
					"cc": ["manager1@example.com"],
					"bcc": []
				},
				{
					"to": ["user2@example.com"],
					"cc": [],
					"bcc": ["archive@example.com"]
				}
			]
		}`

		req := httptest.NewRequest("POST", "/send", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var response dto.EmailSuccessResponse
		json.Unmarshal(body, &response)

		assert.True(t, response.Success)
		assert.Equal(t, "All emails sent successfully", response.Message)
		assert.Equal(t, 2, response.Total)
		assert.Equal(t, 2, response.Sent)
		assert.Equal(t, 0, response.Failed)
		assert.Len(t, response.Results, 2)
		assert.True(t, response.Results[0].Success)
		assert.True(t, response.Results[1].Success)

		mockSender.AssertExpectations(t)
	})

	t.Run("should send plain text email successfully", func(t *testing.T) {
		mockSender := new(MockSender)
		EmailProvider = mockSender

		expectedEmailData := lib.EmailData{
			To:      []string{"user@example.com"},
			Subject: "Plain Text",
			Text:    "This is plain text",
			Cc:      []string{},
			Bcc:     []string{},
		}

		mockSender.On("Send", mock.Anything, expectedEmailData).Return(nil)

		app := setupTestApp()
		app.Post("/send", SendEmail)

		reqBody := `{
			"subject": "Plain Text",
			"body": "This is plain text",
			"is_html": false,
			"recipients": [
				{
					"to": ["user@example.com"],
					"cc": [],
					"bcc": []
				}
			]
		}`

		req := httptest.NewRequest("POST", "/send", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		mockSender.AssertExpectations(t)
	})

	t.Run("should handle partial failure when some recipients fail", func(t *testing.T) {
		mockSender := new(MockSender)
		EmailProvider = mockSender

		expectedEmailData1 := lib.EmailData{
			To:      []string{"user1@example.com"},
			Subject: "Test",
			Html:    "<h1>Test</h1>",
			Cc:      []string{},
			Bcc:     []string{},
		}

		expectedEmailData2 := lib.EmailData{
			To:      []string{"user2@example.com"},
			Subject: "Test",
			Html:    "<h1>Test</h1>",
			Cc:      []string{},
			Bcc:     []string{},
		}

		expectedEmailData3 := lib.EmailData{
			To:      []string{"user3@example.com"},
			Subject: "Test",
			Html:    "<h1>Test</h1>",
			Cc:      []string{},
			Bcc:     []string{},
		}

		mockSender.On("Send", mock.Anything, expectedEmailData1).Return(nil)
		mockSender.On("Send", mock.Anything, expectedEmailData2).Return(errors.New("SMTP error: connection refused"))
		mockSender.On("Send", mock.Anything, expectedEmailData3).Return(nil)

		app := setupTestApp()
		app.Post("/send", SendEmail)

		reqBody := `{
			"subject": "Test",
			"body": "<h1>Test</h1>",
			"is_html": true,
			"recipients": [
				{"to": ["user1@example.com"], "cc": [], "bcc": []},
				{"to": ["user2@example.com"], "cc": [], "bcc": []},
				{"to": ["user3@example.com"], "cc": [], "bcc": []}
			]
		}`

		req := httptest.NewRequest("POST", "/send", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusPartialContent, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var response dto.EmailSuccessResponse
		json.Unmarshal(body, &response)

		assert.True(t, response.Success)
		assert.Equal(t, "Some emails failed to send", response.Message)
		assert.Equal(t, 3, response.Total)
		assert.Equal(t, 2, response.Sent)
		assert.Equal(t, 1, response.Failed)
		assert.Len(t, response.Results, 3)

		// Check individual results
		assert.True(t, response.Results[0].Success)
		assert.Empty(t, response.Results[0].Error)

		assert.False(t, response.Results[1].Success)
		assert.Equal(t, "SMTP error: connection refused", response.Results[1].Error)

		assert.True(t, response.Results[2].Success)
		assert.Empty(t, response.Results[2].Error)

		mockSender.AssertExpectations(t)
	})

	t.Run("should return 500 when all recipients fail", func(t *testing.T) {
		mockSender := new(MockSender)
		EmailProvider = mockSender

		expectedEmailData1 := lib.EmailData{
			To:      []string{"user1@example.com"},
			Subject: "Test",
			Html:    "<h1>Test</h1>",
			Cc:      []string{},
			Bcc:     []string{},
		}

		expectedEmailData2 := lib.EmailData{
			To:      []string{"user2@example.com"},
			Subject: "Test",
			Html:    "<h1>Test</h1>",
			Cc:      []string{},
			Bcc:     []string{},
		}

		mockSender.On("Send", mock.Anything, expectedEmailData1).Return(errors.New("invalid email address"))
		mockSender.On("Send", mock.Anything, expectedEmailData2).Return(errors.New("SMTP server unavailable"))

		app := setupTestApp()
		app.Post("/send", SendEmail)

		reqBody := `{
			"subject": "Test",
			"body": "<h1>Test</h1>",
			"is_html": true,
			"recipients": [
				{"to": ["user1@example.com"], "cc": [], "bcc": []},
				{"to": ["user2@example.com"], "cc": [], "bcc": []}
			]
		}`

		req := httptest.NewRequest("POST", "/send", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var response dto.EmailSuccessResponse
		json.Unmarshal(body, &response)

		assert.False(t, response.Success)
		assert.Equal(t, "Failed to send all emails", response.Message)
		assert.Equal(t, 2, response.Total)
		assert.Equal(t, 0, response.Sent)
		assert.Equal(t, 2, response.Failed)
		assert.Len(t, response.Results, 2)
		assert.False(t, response.Results[0].Success)
		assert.False(t, response.Results[1].Success)

		mockSender.AssertExpectations(t)
	})

	t.Run("should handle recipient with multiple To addresses", func(t *testing.T) {
		mockSender := new(MockSender)
		EmailProvider = mockSender

		expectedEmailData := lib.EmailData{
			To:      []string{"user1@example.com", "user2@example.com"},
			Subject: "Multi-recipient",
			Html:    "<h1>Test</h1>",
			Cc:      []string{},
			Bcc:     []string{},
		}

		mockSender.On("Send", mock.Anything, expectedEmailData).Return(nil)

		app := setupTestApp()
		app.Post("/send", SendEmail)

		reqBody := `{
			"subject": "Multi-recipient",
			"body": "<h1>Test</h1>",
			"is_html": true,
			"recipients": [
				{
					"to": ["user1@example.com", "user2@example.com"],
					"cc": [],
					"bcc": []
				}
			]
		}`

		req := httptest.NewRequest("POST", "/send", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		mockSender.AssertExpectations(t)
	})

	t.Run("should handle recipient with multiple CC and BCC", func(t *testing.T) {
		mockSender := new(MockSender)
		EmailProvider = mockSender

		expectedEmailData := lib.EmailData{
			To:      []string{"user@example.com"},
			Subject: "Multi CC/BCC",
			Html:    "<h1>Test</h1>",
			Cc:      []string{"cc1@example.com", "cc2@example.com"},
			Bcc:     []string{"bcc1@example.com", "bcc2@example.com"},
		}

		mockSender.On("Send", mock.Anything, expectedEmailData).Return(nil)

		app := setupTestApp()
		app.Post("/send", SendEmail)

		reqBody := `{
			"subject": "Multi CC/BCC",
			"body": "<h1>Test</h1>",
			"is_html": true,
			"recipients": [
				{
					"to": ["user@example.com"],
					"cc": ["cc1@example.com", "cc2@example.com"],
					"bcc": ["bcc1@example.com", "bcc2@example.com"]
				}
			]
		}`

		req := httptest.NewRequest("POST", "/send", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		mockSender.AssertExpectations(t)
	})

	t.Run("should return 400 for invalid JSON", func(t *testing.T) {
		app := setupTestApp()
		app.Post("/send", SendEmail)

		reqBody := `{invalid json`

		req := httptest.NewRequest("POST", "/send", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var response dto.ErrorResponse
		json.Unmarshal(body, &response)

		assert.Equal(t, "Invalid request body", response.Error)
	})

	t.Run("should handle empty recipients array", func(t *testing.T) {
		mockSender := new(MockSender)
		EmailProvider = mockSender

		app := setupTestApp()
		app.Post("/send", SendEmail)

		reqBody := `{
			"subject": "Test",
			"body": "<h1>Test</h1>",
			"is_html": true,
			"recipients": []
		}`

		req := httptest.NewRequest("POST", "/send", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var response dto.EmailSuccessResponse
		json.Unmarshal(body, &response)

		assert.False(t, response.Success)
		assert.Equal(t, "Failed to send all emails", response.Message)
		assert.Equal(t, 0, response.Total)
		assert.Equal(t, 0, response.Sent)
		assert.Equal(t, 0, response.Failed)
	})
}

func TestSendTemplateMerge(t *testing.T) {
	t.Run("should merge template and send successfully", func(t *testing.T) {
		mockSender := new(MockSender)
		EmailProvider = mockSender
		TemplateRenderer = &lib.TemplateRenderer{}

		mockSender.On("Send", mock.Anything, mock.MatchedBy(func(data lib.EmailData) bool {
			return data.To[0] == "user@example.com" && 
				   data.Subject == "Welcome" && 
				   len(data.Html) > 0
		})).Return(nil)

		app := setupTestApp()
		app.Post("/send-template-merge", SendTemplateMerge)

		reqBody := `{
			"template_html": "<h1>{{.greeting}}, {{.name}}!</h1><p>Your code is: {{.code}}</p>",
			"translations": {
				"subject": "Welcome",
				"greeting": "Hello"
			},
			"data": {
				"name": "John",
				"code": "123456"
			},
			"to": ["user@example.com"],
			"cc": [],
			"bcc": []
		}`

		req := httptest.NewRequest("POST", "/send-template-merge", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var response dto.EmailSuccessResponse
		json.Unmarshal(body, &response)

		assert.True(t, response.Success)
		assert.Equal(t, 1, response.Total)
		assert.Equal(t, 1, response.Sent)
		assert.Equal(t, 0, response.Failed)

		mockSender.AssertExpectations(t)
	})

	t.Run("should return 400 when subject missing in translations", func(t *testing.T) {
		app := setupTestApp()
		app.Post("/send-template-merge", SendTemplateMerge)

		reqBody := `{
			"template_html": "<h1>{{.greeting}}</h1>",
			"translations": {
				"greeting": "Hello"
			},
			"data": {},
			"to": ["user@example.com"],
			"cc": [],
			"bcc": []
		}`

		req := httptest.NewRequest("POST", "/send-template-merge", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var response dto.ErrorResponse
		json.Unmarshal(body, &response)

		assert.Equal(t, "Missing subject in translations", response.Error)
	})

	t.Run("should handle template rendering error", func(t *testing.T) {
		mockSender := new(MockSender)
		EmailProvider = mockSender
		TemplateRenderer = &lib.TemplateRenderer{}

		app := setupTestApp()
		app.Post("/send-template-merge", SendTemplateMerge)

		// Invalid template syntax
		reqBody := `{
			"template_html": "<h1>{{.greeting</h1>",
			"translations": {
				"subject": "Test",
				"greeting": "Hello"
			},
			"data": {},
			"to": ["user@example.com"],
			"cc": [],
			"bcc": []
		}`

		req := httptest.NewRequest("POST", "/send-template-merge", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var response dto.ErrorResponse
		json.Unmarshal(body, &response)

		assert.Equal(t, "Failed to render email template", response.Error)
	})

	t.Run("should handle email send failure in template merge", func(t *testing.T) {
		mockSender := new(MockSender)
		EmailProvider = mockSender
		TemplateRenderer = &lib.TemplateRenderer{}

		mockSender.On("Send", mock.Anything, mock.Anything).Return(errors.New("SMTP error"))

		app := setupTestApp()
		app.Post("/send-template-merge", SendTemplateMerge)

		reqBody := `{
			"template_html": "<h1>{{.greeting}}</h1>",
			"translations": {
				"subject": "Test",
				"greeting": "Hello"
			},
			"data": {},
			"to": ["user@example.com"],
			"cc": [],
			"bcc": []
		}`

		req := httptest.NewRequest("POST", "/send-template-merge", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var response dto.ErrorResponse
		json.Unmarshal(body, &response)

		assert.Equal(t, "Failed to send email", response.Error)

		mockSender.AssertExpectations(t)
	})

	t.Run("should handle custom data overriding translations", func(t *testing.T) {
		mockSender := new(MockSender)
		EmailProvider = mockSender
		TemplateRenderer = &lib.TemplateRenderer{}

		mockSender.On("Send", mock.Anything, mock.Anything).Return(nil)

		app := setupTestApp()
		app.Post("/send-template-merge", SendTemplateMerge)

		reqBody := `{
			"template_html": "<h1>{{.greeting}}</h1>",
			"translations": {
				"subject": "Test",
				"greeting": "Hello"
			},
			"data": {
				"greeting": "Hi there"
			},
			"to": ["user@example.com"],
			"cc": [],
			"bcc": []
		}`

		req := httptest.NewRequest("POST", "/send-template-merge", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		mockSender.AssertExpectations(t)
	})

	t.Run("should return 400 for invalid JSON in template merge", func(t *testing.T) {
		app := setupTestApp()
		app.Post("/send-template-merge", SendTemplateMerge)

		reqBody := `{invalid json`

		req := httptest.NewRequest("POST", "/send-template-merge", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var response dto.ErrorResponse
		json.Unmarshal(body, &response)

		assert.Equal(t, "Invalid request body", response.Error)
	})
}

func TestHealthCheck(t *testing.T) {
	t.Run("should return healthy status", func(t *testing.T) {
		app := setupTestApp()
		app.Get("/health", HealthCheck)

		req := httptest.NewRequest("GET", "/health", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var response dto.HealthResponse
		json.Unmarshal(body, &response)

		assert.Equal(t, "healthy", response.Status)
		assert.Equal(t, "email", response.Service)
		assert.Equal(t, "1.0", response.Version)
	})
}
