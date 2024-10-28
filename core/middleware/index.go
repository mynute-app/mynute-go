package middleware

import "github.com/gofiber/fiber/v3"

// type IMiddleware interface {
// 	GET() []func(fiber.Ctx) (int, error)
// 	POST() []func(fiber.Ctx) (int, error)
// 	PATCH() []func(fiber.Ctx) (int, error)
// 	DELETE() []func(fiber.Ctx) (int, error)
// 	ForceGET() []func(fiber.Ctx) (int, error)
// 	ForceDELETE() []func(fiber.Ctx) (int, error)
// }

// Central Registry for middleware actions
type Registry struct {
	actions map[string]map[string][]func(fiber.Ctx) (int, error)
}

// Initialize a new registry
func NewMiddlewareRegistry() *Registry {
	return &Registry{actions: make(map[string]map[string][]func(fiber.Ctx) (int, error))}
}

// Register middleware actions by resource and method
func (mr *Registry) RegisterAction(resource, method string, action func(fiber.Ctx) (int, error)) {
	if mr.actions[resource] == nil {
		mr.actions[resource] = make(map[string][]func(fiber.Ctx) (int, error))
	}
	mr.actions[resource][method] = append(mr.actions[resource][method], action)
}

// Retrieve middleware actions for a specific resource and method
func (mr *Registry) GetActions(resource, method string) []func(fiber.Ctx) (int, error) {
	return mr.actions[resource][method]
}
