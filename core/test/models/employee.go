package modelT

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	mJSON "agenda-kaki-go/core/config/db/model/json"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/test/handlers"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Employee struct {
	Created      *model.Employee
	Company      *Company
	Services     []*Service
	Branches     []*Branch
	Appointments []*Appointment
	X_Auth_Token string
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
	if err := ValidateUpdateChanges("Employee", e.Created, changes); err != nil {
		return err
	}
	return nil
}

func (e *Employee) UpdateWorkSchedule(s int, workSchedule []mJSON.WorkSchedule, x_auth_token *string, x_company_id *string) error {
	if workSchedule == nil {
		return fmt.Errorf("work schedule cannot be nil")
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
	var emp *model.Employee
	if err := handler.NewHttpClient().
		Method("PATCH").
		URL(fmt.Sprintf("/employee/%s", e.Created.ID.String())).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Company, cID).
		Header(namespace.HeadersKey.Auth, t).
		Send(map[string]any{"work_schedule": workSchedule}).
		ParseResponse(&emp).
		Error; err != nil {
		return fmt.Errorf("failed to update employee work schedule: %w", err)
	}

	if err := ValidateWorkSchedule(emp.WorkSchedule, e, e.Company); err != nil {
		return fmt.Errorf("invalid work schedule: %w", err)
	}

	e.Created.WorkSchedule = emp.WorkSchedule

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
	if err := e.GetById(200, nil, nil); err != nil {
		return fmt.Errorf("failed to get employee by ID after login: %w", err)
	}
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

func (e *Employee) AddBranch(s int, b *Branch, token *string, x_company_id *string) error {
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
		URL(fmt.Sprintf("/employee/%s/branch/%s", e.Created.ID.String(), b.Created.ID.String())).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, t).
		Header(namespace.HeadersKey.Company, cID).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to add branch to employee: %w", err)
	}
	if err := b.GetById(s, e.Company.Owner.X_Auth_Token, nil); err != nil {
		return fmt.Errorf("failed to get branch by ID after adding to employee: %w", err)
	}
	if err := e.GetById(s, nil, nil); err != nil {
		return fmt.Errorf("failed to get employee by ID after adding branch: %w", err)
	}
	b.Employees = append(b.Employees, e)
	e.Branches = append(e.Branches, b)
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
	e.Services = append(e.Services, service)
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

func ValidateWorkSchedule(ws mJSON.WorkSchedule, employee *Employee, company *Company) error {
	preferredLocation := time.UTC

	workSchedule := employee.Created.WorkSchedule
	weekdaySchedules := map[time.Weekday][]mJSON.WorkRange{
		time.Sunday:    workSchedule.Sunday,
		time.Monday:    workSchedule.Monday,
		time.Tuesday:   workSchedule.Tuesday,
		time.Wednesday: workSchedule.Wednesday,
		time.Thursday:  workSchedule.Thursday,
		time.Friday:    workSchedule.Friday,
		time.Saturday:  workSchedule.Saturday,
	}

	now := time.Now().In(preferredLocation)
	searchStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, preferredLocation)

	branchCache := make(map[string]*model.Branch)
	serviceCache := make(map[string]*model.Service)

	httpClient := handler.NewHttpClient().
		Header(namespace.HeadersKey.Company, employee.Company.Created.ID.String()).
		Header(namespace.HeadersKey.Auth, employee.Company.Owner.X_Auth_Token)

	for dayOffset := range 8 {
		currentDate := searchStart.AddDate(0, 0, dayOffset)
		currentWeekday := currentDate.Weekday()
		workRanges := weekdaySchedules[currentWeekday]

		for iWr, wr := range workRanges {
			if wr.Start == "" || wr.End == "" || wr.BranchID == uuid.Nil {
				return fmt.Errorf("work range %d has invalid data (Start, End, or BranchID missing)", iWr)
			}

			branchID := wr.BranchID.String()
			branch, ok := branchCache[branchID]
			if !ok {
				var b model.Branch
				if err := httpClient.
					Method("GET").
					URL("/branch/" + branchID).
					ExpectedStatus(200).
					Send(nil).
					ParseResponse(&b).Error; err != nil {
					return fmt.Errorf("failed to get branch %s: %w", branchID, err)
				}
				branchCache[branchID] = &b
				branch = &b
			}

			// Check if employee is assigned to the branch
			assignedToBranch := false
			for _, e := range branch.Employees {
				if e.ID == employee.Created.ID {
					assignedToBranch = true
					break
				}
			}

			if !assignedToBranch {
				return fmt.Errorf("employee %s is not assigned to branch %s.\nEmployee.Branches: %+v\nBranches.Employees: %+v", employee.Created.ID, branchID, employee.Created.Branches, branch.Employees)
			}

			startTime, err := parseTimeWithLocation(currentDate, wr.Start, preferredLocation)
			if err != nil {
				return fmt.Errorf("failed to parse start time for work range #%d: %w", iWr, err)
			}
			endTime, err := parseTimeWithLocation(currentDate, wr.End, preferredLocation)
			if err != nil || !startTime.Before(endTime) {
				return fmt.Errorf("invalid time range for work range #%d: %w", iWr, err)
			}

			for _, serviceID := range wr.Services {
				if serviceID == uuid.Nil {
					return fmt.Errorf("work range %d has a nil service ID", iWr)
				}
				sID := serviceID.String()

				service, ok := serviceCache[sID]
				if !ok {
					var s model.Service
					if err := httpClient.
						Method("GET").
						URL("/service/" + sID).
						ExpectedStatus(200).
						Send(nil).
						ParseResponse(&s).Error; err != nil {
						return fmt.Errorf("failed to get service %s: %w", sID, err)
					}
					serviceCache[sID] = &s
					service = &s
				}

				// Check if employee is assigned to the service
				assignedToService := false
				for _, e := range service.Employees {
					if e.ID == employee.Created.ID {
						assignedToService = true
						break
					}
				}
				if !assignedToService {
					return fmt.Errorf("employee %s is not assigned to service %s", employee.Created.ID, sID)
				}

				// Check if branch is assigned to the service
				serviceAvailableAtBranch := false
				for _, s := range branch.Services {
					if s.ID == serviceID {
						serviceAvailableAtBranch = true
						break
					}
				}
				if !serviceAvailableAtBranch {
					return fmt.Errorf("service %s is not available at branch %s", sID, branchID)
				}
			}
		}
	}
	return nil
}
