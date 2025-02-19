package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Central Registry for middleware actions
type Registry struct {
	actions map[string]map[string][]func(*fiber.Ctx) (int, error)
}

type MiddlewareActions struct {
	methods interface{}
	action func(*fiber.Ctx) (int, error)
}

// Initialize a new registry
func NewRegistry() *Registry {
	return &Registry{actions: make(map[string]map[string][]func(*fiber.Ctx) (int, error))}
}

func (mr *Registry) RegisterActions(resource string, actions []MiddlewareActions) {
	for _, act := range actions {
		mr.RegisterAction(resource, act.methods, act.action)
	}
}

// Register middleware actions by resource and method(s)
func (mr *Registry) RegisterAction(resource string, methods interface{}, action func(*fiber.Ctx) (int, error)) {
	if mr.actions[resource] == nil {
		mr.actions[resource] = make(map[string][]func(*fiber.Ctx) (int, error))
	}

	allMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"} // List of all HTTP methods

	registerAllMethods := func() {
		for _, method := range allMethods {
			mr.actions[resource][method] = append(mr.actions[resource][method], action)
		}
	}

	// Handle nil or "ALL" (case-insensitive) to register for all methods
	if methods == nil {
		registerAllMethods()
	} else if s, ok := methods.(string); ok && strings.ToUpper(s) == "ALL" {
		registerAllMethods()
	} else {
		switch m := methods.(type) {
		case string:
			mr.actions[resource][m] = append(mr.actions[resource][m], action)
		case []string:
			if len(m) == 0 {
				registerAllMethods()
			} else {
				for _, method := range m {
					mr.actions[resource][method] = append(mr.actions[resource][method], action)
				}
			}
		default:
			panic("methods must be either nil, a string, or a slice of strings")
		}
	}
}

// Retrieve middleware actions for a specific resource and method
func (mr *Registry) GetActions(resource, method string) []func(*fiber.Ctx) (int, error) {
	return mr.actions[resource][method]
}
