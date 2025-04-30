package model

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Company struct {
	BaseModel
	Name       string     `gorm:"not null;unique" json:"name"`
	TaxID      string     `gorm:"not null;unique" json:"tax_id"`
	SchemaName string     `gorm:"type:varchar(100);not null;uniqueIndex" json:"schema_name"`
	SectorID   *uuid.UUID `json:"sector_id"`
	Sector     *Sector    `gorm:"foreignKey:SectorID;constraint:OnDelete:SET NULL;"`
}

func (c *Company) GenerateSchemaName() string {
	return c.Name + "_" + c.ID.String()
}

func (c *Company) AfterCreate(tx *gorm.DB) error {
	// Generate the schema name after the company is created
	c.SchemaName = c.GenerateSchemaName()
	if err := tx.Save(c).Error; err != nil {
		return err
	}

	// Create the schema in the database
	if err := tx.Exec("CREATE SCHEMA IF NOT EXISTS " + c.SchemaName).Error; err != nil {
		return err
	}
	
	// Set the search path to the new schema
	if err := tx.Exec("SET search_path TO " + c.SchemaName).Error; err != nil {
		return err
	}

	// Migrate tenant specific tables
	for _, model := range TenantModels {
		if err := tx.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate model %T: %w", model, err)
		}
	}

	return nil
}