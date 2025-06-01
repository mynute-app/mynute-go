package modelT

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/test/handlers"
	"fmt"
)

type Branch struct {
	Created    model.Branch
	Auth_token string
	Company    *Company
	Services   []*Service
	Employees  []*Employee
}

func (b *Branch) Create(status int) error {
	if err := handler.NewHttpClient().
		Method("POST").
		URL("/branch").
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Company, b.Company.Created.ID.String()).
		Header(namespace.HeadersKey.Auth, b.Auth_token).
		Send(DTO.CreateBranch{
			Name:         lib.GenerateRandomName("Branch Name"),
			CompanyID:    b.Company.Created.ID,
			Street:       lib.GenerateRandomName("Street"),
			Number:       lib.GenerateRandomStrNumber(3),
			Neighborhood: lib.GenerateRandomName("Neighborhood"),
			ZipCode:      lib.GenerateRandomStrNumber(5),
			City:         lib.GenerateRandomName("City"),
			State:        lib.GenerateRandomName("State"),
			Country:      lib.GenerateRandomName("Country"),
		}).
		ParseResponse(&b.Created).
		Error; err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}
	return nil
}

func (b *Branch) Update(status int, changes map[string]any) error {
	if err := handler.NewHttpClient().
		Method("PATCH").
		URL("/branch/"+fmt.Sprintf("%v", b.Created.ID.String())).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Company, b.Company.Created.ID.String()).
		Header(namespace.HeadersKey.Auth, b.Auth_token).
		Send(changes).
		ParseResponse(&b.Created).Error; err != nil {
		return fmt.Errorf("failed to update branch: %w", err)
	}
	return nil
}

func (b *Branch) GetByName(status int) error {
	if err := handler.NewHttpClient().
		Method("GET").
		URL(fmt.Sprintf("/branch/name/%s", b.Created.Name)).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Company, b.Company.Created.ID.String()).
		Header(namespace.HeadersKey.Auth, b.Auth_token).
		Send(nil).
		ParseResponse(&b.Created).Error; err != nil {
		return fmt.Errorf("failed to get branch by name: %w", err)
	}
	return nil
}

func (b *Branch) GetById(status int) error {
	if err := handler.NewHttpClient().
		Method("GET").
		URL(fmt.Sprintf("/branch/%s", b.Created.ID.String())).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Company, b.Company.Created.ID.String()).
		Header(namespace.HeadersKey.Auth, b.Auth_token).
		Send(nil).
		ParseResponse(&b.Created).Error; err != nil {
		return fmt.Errorf("failed to get branch by id: %w", err)
	}
	return nil
}

func (b *Branch) Delete(status int) error {
	if err := handler.NewHttpClient().
		Method("DELETE").
		URL(fmt.Sprintf("/branch/%s", b.Created.ID.String())).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Company, b.Company.Created.ID.String()).
		Header(namespace.HeadersKey.Auth, b.Auth_token).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to delete branch: %w", err)
	}
	return nil
}

func (b *Branch) AddService(status int, service *Service, token *string) error {
	var t string
	if token == nil && b.Auth_token == "" {
		return fmt.Errorf("no authentication token provided")
	} else if token != nil {
		t = *token
	} else {
		t = b.Auth_token
	}
	if err := handler.NewHttpClient().
		Method("POST").
		URL(fmt.Sprintf("/branch/%s/service/%s", b.Created.ID.String(), service.Created.ID.String())).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Auth, t).
		Header(namespace.HeadersKey.Company, b.Company.Created.ID.String()).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to add service to branch: %w", err)
	}
	return nil
}
