package email

// Example usage of TemplateRenderer
//
// This file demonstrates how to use the email template rendering system
// to send emails with multi-language support using either Resend or MailHog.

/*

// Initialize the template renderer
renderer := email.NewTemplateRenderer("./static", "./translation")

// Example 1: Render login validation email in English (default)
htmlBody, err := renderer.RenderEmail("login_validation", "", email.TemplateData{
	"ValidationCode": "123456",
})
if err != nil {
	// handle error
}

// Example 2: Render login validation email in Portuguese
htmlBody, err := renderer.RenderEmail("login_validation", "pt", email.TemplateData{
	"ValidationCode": "654321",
})
if err != nil {
	// handle error
}

// Example 3: Render login validation email in Spanish
htmlBody, err := renderer.RenderEmail("login_validation", "es", email.TemplateData{
	"ValidationCode": "ABCDEF",
})
if err != nil {
	// handle error
}

// Example 4: Send the rendered email using Resend (production)
provider, err := email.NewProvider("resend")
if err != nil {
	// handle error
}

err = provider.Send(context.Background(), email.EmailData{
	To:       []string{"user@example.com"},
	Subject:  "Login Validation Code",
	Html: htmlBody,
})
if err != nil {
	// handle error
}

// Example 5: Send the rendered email using MailHog (testing/development)
provider, err := email.NewProvider("mailhog")
if err != nil {
	// handle error
}

err = provider.Send(context.Background(), email.EmailData{
	To:       []string{"test@example.com"},
	Subject:  "Login Validation Code",
	Html: htmlBody,
})
if err != nil {
	// handle error
}

// Complete example: Render and send in one flow
func SendLoginValidationEmail(userEmail, code, language string) error {
	// Initialize renderer
	renderer := email.NewTemplateRenderer("./static", "./translation")

	// Render email template
	htmlBody, err := renderer.RenderEmail("login_validation", language, email.TemplateData{
		"ValidationCode": code,
	})
	if err != nil {
		return fmt.Errorf("failed to render email: %w", err)
	}

	// Initialize email provider (use "mailhog" for testing, "resend" for production)
	provider, err := email.NewProvider("resend")
	if err != nil {
		return fmt.Errorf("failed to initialize email provider: %w", err)
	}

	// Send email
	err = provider.Send(context.Background(), email.EmailData{
		To:       []string{userEmail},
		Subject:  "Login Validation Code",
		Html: htmlBody,
	})
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

*/
