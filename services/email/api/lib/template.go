package lib

import (
	"bytes"
	"fmt"
	"html/template"
)

// TemplateData holds the data to be inserted into email templates
type TemplateData map[string]any

// TemplateRenderer handles email template rendering
type TemplateRenderer struct {
	defaultLanguage string
}

// NewTemplateRenderer creates a new template renderer
func NewTemplateRenderer(templateDir, translationDir string) *TemplateRenderer {
	return &TemplateRenderer{
		defaultLanguage: "en",
	}
}

// RenderFromString renders an email template from a string with provided data
// templateHTML: the HTML template as a string
// data: data to merge into the template
func (r *TemplateRenderer) RenderFromString(templateHTML string, data TemplateData) (string, error) {
	// Parse template from string
	tmpl, err := template.New("email").Parse(templateHTML)
	if err != nil {
		return "", fmt.Errorf("failed to parse template string: %w", err)
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
