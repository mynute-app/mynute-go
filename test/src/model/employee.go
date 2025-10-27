package model

import (
	"bytes"
	"fmt"
	DTO "mynute-go/core/src/config/api/dto"
	"mynute-go/core/src/config/db/model"
	"mynute-go/core/src/config/namespace"
	"mynute-go/core/src/lib"
	"mynute-go/core/src/lib/email"
	"mynute-go/test/src/handler"
	"net/url"

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

func (e *Employee) GetID() string        { return e.Created.ID.String() }
func (e *Employee) GetCompanyID() string { return e.Company.Created.ID.String() }
func (e *Employee) GetAuthToken() string { return e.X_Auth_Token }
func (e *Employee) SetWorkRanges(wr []any) error {
	e.Created.WorkSchedule = make([]model.EmployeeWorkRange, len(wr))
	for i, v := range wr {
		if ewr, ok := v.(model.EmployeeWorkRange); !ok {
			return fmt.Errorf("invalid work range type")
		} else {
			e.Created.WorkSchedule[i] = ewr
		}
	}
	return nil
}

func (e *Employee) Create(s int, x_auth_token *string, x_company_id *string) error {
	pswd := lib.GenerateValidPassword()

	t, err := Get_x_auth_token(x_auth_token, &e.X_Auth_Token)
	if err != nil {
		return err
	}
	companyIDStr := e.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
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
			TimeZone:  "America/Sao_Paulo", // Use a valid timezone
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
	t, err := Get_x_auth_token(x_auth_token, &e.X_Auth_Token)
	if err != nil {
		return err
	}
	companyIDStr := e.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
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
	if s > 200 && s < 300 {
		if err := ValidateUpdateChanges("Employee", e.Created, changes); err != nil {
			return err
		}
	}
	return nil
}

func (e *Employee) CreateWorkSchedule(s int, EmployeeWorkSchedule DTO.CreateEmployeeWorkSchedule, x_auth_token *string, x_company_id *string) error {
	if EmployeeWorkSchedule.WorkRanges == nil {
		return fmt.Errorf("work schedule cannot be nil")
	}
	t, err := Get_x_auth_token(x_auth_token, &e.X_Auth_Token)
	if err != nil {
		return err
	}
	companyIDStr := e.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	var updated *model.EmployeeWorkSchedule

	http := handler.NewHttpClient()
	if err := http.
		Method("POST").
		URL(fmt.Sprintf("/employee/%s/work_schedule", e.Created.ID.String())).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Company, cID).
		Header(namespace.HeadersKey.Auth, t).
		Send(EmployeeWorkSchedule).
		ParseResponse(&updated).Error; err != nil {
		return fmt.Errorf("failed to update employee work schedule: %w", err)
	}

	e.Created.WorkSchedule = updated.WorkRanges
	return nil
}

func (e *Employee) GetWorkSchedule(s int, x_auth_token *string, x_company_id *string) error {
	t, err := Get_x_auth_token(x_auth_token, &e.X_Auth_Token)
	if err != nil {
		return err
	}
	companyIDStr := e.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	var schedule *model.EmployeeWorkSchedule
	if err := handler.NewHttpClient().
		Method("GET").
		URL(fmt.Sprintf("/employee/%s/work_schedule", e.Created.ID.String())).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Company, cID).
		Header(namespace.HeadersKey.Auth, t).
		Send(nil).
		ParseResponse(&schedule).
		Error; err != nil {
		return fmt.Errorf("failed to get employee work schedule: %w", err)
	}
	e.Created.WorkSchedule = schedule.WorkRanges
	return nil
}

func (e *Employee) UpdateWorkRange(status int, wrID string, changes map[string]any, x_auth_token *string, x_company_id *string) error {
	companyIDStr := e.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	t, err := Get_x_auth_token(x_auth_token, &e.X_Auth_Token)
	if err != nil {
		return err
	}
	var updated *model.EmployeeWorkSchedule
	if err := handler.NewHttpClient().
		Method("PUT").
		URL(fmt.Sprintf("/employee/%s/work_range/%s", e.Created.ID.String(), wrID)).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Company, cID).
		Header(namespace.HeadersKey.Auth, t).
		Send(changes).
		ParseResponse(&updated).
		Error; err != nil {
		return fmt.Errorf("failed to update branch work range: %w", err)
	}
	e.Created.WorkSchedule = updated.WorkRanges
	return nil
}

func (e *Employee) DeleteWorkRange(status int, wrID string, x_auth_token *string, x_company_id *string) error {
	companyIDStr := e.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	t, err := Get_x_auth_token(x_auth_token, &e.X_Auth_Token)
	if err != nil {
		return err
	}
	var updated *model.EmployeeWorkSchedule
	if err := handler.NewHttpClient().
		Method("DELETE").
		URL(fmt.Sprintf("/employee/%s/work_range/%s", e.Created.ID.String(), wrID)).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Company, cID).
		Header(namespace.HeadersKey.Auth, t).
		Send(nil).
		ParseResponse(&updated).
		Error; err != nil {
		return fmt.Errorf("failed to delete branch work schedule: %w", err)
	}
	e.Created.WorkSchedule = updated.WorkRanges
	return nil
}

func (e *Employee) AddServicesToWorkRange(s int, wrID string, body DTO.EmployeeWorkRangeServices, x_auth_token *string, x_company_id *string) error {
	t, err := Get_x_auth_token(x_auth_token, &e.X_Auth_Token)
	if err != nil {
		return err
	}
	companyIDStr := e.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	var updated *model.EmployeeWorkSchedule
	if err := handler.NewHttpClient().
		Method("POST").
		URL(fmt.Sprintf("/employee/%s/work_range/%s/services", e.Created.ID.String(), wrID)).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Company, cID).
		Header(namespace.HeadersKey.Auth, t).
		Send(body).
		ParseResponse(&updated).
		Error; err != nil {
		return fmt.Errorf("failed to add services to employee work range: %w", err)
	}
	e.Created.WorkSchedule = updated.WorkRanges
	return nil
}

func (e *Employee) RemoveServiceFromWorkRange(s int, wrID string, serviceID string, x_auth_token *string, x_company_id *string) error {
	t, err := Get_x_auth_token(x_auth_token, &e.X_Auth_Token)
	if err != nil {
		return err
	}
	companyIDStr := e.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	var updated *model.EmployeeWorkSchedule
	if err := handler.NewHttpClient().
		Method("DELETE").
		URL(fmt.Sprintf("/employee/%s/work_range/%s/service/%s", e.Created.ID.String(), wrID, serviceID)).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Company, cID).
		Header(namespace.HeadersKey.Auth, t).
		Send(nil).
		ParseResponse(&updated).
		Error; err != nil {
		return fmt.Errorf("failed to remove service from employee work range: %w", err)
	}
	e.Created.WorkSchedule = updated.WorkRanges
	return nil
}

func (e *Employee) GetById(s int, x_auth_token *string, x_company_id *string) error {
	t, err := Get_x_auth_token(x_auth_token, &e.X_Auth_Token)
	if err != nil {
		return err
	}
	companyIDStr := e.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
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
	t, err := Get_x_auth_token(x_auth_token, &e.X_Auth_Token)
	if err != nil {
		return err
	}
	companyIDStr := e.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("GET").
		URL(fmt.Sprintf("/employee/email/%s", url.PathEscape(e.Created.Email))).
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
	t, err := Get_x_auth_token(x_auth_token, &e.X_Auth_Token)
	if err != nil {
		return err
	}
	companyIDStr := e.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
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

func (e *Employee) LoginWith(s int, login_type string, x_company_id *string) error {
	if login_type == "password" {
		return e.LoginWithPassword(s, x_company_id)
	} else if login_type == "email_code" {
		return e.LoginWithEmailCode(s, x_company_id)
	}
	return fmt.Errorf("invalid login type: %s", login_type)
}

func (e *Employee) LoginWithPassword(s int, x_company_id *string) error {
	if err := e.LoginByPassword(s, e.Created.Password, x_company_id); err != nil {
		return fmt.Errorf("failed to login with password: %w", err)
	}
	return nil
}

func (e *Employee) LoginByPassword(s int, password string, x_company_id *string) error {
	login := DTO.LoginEmployee{
		Email:    e.Created.Email,
		Password: password,
	}
	companyIDStr := e.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	http := handler.NewHttpClient()
	if err := http.
		Method("POST").
		URL("/employee/login").
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Company, cID).
		Send(login).Error; err != nil {
		return fmt.Errorf("failed to login employee by password: %w", err)
	}

	if s == 200 {
		auth := http.ResHeaders[namespace.HeadersKey.Auth]
		if len(auth) == 0 {
			return fmt.Errorf("authentication token not found in response headers")
		}
		e.X_Auth_Token = auth[0]
		if err := e.GetById(200, nil, nil); err != nil {
			return fmt.Errorf("failed to get employee by ID after login by password: %w", err)
		}
	}
	return nil
}

func (e *Employee) LoginWithEmailCode(s int, x_company_id *string) error {
	if err := e.SendLoginCode(s, x_company_id); err != nil {
		return fmt.Errorf("failed to send login code: %w", err)
	}
	code, err := e.GetLoginCodeFromEmail()
	if err != nil {
		return fmt.Errorf("failed to get login code from email: %w", err)
	}
	if err := e.LoginByEmailCode(s, code, x_company_id); err != nil {
		return fmt.Errorf("failed to login by email code: %w", err)
	}
	return nil
}

func (e *Employee) LoginByEmailCode(s int, code string, x_company_id *string) error {
	loginData := DTO.LoginByEmailCode{
		Email: e.Created.Email,
		Code:  code,
	}
	companyIDStr := e.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	http := handler.NewHttpClient()
	if err := http.
		Method("POST").
		URL("/employee/login-with-code").
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Company, cID).
		Send(loginData).Error; err != nil {
		return fmt.Errorf("failed to login employee by email code: %w", err)
	}

	if s == 200 {
		auth := http.ResHeaders[namespace.HeadersKey.Auth]
		if len(auth) == 0 {
			return fmt.Errorf("authorization header '%s' not found", namespace.HeadersKey.Auth)
		}
		e.X_Auth_Token = auth[0]
		if err := e.GetById(200, nil, nil); err != nil {
			return fmt.Errorf("failed to get employee by email after login by code: %w", err)
		}
	}
	return nil
}

func (e *Employee) SendLoginCode(s int, x_company_id *string) error {
	// Note: The employee send-login-code endpoint DOES require X-Company-ID header
	// to switch to the correct company schema before querying for the employee
	companyIDStr := e.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}

	http := handler.NewHttpClient()
	if err := http.
		Method("POST").
		URL(fmt.Sprintf("/employee/send-login-code/email/%s?lang=en", url.PathEscape(e.Created.Email))).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Company, cID).
		Send(nil).Error; err != nil {
		return fmt.Errorf("failed to send login code to employee: %w", err)
	}
	return nil
}

func (e *Employee) GetLoginCodeFromEmail() (string, error) {
	// Initialize MailHog client
	mailhog, err := email.MailHog()
	if err != nil {
		return "", err
	}

	// Get all messages to find the login validation email
	messages, err := mailhog.GetMessages()
	if err != nil {
		return "", err
	}

	// Search for the most recent login validation email (searching from newest to oldest)
	var loginMessage *email.MailHogMessage
	for i := len(messages) - 1; i >= 0; i-- {
		msg := &messages[i]

		// Check if this message is for the employee
		isForEmployee := false
		for _, to := range msg.To {
			recipientEmail := fmt.Sprintf("%s@%s", to.Mailbox, to.Domain)
			if recipientEmail == e.Created.Email {
				isForEmployee = true
				break
			}
		}

		if !isForEmployee {
			continue
		}

		// Check if this is a login validation email
		subject := msg.GetSubject()
		if subject == "Your Login Validation Code" ||
			subject == "Seu Código de Validação de Login" ||
			subject == "Su Código de Validación de Inicio de Sesión" {
			loginMessage = msg
			break
		}
	}

	if loginMessage == nil {
		return "", fmt.Errorf("no login validation email found for %s", e.Created.Email)
	}

	// Extract the validation code from the email
	code, err := loginMessage.ExtractValidationCode()
	if err != nil {
		return "", err
	}

	return code, nil
}

func (e *Employee) SendPasswordResetEmail(s int, x_company_id *string) error {
	// Note: The employee reset-password endpoint requires X-Company-ID header
	companyIDStr := e.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}

	http := handler.NewHttpClient()
	if err := http.
		Method("POST").
		URL(fmt.Sprintf("/employee/reset-password/%s?lang=en", url.PathEscape(e.Created.Email))).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Company, cID).
		Send(nil).Error; err != nil {
		return fmt.Errorf("failed to send password reset email to employee: %w", err)
	}
	return nil
}

func (e *Employee) GetNewPasswordFromEmail() (string, error) {
	// Initialize MailHog client
	mailhog, err := email.MailHog()
	if err != nil {
		return "", err
	}

	// Get the latest email sent to the employee
	message, err := mailhog.GetLatestMessageTo(e.Created.Email)
	if err != nil {
		return "", err
	}

	// Verify the email has a subject
	if message.GetSubject() == "" {
		return "", fmt.Errorf("email subject is empty")
	}

	// Extract the new password from the email
	password, err := message.ExtractPassword()
	if err != nil {
		return "", err
	}

	return password, nil
}

func (e *Employee) ResetPasswordByEmail(s int, x_company_id *string) error {
	if err := e.SendPasswordResetEmail(s, x_company_id); err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	newPassword, err := e.GetNewPasswordFromEmail()
	if err != nil {
		return fmt.Errorf("failed to get new password from email: %w", err)
	}

	// Update the password in memory
	e.Created.Password = newPassword

	// Try to login with the new password
	if err := e.LoginByPassword(200, newPassword, x_company_id); err != nil {
		return fmt.Errorf("failed to login with new password: %w", err)
	}

	return nil
}

func (e *Employee) SendVerificationEmail(s int, x_company_id *string) error {
	// Note: The employee send-verification-code endpoint requires X-Company-ID header
	companyIDStr := e.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}

	http := handler.NewHttpClient()
	if err := http.
		Method("POST").
		URL(fmt.Sprintf("/employee/send-verification-code/email/%s/%s?language=en", url.PathEscape(e.Created.Email), cID)).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Company, cID).
		Send(nil).Error; err != nil {
		return fmt.Errorf("failed to send verification email to employee: %w", err)
	}
	return nil
}

func (e *Employee) GetVerificationCodeFromEmail() (string, error) {
	// Initialize MailHog client
	mailhog, err := email.MailHog()
	if err != nil {
		return "", err
	}

	// Get the latest email sent to the employee
	message, err := mailhog.GetLatestMessageTo(e.Created.Email)
	if err != nil {
		return "", err
	}

	// Verify the email has a subject
	if message.GetSubject() == "" {
		return "", fmt.Errorf("email subject is empty")
	}

	// Extract the verification code from the email
	code, err := message.ExtractValidationCode()
	if err != nil {
		return "", err
	}

	return code, nil
}

func (e *Employee) VerifyEmailByCode(s int, code string, x_company_id *string) error {
	companyIDStr := e.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}

	http := handler.NewHttpClient()
	if err := http.
		Method("GET").
		URL(fmt.Sprintf("/employee/verify-email/%s/%s/%s", url.PathEscape(e.Created.Email), code, cID)).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Company, cID).
		Send(nil).Error; err != nil {
		return fmt.Errorf("failed to verify employee email: %w", err)
	}

	if s == 200 {
		// Update the verified status in memory
		e.Created.Verified = true
	}
	return nil
}

func (e *Employee) VerifyEmail(s int, x_company_id *string) error {
	if err := e.SendVerificationEmail(s, x_company_id); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	code, err := e.GetVerificationCodeFromEmail()
	if err != nil {
		return fmt.Errorf("failed to get verification code from email: %w", err)
	}

	if err := e.VerifyEmailByCode(s, code, x_company_id); err != nil {
		return fmt.Errorf("failed to verify email with code: %w", err)
	}

	return nil
}

// func (e *Employee) VerifyEmail(s int, x_company_id *string) error {
// 	companyIDStr := e.Company.Created.ID.String()
// 	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
// 	if err != nil {
// 		return err
// 	}
// 	if err := handler.NewHttpClient().
// 		Method("POST").
// 		URL(fmt.Sprintf("/employee/verify-email/%s/%s", e.Created.Email, "12345")).
// 		ExpectedStatus(s).
// 		Header(namespace.HeadersKey.Company, cID).
// 		Send(nil).
// 		Error; err != nil {
// 		return fmt.Errorf("failed to verify employee email: %w", err)
// 	}
// 	return nil
// }

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
	t, err := Get_x_auth_token(token, &e.X_Auth_Token)
	if err != nil {
		return err
	}
	companyIDStr := e.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
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
	if err := b.GetById(s, b.Company.Owner.X_Auth_Token, nil); err != nil {
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
	t, err := Get_x_auth_token(token, &e.X_Auth_Token)
	if err != nil {
		return err
	}
	companyIDStr := e.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
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
	t, err := Get_x_auth_token(x_auth_token, &e.X_Auth_Token)
	if err != nil {
		return err
	}
	companyIDStr := e.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
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

func (e *Employee) UploadImages(status int, files map[string][]byte, x_auth_token *string, x_company_id *string) error {
	var fileMap = make(handler.Files)
	for field, content := range files {
		fileMap[field] = handler.MyFile{
			Name:    field + "_" + lib.GenerateRandomString(6) + ".png",
			Content: content,
		}
	}

	companyIDStr := e.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}

	t, err := Get_x_auth_token(x_auth_token, &e.X_Auth_Token)
	if err != nil {
		return err
	}

	if err := handler.NewHttpClient().
		Method("PATCH").
		URL(fmt.Sprintf("/employee/%s/design/images", e.Created.ID.String())).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Auth, t).
		Header(namespace.HeadersKey.Company, cID).
		Send(fileMap).
		ParseResponse(&e.Created.Meta.Design.Images).
		Error; err != nil {
		return fmt.Errorf("failed to upload employee images: %w", err)
	}

	return nil
}

func (e *Employee) GetImage(status int, imageURL string, compareImgBytes *[]byte) error {
	if imageURL == "" {
		return fmt.Errorf("image URL cannot be empty")
	}
	http := handler.NewHttpClient()
	http.Method("GET")
	http.URL(imageURL)
	http.ExpectedStatus(status)
	http.Send(nil)
	// Compare the response bytes with the expected image bytes
	if compareImgBytes != nil {
		var response []byte
		http.ParseResponse(&response)
		if len(response) == 0 {
			return fmt.Errorf("received empty response for image (%s)", imageURL)
		} else if len(response) != len(*compareImgBytes) {
			return fmt.Errorf("image size mismatch for %s: expected %d bytes, got %d bytes", imageURL, len(*compareImgBytes), len(response))
		} else if !bytes.Equal(response, *compareImgBytes) {
			return fmt.Errorf("image content mismatch for %s", imageURL)
		}
	}
	return nil
}

func (e *Employee) DeleteImages(status int, image_types []string, x_auth_token string, x_company_id *string) error {
	if len(image_types) == 0 {
		return fmt.Errorf("no image types provided to delete")
	}

	createdEmployeeID := e.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &createdEmployeeID)
	if err != nil {
		return fmt.Errorf("failed to get company ID for deletion: %w", err)
	}

	http := handler.NewHttpClient()

	if err := http.
		Method("DELETE").
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Header(namespace.HeadersKey.Company, cID).
		Error; err != nil {
		return fmt.Errorf("failed to prepare delete images request: %w", err)
	}

	base_url := fmt.Sprintf("/employee/%s/design/images", e.Created.ID.String())
	for _, image_type := range image_types {
		image_url := base_url + "/" + image_type
		http.URL(image_url)
		http.Send(nil)
		http.ParseResponse(&e.Created.Meta.Design.Images)
		if http.Error != nil {
			return fmt.Errorf("failed to delete image %s: %w", image_type, http.Error)
		}
		url := e.Created.Meta.Design.Images.GetImageURL(image_type)
		if url != "" {
			return fmt.Errorf("image %s was not deleted successfully, expected empty URL but got %s", image_type, url)
		}
	}
	return nil
}

func Get_x_auth_token(priority *string, secundary *string) (string, error) {
	if priority != nil {
		return *priority, nil
	} else if secundary != nil {
		return *secundary, nil
	}
	return "", fmt.Errorf("no authentication token provided")
}

func Get_x_company_id(priority *string, secundary *string) (string, error) {
	if priority != nil {
		return *priority, nil
	} else if secundary != nil {
		return *secundary, nil
	}
	return "", fmt.Errorf("no company ID provided")
}

func GetExampleEmployeeWorkSchedule(employeeID uuid.UUID, branchID uuid.UUID, servicesID []DTO.ServiceBase) DTO.CreateEmployeeWorkSchedule {
	return DTO.CreateEmployeeWorkSchedule{
		WorkRanges: []DTO.CreateEmployeeWorkRange{
			{
				EmployeeID:                employeeID,
				BranchID:                  branchID,
				Weekday:                   1,
				StartTime:                 "08:00",
				EndTime:                   "12:00",
				TimeZone:                  "America/Sao_Paulo",
				EmployeeWorkRangeServices: DTO.EmployeeWorkRangeServices{Services: servicesID},
			},
			{
				EmployeeID:                employeeID,
				BranchID:                  branchID,
				Weekday:                   1,
				StartTime:                 "13:00",
				EndTime:                   "17:00",
				TimeZone:                  "America/Sao_Paulo",
				EmployeeWorkRangeServices: DTO.EmployeeWorkRangeServices{Services: servicesID},
			},
			{
				EmployeeID:                employeeID,
				BranchID:                  branchID,
				Weekday:                   2,
				StartTime:                 "08:00",
				EndTime:                   "12:00",
				TimeZone:                  "America/Sao_Paulo",
				EmployeeWorkRangeServices: DTO.EmployeeWorkRangeServices{Services: servicesID},
			},
			{
				EmployeeID:                employeeID,
				BranchID:                  branchID,
				Weekday:                   2,
				StartTime:                 "13:00",
				EndTime:                   "17:00",
				TimeZone:                  "America/Sao_Paulo",
				EmployeeWorkRangeServices: DTO.EmployeeWorkRangeServices{Services: servicesID},
			},
			{
				EmployeeID:                employeeID,
				BranchID:                  branchID,
				Weekday:                   3,
				StartTime:                 "08:00",
				EndTime:                   "12:00",
				TimeZone:                  "America/Sao_Paulo",
				EmployeeWorkRangeServices: DTO.EmployeeWorkRangeServices{Services: servicesID},
			},
			{
				EmployeeID:                employeeID,
				BranchID:                  branchID,
				Weekday:                   3,
				StartTime:                 "13:00",
				EndTime:                   "17:00",
				TimeZone:                  "America/Sao_Paulo",
				EmployeeWorkRangeServices: DTO.EmployeeWorkRangeServices{Services: servicesID},
			},
			{
				EmployeeID:                employeeID,
				BranchID:                  branchID,
				Weekday:                   4,
				StartTime:                 "08:00",
				EndTime:                   "12:00",
				TimeZone:                  "America/Sao_Paulo",
				EmployeeWorkRangeServices: DTO.EmployeeWorkRangeServices{Services: servicesID},
			},
			{
				EmployeeID:                employeeID,
				BranchID:                  branchID,
				Weekday:                   4,
				StartTime:                 "13:00",
				EndTime:                   "17:00",
				TimeZone:                  "America/Sao_Paulo",
				EmployeeWorkRangeServices: DTO.EmployeeWorkRangeServices{Services: servicesID},
			},
			{
				EmployeeID:                employeeID,
				BranchID:                  branchID,
				Weekday:                   5,
				StartTime:                 "08:00",
				EndTime:                   "12:00",
				TimeZone:                  "America/Sao_Paulo",
				EmployeeWorkRangeServices: DTO.EmployeeWorkRangeServices{Services: servicesID},
			},
			{
				EmployeeID:                employeeID,
				BranchID:                  branchID,
				Weekday:                   5,
				StartTime:                 "13:00",
				EndTime:                   "17:00",
				TimeZone:                  "America/Sao_Paulo",
				EmployeeWorkRangeServices: DTO.EmployeeWorkRangeServices{Services: servicesID},
			},
			{
				EmployeeID:                employeeID,
				BranchID:                  branchID,
				Weekday:                   6,
				StartTime:                 "08:00",
				EndTime:                   "12:00",
				TimeZone:                  "America/Sao_Paulo",
				EmployeeWorkRangeServices: DTO.EmployeeWorkRangeServices{Services: servicesID},
			},
		},
	}
}
