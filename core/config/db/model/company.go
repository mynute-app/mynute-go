package model

import "gorm.io/gorm"

// Company represents a company entity

type Company struct {
	gorm.Model
	Name      string     `gorm:"not null;unique" json:"name"`
	TaxID     string     `gorm:"not null;unique" json:"tax_id"`
	SectorID  *uint      `gorm:"index" json:"sector_id"`
	Sector    *Sector    `gorm:"foreignKey:SectorID;constraint:OnDelete:SET NULL;"`
	Employees []Employee `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"employees"` // One-to-many relation with Client
	Branches  []Branch   `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"branches"`  // One-to-many relation with Branch
	Services  []Service  `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"services"`  // One-to-many relation with Service
}
