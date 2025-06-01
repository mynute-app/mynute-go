package modelT

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/test/handlers"
	"fmt"
)

type Employee struct {
	Auth_token string
	Company    *Company
	Created    model.Employee
	Services   []*Service
	Branches   []*Branch
}

func (e *Employee) Create(s int) error {
	pswd := "1SecurePswd!"

	if err := handler.NewHttpClient().
		Method("POST").
		URL("/employee").
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Company, e.Company.Created.ID.String()).
		Header(namespace.HeadersKey.Auth, e.Company.Auth_token).
		Send(DTO.CreateEmployee{
			CompanyID: e.Company.Created.ID,
			Name:      lib.GenerateRandomName("Employee Name"),
			Surname:   lib.GenerateRandomName("Employee Surname"),
			Email:     lib.GenerateRandomEmail("employee"),
			Phone:     lib.GenerateRandomPhoneNumber(),
			Password:  pswd,
		}).
		ParseResponse(&e.Created).
		Error; err != nil {
		return fmt.Errorf("failed to create employee: %w", err)
	}
	e.Created.Password = pswd
	return nil
}

func (e *Employee) Update(s int, changes map[string]any) error {
	if len(changes) == 0 {
		return fmt.Errorf("no changes provided")
	}
	if err := handler.NewHttpClient().
		Method("PATCH").
		URL(fmt.Sprintf("/employee/%s", e.Created.ID.String())).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Company, e.Company.Created.ID.String()).
		Header(namespace.HeadersKey.Auth, e.Company.Auth_token).
		Send(changes).
		ParseResponse(&e.Created).
		Error; err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}
	return nil
}

func (e *Employee) GetById(s int) error {
	if err := handler.NewHttpClient().
		Method("GET").
		URL(fmt.Sprintf("/employee/%s", e.Created.ID.String())).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Company, e.Company.Created.ID.String()).
		Header(namespace.HeadersKey.Auth, e.Company.Auth_token).
		Send(nil).
		ParseResponse(&e.Created).
		Error; err != nil {
		return fmt.Errorf("failed to get employee by ID: %w", err)
	}
	return nil
}

func (e *Employee) GetByEmail(s int) error {
	if err := handler.NewHttpClient().
		Method("GET").
		URL(fmt.Sprintf("/employee/email/%s", e.Created.Email)).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Company, e.Company.Created.ID.String()).
		Header(namespace.HeadersKey.Auth, e.Company.Auth_token).
		Send(nil).
		ParseResponse(&e.Created).
		Error; err != nil {
		return fmt.Errorf("failed to get employee by email: %w", err)
	}
	return nil
}

func (e *Employee) Delete(s int) error {
	if err := handler.NewHttpClient().
		Method("DELETE").
		URL(fmt.Sprintf("/employee/%s", e.Created.ID.String())).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Company, e.Company.Created.ID.String()).
		Header(namespace.HeadersKey.Auth, e.Company.Auth_token).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to delete employee: %w", err)
	}
	return nil
}

func (e *Employee) Login(s int) error {
	http := handler.NewHttpClient()
	login := DTO.LoginEmployee{
		Email:    e.Created.Email,
		Password: e.Created.Password,
	}
	if err := http.
		Method("POST").
		URL("/employee/login").
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Company, e.Created.CompanyID.String()).
		Send(login).Error; err != nil {
		return fmt.Errorf("failed to login employee: %w", err)
	}
	auth := http.ResHeaders[namespace.HeadersKey.Auth]
	if len(auth) == 0 {
		return fmt.Errorf("authentication token not found in response headers")
	}
	e.Auth_token = auth[0]
	return nil
}

func (e *Employee) VerifyEmail(s int) error {
	if err := handler.NewHttpClient().
		Method("POST").
		URL(fmt.Sprintf("/employee/verify-email/%s/%s", e.Created.Email, "12345")).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Company, e.Created.CompanyID.String()).
		Header(namespace.HeadersKey.Auth, e.Company.Auth_token).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to verify employee email: %w", err)
	}
	return nil
}

func (e *Employee) CreateBranch(s int) error {
	Branch := &Branch{}
	Branch.Auth_token = e.Auth_token
	Branch.Company = e.Company
	if err := Branch.Create(s); err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}
	e.Company.Branches = append(e.Company.Branches, Branch)
	return nil
}

func (e *Employee) CreateService(s int) error {
	Service := &Service{}
	Service.Auth_token = e.Auth_token
	Service.Company = e.Company
	if err := Service.Create(s); err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}
	e.Company.Services = append(e.Company.Services, Service)
	return nil
}

func (e *Employee) AddBranch(s int, branch *Branch, token *string) error {
	t, err := get_token(&e.Auth_token, token)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("POST").
		URL(fmt.Sprintf("/employee/%s/branch/%s", e.Created.ID.String(), branch.Created.ID.String())).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, t).
		Header(namespace.HeadersKey.Company, e.Company.Created.ID.String()).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to add branch to employee: %w", err)
	}
	if err := branch.GetById(s); err != nil {
		return fmt.Errorf("failed to get branch by ID after adding to employee: %w", err)
	}
	branch.Employees = append(branch.Employees, e)
	return nil
}

func (e *Employee) AddService(s int, service *Service, token *string) error {
	t, err := get_token(&e.Auth_token, token)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("POST").
		URL(fmt.Sprintf("/employee/%s/service/%s", e.Created.ID.String(), service.Created.ID.String())).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, t).
		Header(namespace.HeadersKey.Company, e.Company.Created.ID.String()).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to add service to employee: %w", err)
	}
	service.Employees = append(service.Employees, e)
	return nil
}

func (e *Employee) AddRole(s int, role *Role, token *string) error {
	t, err := get_token(&e.Auth_token, token)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("POST").
		URL(fmt.Sprintf("/employee/%s/role/%s", e.Created.ID.String(), role.Created.ID.String())).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, t).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to add role to employee: %w", err)
	}
	role.Employees = append(role.Employees, e)
	return nil
}

func get_token(token1 *string, token2 *string) (string, error) {
	if token1 != nil {
		return *token1, nil
	} else if token2 != nil {
		return *token2, nil
	}
	return "", fmt.Errorf("no authentication token provided")
}
