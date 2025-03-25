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

	method := c.Method()
	path := c.Path()
	db := am.Gorm.DB
	userID := claim.ID
	tenantID := claim.CompanyID

	// RBAC Check
	var roles []model.Role
	db.Raw(`SELECT r.* FROM roles r 
		JOIN user_roles ur ON ur.role_id = r.id 
		WHERE ur.user_id = ? AND ur.tenant_id = ?`, userID, tenantID).Scan(&roles)

	for _, role := range roles {
		var perms []model.RolePermission
		db.Where("role_id = ? AND method = ?", role.ID, method).Find(&perms)
		for _, perm := range perms {
			if matchPath(perm.Path, path) {
				return c.Next()
			}
		}
	}

	// ABAC Check
	sub := handler.PolicySubject{
		ID: userID,
		Attrs: map[string]string{
			"user_id":    fmt.Sprintf("%d", userID),
			"role":       claim.Role,
			"company_id": fmt.Sprintf("%d", claim.CompanyID),
		},
	}

	res := handler.PolicyResource{
		Attrs: map[string]string{
			"branch_id": c.Query("branch_id"),
		},
	}

	env := handler.PolicyEnvironment{}

	var rules []model.PolicyRule
	db.Where("company_id = ?", tenantID).Find(&rules)
	engine := handler.Policy(rules)

	if engine.CanAccess(sub, method, path, res, env) {
		return c.Next()
	}

	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"error": "Access denied",
	})
}

func matchPath(rulePath, actualPath string) bool {
	ruleSegments := strings.Split(rulePath, "/")
	actualSegments := strings.Split(actualPath, "/")

	if len(ruleSegments) != len(actualSegments) {
		return false
	}

	for i := range ruleSegments {
		if strings.HasPrefix(ruleSegments[i], ":") {
			continue
		}
		if ruleSegments[i] != actualSegments[i] {
			return false
		}
	}
	return true
}

// func (am *auth_middleware) DenyUnauthorized(c *fiber.Ctx) error {
// 	auth_claims := c.Locals(namespace.RequestKey.Auth_Claims)
// 	claim, ok := auth_claims.(*DTO.Claims)
// 	if !ok {
// 		return lib.Error.Auth.InvalidToken.SendToClient(c)
// 	}
// 	if claim.ID == 0 {
// 		return lib.Error.Auth.InvalidToken.SendToClient(c)
// 	}
// 	if !claim.Verified {
// 		return lib.Error.Client.NotVerified.SendToClient(c)
// 	}
// 	return c.Next()
// }

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
