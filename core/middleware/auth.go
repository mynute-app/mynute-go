package middleware

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"fmt"
	"strconv"

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
	userID := fmt.Sprintf("%d", claim.ID)
	companyID := fmt.Sprintf("%d", claim.CompanyID)

	// RBAC Check
	var roles []model.Role
	db.Raw(`SELECT r.* FROM roles r 
		JOIN user_roles ur ON ur.role_id = r.id 
		WHERE ur.user_id = ? AND ur.company_id = ?`, userID, companyID).Scan(&roles)

	for _, role := range roles {
		var perms []model.RolePermission
		db.Where("role_id = ? AND method = ?", role.ID, method).Find(&perms)
		for _, perm := range perms {
			if lib.MatchPath(perm.Path, path) {
				return c.Next()
			}
		}
	}

	// ABAC Check
	sub := handler.PolicySubject{
		ID: userID,
		Attrs: map[string]string{
			"user_id":    userID,
			"role":       claim.Role,
			"company_id": companyID,
		},
	}

	res := handler.PolicyResource{
		Attrs: map[string]string{
			"company_id": strconv.Itoa(int(claim.CompanyID)),
			"branch_id":  c.Params("branch_id"),
			"employee_id": c.Params("employee_id"),
		},
	}

	env := handler.PolicyEnvironment{}

	var rules []model.PolicyRule
	db.Where("company_id = ?", companyID).Find(&rules)
	engine := handler.Policy(rules)

	if engine.CanAccess(sub, method, path, res, env) {
		return c.Next()
	}

	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"error": "Access denied",
	})
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
