package DTO

type CreateBranch struct {
	CompanyID    uint   `json:"company_id" example:"1"`
	Name         string `json:"name" example:"Main Branch"`
	Street       string `gorm:"not null" json:"street" example:"123 Main St"`
	Number       string `gorm:"not null" json:"number" example:"456"`
	Complement   string `json:"complement" example:"Suite 100"`
	Neighborhood string `gorm:"not null" json:"neighborhood" example:"Downtown"`
	ZipCode      string `gorm:"not null" json:"zip_code" example:"10001"`
	City         string `gorm:"not null" json:"city" example:"New York"`
	State        string `gorm:"not null" json:"state" example:"NY"`
	Country      string `gorm:"not null" json:"country" example:"USA"`
}

type UpdateBranch struct {
	CompanyID uint   `json:"company_id" example:"1"`
	Name      string `json:"name" example:"Main Branch Updated"`
	Street    string `gorm:"not null" json:"street" example:"556 Main St"`
}

type Branch struct {
	ID           uint               `json:"id" example:"1"`
	Employees    []UserPopulated    `json:"employees"`
	Services     []ServicePopulated `json:"services"`
	CompanyID    uint               `json:"company_id" example:"1"`
	Name         string             `json:"name" example:"Main Branch"`
	Street       string             `gorm:"not null" json:"street" example:"123 Main St"`
	Number       string             `gorm:"not null" json:"number" example:"456"`
	Complement   string             `json:"complement" example:"Suite 100"`
	Neighborhood string             `gorm:"not null" json:"neighborhood" example:"Downtown"`
	ZipCode      string             `gorm:"not null" json:"zip_code" example:"10001"`
	City         string             `gorm:"not null" json:"city" example:"New York"`
	State        string             `gorm:"not null" json:"state" example:"NY"`
	Country      string             `gorm:"not null" json:"country" example:"USA"`
}

type BranchPopulated struct {
	ID           uint   `json:"id" example:"1"`
	Name         string `json:"name" example:"Main Branch"`
	Street       string `gorm:"not null" json:"street" example:"123 Main St"`
	Number       string `gorm:"not null" json:"number" example:"456"`
	Complement   string `json:"complement" example:"Suite 100"`
	Neighborhood string `gorm:"not null" json:"neighborhood" example:"Downtown"`
	ZipCode      string `gorm:"not null" json:"zip_code" example:"10001"`
	City         string `gorm:"not null" json:"city" example:"New York"`
	State        string `gorm:"not null" json:"state" example:"NY"`
	Country      string `gorm:"not null" json:"country" example:"USA"`
}
