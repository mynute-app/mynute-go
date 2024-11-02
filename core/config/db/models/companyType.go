package models

type CompanyType struct {
	GeneralResourceInfo
	Name        string `gorm:"not null;unique" json:"name"`
	Description string `json:"description"`
}
