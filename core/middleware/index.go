package middleware

import "github.com/gofiber/fiber/v2"

// Central Registry for middleware actions
type Registry struct {
	actions map[string]map[string][]func(*fiber.Ctx) (int, error)
}

// Initialize a new registry
func NewRegistry() *Registry {
	return &Registry{actions: make(map[string]map[string][]func(*fiber.Ctx) (int, error))}
}

// Register middleware actions by resource and method(s)
func (mr *Registry) RegisterAction(resource string, methods interface{}, action func(*fiber.Ctx) (int, error)) {
	if mr.actions[resource] == nil {
		mr.actions[resource] = make(map[string][]func(*fiber.Ctx) (int, error))
	}

	switch m := methods.(type) {
	case string:
		mr.actions[resource][m] = append(mr.actions[resource][m], action)
	case []string:
		for _, method := range m {
			mr.actions[resource][method] = append(mr.actions[resource][method], action)
		}
	default:
		panic("methods must be either a string or a slice of strings")
	}
}

// Retrieve middleware actions for a specific resource and method
func (mr *Registry) GetActions(resource, method string) []func(*fiber.Ctx) (int, error) {
	return mr.actions[resource][method]
}
