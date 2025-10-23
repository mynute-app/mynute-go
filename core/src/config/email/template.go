package email

import (
	"maps"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

// TemplateData holds the data to be inserted into email templates
type TemplateData map[string]any

// TemplateRenderer handles email template rendering with translations
type TemplateRenderer struct {
	templateDir     string
	translationDir  string
	defaultLanguage string
}

// NewTemplateRenderer creates a new template renderer
func NewTemplateRenderer(templateDir, translationDir string) *TemplateRenderer {
	return &TemplateRenderer{
		templateDir:     templateDir,
		translationDir:  translationDir,
		defaultLanguage: "en",
	}
}

// RenderEmail renders an email template with translations and custom data
// templateName: name of the HTML template file (without extension)
// language: language code (e.g., "en", "pt", "es"). Defaults to "en" if empty
// customData: additional data to merge with translations
func (r *TemplateRenderer) RenderEmail(templateName, language string, customData TemplateData) (string, error) {
	// Default to English if no language specified
	if language == "" {
		language = r.defaultLanguage
	}

	// Load translations
	translations, err := r.loadTranslations(templateName, language)
	if err != nil {
		return "", fmt.Errorf("failed to load translations: %w", err)
	}

	// Merge custom data with translations
	templateData := r.mergeData(translations, customData)

	// Load and parse template
	templatePath := filepath.Join(r.templateDir, templateName+".html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// loadTranslations loads translations from the JSON file for the specified template and language
func (r *TemplateRenderer) loadTranslations(templateName, language string) (map[string]interface{}, error) {
	translationPath := filepath.Join(r.translationDir, templateName+".json")

	// Read translation file
	data, err := os.ReadFile(translationPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read translation file %s: %w", translationPath, err)
	}

	// Parse JSON
	var allTranslations map[string]map[string]interface{}
	if err := json.Unmarshal(data, &allTranslations); err != nil {
		return nil, fmt.Errorf("failed to parse translation JSON: %w", err)
	}

	// Get translations for the specified language
	translations, exists := allTranslations[language]
	if !exists {
		return nil, fmt.Errorf("language '%s' not found in translation file", language)
	}

	return translations, nil
}

// mergeData merges translations with custom data
// Custom data takes precedence over translations
func (r *TemplateRenderer) mergeData(translations map[string]any, customData TemplateData) TemplateData {
	merged := make(TemplateData)

	// Copy translations first
	maps.Copy(merged, translations)

	// Override with custom data
	maps.Copy(merged, customData)

	return merged
}
