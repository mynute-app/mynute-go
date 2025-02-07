package DTO

type Branch struct {
	ID        uint               `json:"id"`
	CompanyID uint               `gorm:"not null"`
	Name      string             `json:"name"`
	Employees []UserPopulated    `json:"employees"`
	Services  []ServicePopulated `json:"services"`
}

type BranchPopulated struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
