package middleware

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	if err := db.Where("method = ? AND path = ?", method, path).Preload("Resources").First(&EndPoint).Error; err != nil || EndPoint.ID == uuid.Nil {
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

	var policies []*model.PolicyRule
	PoliciesWhereClause := "method = ? AND resource_id = ? AND (company_id IS NULL OR company_id = ?)"
	if err := db.Where(PoliciesWhereClause, EndPoint.Method, EndPoint.ResourceID, claim.CompanyID).
		Find(&policies).Error; err != nil {
		return lib.Error.Auth.Unauthorized
	}

	if len(policies) == 0 {
		return lib.Error.Auth.Unauthorized
	}

	subject_data := make(map[string]any)

	if err := db.Table(UserTableName).
		Where("id = ?", claim.ID).
		Take(&subject_data).Error; err != nil {
		return lib.Error.Auth.InvalidToken
	}

	var RequestVal string
	var ResourceReference model.ResourceReference
	for _, ref := range EndPoint.Resource.References {
		switch ref.RequestRef {
		case "query":
			if c.Query(ref.RequestKey) != "" {
				ResourceReference = ref
				RequestVal = c.Query(ref.RequestKey)
				break
			}
		case "body":
			var body map[string]any
			bbytes := c.Request().Body()
			if len(bbytes) == 0 {
				continue
			}
			if err := json.Unmarshal(bbytes, &body); err != nil {
				return err
			}
			if body[ref.RequestKey] != nil && body[ref.RequestKey] != "" {
				ResourceReference = ref
				RequestVal = fmt.Sprintf("%v", body[ref.RequestKey])
				break
			}
		case "header":
			if c.Get(ref.RequestKey) != "" {
				ResourceReference = ref
				RequestVal = c.Get(ref.RequestKey)
				break
			}
		case "path":
			if c.Params(ref.RequestKey) != "" {
				ResourceReference = ref
				RequestVal = c.Params(ref.RequestKey)
				break
			}
		default:
			return fmt.Errorf("invalid request reference type: %s. Endpoint.Resource.ID: %d", ref.RequestRef, EndPoint.Resource.ID)
		}
	}

	if RequestVal == "" {
		return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("Request is malformed. Endpoint.Resource.ID: %d", EndPoint.Resource.ID))
	}

	resource_data := make(map[string]any)

	if err := db.
		Table(EndPoint.Resource.Table).
		Where(ResourceReference.DatabaseKey+" = ?", RequestVal).
		Take(&resource_data).Error; err != nil {
		return lib.Error.Auth.Unauthorized.WithError(err)
	}

	for _, policy := range policies {
		if ok, err := am.PolicyEngine.CanAccess(subject_data, resource_data, policy); err != nil {
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
