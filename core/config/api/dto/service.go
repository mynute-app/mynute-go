package DTO

type Service struct {
	ID          uint     `json:"id" example:"1"`
	CompanyID   uint     `json:"company_id" gorm:"not null;index;foreignKey:CompanyID;references:ID;constraint:OnDelete:CASCADE;" example:"1"`
	Name        string   `json:"name" example:"Premium Consultation"`
	Description string   `json:"description" example:"A 60-minute in-depth business consultation"`
	Price       int32    `json:"price" example:"150"`
	Duration    int      `json:"duration" example:"60"`
	Branches    []Branch `json:"branches"`
	Users       []User   `json:"employees"`
}

type ServicePopulated struct {
	ID          uint   `json:"id" example:"1"`
	Name        string `json:"name" example:"Premium Consultation"`
	Description string `json:"description" example:"A 60-minute in-depth business consultation"`
	Price       int32  `json:"price" example:"150"`
	Duration    int    `json:"duration" example:"60"`
}