package lib

import (
	"fmt"

	"gorm.io/gorm"
)

func ChangeToPublicSchema(db *gorm.DB) error {
	// Change the current schema to public
	if err := db.Exec("SET search_path TO public").Error; err != nil {
		return fmt.Errorf("failed to set search path to public: %w", err)
	}
	return nil
}

func ChangeToTenantSchema(db *gorm.DB, schemaName string) error {
	if schemaName == "" {
		return fmt.Errorf("schema name is empty")
	}
	// Change the current schema to the tenant schema
	if err := db.Exec("SET search_path TO " + schemaName).Error; err != nil {
		return fmt.Errorf("failed to set search path to %s: %w", schemaName, err)
	}
	return nil
}