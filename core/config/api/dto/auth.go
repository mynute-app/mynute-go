package DTO

type Claims struct {
	ID        uint     `json:"id" example:"1"`
	Name      string   `json:"name" example:"John"`
	Surname   string   `json:"surname" example:"Doe"`
	Role      string   `json:"role" example:"client"`
	Email     string   `json:"email" example:"john.doe@example.com"`
	Phone     string   `json:"phone" example:"+15555555555"`
	Tags      []string `json:"tags" example:"[\"tag1\", \"tag2\"]"`
	Verified  bool     `json:"verified" example:"true"`
	CompanyID uint     `json:"company_id" example:"1"`
}
