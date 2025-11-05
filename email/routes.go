package email

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	emailLib "mynute-go/email/lib"

	"github.com/gofiber/fiber/v2"
)

var (
	emailProvider    emailLib.Sender
	templateRenderer *emailLib.TemplateRenderer
)

// initEmailServices initializes the email provider and template renderer
func initEmailServices() error {
	var err error

	// Initialize email provider
	emailProvider, err = emailLib.NewProvider(nil) // Will use APP_ENV to determine provider
	if err != nil {
		return fmt.Errorf("failed to initialize email provider: %w", err)
	}

	// Initialize template renderer
	staticPath := filepath.Join(".", "static", "email")
	translationPath := filepath.Join(".", "translation", "email")
	templateRenderer = emailLib.NewTemplateRenderer(staticPath, translationPath)

	return nil
}

// setupRoutes configures all email service routes
func setupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// Email endpoints
	emails := api.Group("/emails")

	// Send email endpoint
	emails.Post("/send", handleSendEmail)

	// Send template email endpoint
	emails.Post("/send-template", handleSendTemplateEmail)

	// Send bulk emails endpoint
	emails.Post("/send-bulk", handleSendBulkEmail)
}

// handleSendEmail sends a single email
func handleSendEmail(c *fiber.Ctx) error {
	type EmailRequest struct {
		To      string   `json:"to" validate:"required,email"`
		Subject string   `json:"subject" validate:"required"`
		Body    string   `json:"body" validate:"required"`
		CC      []string `json:"cc,omitempty"`
		BCC     []string `json:"bcc,omitempty"`
		IsHTML  bool     `json:"is_html"`
	}

	var req EmailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Create email data
	emailData := emailLib.EmailData{
		To:      []string{req.To},
		Subject: req.Subject,
		Cc:      req.CC,
		Bcc:     req.BCC,
	}

	if req.IsHTML {
		emailData.Html = req.Body
	} else {
		emailData.Text = req.Body
	}

	// Send email
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := emailProvider.Send(ctx, emailData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to send email",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Email sent successfully",
		"to":      req.To,
	})
}

// handleSendTemplateEmail sends an email using a template
func handleSendTemplateEmail(c *fiber.Ctx) error {
	type TemplateEmailRequest struct {
		To           string                 `json:"to" validate:"required,email"`
		TemplateName string                 `json:"template_name" validate:"required"`
		Language     string                 `json:"language" validate:"required"`
		Data         map[string]interface{} `json:"data"`
		CC           []string               `json:"cc,omitempty"`
		BCC          []string               `json:"bcc,omitempty"`
	}

	var req TemplateEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Render the email template
	rendered, err := templateRenderer.RenderEmail(req.TemplateName, req.Language, req.Data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to render email template",
			"details": err.Error(),
		})
	}

	// Create email data
	emailData := emailLib.EmailData{
		To:      []string{req.To},
		Subject: rendered.Subject,
		Html:    rendered.HTMLBody,
		Cc:      req.CC,
		Bcc:     req.BCC,
	}

	// Send email
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := emailProvider.Send(ctx, emailData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to send email",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Template email sent successfully",
		"to":      req.To,
	})
}

// handleSendBulkEmail sends emails to multiple recipients
func handleSendBulkEmail(c *fiber.Ctx) error {
	type BulkEmailRequest struct {
		Recipients []string `json:"recipients" validate:"required,dive,email"`
		Subject    string   `json:"subject" validate:"required"`
		Body       string   `json:"body" validate:"required"`
		IsHTML     bool     `json:"is_html"`
	}

	var req BulkEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Send emails individually (could be optimized with goroutines)
	successCount := 0
	failedRecipients := []string{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	for _, recipient := range req.Recipients {
		emailData := emailLib.EmailData{
			To:      []string{recipient},
			Subject: req.Subject,
		}

		if req.IsHTML {
			emailData.Html = req.Body
		} else {
			emailData.Text = req.Body
		}

		if err := emailProvider.Send(ctx, emailData); err != nil {
			failedRecipients = append(failedRecipients, recipient)
		} else {
			successCount++
		}
	}

	response := fiber.Map{
		"success": successCount > 0,
		"total":   len(req.Recipients),
		"sent":    successCount,
		"failed":  len(failedRecipients),
	}

	if len(failedRecipients) > 0 {
		response["failed_recipients"] = failedRecipients
	}

	statusCode := fiber.StatusOK
	if successCount == 0 {
		statusCode = fiber.StatusInternalServerError
	} else if len(failedRecipients) > 0 {
		statusCode = fiber.StatusPartialContent
	}

	return c.Status(statusCode).JSON(response)
}
