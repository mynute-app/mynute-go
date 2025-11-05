package controller

import (
	"context"
	"mynute-go/services/email/config/dto"
	emailLib "mynute-go/services/email/api/lib"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	EmailProvider    emailLib.Sender
	TemplateRenderer *emailLib.TemplateRenderer
)

// SendEmail godoc
//
//	@Summary		Send a single email
//	@Description	Send a plain text or HTML email to a single recipient
//	@Tags			Email
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.SendEmailRequest		true	"Email request"
//	@Success		200		{object}	dto.EmailSuccessResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/api/v1/emails/send [post]
func SendEmail(c *fiber.Ctx) error {
	var req dto.SendEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
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

	if err := EmailProvider.Send(ctx, emailData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to send email",
			Details: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.EmailSuccessResponse{
		Success: true,
		Message: "Email sent successfully",
		To:      req.To,
	})
}

// SendTemplateEmail godoc
//
//	@Summary		Send an email using a template
//	@Description	Send an email using a predefined template with dynamic data
//	@Tags			Email
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.SendTemplateEmailRequest	true	"Template email request"
//	@Success		200		{object}	dto.EmailSuccessResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/api/v1/emails/send-template [post]
func SendTemplateEmail(c *fiber.Ctx) error {
	var req dto.SendTemplateEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
	}

	// Render the email template
	rendered, err := TemplateRenderer.RenderEmail(req.TemplateName, req.Language, req.Data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to render email template",
			Details: err.Error(),
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

	if err := EmailProvider.Send(ctx, emailData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to send email",
			Details: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.EmailSuccessResponse{
		Success: true,
		Message: "Template email sent successfully",
		To:      req.To,
	})
}

// SendBulkEmail godoc
//
//	@Summary		Send bulk emails
//	@Description	Send the same email to multiple recipients
//	@Tags			Email
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.SendBulkEmailRequest	true	"Bulk email request"
//	@Success		200		{object}	dto.BulkEmailResponse
//	@Failure		206		{object}	dto.BulkEmailResponse		"Partial success"
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/api/v1/emails/send-bulk [post]
func SendBulkEmail(c *fiber.Ctx) error {
	var req dto.SendBulkEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
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

		if err := EmailProvider.Send(ctx, emailData); err != nil {
			failedRecipients = append(failedRecipients, recipient)
		} else {
			successCount++
		}
	}

	response := dto.BulkEmailResponse{
		Success: successCount > 0,
		Total:   len(req.Recipients),
		Sent:    successCount,
		Failed:  len(failedRecipients),
	}

	if len(failedRecipients) > 0 {
		response.FailedRecipients = failedRecipients
	}

	statusCode := fiber.StatusOK
	if successCount == 0 {
		statusCode = fiber.StatusInternalServerError
		return c.Status(statusCode).JSON(dto.ErrorResponse{
			Error:   "Failed to send any emails",
			Details: "All recipients failed",
		})
	} else if len(failedRecipients) > 0 {
		statusCode = fiber.StatusPartialContent
	}

	return c.Status(statusCode).JSON(response)
}

// HealthCheck godoc
//
//	@Summary		Health check
//	@Description	Check if the email service is running
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	dto.HealthResponse
//	@Router			/health [get]
func HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(dto.HealthResponse{
		Status:  "healthy",
		Service: "email",
		Version: "1.0",
	})
}

// SendTemplateMerge godoc
//
//	@Summary		Send email by merging template HTML with translations
//	@Description	Receives template HTML, translations map, and custom data, merges them, and sends the email
//	@Tags			Email
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.SendTemplateMergeRequest	true	"Template merge request"
//	@Success		200		{object}	dto.EmailSuccessResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/api/v1/emails/send-template-merge [post]
func SendTemplateMerge(c *fiber.Ctx) error {
	var req dto.SendTemplateMergeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
	}

	// Get subject from translations
	subject, ok := req.Translations["subject"].(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Missing subject in translations",
			Details: "translations must contain a 'subject' key",
		})
	}

	// Merge translations with custom data (custom data takes precedence)
	mergedData := make(map[string]interface{})
	for k, v := range req.Translations {
		mergedData[k] = v
	}
	for k, v := range req.Data {
		mergedData[k] = v
	}

	// Render the template with merged data
	rendered, err := TemplateRenderer.RenderFromString(req.TemplateHTML, mergedData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to render email template",
			Details: err.Error(),
		})
	}

	// Create email data
	emailData := emailLib.EmailData{
		To:      req.To,
		Subject: subject,
		Html:    rendered,
		Cc:      req.CC,
		Bcc:     req.BCC,
	}

	// Send email
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := EmailProvider.Send(ctx, emailData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Failed to send email",
			Details: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.EmailSuccessResponse{
		Success: true,
		Message: "Template merge email sent successfully",
		To:      req.To[0], // Return first recipient
	})
}
