package dto

// SendEmailRequest represents a request to send emails to one or more recipients
type SendEmailRequest struct {
	Subject    string      `json:"subject" validate:"required" example:"Welcome to Mynute"`
	Body       string      `json:"body" validate:"required" example:"This is the email body"`
	IsHTML     bool        `json:"is_html" example:"true"`
	Recipients []Recipient `json:"recipients" validate:"required,dive" example:"[{\"to\":[\"user@example.com\"]}]"`
}

// Recipient represents an email recipient with optional CC and BCC
type Recipient struct {
	To  []string `json:"to" validate:"required,min=1,dive,email" example:"recipient@example.com"`
	CC  []string `json:"cc,omitempty" validate:"omitempty,dive,email" example:"cc@example.com"`
	BCC []string `json:"bcc,omitempty" validate:"omitempty,dive,email" example:"bcc@example.com"`
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

// EmailSuccessResponse represents a successful email send response
type EmailSuccessResponse struct {
	Success bool     `json:"success" example:"true"`
	Message string   `json:"message" example:"Email sent successfully"`
	Total   int      `json:"total" example:"5"`
	Sent    int      `json:"sent" example:"4"`
	Failed  int      `json:"failed,omitempty" example:"1"`
	Results []Result `json:"results,omitempty"`
}

// Result represents the result of sending to a specific recipient
type Result struct {
	To      []string `json:"to" example:"recipient@example.com"`
	Success bool     `json:"success" example:"true"`
	Error   string   `json:"error,omitempty" example:"Failed to send"`
}

// BulkEmailResponse represents a response for bulk email sending (deprecated, use EmailSuccessResponse)
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
