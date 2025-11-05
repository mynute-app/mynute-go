package lib

import (
	"fmt"
	"mynute-go/services/core/src/config/namespace"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ChangeToPublicSchema(db *gorm.DB) error {
	// Change the current schema to public
	if err := db.Exec("SET search_path TO public").Error; err != nil {
		return Error.General.InternalError.WithError(err)
	}
	var currentSchemaName string
	if err := db.Raw("SELECT current_schema()").Scan(&currentSchemaName).Error; err != nil {
		return Error.General.InternalError.WithError(err)
	}
	if currentSchemaName != "public" {
		return Error.General.InternalError.WithError(fmt.Errorf("failed to change schema to public"))
	}
	return nil
}

func ChangeToCompanySchema(db *gorm.DB, schemaName string) error {
	if schemaName == "" {
		return Error.General.SessionNotFound.WithError(fmt.Errorf("schema name is empty"))
	}
	// Change the current schema to the tenant schema
	changePath := fmt.Sprintf(`SET search_path TO "%s"`, schemaName)
	if err := db.Exec(changePath).Error; err != nil {
		return Error.General.InternalError.WithError(err)
	}
	var currentSchemaName string
	if err := db.Raw("SELECT current_schema()").Scan(&currentSchemaName).Error; err != nil {
		return Error.General.InternalError.WithError(err)
	}
	if currentSchemaName != schemaName {
		return Error.General.InternalError.WithError(fmt.Errorf("failed to change schema to %s", schemaName))
	}
	return nil
}

func ChangeToPublicSchemaByContext(c *fiber.Ctx) error {
	tx, err := Session(c)
	if err != nil {
		return err
	}
	return ChangeToPublicSchema(tx)
}

func ChangeToCompanySchemaByContext(c *fiber.Ctx) error {
	tx, err := Session(c)
	if err != nil {
		return err
	}
	schemaName, err := GetCompanySchemaName(c)
	if err != nil {
		return err
	}
	if schemaName == "" {
		return Error.General.InternalError.WithError(fmt.Errorf("company schema name is empty"))
	}
	return ChangeToCompanySchema(tx, schemaName)
}

/*
* Gets the database session from the fiber context.
Recomended when you need to perform a single database operation.
* @return *gorm.DB - The database session
* @return error - The error if any
* @param c *fiber.Ctx - The fiber context
*/
func Session(c *fiber.Ctx) (*gorm.DB, error) {
	tx, ok := c.Locals(namespace.GeneralKey.DatabaseSession).(*gorm.DB)
	if !ok {
		return nil, Error.General.SessionNotFound
	}
	return tx, nil
}

func GetCompanySchemaName(c *fiber.Ctx) (string, error) {
	schemaName, ok := c.Locals(namespace.GeneralKey.CompanySchema).(string)
	if !ok {
		return "", Error.General.SessionNotFound.WithError(fmt.Errorf("company schema name not found in context"))
	}
	return schemaName, nil
}
