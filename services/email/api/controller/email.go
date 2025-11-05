package controller

import (
	"context"
	"mynute-go/services/email/api/lib"
	"mynute-go/services/email/config/dto"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	EmailProvider    lib.Sender
	TemplateRenderer *lib.TemplateRenderer
)

// SendEmail godoc
//
//	@Summary		Send emails to one or more recipients
//	@Description	Send plain text or HTML emails to single or multiple recipients with individual CC/BCC lists
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	results := []dto.Result{}
	successCount := 0
	failedCount := 0

	// Send email to each recipient
	for _, recipient := range req.Recipients {
		emailData := lib.EmailData{
			To:      recipient.To,
			Subject: req.Subject,
			Cc:      recipient.CC,
			Bcc:     recipient.BCC,
		}

		if req.IsHTML {
			emailData.Html = req.Body
		} else {
			emailData.Text = req.Body
		}

		// Send email
		err := EmailProvider.Send(ctx, emailData)
		
		result := dto.Result{
			To:      recipient.To,
			Success: err == nil,
		}
		
		if err != nil {
			result.Error = err.Error()
			failedCount++
		} else {
			successCount++
		}
		
		results = append(results, result)
	}

	// Determine response
	response := dto.EmailSuccessResponse{
		Success: successCount > 0,
		Total:   len(req.Recipients),
		Sent:    successCount,
		Failed:  failedCount,
		Results: results,
	}

	if successCount == 0 {
		response.Message = "Failed to send all emails"
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	} else if failedCount > 0 {
		response.Message = "Some emails failed to send"
		return c.Status(fiber.StatusPartialContent).JSON(response)
	}

	response.Message = "All emails sent successfully"
	return c.Status(fiber.StatusOK).JSON(response)
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
	emailData := lib.EmailData{
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
		Total:   1,
		Sent:    1,
		Failed:  0,
	})
}
