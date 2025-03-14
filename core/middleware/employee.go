package middleware

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"

	"github.com/gofiber/fiber/v2"
)

type employee_middleware struct {
	Gorm *handler.Gorm
	Auth *auth_middleware
	User *user_middleware
}

func Employee(Gorm *handler.Gorm) *employee_middleware {
	return &employee_middleware{Gorm: Gorm, Auth: Auth(Gorm), User: User(Gorm)}
}

func (em *employee_middleware) Create() []fiber.Handler {
	return []fiber.Handler{
		em.Auth.WhoAreYou,
		em.Auth.DenyUnauthorized,
		lib.SaveBodyOnCtx[DTO.CreateEmployee],
		em.User.MatchUserAndCompany,
	}
}
