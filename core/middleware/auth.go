package middleware

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	if err := db.Where("method = ? AND path = ?", method, path).Preload("Resource").First(&EndPoint).Error; err != nil || EndPoint.ID == uuid.Nil {
		return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("endpoint not found: %s %s", method, path))
	}

	auth_claims := c.Locals(namespace.RequestKey.Auth_Claims)
	claim, ok := auth_claims.(*DTO.Claims)
	if !ok || claim.ID == uuid.Nil || claim.ID.Variant() != uuid.RFC4122 || !claim.Verified {
		return lib.Error.Auth.InvalidToken
	}

	var policies []*model.PolicyRule
	PoliciesWhereClause := "end_point_id = ? AND (company_id IS NULL OR company_id = ?)"
	if err := db.
		Where(PoliciesWhereClause, EndPoint.ID, claim.CompanyID).
		Find(&policies).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("no policies found for endpoint: %s %s and company: %s", method, path, claim.CompanyID.String()))
		}
		return lib.Error.Auth.Unauthorized.WithError(err)
	}

	if len(policies) == 0 {
		return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("no policies found for endpoint: %s %s and company: %s", method, path, claim.CompanyID.String()))
	}

	var user any
	if claim.CompanyID == uuid.Nil {
		user = &model.Client{}
	} else {
		user = &model.Employee{}
	}

	subject_data := make(map[string]any)

	if err := db.
		Model(user).                  // Tell GORM which model (and thus table)
		Preload(clause.Associations). // Preload ALL defined top-level associations
		Where("id = ?", claim.ID).
		Take(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.Auth.Unauthorized
		}
		return lib.Error.General.AuthError.WithError(err)
	}

	jsonData, err := json.Marshal(user) // Convert struct to JSON bytes
	if err != nil {
		// Handle marshaling error (should be rare for valid structs)
		return lib.Error.General.AuthError.WithError(fmt.Errorf("failed to marshal user subject: %w", err))
	}
	err = json.Unmarshal(jsonData, &subject_data) // Convert JSON bytes to map
	if err != nil {
		// Handle unmarshaling error (should be rare)
		return lib.Error.General.AuthError.WithError(fmt.Errorf("failed to unmarshal user subject to map: %w", err))
	}

	var RequestVal string
	var ResourceReference model.ResourceReference
forLoop:
	for _, ref := range EndPoint.Resource.References {
		switch ref.RequestRef {
		case "query":
			if c.Query(ref.RequestKey) != "" {
				ResourceReference = ref
				RequestVal = c.Query(ref.RequestKey)
				break forLoop
			}
		case "body":
			var body map[string]any
			bbytes := c.Request().Body()
			if len(bbytes) == 0 {
				continue forLoop
			}
			if err := json.Unmarshal(bbytes, &body); err != nil {
				return err
			}
			if body[ref.RequestKey] != nil && body[ref.RequestKey] != "" {
				ResourceReference = ref
				RequestVal = fmt.Sprintf("%v", body[ref.RequestKey])
				break forLoop
			}
		case "header":
			if c.Get(ref.RequestKey) != "" {
				ResourceReference = ref
				RequestVal = c.Get(ref.RequestKey)
				break forLoop
			}
		case "path":
			if c.Params(ref.RequestKey) != "" {
				ResourceReference = ref
				RequestVal = c.Params(ref.RequestKey)
				break forLoop
			}
		default:
			return fmt.Errorf("invalid request reference type: %s. Endpoint.Resource.ID: %s", ref.RequestRef, EndPoint.Resource.ID.String())
		}
	}

	if RequestVal == "" {
		return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("request is malformed. Endpoint.Resource.ID: %s", EndPoint.Resource.ID.String()))
	}

	RequestVal, err = url.QueryUnescape(RequestVal)

	if err != nil {
		return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("failed to unescape request value: %w", err))
	}

	var resource any

	switch EndPoint.Resource.Table {
	case "appointments":
		resource = &model.Appointment{}
	case "branches":
		resource = &model.Branch{}
	case "clients":
		resource = &model.Client{}
	case "companies":
		resource = &model.Company{}
	case "employees":
		resource = &model.Employee{}
	case "holidays":
		resource = &model.Holiday{}
	case "policies":
		resource = &model.PolicyRule{}
	case "roles":
		resource = &model.Role{}
	case "sectors":
		resource = &model.Sector{}
	case "services":
		resource = &model.Service{}
	default:
		return lib.Error.General.AuthError.WithError(fmt.Errorf("invalid resource table: %s", EndPoint.Resource.Table))
	}

	if err := db.
		Model(resource). // Tell GORM which model (and thus table)
		Where(ResourceReference.DatabaseKey+" = ?", RequestVal).
		Preload(clause.Associations).
		Take(resource).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("resource not found: [%s] = %s at table `%s`", ResourceReference.DatabaseKey, RequestVal, EndPoint.Resource.Table))
		}
		return lib.Error.General.AuthError.WithError(err)
	}

	resource_data := make(map[string]any)

	jsonData, err = json.Marshal(user) // Convert struct to JSON bytes
	if err != nil {
		// Handle marshaling error (should be rare for valid structs)
		return lib.Error.General.AuthError.WithError(fmt.Errorf("failed to marshal user subject: %w", err))
	}
	err = json.Unmarshal(jsonData, &resource_data) // Convert JSON bytes to map
	if err != nil {
		// Handle unmarshaling error (should be rare)
		return lib.Error.General.AuthError.WithError(fmt.Errorf("failed to unmarshal user subject to map: %w", err))
	}

	for _, policy := range policies {
		decision := am.PolicyEngine.CanAccess(subject_data, resource_data, policy)
		if decision.Error != nil {
			return lib.Error.General.AuthError.WithError(err)
		} else if !ok {
			denied := lib.Error.Auth.Unauthorized.WithError(errors.New("policy engine denied access"))
			denied.WithError(fmt.Errorf("reason: %s", decision.Reason))
			denied.WithError(fmt.Errorf("endpoint: %s %s", method, path))
			denied.WithError(fmt.Errorf("policy.name: %s", policy.Name))
			return denied
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
