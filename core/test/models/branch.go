package modelT

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/test/handlers"
	"fmt"

	"github.com/google/uuid"
)

type Branch struct {
	Created      *model.Branch
	Company      *Company
	Services     []*Service
	Employees    []*Employee
	Appointments []*Appointment
}

func (b *Branch) Create(status int, x_auth_token string, x_company_id *string) error {
	companyIDStr := b.Company.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("POST").
		URL("/branch").
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Company, cID).
		Header(namespace.HeadersKey.Auth, x_auth_token).
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
			TimeZone:     "America/Sao_Paulo",
		}).
		ParseResponse(&b.Created).
		Error; err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}
	return nil
}

func (b *Branch) Update(status int, changes map[string]any, x_auth_token string, x_company_id *string) error {
	companyIDStr := b.Company.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("PATCH").
		URL("/branch/"+fmt.Sprintf("%v", b.Created.ID.String())).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Company, cID).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Send(changes).
		ParseResponse(&b.Created).Error; err != nil {
		return fmt.Errorf("failed to update branch: %w", err)
	}
	if status > 200 && status < 300 {
		if err := ValidateUpdateChanges("Branch", b.Created, changes); err != nil {
			return err
		}
	}
	return nil
}

func (b *Branch) GetByName(status int, x_auth_token string, x_company_id *string) error {
	companyIDStr := b.Company.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("GET").
		URL(fmt.Sprintf("/branch/name/%s", b.Created.Name)).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Company, cID).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Send(nil).
		ParseResponse(&b.Created).Error; err != nil {
		return fmt.Errorf("failed to get branch by name: %w", err)
	}
	return nil
}

func (b *Branch) GetById(status int, x_auth_token string, x_company_id *string) error {
	companyIDStr := b.Company.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("GET").
		URL(fmt.Sprintf("/branch/%s", b.Created.ID.String())).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Company, cID).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Send(nil).
		ParseResponse(&b.Created).Error; err != nil {
		return fmt.Errorf("failed to get branch by id: %w", err)
	}
	return nil
}

func (b *Branch) Delete(status int, x_auth_token string, x_company_id *string) error {
	companyIDStr := b.Company.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("DELETE").
		URL(fmt.Sprintf("/branch/%s", b.Created.ID.String())).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Company, cID).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to delete branch: %w", err)
	}
	return nil
}

func (b *Branch) AddService(status int, service *Service, x_auth_token string, x_company_id *string) error {
	companyIDStr := b.Company.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("POST").
		URL(fmt.Sprintf("/branch/%s/service/%s", b.Created.ID.String(), service.Created.ID.String())).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Header(namespace.HeadersKey.Company, cID).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to add service to branch: %w", err)
	}
	if err := b.GetById(200, b.Company.Owner.X_Auth_Token, nil); err != nil {
		return fmt.Errorf("failed to get branch by ID after adding service: %w", err)
	}
	if err := service.GetById(200, b.Company.Owner.X_Auth_Token, nil); err != nil {
		return fmt.Errorf("failed to get service by ID after adding to branch: %w", err)
	}
	service.Branches = append(service.Branches, b)
	b.Services = append(b.Services, service)
	return nil
}

func (b *Branch) CreateWorkSchedule(status int, schedule DTO.CreateBranchWorkSchedule, x_auth_token string, x_company_id *string) error {
	if schedule.WorkRanges == nil {
		return fmt.Errorf("work schedule cannot be nil")
	}
	companyIDStr := b.Company.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	var updated *model.Branch
	if err := handler.NewHttpClient().
		Method("POST").
		URL(fmt.Sprintf("/branch/%s/work_schedule", b.Created.ID.String())).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Company, cID).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Send(schedule).
		ParseResponse(&updated).
		Error; err != nil {
		return fmt.Errorf("failed to create branch work schedule: %w", err)
	}
	b.Created.BranchWorkSchedule = updated.BranchWorkSchedule
	return nil
}

func GetExampleBranchWorkSchedule(branchID uuid.UUID, services []*Service) DTO.CreateBranchWorkSchedule {
	var servicesID []DTO.ServiceID
	for _, service := range services {
		servicesID = append(servicesID, DTO.ServiceID{ID: service.Created.ID})
	}
	return DTO.CreateBranchWorkSchedule{
		WorkRanges: []DTO.CreateBranchWorkRange{
			{
				BranchID:  branchID,
				Weekday:   1,
				StartTime: "08:00",
				EndTime:   "20:00",
				TimeZone:  "America/Sao_Paulo",
				Services:  servicesID,
			},
			{
				BranchID:  branchID,
				Weekday:   2,
				StartTime: "08:00",
				EndTime:   "20:00",
				TimeZone:  "America/Sao_Paulo",
				Services:  servicesID,
			},
			{
				BranchID:  branchID,
				Weekday:   3,
				StartTime: "08:00",
				EndTime:   "20:00",
				TimeZone:  "America/Sao_Paulo",
				Services:  servicesID,
			},
			{
				BranchID:  branchID,
				Weekday:   4,
				StartTime: "08:00",
				EndTime:   "20:00",
				TimeZone:  "America/Sao_Paulo",
				Services:  servicesID,
			},
			{
				BranchID:  branchID,
				Weekday:   5,
				StartTime: "08:00",
				EndTime:   "20:00",
				TimeZone:  "America/Sao_Paulo",
				Services:  servicesID,
			},
			{
				BranchID:  branchID,
				Weekday:   6,
				StartTime: "08:00",
				EndTime:   "12:00",
				TimeZone:  "America/Sao_Paulo",
				Services:  servicesID,
			},
			{
				BranchID:  branchID,
				Weekday:   0,
				StartTime: "08:00",
				EndTime:   "12:00",
				TimeZone:  "America/Sao_Paulo",
				Services:  servicesID,
			},
		},
	}
}
