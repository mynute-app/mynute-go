package uploader

import (
	"os"
	"strings"
	"testing"
)

// mocks
var pngData = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A} // PNG
var pdfData = []byte("%PDF-1.4")
var invalidData = []byte("Hello, not an image")

func TestGenerateUniqueFilename(t *testing.T) {
	name := GenerateUniqueFilename("logo.png")
	if !strings.HasSuffix(name, ".png") {
		t.Errorf("expected .png suffix, got: %s", name)
	}
}

func TestExtractFilenameFromURL(t *testing.T) {
	url := "https://cdn.site.com/assets/uploads/logo.png"
	expected := "logo.png"
	got := ExtractFilenameFromURL(url)
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func TestLocalUploader_SaveImage_Success(t *testing.T) {
	l := &LocalUploader{}
	url, err := l.Save("image", pngData, "logo.png")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	filename := ExtractFilenameFromURL(url)
	defer os.Remove("./uploads/" + filename)
	if !strings.HasSuffix(url, ".png") {
		t.Errorf("expected .png url, got %s", url)
	}
}

func TestLocalUploader_SaveImage_Invalid(t *testing.T) {
	l := &LocalUploader{}
	_, err := l.Save("image", invalidData, "text.txt")
	if err == nil {
		t.Fatal("expected error for invalid image, got nil")
	}
}

func TestLocalUploader_SavePDF_Success(t *testing.T) {
	l := &LocalUploader{}
	url, err := l.Save("pdf", pdfData, "doc.pdf")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	filename := ExtractFilenameFromURL(url)
	defer os.Remove("./uploads/" + filename)
	if !strings.HasSuffix(url, ".pdf") {
		t.Errorf("expected .pdf url, got %s", url)
	}
}

func TestLocalUploader_Delete_Success(t *testing.T) {
	l := &LocalUploader{}
	// Save something first
	url, err := l.Save("image", pngData, "todelete.png")
	if err != nil {
		t.Fatalf("failed to save: %v", err)
	}
	if err := l.Delete(url); err != nil {
		t.Fatalf("expected successful delete, got error: %v", err)
	}
}

func TestLocalUploader_Replace(t *testing.T) {
	l := &LocalUploader{}
	// Save first file
	oldURL, err := l.Save("image", pngData, "replace-me.png")
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}
	// Replace with a new one
	newURL, err := l.Replace("image", oldURL, pngData, "replaced.png")
	if err != nil {
		t.Fatalf("replace failed: %v", err)
	}
	oldFilename := ExtractFilenameFromURL(oldURL)
	newFilename := ExtractFilenameFromURL(newURL)
	defer os.Remove("./uploads/" + newFilename)

	// Ensure old file is gone, new file exists
	if _, err := os.Stat("./uploads/" + oldFilename); err == nil {
		t.Errorf("expected old file %s to be deleted", oldFilename)
	}
	if _, err := os.Stat("./uploads/" + newFilename); os.IsNotExist(err) {
		t.Errorf("expected new file %s to exist", newFilename)
	}
}