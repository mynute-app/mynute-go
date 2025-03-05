package middleware

import "agenda-kaki-go/core/handler"

type EmployeeMiddlewareActions struct {
	Gorm *handler.Gorm
}

func Employee(Gorm *handler.Gorm) *Registry {
	registry := NewRegistry()
	// employee := &EmployeeMiddlewareActions{Gorm: Gorm}
	return registry
}