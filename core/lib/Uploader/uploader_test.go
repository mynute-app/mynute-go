// File: uploader/uploader_test.go
package myUploader

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib" // Assuming lib.FindProjectRoot is here
	"agenda-kaki-go/core/lib/FileBytes"
	"fmt"
	"os"
	"path/filepath" // Import for filepath operations
	"strings"
	"testing"

	"github.com/google/uuid"
)

var pdfData = []byte("%PDF-1.4")
var invalidData = []byte("Hello, not an image")

// Helper to get root_path for tests
func getTestRootPath(t *testing.T) string {
	t.Helper()
	rootPath, err := lib.FindProjectRoot()
	if err != nil {
		t.Fatalf("Setup: Failed to find project root: %v", err)
	}
	return rootPath
}

func TestGenerateUniqueFilename(t *testing.T) {
	id := uuid.New()
	name := GenerateUniqueFilename("client", id.String(), "logo.png")
	expectedSuffix := ".png"
	if !strings.HasSuffix(name, expectedSuffix) {
		t.Errorf("expected suffix %s, got: %s", expectedSuffix, name)
	}
	if !strings.Contains(name, id.String()) {
		t.Errorf("expected path to contain UUID %s, got: %s", id.String(), name)
	}
	// GenerateUniqueFilename returns a path like "entity/entityID/name_uuid.ext"
	// So the prefix should be "client/" + id.String() + "/"
	expectedPrefix := filepath.Join("client", id.String()) // Use filepath.Join for OS-agnostic paths
	// Convert to forward slashes for comparison as GenerateUniqueFilename uses /
	if !strings.HasPrefix(filepath.ToSlash(name), filepath.ToSlash(expectedPrefix)+"/") {
		t.Errorf("expected prefix to be '%s/', got: '%s'", expectedPrefix, name)
	}
}

func TestExtractFilenameFromURL(t *testing.T) {
	// This test seems fine as is, assuming namespace.UploadsFolder is correctly defined.
	// For robustness, use actual value if namespace.UploadsFolder can change.
	// Assuming namespace.UploadsFolder = "local_uploads" for this example URL construction.
	url := fmt.Sprintf("https://cdn.site.com/assets/%s/client/some_uuid/logo_another_uuid.png", "local_uploads") // Example usage
	expected := "logo_another_uuid.png"
	got := ExtractFilenameFromURL(url)
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}

	urlLocal := fmt.Sprintf("/static/%s/client/some_uuid/doc_another_uuid.pdf", "local_uploads")
	expectedLocal := "doc_another_uuid.pdf"
	gotLocal := ExtractFilenameFromURL(urlLocal)
	if gotLocal != expectedLocal {
		t.Errorf("expected %s for local URL, got %s", expectedLocal, gotLocal)
	}
}

func TestLocalUploader_SaveImage_Success(t *testing.T) {
	rootPath := getTestRootPath(t)
	id := uuid.New()
	l := &localUploader{Entity: "client", EntityID: id.String()}

	t.Cleanup(func() {
		cleanupDir := filepath.Join(rootPath, namespace.UploadsFolder, l.Entity, l.EntityID)
		if err := os.RemoveAll(cleanupDir); err != nil {
			t.Logf("Cleanup: failed to remove directory %s: %v", cleanupDir, err)
		}
	})

	originalFilename := "logo.png"
	url, err := l.Save("image", FileBytes.PNG_FILE_1, originalFilename)
	if err != nil {
		t.Fatalf("l.Save() error = %v, want nil", err)
	}
	if !strings.HasSuffix(url, expectedSuffixFromFilename(originalFilename)) {
		t.Errorf("expected URL to end with .png suffix, got %s", url)
	}

	// Verify file existence
	urlPrefix := "/static/" + filepath.ToSlash(namespace.UploadsFolder) + "/"
	if !strings.HasPrefix(url, urlPrefix) {
		t.Fatalf("URL '%s' does not have expected prefix '%s'", url, urlPrefix)
	}
	scopedPath := strings.TrimPrefix(url, urlPrefix)
	fullSavedPath := filepath.Join(rootPath, namespace.UploadsFolder, filepath.FromSlash(scopedPath)) // Ensure OS-specific path for Stat

	if _, statErr := os.Stat(fullSavedPath); os.IsNotExist(statErr) {
		t.Fatalf("expected file to exist at %s, but it does not: %v", fullSavedPath, statErr)
	}
}

func TestLocalUploader_SaveImage_Invalid(t *testing.T) {
	// This test should already pass if the validation logic is correct.
	// No changes needed based on the error messages.
	l := &localUploader{Entity: "client", EntityID: uuid.New().String()}
	_, err := l.Save("image", invalidData, "text.txt")
	if err == nil {
		t.Fatal("expected error for invalid image, got nil")
	}
}

func TestLocalUploader_SavePDF_Success(t *testing.T) {
	rootPath := getTestRootPath(t)
	id := uuid.New()
	l := &localUploader{Entity: "client", EntityID: id.String()}

	t.Cleanup(func() {
		cleanupDir := filepath.Join(rootPath, namespace.UploadsFolder, l.Entity, l.EntityID)
		if err := os.RemoveAll(cleanupDir); err != nil {
			t.Logf("Cleanup: failed to remove directory %s: %v", cleanupDir, err)
		}
	})

	originalFilename := "doc.pdf"
	url, err := l.Save("pdf", pdfData, originalFilename)
	if err != nil {
		t.Fatalf("l.Save() error = %v, want nil", err)
	}
	if !strings.HasSuffix(url, expectedSuffixFromFilename(originalFilename)) {
		t.Errorf("expected URL to end with .pdf suffix, got %s", url)
	}

	urlPrefix := "/static/" + filepath.ToSlash(namespace.UploadsFolder) + "/"
	if !strings.HasPrefix(url, urlPrefix) {
		t.Fatalf("URL '%s' does not have expected prefix '%s'", url, urlPrefix)
	}
	scopedPath := strings.TrimPrefix(url, urlPrefix)
	fullSavedPath := filepath.Join(rootPath, namespace.UploadsFolder, filepath.FromSlash(scopedPath))

	if _, statErr := os.Stat(fullSavedPath); os.IsNotExist(statErr) {
		t.Fatalf("expected file to exist at %s, but it does not: %v", fullSavedPath, statErr)
	}
}

func TestLocalUploader_Delete_Success(t *testing.T) {
	rootPath := getTestRootPath(t)
	id := uuid.New()
	l := &localUploader{Entity: "client", EntityID: id.String()}

	t.Cleanup(func() {
		cleanupDir := filepath.Join(rootPath, namespace.UploadsFolder, l.Entity, l.EntityID)
		// If os.RemoveAll fails, it's usually fine as the file might have been deleted by the test.
		os.RemoveAll(cleanupDir)
	})

	originalFilename := "todelete.png"
	url, err := l.Save("image", FileBytes.PNG_FILE_1, originalFilename)
	if err != nil {
		t.Fatalf("l.Save() failed: %v", err)
	}

	urlPrefix := "/static/" + filepath.ToSlash(namespace.UploadsFolder) + "/"
	if !strings.HasPrefix(url, urlPrefix) {
		t.Fatalf("URL '%s' does not have expected prefix '%s'", url, urlPrefix)
	}
	scopedPath := strings.TrimPrefix(url, urlPrefix)
	filePathToVerify := filepath.Join(rootPath, namespace.UploadsFolder, filepath.FromSlash(scopedPath))

	// Ensure file exists before delete attempt
	if _, statErr := os.Stat(filePathToVerify); os.IsNotExist(statErr) {
		t.Fatalf("file %s was not created by Save, cannot test Delete: %v", filePathToVerify, statErr)
	}

	if err := l.Delete(url); err != nil {
		t.Fatalf("l.Delete(%s) error = %v, want nil", url, err)
	}

	// Verify file no longer exists
	if _, statErr := os.Stat(filePathToVerify); !os.IsNotExist(statErr) {
		if statErr == nil {
			t.Errorf("expected file %s to be deleted, but it still exists", filePathToVerify)
		} else {
			t.Errorf("os.Stat(%s) after delete returned unexpected error: %v", filePathToVerify, statErr)
		}
	}
}

func TestLocalUploader_Replace(t *testing.T) {
	rootPath := getTestRootPath(t)
	id := uuid.New()
	l := &localUploader{Entity: "client", EntityID: id.String()}

	t.Cleanup(func() {
		cleanupDir := filepath.Join(rootPath, namespace.UploadsFolder, l.Entity, l.EntityID)
		os.RemoveAll(cleanupDir) // Attempt cleanup, ignore error if dir/files already gone
	})

	oldOriginalFilename := "replace-me.png"
	oldURL, err := l.Save("image", FileBytes.PNG_FILE_1, oldOriginalFilename)
	if err != nil {
		t.Fatalf("initial l.Save() for old file failed: %v", err)
	}

	urlPrefix := "/static/" + filepath.ToSlash(namespace.UploadsFolder) + "/"
	if !strings.HasPrefix(oldURL, urlPrefix) {
		t.Fatalf("Old URL '%s' does not have expected prefix '%s'", oldURL, urlPrefix)
	}
	oldScopedPath := strings.TrimPrefix(oldURL, urlPrefix)
	oldFilePath := filepath.Join(rootPath, namespace.UploadsFolder, filepath.FromSlash(oldScopedPath))

	// Ensure old file exists
	if _, statErr := os.Stat(oldFilePath); os.IsNotExist(statErr) {
		t.Fatalf("old file %s was not created by Save, cannot test Replace: %v", oldFilePath, statErr)
	}

	newOriginalFilename := "replaced.png"
	newURL, err := l.Replace("image", oldURL, FileBytes.PNG_FILE_2, newOriginalFilename)
	if err != nil {
		t.Fatalf("l.Replace() failed: %v", err)
	}

	// Verify old file is deleted
	if _, statErr := os.Stat(oldFilePath); !os.IsNotExist(statErr) {
		if statErr == nil {
			t.Errorf("expected old file %s to be deleted by Replace, but it still exists", oldFilePath)
		} else {
			t.Errorf("os.Stat(%s) for old file after replace returned unexpected error: %v", oldFilePath, statErr)
		}
	}

	// Verify new file exists
	if !strings.HasPrefix(newURL, urlPrefix) {
		t.Fatalf("New URL '%s' does not have expected prefix '%s'", newURL, urlPrefix)
	}
	newScopedPath := strings.TrimPrefix(newURL, urlPrefix)
	newFilePath := filepath.Join(rootPath, namespace.UploadsFolder, filepath.FromSlash(newScopedPath))

	if _, statErr := os.Stat(newFilePath); os.IsNotExist(statErr) {
		t.Fatalf("expected new file %s to exist after Replace, but it does not: %v", newFilePath, statErr)
	}

	// Original checks
	if oldURL == newURL {
		t.Errorf("expected new URL to differ from old one (old: %s, new: %s)", oldURL, newURL)
	}
	if !strings.Contains(newURL, id.String()) {
		t.Errorf("expected new URL to contain entity ID (%s), got: %s", id.String(), newURL)
	}
	if !strings.HasSuffix(newURL, expectedSuffixFromFilename(newOriginalFilename)) {
		t.Errorf("expected new URL to end with %s suffix, got %s", expectedSuffixFromFilename(newOriginalFilename), newURL)
	}
	// Check if the new URL path contains the base name of the new file (without extension and UUID)
	newBaseName := strings.TrimSuffix(newOriginalFilename, filepath.Ext(newOriginalFilename))
	if !strings.Contains(newURL, newBaseName) {
		t.Errorf("expected new URL to contain new base filename '%s', got: %s", newBaseName, newURL)
	}
}

// Helper to get expected suffix from original filename (e.g. ".png")
func expectedSuffixFromFilename(filename string) string {
	return filepath.Ext(filename)
}
