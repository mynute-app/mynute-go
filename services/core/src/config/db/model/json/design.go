package mJSON

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"mynute-go/services/core/src/lib"
	myUploader "mynute-go/services/core/src/lib/cloud_uploader"
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
	Profile    Image `json:"profile"`
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

func (is *Images) Save(img_type string, caller_entity string, caller_id string, file *multipart.FileHeader) (string, error) {
	switch img_type {
	case "profile":
		return is.Profile.Save(caller_entity, caller_id, file)
	case "logo":
		return is.Logo.Save(caller_entity, caller_id, file)
	case "banner":
		return is.Banner.Save(caller_entity, caller_id, file)
	case "background":
		return is.Background.Save(caller_entity, caller_id, file)
	case "favicon":
		return is.Favicon.Save(caller_entity, caller_id, file)
	default:
		return "", lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid image type: %s", img_type))
	}
}

func (is *Images) Delete(img_type string, caller_entity, caller_id string) error {
	up, err := myUploader.FileUploader(caller_entity, caller_id)
	if err != nil {
		return err
	}
	img_types_url := map[string]*Image{
		"profile":    &is.Profile,
		"logo":       &is.Logo,
		"banner":     &is.Banner,
		"favicon":    &is.Favicon,
		"background": &is.Background,
	}
	target, ok := img_types_url[img_type]
	if !ok {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid image type: %s", img_type))
	}
	if err := up.Delete(target.URL); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	target.URL = "" // Clear the URL after deletion
	target.Alt = ""
	target.Title = ""
	target.Caption = ""
	return nil
}

func (i *Image) Save(caller_entity, caller_id string, file *multipart.FileHeader) (string, error) {
	if file == nil {
		return "", lib.Error.General.InternalError.WithError(fmt.Errorf("file header passed is nil trying to save image for %s with ID %s", caller_entity, caller_id))
	}

	f, err := file.Open()
	if err != nil {
		return "", lib.Error.General.InternalError.WithError(fmt.Errorf("failed to open file %s for caller %s with ID %s: %w", file.Filename, caller_entity, caller_id, err))
	}
	defer f.Close()

	newFile := make([]byte, file.Size)
	_, err = f.Read(newFile)
	if err != nil {
		return "", lib.Error.General.InternalError.WithError(fmt.Errorf("failed to read file %s for caller %s with ID %s: %w", file.Filename, caller_entity, caller_id, err))
	}

	up, err := myUploader.FileUploader(caller_entity, caller_id)
	if err != nil {
		return "", lib.Error.General.InternalError.WithError(fmt.Errorf("failed to get file uploader for caller %s with ID %s: %w", caller_entity, caller_id, err))
	}

	newUrl, err := up.Replace("image", i.URL, newFile, file.Filename)
	if err != nil {
		return "", lib.Error.General.InternalError.WithError(fmt.Errorf("failed to replace image %s for caller %s with ID %s: %w", i.URL, caller_entity, caller_id, err))
	}
	i.URL = newUrl
	i.Alt = file.Filename
	i.Title = file.Filename
	i.Caption = file.Filename
	return i.URL, nil
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

func (c Colors) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *Colors) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan Colors: expected []byte")
	}
	return json.Unmarshal(bytes, c)
}

func (i Image) Value() (driver.Value, error) {
	return json.Marshal(i)
}

func (i *Image) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan Image: expected []byte")
	}
	return json.Unmarshal(bytes, i)
}

func (is Images) Value() (driver.Value, error) {
	return json.Marshal(is)
}

func (is *Images) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan Images: expected []byte")
	}
	return json.Unmarshal(bytes, is)
}
