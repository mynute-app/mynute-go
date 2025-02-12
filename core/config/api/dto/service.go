package DTO

type Service struct {
	ID          uint     `json:"id"`
	CompanyID   uint     `gorm:"not null;index;foreignKey:CompanyID;references:ID"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       int32    `json:"price"`
	Duration    int      `json:"duration"`
	Branches    []Branch `json:"branches"`
	Users       []User   `json:"employees"`
}

type ServicePopulated struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int32  `json:"price"`
	Duration    int    `json:"duration"`
}
