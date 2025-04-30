package model

import (
	"github.com/google/uuid"
)

type Company struct {
	BaseModel
	Name       string     `gorm:"not null;unique" json:"name"`
	TaxID      string     `gorm:"not null;unique" json:"tax_id"`
	SchemaName string     `gorm:"type:varchar(100);not null;uniqueIndex" json:"schema_name"`
	SectorID   *uuid.UUID `json:"sector_id"`
	Sector     *Sector    `gorm:"foreignKey:SectorID;constraint:OnDelete:SET NULL;"`
}
