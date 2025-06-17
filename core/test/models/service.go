package modelT

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/test/handlers"
	"fmt"
)

type Service struct {
	Created   *model.Service
	Company   *Company
	Employees []*Employee
	Branches  []*Branch
}

func (s *Service) Create(status int, x_auth_token string, x_company_id *string) error {
	companyIDStr := s.Company.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("POST").
		URL("/service").
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Company, cID).
		Header(namespace.HeadersKey.Auth, x_auth_token).
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

func (s *Service) Update(status int, changes map[string]any, x_auth_token string, x_company_id *string) error {
	if len(changes) == 0 {
		return fmt.Errorf("no changes provided")
	}
	companyIDStr := s.Company.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("PATCH").
		URL("/service/"+fmt.Sprintf("%v", s.Created.ID.String())).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Company, cID).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Send(changes).
		ParseResponse(&s.Created).
		Error; err != nil {
		return fmt.Errorf("failed to update service: %w", err)
	}
	if status > 200 && status < 300 {
		if err := ValidateUpdateChanges("Service", s.Created, changes); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) GetById(status int, x_auth_token string, x_company_id *string) error {
	companyIDStr := s.Company.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("GET").
		URL("/service/"+fmt.Sprintf("%v", s.Created.ID.String())).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Company, cID).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Send(nil).
		ParseResponse(&s.Created).
		Error; err != nil {
		return fmt.Errorf("failed to get service by ID: %w", err)
	}
	return nil
}

func (s *Service) GetByName(status int, x_auth_token string, x_company_id *string) error {
	companyIDStr := s.Company.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("GET").
		URL("/service/name/"+s.Created.Name).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Company, cID).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Send(nil).
		ParseResponse(&s.Created).
		Error; err != nil {
		return fmt.Errorf("failed to get service by name: %w", err)
	}
	return nil
}

func (s *Service) Delete(status int, x_auth_token string, x_company_id *string) error {
	companyIDStr := s.Company.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("DELETE").
		URL("/service/"+fmt.Sprintf("%v", s.Created.ID.String())).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Company, cID).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}
	return nil
}
