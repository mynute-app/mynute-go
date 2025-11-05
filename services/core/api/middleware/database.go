package middleware

import (
	"fmt"
	"mynute-go/services/core/api/lib"
	"mynute-go/services/core/config/db/model"
	"mynute-go/services/core/config/namespace"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

/*
	* Middleware to save the database session in the fiber context.
	It should be used in the main function to set the session for the current request.
*/
// @return fiber.Handler - The middleware function
func SavePublicSession(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tx, err := MakeSession(db, c)
		if err != nil {
			return err
		}
		c.Locals(namespace.GeneralKey.DatabaseSession, tx)
		return c.Next()
	}
}

/*
	* Middleware to save the tenant (company) database session in the fiber context.
	It should be used in the main function to set the session for the current request.
*/
// @return fiber.Handler - The middleware function
func SaveCompanySession(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		companyID := c.Get(namespace.HeadersKey.Company)

		if companyID == "" {
			return lib.Error.Auth.CompanyHeaderMissing
		}

		tx, err := MakeSession(db, c)
		if err != nil {
			return err
		}

		if err := lib.ChangeToPublicSchema(tx); err != nil {
			return lib.Error.General.InternalError.WithError(err)
		}

		var SchemaName string

		if err := tx.Model(&model.Company{}).
			Where("id = ?", companyID).
			Pluck("schema_name", &SchemaName).
			Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return lib.Error.Company.NotFound
			}
			return lib.Error.General.AuthError.WithError(err)
		}

		if SchemaName == "" {
			// Check if the company exists
			var qttyOfCompaniesFound int64
			if err := tx.Model(&model.Company{}).
				Where("id = ?", companyID).
				Count(&qttyOfCompaniesFound).Error; err != nil {
				return lib.Error.Company.NotFound.WithError(fmt.Errorf("error while trying to find company with ID %s: %w", companyID, err))
			}
			if qttyOfCompaniesFound == 0 {
				return lib.Error.Company.NotFound.WithError(fmt.Errorf("company with ID %s not found", companyID))
			}
			return lib.Error.General.InternalError.WithError(fmt.Errorf("company with ID %s does not have a schema name", companyID))
		}

		if err := lib.ChangeToCompanySchema(tx, SchemaName); err != nil {
			return lib.Error.General.InternalError.WithError(err)
		}
		c.Locals(namespace.GeneralKey.CompanySchema, SchemaName)
		c.Locals(namespace.GeneralKey.DatabaseSession, tx)
		return c.Next()
	}
}

/*
	* MakeSession is a middleware to create a new database session for the current request.
	It should be used in the main function to set the session for the current request.
*/
// @return func(c *fiber.Ctx) error - The middleware function
func MakeSession(db *gorm.DB, c *fiber.Ctx) (*gorm.DB, error) {
	tx := db.Session(&gorm.Session{NewDB: true, Context: c.Context()})
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func ChangeToPublicSchema(c *fiber.Ctx) error {
	if err := lib.ChangeToPublicSchemaByContext(c); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return c.Next()
}

func ChangeToCompanySchema(c *fiber.Ctx) error {
	if err := lib.ChangeToCompanySchemaByContext(c); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return c.Next()
}
