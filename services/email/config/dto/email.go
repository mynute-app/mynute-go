package dto

// SendEmailRequest represents a request to send a single email
type SendEmailRequest struct {
	To      string   `json:"to" validate:"required,email" example:"recipient@example.com"`
	Subject string   `json:"subject" validate:"required" example:"Welcome to Mynute"`
	Body    string   `json:"body" validate:"required" example:"This is the email body"`
	CC      []string `json:"cc,omitempty" example:"cc1@example.com,cc2@example.com"`
	BCC     []string `json:"bcc,omitempty" example:"bcc@example.com"`
	IsHTML  bool     `json:"is_html" example:"true"`
}

// SendTemplateEmailRequest represents a request to send an email using a template
type SendTemplateEmailRequest struct {
	To           string                 `json:"to" validate:"required,email" example:"recipient@example.com"`
	TemplateName string                 `json:"template_name" validate:"required" example:"login_validation_code"`
	Language     string                 `json:"language" validate:"required" example:"en"`
	Data         map[string]interface{} `json:"data"`
	CC           []string               `json:"cc,omitempty" example:"cc@example.com"`
	BCC          []string               `json:"bcc,omitempty" example:"bcc@example.com"`
}

// SendTemplateMergeRequest represents a request to send an email by merging template HTML with translations
// This is used when the calling service (e.g., Core) sends the template HTML and translations directly
type SendTemplateMergeRequest struct {
	To           []string               `json:"to" validate:"required,dive,email" example:"recipient@example.com"`
	TemplateHTML string                 `json:"template_html" validate:"required" example:"<html><body>{{.greeting}}</body></html>"`
	Translations map[string]interface{} `json:"translations" validate:"required" example:"{\"subject\":\"Welcome\",\"greeting\":\"Hello\"}"`
	Data         map[string]interface{} `json:"data,omitempty" example:"{\"name\":\"John\"}"`
	CC           []string               `json:"cc,omitempty" example:"cc@example.com"`
	BCC          []string               `json:"bcc,omitempty" example:"bcc@example.com"`
}

// SendBulkEmailRequest represents a request to send emails to multiple recipients
type SendBulkEmailRequest struct {
	Recipients []string `json:"recipients" validate:"required,dive,email" example:"user1@example.com,user2@example.com"`
	Subject    string   `json:"subject" validate:"required" example:"Newsletter"`
	Body       string   `json:"body" validate:"required" example:"Monthly newsletter content"`
	IsHTML     bool     `json:"is_html" example:"true"`
}

// EmailSuccessResponse represents a successful email send response
type EmailSuccessResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Email sent successfully"`
	To      string `json:"to,omitempty" example:"recipient@example.com"`
}

// BulkEmailResponse represents a response for bulk email sending
type BulkEmailResponse struct {
	Success          bool     `json:"success" example:"true"`
	Total            int      `json:"total" example:"10"`
	Sent             int      `json:"sent" example:"9"`
	Failed           int      `json:"failed" example:"1"`
	FailedRecipients []string `json:"failed_recipients,omitempty" example:"failed@example.com"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid request body"`
	Details string `json:"details,omitempty" example:"Missing required field: to"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status  string `json:"status" example:"healthy"`
	Service string `json:"service" example:"email"`
	Version string `json:"version" example:"1.0"`
}
