package modelT

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/test/handlers"
	"fmt"
)

type Service struct {
	Created    DTO.Service
	Auth_token string
	Company    *Company
	Employees  []*Employee
	Branches   []*Branch
}

func (s *Service) Create(status int) error {
	if err := handler.NewHttpClient().
		Method("POST").
		URL("/service").
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Company, s.Company.Created.ID.String()).
		Header(namespace.HeadersKey.Auth, s.Auth_token).
		Send(DTO.CreateService{
			Name:        lib.GenerateRandomName("Service"),
			Description: lib.GenerateRandomName("Description"),
			CompanyID:   s.Company.Created.ID,
			Price:       int32(lib.GenerateRandomInt(3)),
			Duration:    60,
		}).
		ParseResponse(&s.Created).
		Error; err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}
	return nil
}

func (s *Service) Update(status int, changes map[string]any) error {
	if len(changes) == 0 {
		return fmt.Errorf("no changes provided")
	}
	if err := handler.NewHttpClient().
		Method("PATCH").
		URL("/service/"+fmt.Sprintf("%v", s.Created.ID.String())).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Company, s.Company.Created.ID.String()).
		Header(namespace.HeadersKey.Auth, s.Auth_token).
		Send(changes).
		ParseResponse(&s.Created).
		Error; err != nil {
		return fmt.Errorf("failed to update service: %w", err)
	}
	return nil
}

func (s *Service) GetById(status int, token *string) error {
	var t string
	if token == nil && s.Auth_token == "" {
		return fmt.Errorf("no authentication token provided")
	} else if token != nil {
		t = *token
	} else {
		t = s.Auth_token
	}
	if err := handler.NewHttpClient().
		Method("GET").
		URL("/service/"+fmt.Sprintf("%v", s.Created.ID.String())).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Company, s.Company.Created.ID.String()).
		Header(namespace.HeadersKey.Auth, t).
		Send(nil).
		ParseResponse(&s.Created).
		Error; err != nil {
		return fmt.Errorf("failed to get service by ID: %w", err)
	}
	return nil
}

func (s *Service) GetByName(status int) error {
	if err := handler.NewHttpClient().
		Method("GET").
		URL("/service/name/"+s.Created.Name).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Company, s.Company.Created.ID.String()).
		Header(namespace.HeadersKey.Auth, s.Auth_token).
		Send(nil).
		ParseResponse(&s.Created).
		Error; err != nil {
		return fmt.Errorf("failed to get service by name: %w", err)
	}
	return nil
}

func (s *Service) Delete(status int) error {
	if err := handler.NewHttpClient().
		Method("DELETE").
		URL("/service/"+fmt.Sprintf("%v", s.Created.ID.String())).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Company, s.Company.Created.ID.String()).
		Header(namespace.HeadersKey.Auth, s.Auth_token).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}
	return nil
}
