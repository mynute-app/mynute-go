package model

import (
	"bytes"
	"fmt"
	DTO "mynute-go/core/src/config/api/dto"
	"mynute-go/core/src/config/db/model"
	"mynute-go/core/src/config/namespace"
	"mynute-go/core/src/lib"
	"mynute-go/core/src/lib/email"
	FileBytes "mynute-go/core/src/lib/file_bytes"
	"mynute-go/test/src/handler"
	"net/url"
	"reflect"
)

type Client struct {
	Created      *model.Client
	Appointments []*Appointment
	X_Auth_Token string
}

func (u *Client) Set() error {
	if err := u.Create(200); err != nil {
		return err
	}

	// 50/50 chance to verify email either by VerifyEmail or LoginWithEmailCode
	// Both methods verify the email, but LoginWithEmailCode also logs in
	if lib.GenerateRandomIntFromRange(0, 1) == 0 {
		// Option 1: Verify email and then login
		if err := u.VerifyEmail(200); err != nil {
			return err
		}
		if err := u.LoginWithPassword(200); err != nil {
			return err
		}
	} else {
		// Option 2: Login with email code (also verifies email)
		if err := u.LoginWithEmailCode(200); err != nil {
			return err
		}
	}

	if err := u.UploadImages(200, map[string][]byte{
		"profile": FileBytes.PNG_FILE_1,
	}, nil); err != nil {
		return err
	}
	if err := u.GetByEmail(200); err != nil {
		return err
	}
	return nil
}

func (u *Client) Create(s int) error {
	pswd := lib.GenerateValidPassword()
	if err := handler.NewHttpClient().
		Method("POST").
		URL("/client").
		ExpectedStatus(s).
		Send(DTO.CreateClient{
			Email:    lib.GenerateRandomEmail("client"),
			Name:     lib.GenerateRandomName("Client Name"),
			Surname:  lib.GenerateRandomName("Client Surname"),
			Password: pswd,
			Phone:    lib.GenerateRandomPhoneNumber(),
		}).ParseResponse(&u.Created).Error; err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	u.Created.Password = pswd
	return nil
}

func (u *Client) Update(s int, changes map[string]any) error {
	if err := handler.NewHttpClient().
		Method("PATCH").
		URL("/client/"+fmt.Sprintf("%v", u.Created.ID.String())).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, u.X_Auth_Token).
		Send(changes).
		ParseResponse(&u.Created).
		Error; err != nil {
		return fmt.Errorf("failed to update client: %w", err)
	}
	if s > 200 && s < 300 {
		if err := ValidateUpdateChanges("Client", u.Created, changes); err != nil {
			return err
		}
	}

	return nil
}

func (u *Client) GetByEmail(s int) error {
	if err := handler.NewHttpClient().
		Method("GET").
		URL("/client/email/" + u.Created.Email).
		ExpectedStatus(s).
		Send(nil).
		ParseResponse(&u.Created).Error; err != nil {
		return fmt.Errorf("failed to get client by email: %w", err)
	}
	return nil
}

func (u *Client) Delete(s int) error {
	if err := handler.NewHttpClient().
		Method("DELETE").
		URL(fmt.Sprintf("/client/%v", u.Created.ID.String())).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, u.X_Auth_Token).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to delete client: %w", err)
	}
	return nil
}

func (u *Client) Login(s int, login_type string) error {
	if login_type == "password" {
		return u.LoginWithPassword(s)
	} else if login_type == "email_code" {
		return u.LoginWithEmailCode(s)
	}
	return fmt.Errorf("invalid login type: %s", login_type)
}

func (u *Client) LoginWithPassword(s int) error {
	if err := u.LoginByPassword(s, u.Created.Password); err != nil {
		return fmt.Errorf("failed to login with password: %w", err)
	}
	return nil
}

func (u *Client) LoginByPassword(s int, password string) error {
	login := DTO.LoginClient{
		Email:    u.Created.Email,
		Password: password,
	}
	http := handler.NewHttpClient()
	if err := http.
		Method("POST").
		URL("/client/login").
		ExpectedStatus(s).
		Send(login).Error; err != nil {
		return fmt.Errorf("failed to login client by password: %w", err)
	}

	if s == 200 {
		auth := http.ResHeaders[namespace.HeadersKey.Auth]
		if len(auth) == 0 {
			return fmt.Errorf("authorization header '%s' not found", namespace.HeadersKey.Auth)
		}
		u.X_Auth_Token = auth[0]
		if err := u.GetByEmail(200); err != nil {
			return fmt.Errorf("failed to get client by email after login by password: %w", err)
		}
	}
	return nil
}

func (u *Client) LoginWithEmailCode(s int) error {
	if err := u.SendLoginCode(s); err != nil {
		return fmt.Errorf("failed to send login code: %w", err)
	}
	code, err := u.GetLoginCodeFromEmail()
	if err != nil {
		return fmt.Errorf("failed to get login code from email: %w", err)
	}
	if err := u.LoginByEmailCode(s, code); err != nil {
		return fmt.Errorf("failed to login by email code: %w", err)
	}
	return nil
}

func (u *Client) LoginByEmailCode(s int, code string) error {
	loginData := DTO.LoginByEmailCode{
		Email: u.Created.Email,
		Code:  code,
	}
	http := handler.NewHttpClient()
	if err := http.
		Method("POST").
		URL("/client/login-with-code").
		ExpectedStatus(s).
		Send(loginData).Error; err != nil {
		return fmt.Errorf("failed to login client by email code: %w", err)
	}

	if s == 200 {
		auth := http.ResHeaders[namespace.HeadersKey.Auth]
		if len(auth) == 0 {
			return fmt.Errorf("authorization header '%s' not found", namespace.HeadersKey.Auth)
		}
		u.X_Auth_Token = auth[0]
		if err := u.GetByEmail(200); err != nil {
			return fmt.Errorf("failed to get client by email after login by code: %w", err)
		}
	}
	return nil
}

func (u *Client) SendLoginCode(s int) error {
	http := handler.NewHttpClient()
	if err := http.
		Method("POST").
		URL(fmt.Sprintf("/client/send-login-code/email/%s?lang=en", url.PathEscape(u.Created.Email))).
		ExpectedStatus(s).
		Send(nil).Error; err != nil {
		return fmt.Errorf("failed to send login code to client: %w", err)
	}
	return nil
}

func (u *Client) GetLoginCodeFromEmail() (string, error) {
	// Initialize MailHog client
	mailhog, err := email.MailHog()
	if err != nil {
		return "", err
	}

	// Get the latest email sent to the client
	message, err := mailhog.GetLatestMessageTo(u.Created.Email)
	if err != nil {
		return "", err
	}

	// Verify the email has a subject
	if message.GetSubject() == "" {
		return "", fmt.Errorf("email subject is empty")
	}

	// Extract the validation code from the email
	code, err := message.ExtractValidationCode()
	if err != nil {
		return "", err
	}

	return code, nil
}

func (u *Client) SendPasswordResetEmail(s int) error {
	http := handler.NewHttpClient()
	if err := http.
		Method("POST").
		URL(fmt.Sprintf("/client/reset-password/%s?lang=en", url.PathEscape(u.Created.Email))).
		ExpectedStatus(s).
		Send(nil).Error; err != nil {
		return fmt.Errorf("failed to send password reset email to client: %w", err)
	}
	return nil
}

func (u *Client) GetNewPasswordFromEmail() (string, error) {
	// Initialize MailHog client
	mailhog, err := email.MailHog()
	if err != nil {
		return "", err
	}

	// Get the latest email sent to the client
	message, err := mailhog.GetLatestMessageTo(u.Created.Email)
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
		return "", fmt.Errorf("failed to extract password: %w", err)
	}

	return password, nil
}

func (u *Client) ResetPasswordByEmail(s int) error {
	if err := u.SendPasswordResetEmail(s); err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	newPassword, err := u.GetNewPasswordFromEmail()
	if err != nil {
		return fmt.Errorf("failed to get new password from email: %w", err)
	}

	// Update the password in memory
	u.Created.Password = newPassword

	// Try to login with the new password
	if err := u.LoginByPassword(200, newPassword); err != nil {
		return fmt.Errorf("failed to login with new password: %w", err)
	}

	return nil
}

func (u *Client) SendVerificationEmail(s int) error {
	http := handler.NewHttpClient()
	if err := http.
		Method("POST").
		URL(fmt.Sprintf("/client/send-verification-code/email/%s?language=en", url.PathEscape(u.Created.Email))).
		ExpectedStatus(s).
		Send(nil).Error; err != nil {
		return fmt.Errorf("failed to send verification email to client: %w", err)
	}
	return nil
}

func (u *Client) GetVerificationCodeFromEmail() (string, error) {
	// Initialize MailHog client
	mailhog, err := email.MailHog()
	if err != nil {
		return "", err
	}

	// Get the latest email sent to the client
	message, err := mailhog.GetLatestMessageTo(u.Created.Email)
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

func (u *Client) VerifyEmailByCode(s int, code string) error {
	http := handler.NewHttpClient()
	if err := http.
		Method("GET").
		URL(fmt.Sprintf("/client/verify-email/%s/%s", url.PathEscape(u.Created.Email), code)).
		ExpectedStatus(s).
		Send(nil).Error; err != nil {
		return fmt.Errorf("failed to verify client email: %w", err)
	}

	if s == 200 {
		// Update the verified status in memory
		u.Created.Verified = true
		if err := u.GetByEmail(200); err != nil {
			return fmt.Errorf("failed to get client by email after verification: %w", err)
		}
	}
	return nil
}

func (u *Client) VerifyEmail(s int) error {
	if err := u.SendVerificationEmail(s); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	code, err := u.GetVerificationCodeFromEmail()
	if err != nil {
		return fmt.Errorf("failed to get verification code from email: %w", err)
	}

	if err := u.VerifyEmailByCode(s, code); err != nil {
		return fmt.Errorf("failed to verify email with code: %w", err)
	}

	return nil
}

func ValidateUpdateChanges(modelName string, v any, changes map[string]any) error {
	mappy, err := lib.StructToMap(v)
	if err != nil {
		return fmt.Errorf("failed to convert %s struct to map: %w", modelName, err)
	}

	for key, expected := range changes {
		// Se o expected for struct, transforma em map
		if reflect.TypeOf(expected).Kind() == reflect.Struct {
			expected, err = lib.StructToMap(expected)
			if err != nil {
				return fmt.Errorf("failed to convert expected value for key '%s' to map: %w", key, err)
			}
		}

		actual := mappy[key]

		if !reflect.DeepEqual(actual, expected) {
			return fmt.Errorf("%s %s was not updated: expected '%#v' but got '%#v'", modelName, key, expected, actual)
		}
	}

	return nil
}

func (c *Client) UploadImages(status int, files map[string][]byte, x_auth_token *string) error {
	var fileMap = make(handler.Files)
	for field, content := range files {
		fileMap[field] = handler.MyFile{
			Name:    field + "_" + lib.GenerateRandomString(6) + ".png",
			Content: content,
		}
	}

	t, err := Get_x_auth_token(x_auth_token, &c.X_Auth_Token)
	if err != nil {
		return err
	}

	if err := handler.NewHttpClient().
		Method("PATCH").
		URL(fmt.Sprintf("/client/%s/design/images", c.Created.ID.String())).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Auth, t).
		Send(fileMap).
		ParseResponse(&c.Created.Meta.Design.Images).
		Error; err != nil {
		return fmt.Errorf("failed to upload client images: %w", err)
	}

	return nil
}

func (c *Client) GetImage(status int, imageURL string, compareImgBytes *[]byte) error {
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

func (c *Client) DeleteImages(status int, image_types []string, x_auth_token *string) error {
	if len(image_types) == 0 {
		return fmt.Errorf("no image types provided to delete")
	}

	t, err := Get_x_auth_token(x_auth_token, &c.X_Auth_Token)
	if err != nil {
		return err
	}

	http := handler.NewHttpClient()

	if err := http.
		Method("DELETE").
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Auth, t).
		Error; err != nil {
		return fmt.Errorf("failed to prepare delete images request: %w", err)
	}

	base_url := fmt.Sprintf("/client/%s/design/images", c.Created.ID.String())
	for _, image_type := range image_types {
		image_url := base_url + "/" + image_type
		http.URL(image_url)
		http.Send(nil)
		http.ParseResponse(&c.Created.Meta.Design.Images)
		if http.Error != nil {
			return fmt.Errorf("failed to delete image %s: %w", image_type, http.Error)
		}
		url := c.Created.Meta.Design.Images.GetImageURL(image_type)
		if url != "" {
			return fmt.Errorf("image %s was not deleted successfully, expected empty URL but got %s", image_type, url)
		}
	}
	return nil
}

// GetAppointments retrieves the appointments for this client with pagination and filters
func (c *Client) GetAppointments(status int, page int, pageSize int, startDate string, endDate string, cancelled string, timezone string, x_auth_token *string, x_company_id *string) (*DTO.AppointmentList, error) {
	t, err := Get_x_auth_token(x_auth_token, &c.X_Auth_Token)
	if err != nil {
		return nil, err
	}

	var appointmentList DTO.AppointmentList
	urlStr := fmt.Sprintf("/client/%s/appointments?page=%d&page_size=%d", c.Created.ID.String(), page, pageSize)

	if timezone != "" {
		urlStr += fmt.Sprintf("&timezone=%s", url.QueryEscape(timezone))
	}
	if startDate != "" {
		urlStr += fmt.Sprintf("&start_date=%s", startDate)
	}
	if endDate != "" {
		urlStr += fmt.Sprintf("&end_date=%s", endDate)
	}
	if cancelled != "" {
		urlStr += fmt.Sprintf("&cancelled=%s", cancelled)
	}

	req := handler.NewHttpClient().
		Method("GET").
		URL(urlStr).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Auth, t)

	// Add X-Company-ID header if provided
	if x_company_id != nil && *x_company_id != "" {
		req.Header(namespace.HeadersKey.Company, *x_company_id)
	}

	if err := req.Send(nil).
		ParseResponse(&appointmentList).
		Error; err != nil {
		return nil, fmt.Errorf("failed to get client appointments: %w", err)
	}

	return &appointmentList, nil
}
