package lib

import (
	"agenda-kaki-go/core/config/namespace"
	"fmt"

	"github.com/gofiber/fiber/v2"
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
	changePath := fmt.Sprintf(`SET search_path TO "%s"`, schemaName)
	if err := db.Exec(changePath).Error; err != nil {
		return fmt.Errorf("failed to set search path to %s: %w", schemaName, err)
	}
	return nil
}

/*
 * Gets the database session from the fiber context.
 Recomended when you need to perform a single database operation.
 * @return *gorm.DB - The database session
 * @return error - The error if any
*/
// @param c *fiber.Ctx - The fiber context
func Session(c *fiber.Ctx) (*gorm.DB, error) {
	tx, ok := c.Locals(namespace.GeneralKey.DatabaseSession).(*gorm.DB)
	if !ok {
		return nil, Error.General.SessionNotFound
	}
	return tx, nil
}