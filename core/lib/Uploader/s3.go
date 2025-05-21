package uploader

import "fmt"

type S3Uploader struct{}

func (s *S3Uploader) Save(fileType string, file []byte, filename string) (string, error) {
	strategy, err := getStrategy(s, fileType)
	if err != nil {
		return "", err
	}
	return strategy(file, filename)
}

func (s *S3Uploader) Delete(fileURL string) error {
	filename := ExtractFilenameFromURL(fileURL)
	if filename == "" {
		return fmt.Errorf("invalid file URL: %s", fileURL)
	}
	// Ex: s3Client.DeleteObject(bucket, filename)
	return nil
}

func (s *S3Uploader) Replace(fileType string, oldURL string, newFile []byte, originalFilename string) (string, error) {
	if err := s.Delete(oldURL); err != nil {
		return "", err
	}
	return s.Save(fileType, newFile, originalFilename)
}

func (s *S3Uploader) save(file []byte, originalFilename string) (string, error) {
	uniqueName := GenerateUniqueFilename(originalFilename)
	// Ex: s3Client.PutObject(bucket, uniqueName, file)
	url := fmt.Sprintf("https://your-bucket.s3.amazonaws.com/%s", uniqueName)
	return url, nil
}