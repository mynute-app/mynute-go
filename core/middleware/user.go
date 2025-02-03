package middleware

import (
	"agenda-kaki-go/core/handlers"
)

type EmployeeMiddlewareActions struct {
	Gorm *handlers.Gorm
}

func User(Gorm *handlers.Gorm) *Registry {
	// user := &EmployeeMiddlewareActions{Gorm: Gorm}
	registry := NewRegistry()

	// registry.RegisterAction(namespace.UserKey.Name, "POST", user.Create)

	return registry
}

// func (em *EmployeeMiddlewareActions) Create(c fiber.Ctx) (int, error) {
// 	user, err := lib.GetFromCtx[*models.User](c, namespace.GeneralKey.Model)
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
