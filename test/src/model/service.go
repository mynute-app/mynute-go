package model

import (
	"bytes"
	"fmt"
	DTO "mynute-go/core/src/config/api/dto"
	"mynute-go/core/src/config/db/model"
	"mynute-go/core/src/config/namespace"
	"mynute-go/core/src/lib"
	"mynute-go/test/src/handler"
	"time"
)

type Service struct {
	Created   *model.Service
	Company   *Company
	Employees []*Employee
	Branches  []*Branch
}

func (s *Service) Create(status int, x_auth_token string, x_company_id *string) error {
	companyIDStr := s.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
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
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
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
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
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
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
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
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
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

func (s *Service) UploadImages(status int, files map[string][]byte, x_auth_token string, x_company_id *string) error {
	var fileMap = make(handler.Files)
	for field, content := range files {
		fileMap[field] = handler.MyFile{
			Name:    field + "_" + lib.GenerateRandomString(6) + ".png",
			Content: content,
		}
	}

	companyIDStr := s.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}

	if err := handler.NewHttpClient().
		Method("PATCH").
		URL(fmt.Sprintf("/service/%s/design/images", s.Created.ID.String())).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Header(namespace.HeadersKey.Company, cID).
		Send(fileMap).
		ParseResponse(&s.Created.Design.Images).
		Error; err != nil {
		return fmt.Errorf("failed to upload service images: %w", err)
	}

	return nil
}

func (s *Service) DeleteImages(status int, image_types []string, x_auth_token string, x_company_id *string) error {
	if len(image_types) == 0 {
		return fmt.Errorf("no images provided to delete")
	}

	createdCompanyID := s.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &createdCompanyID)
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

	base_url := fmt.Sprintf("/service/%s/design/images", s.Created.ID.String())
	for _, image_type := range image_types {
		image_url := base_url + "/" + image_type
		http.URL(image_url)
		http.Send(nil)
		http.ParseResponse(&s.Created.Design.Images)
		if http.Error != nil {
			return fmt.Errorf("failed to delete image %s: %w", image_type, http.Error)
		}
		url := s.Created.Design.Images.GetImageURL(image_type)
		if url != "" {
			return fmt.Errorf("image %s was not deleted successfully, expected empty URL but got %s", image_type, url)
		}
	}
	return nil
}

func (s *Service) GetImage(status int, imageURL string, compareImgBytes *[]byte) error {
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

func (s *Service) GetAvailability(status int, x_company_id *string, from, to int) error {
	companyIDStr := s.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	http := handler.NewHttpClient()
	http.Method("GET")
	http.ExpectedStatus(status)
	url := fmt.Sprintf("/service/%s/availability?date_forward_start=%d&date_forward_end=%d", s.Created.ID.String(), from, to)
	http.URL(url)
	http.Header(namespace.HeadersKey.Company, cID)
	http.Send(nil)

	if http.Error != nil {
		return fmt.Errorf("failed to get service availability: %w", http.Error)
	}

	return nil
}

func (s *Service) FindValidRandomAppointmentSlot() (*FoundAppointmentSlot, error) {
	cID := s.Company.Created.ID.String()
	x_auth_token := s.Company.Owner.X_Auth_Token

	http := handler.NewHttpClient()
	http.Method("GET")
	http.ExpectedStatus(200)
	url := fmt.Sprintf("/service/%s/availability?date_forward_start=%d&date_forward_end=%d", s.Created.ID.String(), time.Now().Unix(), time.Now().Add(30*24*time.Hour).Unix())
	http.URL(url)
	http.Header(namespace.HeadersKey.Company, cID)
	http.Header(namespace.HeadersKey.Auth, x_auth_token)
	http.Send(nil)
	var availability DTO.ServiceAvailability
	http.ParseResponse(&availability)
	if http.Error != nil {
		return nil, fmt.Errorf("failed to get service availability: %w", http.Error)
	}
	if len(availability.AvailableDates) == 0 {
		return nil, fmt.Errorf("no available slots found for service %s", s.Created.Name)
	}
	// Pick a random available date
	randomAvailableDate := availability.AvailableDates[lib.GenerateRandomInt(len(availability.AvailableDates))]
	BranchID := randomAvailableDate.BranchID.String()
	dateStr := randomAvailableDate.Date
	randomAvailableTime := randomAvailableDate.AvailableTimes[lib.GenerateRandomInt(len(randomAvailableDate.AvailableTimes))]
	timeStr := randomAvailableTime.Time
	StartTimeRFC3339 := fmt.Sprintf("%sT%s:00Z", dateStr, timeStr)
	// Find the Branch TimeZone
	TimeZone := ""
	for _, b := range availability.BranchInfo {
		if b.ID.String() == BranchID {
			TimeZone = b.TimeZone
			break
		}
	}
	if TimeZone == "" {
		return nil, fmt.Errorf("branch time zone not found for branch ID %s", BranchID)
	}

	return &FoundAppointmentSlot{
		StartTimeRFC3339: StartTimeRFC3339,
		BranchID:         BranchID,
		ServiceID:        s.Created.ID.String(),
		TimeZone:         TimeZone,
	}, nil
}
