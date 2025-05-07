package middleware

import (
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type database struct {
	Gorm *gorm.DB
}

func DatabaseFactory(db *gorm.DB) *database {
	return &database{Gorm: db}
}

/*
	* SaveSession is a middleware to save the database session in the fiber context.
	It should be used in the main function to set the session for the current request.
*/
// @return func(c *fiber.Ctx) error - The middleware function
func (db *database) SavePublicSession(c *fiber.Ctx) error {
	c.Locals(namespace.GeneralKey.DatabaseSession, db.Gorm)
	return nil
}

/*
	* SaveTenantSession is a middleware to save the tenant database session in the fiber context.
	It should be used in the main function to set the session for the current request.
*/
// @return func(c *fiber.Ctx) error - The middleware function
func (db *database) SaveTenantSession(c *fiber.Ctx) error {
	companyID := c.Get(namespace.HeadersKey.Company)
	if companyID == "" {
		return lib.Error.Auth.CompanyHeaderMissing
	}
	tx, err := MakeSession(db.Gorm, c)
	if err != nil {
		return err
	}
	var Company model.Company
	if err := tx.Where("id = ?", companyID).First(&Company).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.Company.NotFound
		}
		return lib.Error.General.AuthError.WithError(err)
	}
	if err := lib.ChangeToTenantSchema(tx, Company.SchemaName); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	c.Locals(namespace.GeneralKey.DatabaseSession, tx)
	c.Locals(namespace.GeneralKey.Company, &Company)
	return nil
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
