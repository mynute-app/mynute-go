package middleware

import (
	"os"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/gofiber/fiber/v2"
)

// TypeScriptTranspiler creates a middleware that transpiles TypeScript files on-the-fly
func TypeScriptTranspiler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Only process .ts files
		if !strings.HasSuffix(c.Path(), ".ts") {
			return c.Next()
		}

		// Build the file path
		filePath := "./admin" + strings.TrimPrefix(c.Path(), "/admin")

		// Read the TypeScript file
		tsCode, err := os.ReadFile(filePath)
		if err != nil {
			return c.Status(404).SendString("File not found")
		}

		// Transpile TypeScript to JavaScript using esbuild
		result := api.Transform(string(tsCode), api.TransformOptions{
			Loader: api.LoaderTS,
			Format: api.FormatESModule,
			Target: api.ES2020,
		})

		if len(result.Errors) > 0 {
			return c.Status(500).SendString("TypeScript compilation error: " + result.Errors[0].Text)
		}

		// Set correct MIME type and send transpiled JavaScript
		c.Set("Content-Type", "application/javascript; charset=utf-8")
		return c.Send(result.Code)
	}
}

