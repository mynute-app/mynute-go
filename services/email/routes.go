package email

import (
	"fmt"

	"mynute-go/services/email/api/controller"
	emailLib "mynute-go/services/email/api/lib"

	"github.com/gofiber/fiber/v2"
)

// initEmailServices initializes the email provider and template renderer
func initEmailServices() error {
	var err error

	// Initialize email provider
	controller.EmailProvider, err = emailLib.NewProvider(nil) // Will use APP_ENV to determine provider
	if err != nil {
		return fmt.Errorf("failed to initialize email provider: %w", err)
	}

	// Initialize template renderer (for send-template-merge endpoint)
	controller.TemplateRenderer = emailLib.NewTemplateRenderer("", "")

	return nil
}

// setupRoutes configures all email service routes
func setupRoutes(app *fiber.App) {
	// Health check endpoint (no /api prefix)
	app.Get("/health", controller.HealthCheck)

	// API routes
	api := app.Group("/api/v1")

	// Email endpoints
	emails := api.Group("/emails")
	emails.Post("/send", controller.SendEmail)
	emails.Post("/send-template-merge", controller.SendTemplateMerge)
	emails.Post("/send-bulk", controller.SendBulkEmail)
}
