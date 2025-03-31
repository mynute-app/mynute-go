package middleware

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type auth_middleware struct {
	Gorm         *handler.Gorm
	PolicyEngine *handler.Policy
}

func Auth(Gorm *handler.Gorm) *auth_middleware {
	return &auth_middleware{Gorm: Gorm, PolicyEngine: handler.NewPolicyEngine(Gorm.DB)}
}

func (am *auth_middleware) DenyUnauthorized(c *fiber.Ctx) error {
	db := am.Gorm.DB
	method := c.Method()
	path := c.Route().Path

	var Resource model.Resource
	if err := db.Where("method = ? AND path = ?", method, path).First(&Resource).Error; err != nil || Resource.ID == 0 {
		return lib.Error.Auth.Unauthorized
	}

	auth_claims := c.Locals(namespace.RequestKey.Auth_Claims)
	claim, ok := auth_claims.(*DTO.Claims)
	if !ok || claim.ID == 0 || !claim.Verified {
		return lib.Error.Auth.InvalidToken
	}

	var UserTableName string
	if claim.CompanyID == 0 {
		UserTableName = "clients"
	} else {
		UserTableName = "employees"
	}

	subject := make(map[string]any)

	if err := db.Table(UserTableName).
		Where("id = ?", claim.ID).
		Take(&subject).Error; err != nil {
		return lib.Error.Auth.InvalidToken
	}

	RegistryTable := Resource.RefFromTable
	RegistryParamKey := Resource.RefFromKey
	RegistryParamAt := Resource.RefKeyValueAt
	var RegistryParamVal any
	switch RegistryParamAt {
	case "query":
		RegistryParamVal = c.Query(RegistryParamKey)
	case "header":
		RegistryParamVal = c.Get(RegistryParamKey)
	case "path":
		RegistryParamVal = c.Params(RegistryParamKey)
	default:
		panic("Invalid registry param at")
	}

	WhereClause := fmt.Sprintf("%s = ?", RegistryParamKey)

	registry := make(map[string]any)
	if err := db.Table(RegistryTable).
		Where(WhereClause, RegistryParamVal). // Need to replace with actual resource ID from the path
		Take(&registry).Error; err != nil {
		return lib.Error.Auth.Unauthorized
	}

	var policies []*model.PolicyRule
	if err := db.Where("resource_id = ? AND (company_id IS NULL OR company_id = ?)", Resource.ID, claim.CompanyID).
		Find(&policies).Error; err != nil {
		return lib.Error.Auth.Unauthorized
	}

	if ok, err := am.PolicyEngine.CanAccess(subject, registry, policies); err != nil {
		return err
	} else if !ok {
		return lib.Error.Auth.Unauthorized
	}

	return c.Next()
}

func (am *auth_middleware) WhoAreYou(c *fiber.Ctx) error {
	authorization := c.Get("Authorization")
	if authorization == "" {
		return c.Next()
	}
	user, err := handler.JWT(c).WhoAreYou()
	if err != nil {
		return err
	} else if user == nil {
		return c.Next()
	}
	c.Locals(namespace.RequestKey.Auth_Claims, user)
	return c.Next()
}
