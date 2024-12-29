package models

// Company holds an array of CompanyTypes.
type Company struct {
	GeneralResourceInfo
	CompanyID    uint          `gorm:"primaryKey" json:"company_id"` // Primary key
	Name         string        `gorm:"not null;unique" json:"name"`
	TaxID        string        `gorm:"not null;unique" json:"tax_id"`
	CompanyTypes []CompanyType `gorm:"many2many:company_company_types;constraint:OnDelete:CASCADE" json:"company_types"`
	Employees    []User        `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"employees"` // one-to-many relation with User
	Branches     []Branch      `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"branches"`  // one-to-many relation with Branch
	Services     []Service     `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"services"`  // one-to-many relation with Service
}
