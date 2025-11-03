package middleware

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// LiveReloadConfig holds configuration for live reload
type LiveReloadConfig struct {
	Enabled  bool
	WatchDir string
}

// FileChangeEvent represents a file change notification
type FileChangeEvent struct {
	File      string    `json:"file"`
	Timestamp time.Time `json:"timestamp"`
}

// computeDirectoryHash calculates MD5 hash of all files in admin/src directory
func computeDirectoryHash(dir string) (string, error) {
	hash := md5.New()

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-TS files
		if d.IsDir() || !strings.HasSuffix(path, ".ts") {
			return nil
		}

		// Get file info
		info, err := d.Info()
		if err != nil {
			return err
		}

		// Add file path and modified time to hash
		fmt.Fprintf(hash, "%s:%d:", path, info.ModTime().Unix())

		return nil
	})

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// LiveReloadHash returns a hash of all files in admin/src for polling
func LiveReloadHash(config LiveReloadConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !config.Enabled {
			return c.Status(404).JSON(fiber.Map{"error": "Live reload not enabled"})
		}

		hash, err := computeDirectoryHash(config.WatchDir)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to compute hash"})
		}

		return c.JSON(fiber.Map{
			"hash":      hash,
			"timestamp": time.Now().Unix(),
		})
	}
}

// LiveReloadWatch implements Server-Sent Events for file watching
func LiveReloadWatch(config LiveReloadConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !config.Enabled {
			return c.Status(404).SendString("Live reload not enabled")
		}

		// Set SSE headers
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("X-Accel-Buffering", "no")

		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			lastHash := ""
			ticker := time.NewTicker(1 * time.Second)
			defer ticker.Stop()

			// Send initial ping
			fmt.Fprintf(w, "data: {\"status\":\"connected\"}\n\n")
			w.Flush()

			for {
				select {
				case <-ticker.C:
					// Compute current hash
					currentHash, err := computeDirectoryHash(config.WatchDir)
					if err != nil {
						continue
					}

					// Initialize on first run
					if lastHash == "" {
						lastHash = currentHash
						continue
					}

					// Detect change
					if currentHash != lastHash {
						event := FileChangeEvent{
							File:      config.WatchDir,
							Timestamp: time.Now(),
						}

						eventJSON, _ := json.Marshal(event)
						fmt.Fprintf(w, "data: %s\n\n", eventJSON)
						w.Flush()

						lastHash = currentHash
					}

				case <-c.Context().Done():
					// Client disconnected
					return
				}
			}
		})

		return nil
	}
}

// SetupLiveReload adds live reload endpoints to the app
func SetupLiveReload(app *fiber.App) {
	// Only enable in development
	env := os.Getenv("APP_ENV")
	if env != "dev" {
		return
	}

	config := LiveReloadConfig{
		Enabled:  true,
		WatchDir: "./admin/src",
	}

	// Register endpoints
	app.Get("/admin/dev/hash", LiveReloadHash(config))
	app.Get("/admin/dev/watch", LiveReloadWatch(config))

	fmt.Println("🔄 Live reload enabled - watching ./admin/src")
}
