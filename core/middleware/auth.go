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

	var EndPoint model.EndPoint
	if err := db.Where("method = ? AND path = ?", method, path).First(&EndPoint).Error; err != nil || EndPoint.ID == 0 {
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

	var policies []*model.PolicyRule
	PoliciesWhereClause := "id = ? AND (company_id IS NULL OR company_id = ?)"
	if err := db.Where(PoliciesWhereClause, EndPoint.ID, claim.CompanyID).
		Find(&policies).Error; err != nil {
		return lib.Error.Auth.Unauthorized
	}

	if len(policies) == 0 {
		return lib.Error.Auth.Unauthorized
	}

	getResource := func(table, key, valueAt string) (map[string]any, error) {
		var value any
		switch valueAt {
		case "query":
			value = c.Query(key)
		case "header":
			value = c.Get(key)
		case "path":
			value = c.Params(key)
		case "body":
			var body map[string]any
			if err := c.BodyParser(&body); err != nil {
				return nil, err
			}
			if val, ok := body[key]; ok {
				value = val
			} else {
				return nil, lib.Error.Auth.Unauthorized
			}
		default:
			panic("Invalid Resource param at")
		}

		ResourceWhereClause := fmt.Sprintf("%s = ?", key)

		resource := make(map[string]any)
		if err := db.Table(table).
			Where(ResourceWhereClause, value).
			Take(&resource).Error; err != nil {
			return nil, lib.Error.Auth.Unauthorized
		}

		if len(resource) == 0 {
			return nil, lib.Error.Auth.Unauthorized
		}

		return resource, nil
	}

	for _, policy := range policies {
		table := policy.ResourceTable
		key := policy.ResourceKey
		valueAt := policy.ResourceValueAt
		resource, err := getResource(table, key, valueAt)
		if err != nil {
			return err
		}
		if ok, err := am.PolicyEngine.CanAccess(subject, resource, policy); err != nil {
			return err
		} else if !ok {
			return lib.Error.Auth.Unauthorized
		}
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
