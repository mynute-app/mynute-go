package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTemplateRenderer(t *testing.T) {
	renderer := NewTemplateRenderer("", "")
	assert.NotNil(t, renderer)
	assert.Equal(t, "en", renderer.defaultLanguage)
}

func TestTemplateRenderer_RenderFromString(t *testing.T) {
	renderer := NewTemplateRenderer("", "")

	t.Run("should render simple template", func(t *testing.T) {
		templateHTML := `<html><body><h1>{{.title}}</h1></body></html>`
		data := TemplateData{
			"title": "Welcome",
		}

		result, err := renderer.RenderFromString(templateHTML, data)

		assert.NoError(t, err)
		assert.Contains(t, result, "<h1>Welcome</h1>")
	})

	t.Run("should render template with multiple variables", func(t *testing.T) {
		templateHTML := `<html><body><h1>{{.greeting}}</h1><p>{{.message}}</p><p>Code: {{.code}}</p></body></html>`
		data := TemplateData{
			"greeting": "Hello",
			"message":  "Welcome to our service",
			"code":     "12345",
		}

		result, err := renderer.RenderFromString(templateHTML, data)

		assert.NoError(t, err)
		assert.Contains(t, result, "<h1>Hello</h1>")
		assert.Contains(t, result, "<p>Welcome to our service</p>")
		assert.Contains(t, result, "<p>Code: 12345</p>")
	})

	t.Run("should handle translations merged with custom data", func(t *testing.T) {
		templateHTML := `<html><body><h1>{{.subject}}</h1><p>{{.greeting}}</p><p>User: {{.username}}</p></body></html>`

		// Simulate what the controller does: merge translations + custom data
		translations := map[string]interface{}{
			"subject":  "Email Verification",
			"greeting": "Hello",
		}
		customData := map[string]interface{}{
			"username": "John Doe",
		}

		// Merge them
		mergedData := make(TemplateData)
		for k, v := range translations {
			mergedData[k] = v
		}
		for k, v := range customData {
			mergedData[k] = v
		}

		result, err := renderer.RenderFromString(templateHTML, mergedData)

		assert.NoError(t, err)
		assert.Contains(t, result, "<h1>Email Verification</h1>")
		assert.Contains(t, result, "<p>Hello</p>")
		assert.Contains(t, result, "<p>User: John Doe</p>")
	})

	t.Run("should allow custom data to override translations", func(t *testing.T) {
		templateHTML := `<html><body><h1>{{.title}}</h1></body></html>`

		// Simulate override scenario
		translations := map[string]interface{}{
			"title": "Default Title",
		}
		customData := map[string]interface{}{
			"title": "Custom Title",
		}

		// Merge with custom data taking precedence
		mergedData := make(TemplateData)
		for k, v := range translations {
			mergedData[k] = v
		}
		for k, v := range customData {
			mergedData[k] = v
		}

		result, err := renderer.RenderFromString(templateHTML, mergedData)

		assert.NoError(t, err)
		assert.Contains(t, result, "<h1>Custom Title</h1>")
		assert.NotContains(t, result, "Default Title")
	})

	t.Run("should return error for invalid template syntax", func(t *testing.T) {
		templateHTML := `<html><body>{{.invalid syntax}}</body></html>`
		data := TemplateData{}

		result, err := renderer.RenderFromString(templateHTML, data)

		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "failed to parse template string")
	})

	t.Run("should handle empty template", func(t *testing.T) {
		templateHTML := ``
		data := TemplateData{}

		result, err := renderer.RenderFromString(templateHTML, data)

		assert.NoError(t, err)
		assert.Equal(t, "", result)
	})

	t.Run("should handle template with no variables", func(t *testing.T) {
		templateHTML := `<html><body><h1>Static Content</h1></body></html>`
		data := TemplateData{}

		result, err := renderer.RenderFromString(templateHTML, data)

		assert.NoError(t, err)
		assert.Equal(t, templateHTML, result)
	})

	t.Run("should render complex email template from Core service", func(t *testing.T) {
		// Simulate a real email template that Core would send
		templateHTML := `<!DOCTYPE html>
<html>
<head><title>{{.subject}}</title></head>
<body>
	<h1>{{.heading}}</h1>
	<p>{{.instruction}}</p>
	<p>Your verification code is: <strong>{{.code}}</strong></p>
	<a href="{{.link}}">{{.button_text}}</a>
	<p>{{.footer_notice}}</p>
</body>
</html>`

		// Translations from Core service
		translations := map[string]interface{}{
			"subject":       "Verify Your Email",
			"heading":       "Email Verification",
			"instruction":   "Please use the code below to verify your email",
			"button_text":   "Verify Email",
			"footer_notice": "This is an automated message",
		}

		// Custom data from Core service
		customData := map[string]interface{}{
			"code": "ABC123",
			"link": "https://example.com/verify?code=ABC123",
		}

		// Merge
		mergedData := make(TemplateData)
		for k, v := range translations {
			mergedData[k] = v
		}
		for k, v := range customData {
			mergedData[k] = v
		}

		result, err := renderer.RenderFromString(templateHTML, mergedData)

		require.NoError(t, err)
		assert.Contains(t, result, "<title>Verify Your Email</title>")
		assert.Contains(t, result, "<h1>Email Verification</h1>")
		assert.Contains(t, result, "<strong>ABC123</strong>")
		assert.Contains(t, result, `href="https://example.com/verify?code=ABC123"`)
		assert.Contains(t, result, "Verify Email")
		assert.Contains(t, result, "This is an automated message")
	})
}
