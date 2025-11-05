package model

// AuthDBModels are authentication/authorization models for the auth database
var AuthDBModels = []interface{}{
	&EndPoint{},
	&PolicyRule{},
	&Resource{},
	&Property{},
	&User{}, // Unified user model for auth
}
