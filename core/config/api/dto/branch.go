package DTO

type Branch struct {
	ID           uint               `json:"id"`
	CompanyID    uint               `gorm:"not null" json:"company_id"`
	Name         string             `json:"name"`
	Employees    []UserPopulated    `json:"employees"`
	Services     []ServicePopulated `json:"services"`
	Street       string             `gorm:"not null" json:"street"`
	Number       string             `gorm:"not null" json:"number"`
	Complement   string             `json:"complement"`
	Neighborhood string             `gorm:"not null" json:"neighborhood"`
	ZipCode      string             `gorm:"not null" json:"zip_code"`
	City         string             `gorm:"not null" json:"city"`
	State        string             `gorm:"not null" json:"state"`
	Country      string             `gorm:"not null" json:"country"`
}

type BranchPopulated struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Street       string `gorm:"not null" json:"street"`
	Number       string `gorm:"not null" json:"number"`
	Complement   string `json:"complement"`
	Neighborhood string `gorm:"not null" json:"neighborhood"`
	ZipCode      string `gorm:"not null" json:"zip_code"`
	City         string `gorm:"not null" json:"city"`
	State        string `gorm:"not null" json:"state"`
	Country      string `gorm:"not null" json:"country"`
}
