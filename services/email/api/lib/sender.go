package lib

import (
	"context"
	"fmt"
	"os"
)

// Attachment is the public struct used for adding attachments to emails
type Attachment struct {
	// Content is the binary content of the attachment to use when a Path
	// is not available.
	Content []byte

	// Filename that will appear in the email.
	// Make sure you pick the correct extension otherwise preview
	// may not work as expected
	Filename string

	// Path where the attachment file is hosted instead of providing the
	// content directly.
	Path string

	// Content type for the attachment, if not set will be derived from
	// the filename property
	ContentType string

	// Optional content ID for the attachment, to be used as a reference in the HTML content.
	// If set, this attachment will be sent as an inline attachment and you can reference it
	// in the HTML content using the `cid:` prefix.
	ContentId string
}

type Tag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// EmailData holds all the necessary information for an email.
type EmailData struct {
	From        string            `json:"from"`
	To          []string          `json:"to"`
	Subject     string            `json:"subject"`
	Bcc         []string          `json:"bcc,omitempty"`
	Cc          []string          `json:"cc,omitempty"`
	ReplyTo     string            `json:"reply_to,omitempty"`
	Html        string            `json:"html,omitempty"`
	Text        string            `json:"text,omitempty"`
	Tags        []Tag             `json:"tags,omitempty"`
	Attachments []*Attachment     `json:"attachments,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	ScheduledAt string            `json:"scheduled_at,omitempty"`
}

// Sender defines the contract for any email sending service.
// This allows for easy mocking and swapping of implementations (e.g., for testing).
type Sender interface {
	Send(ctx context.Context, data EmailData) error
}

type ProviderOpts struct {
	Provider string
}

func NewProvider(opts *ProviderOpts) (Sender, error) {
	if opts == nil {
		APP_ENV := os.Getenv("APP_ENV")
		switch APP_ENV {
		case "dev", "test":
			return MailHog()
		case "prod":
			return Resend()
		}
		return nil, fmt.Errorf("email provider not specified and APP_ENV (%s) is not set to a known value", APP_ENV)
	}
	switch opts.Provider {
	case "resend":
		return Resend()
	case "mailhog":
		return MailHog()
	default:
		return nil, fmt.Errorf("email provider (%s) not implemented", opts.Provider)
	}
}
