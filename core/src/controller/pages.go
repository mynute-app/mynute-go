package controller

import (
	"encoding/json"
	"fmt"
	"mynute-go/core/src/lib"
	"os"

	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	root, err := lib.FindProjectRoot()
	if err != nil {
		return err
	}
	return c.SendFile(root + "/static/page/index.html")
}

func VerifyEmailPage(c *fiber.Ctx) error {
	root, err := lib.FindProjectRoot()
	if err != nil {
		return err
	}

	// Simply serve the HTML file without template processing
	return c.SendFile(root + "/static/page/verify-email.html")
}

// GetPageTranslations returns translations for a specific page
func GetPageTranslations(c *fiber.Ctx) error {
	root, err := lib.FindProjectRoot()
	if err != nil {
		return err
	}

	page := c.Params("page")
	lang := c.Query("language", "en")

	// Read the translation file
	translationPath := fmt.Sprintf("%s/translation/page/%s.json", root, page)
	data, err := os.ReadFile(translationPath)
	if err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("failed to read translation file: %w", err))
	}

	// Parse the JSON
	var translations map[string]map[string]string
	if err := json.Unmarshal(data, &translations); err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("failed to parse translation file: %w", err))
	}

	// Get the requested language, default to "en" if not found
	langData, exists := translations[lang]
	if !exists {
		langData = translations["en"]
	}

	return c.JSON(langData)
}
