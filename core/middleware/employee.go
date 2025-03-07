package middleware

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"

	"github.com/gofiber/fiber/v2"
)

type employee_middleware struct {
	Gorm *handler.Gorm
	Auth *auth_middleware
}

func Employee(Gorm *handler.Gorm) *employee_middleware {
	return &employee_middleware{Gorm: Gorm, Auth: Auth(Gorm)}
}

func (em *employee_middleware) CreateEmployee() []fiber.Handler {
	return []fiber.Handler{
		em.Auth.WhoAreYou,
		em.Auth.DenyUnauthorized,
		lib.SaveBodyOnCtx[DTO.CreateEmployee],
		em.FindUser,
		em.LinkEmployeeWithUser,
	}
}

func (em *employee_middleware) FindUser(c *fiber.Ctx) error {
	body, err := lib.GetBodyFromCtx[*DTO.CreateEmployee](c)
	if err != nil {
		return err
	}
	user := &model.User{}
	if err := em.Gorm.DB.Where("email = ?", body.Email).First(user).Error; err != nil {
		return err
	}
	c.Locals(namespace.UserKey.Model, user)
	return c.Next()
}

func (em *employee_middleware) LinkEmployeeWithUser(c *fiber.Ctx) error {
	body, err := lib.GetBodyFromCtx[*DTO.CreateEmployee](c)
	if err != nil {
		return err
	}
	user, err := lib.GetFromCtx[*model.User](c, namespace.UserKey.Model)
	if err != nil {
		return err
	}
	if user.ID != 0 {
		body.UserID = user.ID
		return c.Next()
	}
	lib.ParseToDTO(body, user)
	if err := em.Gorm.DB.Create(user).Error; err != nil {
		return err
	}
	body.UserID = user.ID
	return c.Next()
}
