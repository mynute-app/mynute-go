package email

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTemplateRenderer(t *testing.T) {
	renderer := NewTemplateRenderer("templates", "translations")
	assert.NotNil(t, renderer)
	assert.Equal(t, "templates", renderer.templateDir)
	assert.Equal(t, "translations", renderer.translationDir)
	assert.Equal(t, "en", renderer.defaultLanguage)
}

func TestTemplateRenderer_RenderEmail(t *testing.T) {
	// Create temporary directories for test files
	tempDir := t.TempDir()
	templateDir := filepath.Join(tempDir, "templates")
	translationDir := filepath.Join(tempDir, "translations")

	require.NoError(t, os.MkdirAll(templateDir, 0755))
	require.NoError(t, os.MkdirAll(translationDir, 0755))

	// Create test template
	templateContent := `<html><body><h1>{{.title}}</h1><p>{{.message}}</p><p>Code: {{.code}}</p></body></html>`
	templatePath := filepath.Join(templateDir, "test_email.html")
	require.NoError(t, os.WriteFile(templatePath, []byte(templateContent), 0644))

	// Create test translations
	translationContent := `{
		"en": {
			"title": "Welcome",
			"message": "Hello World"
		},
		"pt": {
			"title": "Bem-vindo",
			"message": "Olá Mundo"
		},
		"es": {
			"title": "Bienvenido",
			"message": "Hola Mundo"
		}
	}`
	translationPath := filepath.Join(translationDir, "test_email.json")
	require.NoError(t, os.WriteFile(translationPath, []byte(translationContent), 0644))

	renderer := NewTemplateRenderer(templateDir, translationDir)

	t.Run("should render email with English translations by default", func(t *testing.T) {
		html, err := renderer.RenderEmail("test_email", "", TemplateData{
			"code": "12345",
		})

		assert.NoError(t, err)
		assert.Contains(t, html, "<h1>Welcome</h1>")
		assert.Contains(t, html, "<p>Hello World</p>")
		assert.Contains(t, html, "<p>Code: 12345</p>")
	})

	t.Run("should render email with Portuguese translations", func(t *testing.T) {
		html, err := renderer.RenderEmail("test_email", "pt", TemplateData{
			"code": "67890",
		})

		assert.NoError(t, err)
		assert.Contains(t, html, "<h1>Bem-vindo</h1>")
		assert.Contains(t, html, "<p>Olá Mundo</p>")
		assert.Contains(t, html, "<p>Code: 67890</p>")
	})

	t.Run("should render email with Spanish translations", func(t *testing.T) {
		html, err := renderer.RenderEmail("test_email", "es", TemplateData{
			"code": "ABCDE",
		})

		assert.NoError(t, err)
		assert.Contains(t, html, "<h1>Bienvenido</h1>")
		assert.Contains(t, html, "<p>Hola Mundo</p>")
		assert.Contains(t, html, "<p>Code: ABCDE</p>")
	})

	t.Run("should allow custom data to override translations", func(t *testing.T) {
		html, err := renderer.RenderEmail("test_email", "en", TemplateData{
			"title": "Custom Title",
			"code":  "99999",
		})

		assert.NoError(t, err)
		assert.Contains(t, html, "<h1>Custom Title</h1>")
		assert.Contains(t, html, "<p>Hello World</p>")
		assert.Contains(t, html, "<p>Code: 99999</p>")
	})

	t.Run("should return error for unsupported language", func(t *testing.T) {
		html, err := renderer.RenderEmail("test_email", "fr", TemplateData{
			"code": "12345",
		})

		assert.Error(t, err)
		assert.Empty(t, html)
		assert.Contains(t, err.Error(), "language 'fr' not found")
	})

	t.Run("should return error for non-existent template", func(t *testing.T) {
		// Create translation file but not template file
		noTemplateTrans := filepath.Join(translationDir, "non_existent.json")
		require.NoError(t, os.WriteFile(noTemplateTrans, []byte(`{"en": {"test": "value"}}`), 0644))

		html, err := renderer.RenderEmail("non_existent", "en", TemplateData{})

		assert.Error(t, err)
		assert.Empty(t, html)
		assert.Contains(t, err.Error(), "failed to parse template")
	})

	t.Run("should return error for non-existent translation file", func(t *testing.T) {
		// Create a template without corresponding translations
		noTransPath := filepath.Join(templateDir, "no_trans.html")
		require.NoError(t, os.WriteFile(noTransPath, []byte("<html><body>{{.test}}</body></html>"), 0644))

		html, err := renderer.RenderEmail("no_trans", "en", TemplateData{})

		assert.Error(t, err)
		assert.Empty(t, html)
		assert.Contains(t, err.Error(), "failed to load translations")
	})
}

func TestTemplateRenderer_LoadTranslations(t *testing.T) {
	tempDir := t.TempDir()
	translationDir := filepath.Join(tempDir, "translations")
	require.NoError(t, os.MkdirAll(translationDir, 0755))

	// Create test translation file
	translationContent := `{
		"en": {
			"key1": "value1",
			"key2": "value2"
		},
		"pt": {
			"key1": "valor1",
			"key2": "valor2"
		}
	}`
	translationPath := filepath.Join(translationDir, "test.json")
	require.NoError(t, os.WriteFile(translationPath, []byte(translationContent), 0644))

	renderer := NewTemplateRenderer("", translationDir)

	t.Run("should load English translations", func(t *testing.T) {
		translations, err := renderer.loadTranslations("test", "en")

		assert.NoError(t, err)
		assert.Equal(t, "value1", translations["key1"])
		assert.Equal(t, "value2", translations["key2"])
	})

	t.Run("should load Portuguese translations", func(t *testing.T) {
		translations, err := renderer.loadTranslations("test", "pt")

		assert.NoError(t, err)
		assert.Equal(t, "valor1", translations["key1"])
		assert.Equal(t, "valor2", translations["key2"])
	})

	t.Run("should return error for non-existent language", func(t *testing.T) {
		translations, err := renderer.loadTranslations("test", "de")

		assert.Error(t, err)
		assert.Nil(t, translations)
		assert.Contains(t, err.Error(), "language 'de' not found")
	})

	t.Run("should return error for invalid JSON", func(t *testing.T) {
		invalidPath := filepath.Join(translationDir, "invalid.json")
		require.NoError(t, os.WriteFile(invalidPath, []byte("invalid json"), 0644))

		translations, err := renderer.loadTranslations("invalid", "en")

		assert.Error(t, err)
		assert.Nil(t, translations)
		assert.Contains(t, err.Error(), "failed to parse translation JSON")
	})
}

func TestTemplateRenderer_MergeData(t *testing.T) {
	renderer := NewTemplateRenderer("", "")

	t.Run("should merge translations with empty custom data", func(t *testing.T) {
		translations := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		}
		customData := TemplateData{}

		merged := renderer.mergeData(translations, customData)

		assert.Equal(t, "value1", merged["key1"])
		assert.Equal(t, "value2", merged["key2"])
	})

	t.Run("should merge translations with custom data", func(t *testing.T) {
		translations := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		}
		customData := TemplateData{
			"key3": "value3",
		}

		merged := renderer.mergeData(translations, customData)

		assert.Equal(t, "value1", merged["key1"])
		assert.Equal(t, "value2", merged["key2"])
		assert.Equal(t, "value3", merged["key3"])
	})

	t.Run("should allow custom data to override translations", func(t *testing.T) {
		translations := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		}
		customData := TemplateData{
			"key2": "overridden",
			"key3": "value3",
		}

		merged := renderer.mergeData(translations, customData)

		assert.Equal(t, "value1", merged["key1"])
		assert.Equal(t, "overridden", merged["key2"])
		assert.Equal(t, "value3", merged["key3"])
	})
}

func TestTemplateRenderer_RealLoginValidation(t *testing.T) {
	// Test with actual login_validation files if they exist
	// Use absolute path construction from test location
	staticPath := filepath.Join("..", "..", "..", "static")
	translationPath := filepath.Join("..", "..", "..", "translation")

	// Check if files exist before running tests
	if _, err := os.Stat(filepath.Join(translationPath, "login_validation.json")); os.IsNotExist(err) {
		t.Skip("Skipping test - login_validation.json not found in translation directory")
	}
	if _, err := os.Stat(filepath.Join(staticPath, "login_validation.html")); os.IsNotExist(err) {
		t.Skip("Skipping test - login_validation.html not found in static directory")
	}

	renderer := NewTemplateRenderer(staticPath, translationPath)

	t.Run("should render login_validation email in English", func(t *testing.T) {
		html, err := renderer.RenderEmail("login_validation", "en", TemplateData{
			"ValidationCode": "123456",
		})

		require.NoError(t, err)
		assert.Contains(t, html, "Login Validation Code")
		assert.Contains(t, html, "123456")
		assert.Contains(t, html, "Your Login Validation Code")
	})

	t.Run("should render login_validation email in Portuguese", func(t *testing.T) {
		html, err := renderer.RenderEmail("login_validation", "pt", TemplateData{
			"ValidationCode": "654321",
		})

		require.NoError(t, err)
		assert.Contains(t, html, "Código de Validação de Login")
		assert.Contains(t, html, "654321")
		assert.Contains(t, html, "Seu Código de Validação de Login")
	})

	t.Run("should render login_validation email in Spanish", func(t *testing.T) {
		html, err := renderer.RenderEmail("login_validation", "es", TemplateData{
			"ValidationCode": "ABCDEF",
		})

		require.NoError(t, err)
		assert.Contains(t, html, "Código de Validación de Inicio de Sesión")
		assert.Contains(t, html, "ABCDEF")
		assert.Contains(t, html, "Su Código de Validación de Inicio de Sesión")
	})
}
