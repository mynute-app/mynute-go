package mJSON

import (
	uploader "agenda-kaki-go/core/lib/Uploader"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type DesignConfig struct {
	Colors    Colors `json:"colors"`
	Images    Images `json:"images"`
	Fonts     Fonts  `json:"fonts"`
	DarkMode  bool   `json:"dark_mode"`
	CustomCSS string `json:"custom_css"`
}

type Fonts []Font
type Font struct {
	Family string `json:"family"`
	Size   string `json:"size"`
	Weight string `json:"weight"`
	Color  string `json:"color"`
}

type Colors struct {
	Primary    string `json:"primary"`
	Secondary  string `json:"secondary"`
	Tertiary   string `json:"tertiary"`
	Quaternary string `json:"quaternary"`
}

type Images struct {
	LogoURL       string `json:"logo_url"`
	BannerURL     string `json:"banner_url"`
	BackgroundURL string `json:"background_url"`
	FaviconURL    string `json:"favicon_url"`
}

func SaveDesignImage(oldURL string, newFile []byte, originalFilename string) (string, error) {
	up, err := uploader.FileUploader()
	if err != nil {
		return "", err
	}
	return up.Replace("image", oldURL, newFile, originalFilename)
}


func (d DesignConfig) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *DesignConfig) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan DesignConfig: expected []byte")
	}
	return json.Unmarshal(bytes, d)
}

func (i *Images) GetLogoURL() string {
	if i.LogoURL == "" {
		return "/assets/images/standard_logo.png"
	}
	return i.LogoURL
}

func (i *Images) GetBannerURL() string {
	if i.BannerURL == "" {
		return "/assets/images/standard_banner.png"
	}
	return i.BannerURL
}

func (i *Images) GetBackgroundURL() string {
	if i.BackgroundURL == "" {
		return "/assets/images/standard_background.png"
	}
	return i.BackgroundURL
}

func (i *Images) GetFaviconURL() string {
	if i.FaviconURL == "" {
		return "/assets/images/standard_favicon.png"
	}
	return i.FaviconURL
}

func (i *Images) GetURLs() Images {
	return Images{
		LogoURL:       i.GetLogoURL(),
		BannerURL:     i.GetBannerURL(),
		BackgroundURL: i.GetBackgroundURL(),
		FaviconURL:    i.GetFaviconURL(),
	}
}
