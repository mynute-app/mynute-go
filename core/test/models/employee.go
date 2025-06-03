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
	X_Auth_Token string
	Company      *Company
	Created      model.Employee
	Services     []*Service
	Branches     []*Branch
}

func (e *Employee) Create(s int, x_auth_token *string, x_company_id *string) error {
	pswd := "1SecurePswd!"

	t, err := get_token(x_auth_token, &e.X_Auth_Token)
	if err != nil {
		return err
	}
	companyIDStr := e.Company.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("POST").
		URL("/employee").
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Company, cID).
		Header(namespace.HeadersKey.Auth, t).
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

func (e *Employee) Update(s int, changes map[string]any, x_auth_token *string, x_company_id *string) error {
	if len(changes) == 0 {
		return fmt.Errorf("no changes provided")
	}
	t, err := get_token(x_auth_token, &e.X_Auth_Token)
	if err != nil {
		return err
	}
	companyIDStr := e.Company.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("PATCH").
		URL(fmt.Sprintf("/employee/%s", e.Created.ID.String())).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Company, cID).
		Header(namespace.HeadersKey.Auth, t).
		Send(changes).
		ParseResponse(&e.Created).
		Error; err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}
	return nil
}

func (e *Employee) GetById(s int, x_auth_token *string, x_company_id *string) error {
	t, err := get_token(x_auth_token, &e.X_Auth_Token)
	if err != nil {
		return err
	}
	companyIDStr := e.Company.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("GET").
		URL(fmt.Sprintf("/employee/%s", e.Created.ID.String())).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Company, cID).
		Header(namespace.HeadersKey.Auth, t).
		Send(nil).
		ParseResponse(&e.Created).
		Error; err != nil {
		return fmt.Errorf("failed to get employee by ID: %w", err)
	}
	return nil
}

func (e *Employee) GetByEmail(s int, x_auth_token *string, x_company_id *string) error {
	t, err := get_token(x_auth_token, &e.X_Auth_Token)
	if err != nil {
		return err
	}
	companyIDStr := e.Company.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("GET").
		URL(fmt.Sprintf("/employee/email/%s", e.Created.Email)).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Company, cID).
		Header(namespace.HeadersKey.Auth, t).
		Send(nil).
		ParseResponse(&e.Created).
		Error; err != nil {
		return fmt.Errorf("failed to get employee by email: %w", err)
	}
	return nil
}

func (e *Employee) Delete(s int, x_auth_token *string, x_company_id *string) error {
	t, err := get_token(x_auth_token, &e.X_Auth_Token)
	if err != nil {
		return err
	}
	companyIDStr := e.Company.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("DELETE").
		URL(fmt.Sprintf("/employee/%s", e.Created.ID.String())).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Company, cID).
		Header(namespace.HeadersKey.Auth, t).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to delete employee: %w", err)
	}
	return nil
}

func (e *Employee) Login(s int, x_company_id *string) error {
	http := handler.NewHttpClient()
	login := DTO.LoginEmployee{
		Email:    e.Created.Email,
		Password: e.Created.Password,
	}
	companyIDStr := e.Company.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := http.
		Method("POST").
		URL("/employee/login").
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Company, cID).
		Send(login).Error; err != nil {
		return fmt.Errorf("failed to login employee: %w", err)
	}
	auth := http.ResHeaders[namespace.HeadersKey.Auth]
	if len(auth) == 0 {
		return fmt.Errorf("authentication token not found in response headers")
	}
	e.X_Auth_Token = auth[0]
	return nil
}

func (e *Employee) VerifyEmail(s int, x_company_id *string) error {
	companyIDStr := e.Company.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("POST").
		URL(fmt.Sprintf("/employee/verify-email/%s/%s", e.Created.Email, "12345")).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Company, cID).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to verify employee email: %w", err)
	}
	return nil
}

func (e *Employee) CreateBranch(s int) error {
	Branch := &Branch{}
	Branch.Company = e.Company
	if err := Branch.Create(s, e.X_Auth_Token, nil); err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}
	e.Company.Branches = append(e.Company.Branches, Branch)
	return nil
}

func (e *Employee) CreateService(s int) error {
	Service := &Service{}
	Service.Company = e.Company
	if err := Service.Create(s, e.X_Auth_Token, nil); err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}
	e.Company.Services = append(e.Company.Services, Service)
	return nil
}

func (e *Employee) AddBranch(s int, branch *Branch, token *string, x_company_id *string) error {
	t, err := get_token(token, &e.X_Auth_Token)
	if err != nil {
		return err
	}
	companyIDStr := e.Company.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("POST").
		URL(fmt.Sprintf("/employee/%s/branch/%s", e.Created.ID.String(), branch.Created.ID.String())).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, t).
		Header(namespace.HeadersKey.Company, cID).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to add branch to employee: %w", err)
	}
	if err := branch.GetById(s, e.Company.Owner.X_Auth_Token, nil); err != nil {
		return fmt.Errorf("failed to get branch by ID after adding to employee: %w", err)
	}
	if err := e.GetById(s, nil, nil); err != nil {
		return fmt.Errorf("failed to get employee by ID after adding branch: %w", err)
	}
	branch.Employees = append(branch.Employees, e)
	e.Branches = append(e.Branches, branch)
	return nil
}

func (e *Employee) AddService(s int, service *Service, token *string, x_company_id *string) error {
	t, err := get_token(token, &e.X_Auth_Token)
	if err != nil {
		return err
	}
	companyIDStr := e.Company.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("POST").
		URL(fmt.Sprintf("/employee/%s/service/%s", e.Created.ID.String(), service.Created.ID.String())).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, t).
		Header(namespace.HeadersKey.Company, cID).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to add service to employee: %w", err)
	}
	if err := service.GetById(s, e.Company.Owner.X_Auth_Token, nil); err != nil {
		return fmt.Errorf("failed to get service by ID after adding to employee: %w", err)
	}
	if err := e.GetById(s, nil, nil); err != nil {
		return fmt.Errorf("failed to get employee by ID after adding service: %w", err)
	}
	service.Employees = append(service.Employees, e)
	return nil
}

func (e *Employee) AddRole(s int, role *Role, x_auth_token *string, x_company_id *string) error {
	t, err := get_token(x_auth_token, &e.X_Auth_Token)
	if err != nil {
		return err
	}
	companyIDStr := e.Company.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("POST").
		URL(fmt.Sprintf("/employee/%s/role/%s", e.Created.ID.String(), role.Created.ID.String())).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, t).
		Header(namespace.HeadersKey.Company, cID).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to add role to employee: %w", err)
	}
	role.Employees = append(role.Employees, e)
	return nil
}

func get_token(priority *string, secundary *string) (string, error) {
	if priority != nil {
		return *priority, nil
	} else if secundary != nil {
		return *secundary, nil
	}
	return "", fmt.Errorf("no authentication token provided")
}

func get_x_company_id(priority *string, secundary *string) (string, error) {
	if priority != nil {
		return *priority, nil
	} else if secundary != nil {
		return *secundary, nil
	}
	return "", fmt.Errorf("no company ID provided")
}
