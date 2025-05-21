// File: uploader/uploader_test.go
package uploader

import (
	"agenda-kaki-go/core/config/namespace"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
)

var pngData = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A} // PNG
var pdfData = []byte("%PDF-1.4")
var invalidData = []byte("Hello, not an image")

func TestGenerateUniqueFilename(t *testing.T) {
	id := uuid.New()
	name := GenerateUniqueFilename("client", id.String(), "logo.png")
	if !strings.HasSuffix(name, ".png") {
		t.Errorf("expected .png suffix, got: %s", name)
	}
	if !strings.Contains(name, id.String()) {
		t.Errorf("expected path to contain UUID: %s", name)
	}
	if !strings.HasPrefix(name, "client/") {
		t.Errorf("expected prefix to be 'client/', got: %s", name)
	}
}

func TestExtractFilenameFromURL(t *testing.T) {
	url := fmt.Sprintf("https://cdn.site.com/assets/%s/logo.png", namespace.UploadsFolder)
	expected := "logo.png"
	got := ExtractFilenameFromURL(url)
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func TestLocalUploader_SaveImage_Success(t *testing.T) {
	id := uuid.New()
	l := &LocalUploader{Entity: "client", EntityID: id.String()}
	url, err := l.Save("image", pngData, "logo.png")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.HasSuffix(url, ".png") {
		t.Errorf("expected .png url, got %s", url)
	}
	if err := os.RemoveAll(fmt.Sprintf("./%s/client/", namespace.UploadsFolder)); err != nil {
		t.Fatalf("failed to clean up: %v", err)
	}
	if _, err := os.Stat(fmt.Sprintf("./%s/client/%s/logo.png", namespace.UploadsFolder, id.String())); os.IsNotExist(err) {
		t.Fatalf("expected file to exist, got error: %v", err)
	}
}

func TestLocalUploader_SaveImage_Invalid(t *testing.T) {
	l := &LocalUploader{Entity: "client", EntityID: uuid.New().String()}
	_, err := l.Save("image", invalidData, "text.txt")
	if err == nil {
		t.Fatal("expected error for invalid image, got nil")
	}
}

func TestLocalUploader_SavePDF_Success(t *testing.T) {
	id := uuid.New()
	l := &LocalUploader{Entity: "client", EntityID: id.String()}
	url, err := l.Save("pdf", pdfData, "doc.pdf")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.HasSuffix(url, ".pdf") {
		t.Errorf("expected .pdf url, got %s", url)
	}
	if err := os.RemoveAll(fmt.Sprintf("./%s/client/", namespace.UploadsFolder)); err != nil {
		t.Fatalf("failed to clean up: %v", err)
	}
	if _, err := os.Stat(fmt.Sprintf("./%s/client/%s/doc.pdf", namespace.UploadsFolder, id.String())); os.IsNotExist(err) {
		t.Fatalf("expected file to exist, got error: %v", err)
	}
}

func TestLocalUploader_Delete_Success(t *testing.T) {
	id := uuid.New()
	l := &LocalUploader{Entity: "client", EntityID: id.String()}
	url, err := l.Save("image", pngData, "todelete.png")
	if err != nil {
		t.Fatalf("failed to save: %v", err)
	}
	if err := l.Delete(url); err != nil {
		t.Fatalf("expected successful delete, got error: %v", err)
	}
}

func TestLocalUploader_Replace(t *testing.T) {
	id := uuid.New()
	l := &LocalUploader{Entity: "client", EntityID: id.String()}
	oldURL, err := l.Save("image", pngData, "replace-me.png")
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}
	newURL, err := l.Replace("image", oldURL, pngData, "replaced.png")
	if err != nil {
		t.Fatalf("replace failed: %v", err)
	}
	if err := os.RemoveAll(fmt.Sprintf("./%s/client/%s/", namespace.UploadsFolder, id.String())); err != nil {
		t.Fatalf("failed to clean up: %v", err)
	}
	if oldURL == newURL {
		t.Errorf("expected new URL to differ from old one")
	}
	if !strings.Contains(newURL, id.String()) {
		t.Errorf("expected new URL to contain entity ID")
	}
}