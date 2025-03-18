package middleware

import (
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type company_middleware struct {
	Gorm *handler.Gorm
}

func Company(Gorm *handler.Gorm) *company_middleware {
	return &company_middleware{Gorm: Gorm}
}

func (cm *company_middleware) CreateCompany() []fiber.Handler {
	auth := Auth(cm.Gorm)
	return []fiber.Handler{
		auth.WhoAreYou,
		auth.DenyClaimless,
		lib.SaveBodyOnCtx[model.Company],
		cm.DenyCnpj,
		cm.ValidateProps,
	}
}

func (cm *company_middleware) GetCompany(c *fiber.Ctx) error {
	var company model.Company
	res := &lib.SendResponse{Ctx: c}
	companyID := c.Params(namespace.QueryKey.CompanyId)

	if companyID == "" {
		return errors.New("missing companyId")
	}

	if err := cm.Gorm.GetOneBy("id", companyID, &company, nil); err != nil {
		return res.Http400(err)
	}

	c.Locals(namespace.CompanyKey.Model, &company)
	return c.Next()
}

func (cm *company_middleware) ValidateProps(c *fiber.Ctx) error {
	company, err := lib.GetBodyFromCtx[*model.Company](c)
	if err != nil {
		return err
	}
	res := &lib.SendResponse{Ctx: c}
	err = lib.ValidateName(company.Name, "company")
	if err != nil {
		return res.Http400(err)
	} else if !lib.ValidateTaxID(company.TaxID) {
		return res.Http400(errors.New("invalid tax ID"))
	}
	return c.Next()
}

func (cm *company_middleware) DenyCnpj(c *fiber.Ctx) error {
	company, err := lib.GetBodyFromCtx[*model.Company](c)
	if err != nil {
		return err
	}
	companies := []model.Company{}
	cm.Gorm.DB.Where("tax_id = ?", company.TaxID).Find(&companies)
	if len(companies) > 0 {
		return lib.Error.Company.CnpjAlreadyExists.SendToClient(c)
	}
	return c.Next()
}