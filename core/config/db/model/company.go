package model

import "gorm.io/gorm"

// Company represents a company entity

type Company struct {
	gorm.Model
	Name      string     `gorm:"not null;unique" json:"name"`
	TaxID     string     `gorm:"not null;unique" json:"tax_id"`
	SectorID  *uint      `gorm:"index" json:"sector_id"`
	Sector    *Sector    `gorm:"foreignKey:SectorID;constraint:OnDelete:SET NULL;"`
	Employees []Employee `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"employees"` // One-to-many relation with User
	Branches  []Branch   `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"branches"`  // One-to-many relation with Branch
	Services  []Service  `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"services"`  // One-to-many relation with Service
}

type EntityPermissions struct {
	EntityName string   `json:"entity_name"`
	Create     []string `json:"create"`
	ReadAny    []string `json:"read_any"`
	ReadOwn    []string `json:"read_own"`
	UpdateAny  []string `json:"update_any"`
	UpdateOwn  []string `json:"update_own"`
	DeleteAny  []string `json:"delete_any"`
	DeleteOwn  []string `json:"delete_own"`
}

type ItemPermissions struct {
	Create    []string `json:"create"`
	ReadAny   []string `json:"read_any"`
	ReadOwn   []string `json:"read_own"`
	UpdateAny []string `json:"update_any"`
	UpdateOwn []string `json:"update_own"`
	DeleteAny []string `json:"delete_any"`
	DeleteOwn []string `json:"delete_own"`
}
