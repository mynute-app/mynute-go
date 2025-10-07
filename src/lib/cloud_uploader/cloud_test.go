// File: uploader/cloud_test.go
package myUploader

import (
	"context"
	FileBytes "mynute-go/src/lib/file_bytes"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type mockS3Client struct {
	PutCalled    bool
	DeleteCalled bool
	LastKey      string
}

func (m *mockS3Client) PutObject(ctx context.Context, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	m.PutCalled = true
	m.LastKey = *input.Key
	return &s3.PutObjectOutput{}, nil
}

func (m *mockS3Client) DeleteObject(ctx context.Context, input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	m.DeleteCalled = true
	m.LastKey = *input.Key
	return &s3.DeleteObjectOutput{}, nil
}

func TestCloudUploader_Save_And_Delete(t *testing.T) {
	id := uuid.New().String()
	mockClient := &mockS3Client{}
	uploader := NewCloudUploader("client", id, mockClient, "my-bucket", "https://cdn.test.com")

	// Save test
	file := FileBytes.PNG_FILE_1
	filename := "test-image.png"
	url, err := uploader.Save("image", file, filename)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if !mockClient.PutCalled {
		t.Error("expected PutObject to be called")
	}
	if !strings.HasPrefix(url, "https://cdn.test.com/my-bucket/client/"+id) {
		t.Errorf("expected URL to contain client/%s, got %s", id, url)
	}

	// Delete test
	err = uploader.Delete(url)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if !mockClient.DeleteCalled {
		t.Error("expected DeleteObject to be called")
	}
}

func TestCloudUploader_Replace(t *testing.T) {
	id := uuid.New().String()
	mockClient := &mockS3Client{}
	uploader := NewCloudUploader("client", id, mockClient, "my-bucket", "https://cdn.test.com")

	oldURL := "https://cdn.test.com/client/" + id + "/old_image_1234.png"
	newFile := FileBytes.PNG_FILE_2
	newName := "new-image.png"

	url, err := uploader.Replace("image", oldURL, newFile, newName)
	if err != nil {
		t.Fatalf("Replace failed: %v", err)
	}
	if !mockClient.PutCalled || !mockClient.DeleteCalled {
		t.Error("expected both PutObject and DeleteObject to be called")
	}
	if !strings.Contains(url, newName[:3]) {
		t.Errorf("expected new URL to contain base filename '%s', got: %s", newName[:3], url)
	}
}
