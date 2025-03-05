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
}

func Employee(Gorm *handler.Gorm) *employee_middleware {
	return &employee_middleware{Gorm: Gorm}
}

func (em *employee_middleware) SaveEmployeeCreateBody(c *fiber.Ctx) error {
	body := &DTO.CreateEmployee{}
	err := lib.BodyParser(c.Body(), body)
	if err != nil {
		return err
	}
	c.Locals(namespace.RequestKey.Body_Parsed, body)
	return c.Next()
}

func (em *employee_middleware) FindUserWhenCreatingEmployee(c *fiber.Ctx) error {
	body, err := lib.GetFromCtx[*DTO.CreateEmployee](c, namespace.RequestKey.Body_Parsed)
	if err != nil {
		return err
	}
	user := &model.User{}
	if err := em.Gorm.DB.Where("email = ?", body.Email).First(user).Error; err != nil {
		return err
	}
	if user.ID != 0 {
		body.UserID = user.ID
		c.Locals(namespace.RequestKey.Body_Parsed, body)
		return c.Next()
	}
	return c.Next()
}

func (em *employee_middleware) SetEmployeeUserAccount(c *fiber.Ctx) error {
	body, err := lib.GetFromCtx[*DTO.CreateEmployee](c, namespace.RequestKey.Body_Parsed)
	if err != nil {
		return err
	}
	if body.UserID != 0 {
		return c.Next()
	}
	user := &model.User{}
	lib.ParseToDTO(body, user)
	if err := em.Gorm.DB.Create(user).Error; err != nil {
		return err
	}
	body.UserID = user.ID
	return c.Next()
}
