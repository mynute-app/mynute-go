package model

import (
	"agenda-kaki-go/core/lib"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Company represents a company entity
type CompanyBase struct {
	Name  string `gorm:"not null;unique" json:"name"`
	TaxID string `gorm:"not null;unique" json:"tax_id"`
}

type CompanyMeta struct {
	BaseModel
	CompanyBase
	SchemaName string `gorm:"type:varchar(100);not null;uniqueIndex" json:"schema_name"`
}

func (CompanyMeta) TableName() string { return "general.companies" }

type Company struct {
	CompanyBase
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;<-:create" json:"id"`
	Employees []Employee `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"employees"` // One-to-many relation with Client
	Branches  []Branch   `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"branches"`  // One-to-many relation with Branch
	Services  []Service  `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"services"`  // One-to-many relation with Service
	SectorID  *uuid.UUID `json:"sector_id"`
	Sector    *Sector    `gorm:"foreignKey:SectorID;constraint:OnDelete:SET NULL;"`
}

func (Company) TableName() string { return "company" }

func (c *Company) BeforeUpdate(tx *gorm.DB) (err error) {
	if tx.Statement.Changed("ID") {
		return lib.Error.Company.IdUpdateForbidden
	}
	return nil
}
