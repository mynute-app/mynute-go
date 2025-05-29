package mJSON

import (
	"agenda-kaki-go/core/lib/Uploader"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type DesignConfig struct {
	Colors Colors `json:"colors"`
	Images Images `json:"images"`
}

type Colors struct {
	Primary    string `json:"primary"`
	Secondary  string `json:"secondary"`
	Tertiary   string `json:"tertiary"`
	Quaternary string `json:"quaternary"`
}

type Images struct {
	Logo       Image `json:"logo"`
	Banner     Image `json:"banner"`
	Background Image `json:"background"`
	Favicon    Image `json:"favicon"`
}

type Image struct {
	Alt     string `json:"alt"`
	Title   string `json:"title"`
	Caption string `json:"caption"`
	URL     string `json:"url"`
}

func (i *Images) GetImageURL(imageType string) string {
	switch imageType {
	case "logo":
		return i.Logo.URL
	case "banner":
		return i.Banner.URL
	case "background":
		return i.Background.URL
	case "favicon":
		return i.Favicon.URL
	default:
		return ""
	}
}

func (d *DesignConfig) SaveImage(caller_entity, caller_id, oldURL, originalFilename string, newFile []byte) (string, error) {
	up, err := myUploader.FileUploader(caller_entity, caller_id)
	if err != nil {
		return "", err
	}
	return up.Replace("image", oldURL, newFile, originalFilename)
}

func (d *DesignConfig) DeleteImage(caller_entity, caller_id, oldURL string) error {
	up, err := myUploader.FileUploader(caller_entity, caller_id)
	if err != nil {
		return err
	}
	return up.Delete(oldURL)
}

func (d *DesignConfig) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *DesignConfig) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan DesignConfig: expected []byte")
	}
	return json.Unmarshal(bytes, d)
}
