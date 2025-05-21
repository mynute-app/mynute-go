package dJSON

type Design struct {
	Colors Colors `json:"colors"`
	Images Images `json:"images"`
}

type Colors struct {
	Primary   string `json:"primary" example:"#FF5733"`
	Secondary string `json:"secondary" example:"#33FF57"`
	Tertiary  string `json:"tertiary" example:"#3357FF"`
	Quaternary string `json:"quaternary" example:"#FF33A1"`
}

type Images struct {
	LogoURL       string `json:"logo_url" example:"https://example.com/logo.png"`
	BannerURL     string `json:"banner_url" example:"https://example.com/banner.png"`
	BackgroundURL string `json:"background_url" example:"https://example.com/background.png"`
	FaviconURL    string `json:"favicon_url" example:"https://example.com/favicon.png"`
}