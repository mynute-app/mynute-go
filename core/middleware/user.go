package middleware

import (
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type user_middleware struct {
	Gorm *handler.Gorm
}

func User(Gorm *handler.Gorm) *user_middleware {
	return &user_middleware{Gorm: Gorm}
}

func (um *user_middleware) Create() []fiber.Handler {
	return []fiber.Handler{
		lib.SaveBodyOnCtx[model.User],
		um.VerifyEmailExists,
		um.ValidateProps,
		um.HashPassword,
	}
}

func (um *user_middleware) VerifyEmailExists(c *fiber.Ctx) error {
	body, err := lib.GetBodyFromCtx[*model.User](c)
	if err != nil {
		return err
	}
	users := &[]model.User{}
	if err := um.Gorm.DB.Where("email = ?", body.Email).Find(users).Error; err != nil {
		return err
	}
	if len(*users) > 0 {
		return lib.MyErrors.EmailExists.SendToClient(c)
	}
	fmt.Printf("Email %v is unique\n", body.Email)
	return c.Next()
}

func (um *user_middleware) ValidateProps(c *fiber.Ctx) error {
	body, err := lib.GetBodyFromCtx[*model.User](c)
	if err != nil {
		return err
	}
	if err := lib.ValidateName(body.Name, "user"); err != nil {
		return lib.MyErrors.InvalidUserName.SendToClient(c)
	}
	if valid := lib.ValidateEmail(body.Email); !valid {
		return lib.MyErrors.InvalidEmail.SendToClient(c)
	}
	return c.Next()
}

func (um *user_middleware) HashPassword(c *fiber.Ctx) error {
	body, err := lib.GetBodyFromCtx[*model.User](c)
	if err != nil {
		return err
	}
	hashed, err := handler.HashPassword(body.Password)
	if err != nil {
		return err
	}
	fmt.Printf("Hashed password: %s\n", hashed)
	body.Password = hashed
	return c.Next()
}

// type UserMiddlewareActions struct {
// 	Gorm *handler.Gorm
// }

// func User(Gorm *handler.Gorm) *Registry {
// 	// user := &UserMiddlewareActions{Gorm: Gorm}
// 	registry := NewRegistry()
// 	// user := &UserMiddlewareActions{Gorm: Gorm}
// 	// registry.RegisterAction(namespace.UserKey.Name, "POST", user.Create)

// 	return registry
// }

// func (em *UserMiddlewareActions) GetUserByEmail(c *fiber.Ctx) error {
// 	user := &model.User{}
// 	if err := em.Gorm.DB.Where("email = ?", c.Params("email")).First(user).Error; err != nil {
// 		return err
// 	}
// 	c.Locals(namespace.UserKey.Model, user)
// 	return c.Next()
// }

// func (em *UserMiddlewareActions) Create(c *fiber.Ctx) (int, error) {
// 	user, err := lib.GetFromCtx[*model.User](c, namespace.GeneralKey.Model)
// 	if err != nil {
// 		return 500, err
// 	}
// 	// Perform validation
// 	if err := lib.ValidateName(user.Name, "user"); err != nil {
// 		return 400, err
// 	}
// 	// Proceed to the next middleware or handler
// 	return 0, nil
// }
