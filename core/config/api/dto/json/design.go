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
	Logo       Image `json:"logo"`
	Banner     Image `json:"banner"`
	Background Image `json:"background"`
	Favicon    Image `json:"favicon"`
}

type Image struct {
	Alt 	 string `json:"alt" example:"Image of something"`
	Title    string `json:"title" example:"Title of this image"`
	Caption  string `json:"caption" example:"This is the image we talk about"`
	URL       string `json:"url" example:"https://example.com/image.png"`
}