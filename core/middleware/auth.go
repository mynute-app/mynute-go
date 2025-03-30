package middleware

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type auth_middleware struct {
	Gorm *handler.Gorm
}

func Auth(Gorm *handler.Gorm) *auth_middleware {
	return &auth_middleware{Gorm: Gorm}
}

func (am *auth_middleware) DenyUnauthorized(c *fiber.Ctx) error {
	auth_claims := c.Locals(namespace.RequestKey.Auth_Claims)
	claim, ok := auth_claims.(*DTO.Claims)
	if !ok || claim.ID == 0 || !claim.Verified {
		return lib.Error.Auth.InvalidToken.SendToClient(c)
	}

	db := am.Gorm.DB
	userID := claim.ID
	companyID := claim.CompanyID
	method := c.Method()
	path := c.Path()

	// 1. Verificar RBAC
	var matchedResource model.Resource
	if err := db.Where("method = ? AND path = ?", method, path).First(&matchedResource).Error; err != nil {
		return lib.Error.Auth.Unauthorized
	}

	var count int64
	err := db.Raw(`
		SELECT COUNT(*)
		FROM roles r
		JOIN employee_roles er ON er.role_id = r.id
		JOIN role_routes rr ON rr.role_id = r.id
		WHERE er.user_id = ? AND er.company_id = ? AND rr.route_id = ?
	`, userID, companyID, matchedResource.ID).Scan(&count).Error

	if err != nil || count == 0 {
		return lib.Error.Auth.Unauthorized
	}

	// 2. Verificar ABAC
	var rules []model.PolicyRule
	if err := db.Where("method = ? AND path = ? AND company_id = ?", method, path, companyID).
		Find(&rules).Error; err != nil {
		return lib.Error.Auth.Unauthorized
	}

	abac := handler.Policy(rules)

	if abac.CanAccess() {
		return c.Next()
	}

	return lib.Error.Auth.Unauthorized
}

func matchCondition(actual string, op string, expected string) bool {
	switch op {
	case "equal":
		return actual == expected
	case "contains":
		return strings.Contains(actual, expected)
	default:
		return false
	}
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
