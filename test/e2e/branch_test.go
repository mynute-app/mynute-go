package e2e_test

import (
	"fmt"

	"mynute-go/core"
	DTO "mynute-go/core/src/config/api/dto"
	coreModel "mynute-go/core/src/config/db/model"
	FileBytes "mynute-go/core/src/lib/file_bytes"
	"mynute-go/test/src/handler"
	"mynute-go/test/src/model"

	"testing"

	"github.com/google/uuid"
)

func Test_Branch(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handler.NewTestErrorHandler(t)

	company := &model.Company{}
	tt.Describe("Company creation").Test(company.Create(200))

	branch := &model.Branch{}
	branch.Company = company
	tt.Describe("Branch creation").Test(branch.Create(200, company.Owner.X_Auth_Token, nil))

	tt.Describe("Branch update").Test(branch.Update(200, map[string]any{
		"name": branch.Created.Name,
	}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Branch get by ID").Test(branch.GetById(200, company.Owner.X_Auth_Token, nil))
	// tt.Describe("Branch get by name").Test(branch.GetByName(200, company.Owner.X_Auth_Token, nil))

	service := &model.Service{}
	service.Company = company

	tt.Describe("Service creation").Test(service.Create(200, company.Owner.X_Auth_Token, nil))
	servicesID := []DTO.ServiceBase{{ID: service.Created.ID}}
	BranchWorkSchedule := model.GetExampleBranchWorkSchedule(branch.Created.ID, servicesID)
	tt.Describe("Branch work schedule fail creation").Test(branch.CreateWorkSchedule(400, BranchWorkSchedule, company.Owner.X_Auth_Token, nil))
	tt.Describe("Adding service to branch").Test(branch.AddService(200, service, company.Owner.X_Auth_Token, nil))
	tt.Describe("Branch work schedule success creation").Test(branch.CreateWorkSchedule(200, BranchWorkSchedule, company.Owner.X_Auth_Token, nil))
	tt.Describe("Get Branch work schedule success").Test(branch.GetWorkSchedule(200, "", nil))

	wr := branch.Created.WorkSchedule[0]
	tt.Describe("Updating fail branch work schedule").Test(branch.UpdateWorkRange(400, wr.ID.String(), map[string]any{
		"start_time": "06:00",
		"end_time":   "20:00",
		"time_zone":  "America/Sao_Paulo",
	}, company.Owner.X_Auth_Token, nil))
	tt.Describe("Updating success branch work schedule").Test(branch.UpdateWorkRange(200, wr.ID.String(), map[string]any{
		"start_time": "06:00",
		"end_time":   "20:00",
		"time_zone":  "America/Sao_Paulo",
		"weekday":    1,
	}, company.Owner.X_Auth_Token, nil))

	removeAllServicesFromWorkRange := func(work_range coreModel.BranchWorkRange) error {
		for _, service := range work_range.Services {
			if err := branch.RemoveServiceFromWorkRange(200, work_range.ID.String(), service.ID.String(), company.Owner.X_Auth_Token, nil); err != nil {
				return err
			}
		}
		return nil
	}

	tt.Describe("Removing all service from branch work range").Test(removeAllServicesFromWorkRange(wr))

	checkIfAllServicesRemoved := func(work_range coreModel.BranchWorkRange) error {
		for _, bwr := range branch.Created.WorkSchedule {
			if bwr.ID == work_range.ID && len(bwr.Services) > 0 {
				return fmt.Errorf("Branch work range %s still has services associated: %v", work_range.ID, bwr.Services)
			}
		}
		return nil
	}

	tt.Describe("Checking if all services were removed from branch work range").Test(checkIfAllServicesRemoved(wr))

	AddAllServicesBackToWorkRange := func(work_range coreModel.BranchWorkRange) error {
		var services DTO.BranchWorkRangeServices
		for _, service := range work_range.Services {
			services.Services = append(services.Services, DTO.ServiceBase{ID: service.ID})
		}
		return branch.AddServicesToWorkRange(200, work_range.ID.String(), services, company.Owner.X_Auth_Token, nil)
	}

	tt.Describe("Adding all services back to branch work range").Test(AddAllServicesBackToWorkRange(wr))

	wrService := wr.Services[0]

	tt.Describe("Add the same service again to branch work range").Test(branch.AddServicesToWorkRange(200, wr.ID.String(), DTO.BranchWorkRangeServices{
		Services: []DTO.ServiceBase{{ID: wrService.ID}},
	}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Check if the number of services in branch work range is still the same").Test(func() error {
		if len(branch.Created.WorkSchedule[0].Services) != len(wr.Services) {
			return fmt.Errorf("Expected %d services, got %d", len(wr.Services), len(branch.Created.WorkSchedule[0].Services))
		}
		return nil
	}())

	tt.Describe("Deleting branch work range").Test(branch.DeleteWorkRange(200, wr.ID.String(), company.Owner.X_Auth_Token, nil))
	tt.Describe("Adding branch to Owner").Test(company.Owner.AddBranch(200, branch, nil, nil))
	tt.Describe("Getting company by ID").Test(company.GetById(200, company.Owner.X_Auth_Token, nil))

	tt.Describe("Changing branch company_id").Test(branch.Update(400, map[string]any{
		"company_id": uuid.New().String(),
	}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Upload profile image").Test(branch.UploadImages(200, map[string][]byte{
		"profile": FileBytes.PNG_FILE_1,
	}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Get profile image").Test(branch.GetImage(200, branch.Created.Design.Images.Profile.URL, &FileBytes.PNG_FILE_1))

	tt.Describe("Overwrite profile image").Test(branch.UploadImages(200, map[string][]byte{
		"profile": FileBytes.PNG_FILE_3,
	}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Get overwritten profile image").Test(branch.GetImage(200, branch.Created.Design.Images.Profile.URL, &FileBytes.PNG_FILE_3))

	// Test that GET branch returns images in design.images
	tt.Describe("Get branch by ID and verify images are returned").Test(func() error {
		if err := branch.GetById(200, company.Owner.X_Auth_Token, nil); err != nil {
			return err
		}
		if branch.Created.Design.Images.Profile.URL == "" {
			return fmt.Errorf("Expected profile image URL to be returned in GET response, but got empty string")
		}
		return nil
	}())

	// Upload multiple images to test all image types
	tt.Describe("Upload multiple images (logo, banner, background)").Test(branch.UploadImages(200, map[string][]byte{
		"logo":       FileBytes.PNG_FILE_1,
		"banner":     FileBytes.PNG_FILE_2,
		"background": FileBytes.PNG_FILE_3,
	}, company.Owner.X_Auth_Token, nil))

	// Verify all images are returned in GET response
	tt.Describe("Get branch and verify all uploaded images are returned").Test(func() error {
		if err := branch.GetById(200, company.Owner.X_Auth_Token, nil); err != nil {
			return err
		}
		images := branch.Created.Design.Images
		if images.Profile.URL == "" {
			return fmt.Errorf("Expected profile image URL in response")
		}
		if images.Logo.URL == "" {
			return fmt.Errorf("Expected logo image URL in response")
		}
		if images.Banner.URL == "" {
			return fmt.Errorf("Expected banner image URL in response")
		}
		if images.Background.URL == "" {
			return fmt.Errorf("Expected background image URL in response")
		}
		return nil
	}())

	img_url := branch.Created.Design.Images.Profile.URL

	tt.Describe("Delete profile image").Test(branch.DeleteImages(200, []string{"profile"}, company.Owner.X_Auth_Token, nil))

	tt.Describe("Get deleted profile image").Test(branch.GetImage(404, img_url, nil))

	// Test branch appointments endpoint
	t.Run("Branch Appointments", func(t *testing.T) {
		// Use the main test branch for appointments endpoint testing
		// This branch already exists and has a work schedule from the main test

		// Create client for appointments
		client := &model.Client{}
		tt.Describe("Client creation").Test(client.Set())

		// Test getting appointments when branch has no appointments
		tt.Describe("Get empty appointments list").Test(func() error {
			appointmentList, err := branch.GetAppointments(200, 1, 10, "", "", "", "UTC", company.Owner.X_Auth_Token, nil)
			if err != nil {
				return err
			}
			if len(appointmentList.Appointments) != 0 {
				return fmt.Errorf("expected 0 appointments, got %d", len(appointmentList.Appointments))
			}
			if appointmentList.TotalCount != 0 {
				return fmt.Errorf("expected total count 0, got %d", appointmentList.TotalCount)
			}
			if appointmentList.Page != 1 {
				return fmt.Errorf("expected page 1, got %d", appointmentList.Page)
			}
			if appointmentList.PageSize != 10 {
				return fmt.Errorf("expected page size 10, got %d", appointmentList.PageSize)
			}
			if len(appointmentList.ClientInfo) != 0 {
				return fmt.Errorf("expected 0 clients in ClientInfo, got %d", len(appointmentList.ClientInfo))
			}
			return nil
		}())

		// Test pagination with different page sizes (still empty)
		tt.Describe("Get empty appointments with page size 5").Test(func() error {
			appointmentList, err := branch.GetAppointments(200, 1, 5, "", "", "", "UTC", company.Owner.X_Auth_Token, nil)
			if err != nil {
				return err
			}
			if len(appointmentList.Appointments) != 0 {
				return fmt.Errorf("expected 0 appointments, got %d", len(appointmentList.Appointments))
			}
			if appointmentList.TotalCount != 0 {
				return fmt.Errorf("expected total count 0, got %d", appointmentList.TotalCount)
			}
			if appointmentList.PageSize != 5 {
				return fmt.Errorf("expected page size 5, got %d", appointmentList.PageSize)
			}
			return nil
		}())

		// Test different pagination parameters
		tt.Describe("Get appointments page 2 (empty)").Test(func() error {
			appointmentList, err := branch.GetAppointments(200, 2, 10, "", "", "", "UTC", company.Owner.X_Auth_Token, nil)
			if err != nil {
				return err
			}
			if len(appointmentList.Appointments) != 0 {
				return fmt.Errorf("expected 0 appointments on page 2, got %d", len(appointmentList.Appointments))
			}
			if appointmentList.Page != 2 {
				return fmt.Errorf("expected page 2, got %d", appointmentList.Page)
			}
			return nil
		}())

		// Test missing timezone parameter
		tt.Describe("Test missing timezone parameter").Test(func() error {
			_, err := branch.GetAppointments(400, 1, 10, "", "", "", "", company.Owner.X_Auth_Token, nil)
			if err != nil {
				return err
			}
			return nil
		}())

		// Test date range filtering - invalid date format
		tt.Describe("Test invalid start_date format").Test(func() error {
			_, err := branch.GetAppointments(400, 1, 10, "2025-04-21", "", "", "UTC", company.Owner.X_Auth_Token, nil)
			if err != nil {
				return err
			}
			return nil
		}())

		// Test date range filtering - valid date range
		tt.Describe("Test valid date range filter").Test(func() error {
			appointmentList, err := branch.GetAppointments(200, 1, 10, "21/04/2025", "31/05/2025", "", "UTC", company.Owner.X_Auth_Token, nil)
			if err != nil {
				return err
			}
			if len(appointmentList.Appointments) != 0 {
				return fmt.Errorf("expected 0 appointments with date filter, got %d", len(appointmentList.Appointments))
			}
			return nil
		}())

		// Test date range filtering - range exceeds 90 days
		tt.Describe("Test date range exceeds 90 days").Test(func() error {
			_, err := branch.GetAppointments(400, 1, 10, "01/01/2025", "15/04/2025", "", "UTC", company.Owner.X_Auth_Token, nil)
			if err != nil {
				return err
			}
			return nil
		}())

		// Test date range filtering - end date before start date
		tt.Describe("Test end_date before start_date").Test(func() error {
			_, err := branch.GetAppointments(400, 1, 10, "31/05/2025", "21/04/2025", "", "UTC", company.Owner.X_Auth_Token, nil)
			if err != nil {
				return err
			}
			return nil
		}())

		// Test cancelled filter - true
		tt.Describe("Test cancelled filter true").Test(func() error {
			appointmentList, err := branch.GetAppointments(200, 1, 10, "", "", "true", "UTC", company.Owner.X_Auth_Token, nil)
			if err != nil {
				return err
			}
			if len(appointmentList.Appointments) != 0 {
				return fmt.Errorf("expected 0 cancelled appointments, got %d", len(appointmentList.Appointments))
			}
			return nil
		}())

		// Test cancelled filter - false
		tt.Describe("Test cancelled filter false").Test(func() error {
			appointmentList, err := branch.GetAppointments(200, 1, 10, "", "", "false", "UTC", company.Owner.X_Auth_Token, nil)
			if err != nil {
				return err
			}
			if len(appointmentList.Appointments) != 0 {
				return fmt.Errorf("expected 0 non-cancelled appointments, got %d", len(appointmentList.Appointments))
			}
			return nil
		}())

		// Test cancelled filter - invalid value
		tt.Describe("Test invalid cancelled filter value").Test(func() error {
			_, err := branch.GetAppointments(400, 1, 10, "", "", "invalid", "UTC", company.Owner.X_Auth_Token, nil)
			if err != nil {
				return err
			}
			return nil
		}())

		// Note: Testing 404 for non-existent branch is complex due to company schema isolation
		// The endpoint correctly returns 404 for branches in other company schemas
	})

	tt.Describe("Branch deletion").Test(branch.Delete(200, company.Owner.X_Auth_Token, nil))

	tt.Describe("Get deleted branch by ID").Test(branch.GetById(404, company.Owner.X_Auth_Token, nil))
}
