package middleware

import (
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type tenant_middleware struct {
	Gorm *gorm.DB
}

func Tenant(db *gorm.DB) *tenant_middleware {
	return &tenant_middleware{
		Gorm: db,
	}
}

func (tm *tenant_middleware) Validate(c *fiber.Ctx) error {
	companyID := c.Get(namespace.HeadersKey.Company)
	if companyID == "" {
		return lib.Error.Auth.CompanyHeaderMissing
	}
	tx := tm.Gorm.Session(&gorm.Session{NewDB: true, Context: c.Context()})
	if tx.Error != nil {
		return tx.Error
	}
	var company model.Company
	if CompanyUUID, err := uuid.Parse(companyID); err != nil {
		return lib.Error.Auth.CompanyHeaderInvalid
	} else {
		company.ID = CompanyUUID
	}
	if err := company.Refresh(tx); err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.Company.NotFound
		}
		return lib.Error.General.AuthError.WithError(err)
	}
	c.Locals(namespace.HeadersKey.Company, companyID)
	c.Locals(namespace.GeneralKey.DatabaseSession, tx)
	c.Locals(namespace.GeneralKey.Company, &company)
	return c.Next()
}