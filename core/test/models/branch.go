package modelT

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/test/handlers"
	"bytes"
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

func (b *Branch) GetID() string        { return b.Created.ID.String() }
func (b *Branch) GetCompanyID() string { return b.Company.Created.ID.String() }
func (b *Branch) GetAuthToken() string { return "" }
func (b *Branch) SetWorkRanges(wr []any) error {
	b.Created.BranchWorkSchedule = make([]model.BranchWorkRange, len(wr))
	for i, v := range wr {
		if ewr, ok := v.(model.BranchWorkRange); !ok {
			return fmt.Errorf("invalid work range type")
		} else {
			b.Created.BranchWorkSchedule[i] = ewr
		}
	}
	return nil
}

func (b *Branch) Create(status int, x_auth_token string, x_company_id *string) error {
	companyIDStr := b.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
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
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
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
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
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
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
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
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
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

func (b *Branch) UploadImages(status int, files map[string][]byte, x_auth_token string, x_company_id *string) error {
	var fileMap = make(handler.Files)
	for field, content := range files {
		fileMap[field] = handler.MyFile{
			Name:    field + "_" + lib.GenerateRandomString(6) + ".png",
			Content: content,
		}
	}

	companyIDStr := b.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}

	if err := handler.NewHttpClient().
		Method("PATCH").
		URL(fmt.Sprintf("/branch/%s/design/images", b.Created.ID.String())).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Header(namespace.HeadersKey.Company, cID).
		Send(fileMap).
		ParseResponse(&b.Created.Design.Images).
		Error; err != nil {
		return fmt.Errorf("failed to upload branch images: %w", err)
	}

	return nil
}

func (b *Branch) DeleteImages(status int, image_types []string, x_auth_token string, x_company_id *string) error {
	if len(image_types) == 0 {
		return fmt.Errorf("no images provided to delete")
	}

	createdCompanyID := b.Company.Created.ID.String()
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

	base_url := fmt.Sprintf("/branch/%s/design/images", b.Created.ID.String())
	for _, image_type := range image_types {
		image_url := base_url + "/" + image_type
		http.URL(image_url)
		http.Send(nil)
		http.ParseResponse(&b.Created.Design.Images)
		if http.Error != nil {
			return fmt.Errorf("failed to delete image %s: %w", image_type, http.Error)
		}
		url := b.Created.Design.Images.GetImageURL(image_type)
		if url != "" {
			return fmt.Errorf("image %s was not deleted successfully, expected empty URL but got %s", image_type, url)
		}
	}
	return nil
}

func (b *Branch) GetImage(status int, imageURL string, compareImgBytes *[]byte) error {
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

func (b *Branch) AddService(status int, service *Service, x_auth_token string, x_company_id *string) error {
	companyIDStr := b.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
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
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	var updated *model.BranchWorkSchedule
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
	b.Created.BranchWorkSchedule = updated.WorkRanges
	return nil
}

func (b *Branch) UpdateWorkRange(status int, wrID string, changes map[string]any, x_auth_token string, x_company_id *string) error {
	companyIDStr := b.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	var updated *model.BranchWorkSchedule
	if err := handler.NewHttpClient().
		Method("PUT").
		URL(fmt.Sprintf("/branch/%s/work_range/%s", b.Created.ID.String(), wrID)).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Company, cID).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Send(changes).
		ParseResponse(&updated).
		Error; err != nil {
		return fmt.Errorf("failed to update branch work range: %w", err)
	}
	b.Created.BranchWorkSchedule = updated.WorkRanges
	return nil
}

func (b *Branch) DeleteWorkRange(status int, wrID string, x_auth_token string, x_company_id *string) error {
	companyIDStr := b.Company.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	var updated *model.BranchWorkSchedule
	if err := handler.NewHttpClient().
		Method("DELETE").
		URL(fmt.Sprintf("/branch/%s/work_range/%s", b.Created.ID.String(), wrID)).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Company, cID).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Send(nil).
		ParseResponse(&updated).
		Error; err != nil {
		return fmt.Errorf("failed to delete branch work schedule: %w", err)
	}
	b.Created.BranchWorkSchedule = updated.WorkRanges
	return nil
}

func GetExampleBranchWorkSchedule(branchID uuid.UUID, servicesID []DTO.ServiceID) DTO.CreateBranchWorkSchedule {
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
