package e2e_test

import (
	"fmt"

	"mynute-go/core"
	DTO "mynute-go/core/src/config/api/dto"
	coreModel "mynute-go/core/src/config/db/model"
	"mynute-go/core/src/lib"
	FileBytes "mynute-go/core/src/lib/file_bytes"
	"mynute-go/test/src/handler"
	"mynute-go/test/src/model"
	"testing"

	"github.com/google/uuid"
)

func Test_Employee(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handler.NewTestErrorHandler(t)
	c := &model.Company{}

	tt.Describe("Company setup").Test(c.Set()) // Cria company, employees, branches, services

	employee := c.Employees[0]
	branch := c.Branches[0]
	service := c.Services[0]

	tt.Describe("Changing employee id").Test(employee.Update(400, map[string]any{
		"Id":   "00000000-0000-0000-0000-000000000001",
		"name": "Updated Employee Name xDDDD",
	}, nil, nil))

	tt.Describe("Changing employee company_id").Test(employee.Update(400, map[string]any{
		"company_id": uuid.New().String(),
	}, &c.Owner.X_Auth_Token, nil))

	tt.Describe("Employee update").Test(employee.Update(200, map[string]any{
		"name": "Updated Employee Name xDDDD",
	}, nil, nil))

	tt.Describe("Employee fail to update").Test(employee.Update(400, map[string]any{
		"name":     "Should Fail Update on Employee Name",
		"password": "newpswrd123",
	}, nil, nil))

	new_password := lib.GenerateValidPassword()

	tt.Describe("Employee update password").Test(employee.Update(200, map[string]any{
		"name":     "Should Succeed Update on Employee Name",
		"password": new_password,
	}, nil, nil))

	// Test that old password fails after password change
	tt.Describe("Employee fail to login with old password").Test(employee.LoginByPassword(401, employee.Created.Password, nil))

	// Update password in memory and test new password works
	employee.Created.Password = new_password
	tt.Describe("Employee login with new password").Test(employee.LoginByPassword(200, new_password, nil))

	// Test password reset by email
	tt.Describe("Reset employee password by email").Test(employee.ResetPasswordByEmail(200, nil))

	// Test that the old password no longer works
	tt.Describe("Employee login with old password fails").Test(employee.LoginByPassword(401, new_password, nil))

	// Test that new password from email works
	tt.Describe("Employee login with password from email").Test(employee.LoginByPassword(200, employee.Created.Password, nil))

	ServicesID := []DTO.ServiceBase{
		{ID: service.Created.ID},
	}

	EmployeeWorkSchedule := model.GetExampleEmployeeWorkSchedule(employee.Created.ID, branch.Created.ID, ServicesID)

	tt.Describe("Employee create work schedule incorrectly").Test(employee.CreateWorkSchedule(400, EmployeeWorkSchedule, nil, nil))

	tt.Describe("Add service to employee").Test(employee.AddService(200, service, &c.Owner.X_Auth_Token, nil))

	tt.Describe("Employee create work schedule incorrectly").Test(employee.CreateWorkSchedule(400, EmployeeWorkSchedule, nil, nil))

	tt.Describe("Add branch to employee").Test(employee.AddBranch(200, branch, &c.Owner.X_Auth_Token, nil))

	tt.Describe("Employee create work schedule successfully").Test(employee.CreateWorkSchedule(200, EmployeeWorkSchedule, nil, nil))

	tt.Describe("Get Employee work schedule successfully").Test(employee.GetWorkSchedule(200, nil, nil))

	wr := employee.Created.WorkSchedule[0]

	tt.Describe("Updating fail branch work schedule").Test(employee.UpdateWorkRange(400, wr.ID.String(), map[string]any{
		"start_time": "06:00",
		"end_time":   "20:00",
		"time_zone":  "America/Sao_Paulo",
	}, nil, nil))

	tt.Describe("Updating success branch work schedule").Test(employee.UpdateWorkRange(400, wr.ID.String(), map[string]any{
		"start_time": "09:00",
		"end_time":   "18:00",
		"time_zone":  "America/Sao_Paulo",
		"weekday":    1,
	}, nil, nil))

	tt.Describe("Updating success branch work schedule").Test(employee.UpdateWorkRange(200, wr.ID.String(), map[string]any{
		"start_time": "09:30",
		"end_time":   "11:00",
		"time_zone":  "America/Sao_Paulo",
		"weekday":    1,
	}, nil, nil))

	removeAllServicesFromWorkRange := func(work_range coreModel.EmployeeWorkRange) error {
		for _, service := range work_range.Services {
			if err := employee.RemoveServiceFromWorkRange(200, work_range.ID.String(), service.ID.String(), nil, nil); err != nil {
				return err
			}
		}
		return nil
	}

	tt.Describe("Removing all services from employee work range").Test(removeAllServicesFromWorkRange(wr))

	checkIfAllServicesRemoved := func(work_range coreModel.EmployeeWorkRange) error {
		for _, ewr := range employee.Created.WorkSchedule {
			if ewr.ID == work_range.ID && len(ewr.Services) > 0 {
				return fmt.Errorf("Employee work range %s still has services associated: %v", work_range.ID, ewr.Services)
			}
		}
		return nil
	}

	tt.Describe("Checking if all services were removed from employee work range").Test(checkIfAllServicesRemoved(wr))

	AddAllServicesBackToWorkRange := func(work_range coreModel.EmployeeWorkRange) error {
		var services DTO.EmployeeWorkRangeServices
		for _, service := range work_range.Services {
			services.Services = append(services.Services, DTO.ServiceBase{ID: service.ID})
		}
		return employee.AddServicesToWorkRange(200, work_range.ID.String(), services, nil, nil)
	}

	tt.Describe("Adding all services back to employee work range").Test(AddAllServicesBackToWorkRange(wr))

	wrService := wr.Services[0]

	tt.Describe("Add the same service again to employee work range").Test(employee.AddServicesToWorkRange(200, wr.ID.String(), DTO.EmployeeWorkRangeServices{
		Services: []DTO.ServiceBase{{ID: wrService.ID}},
	}, nil, nil))

	tt.Describe("Check if the number of services in employee work range is still the same").Test(func() error {
		if len(employee.Created.WorkSchedule[0].Services) != len(wr.Services) {
			return fmt.Errorf("Expected %d services, got %d", len(wr.Services), len(employee.Created.WorkSchedule[0].Services))
		}
		return nil
	}())

	tt.Describe("Deleting branch work range").Test(employee.DeleteWorkRange(200, wr.ID.String(), nil, nil))

	tt.Describe("Upload profile image").Test(employee.UploadImages(200, map[string][]byte{
		"profile": FileBytes.PNG_FILE_1,
	}, nil, nil))

	tt.Describe("Get profile image").Test(employee.GetImage(200, employee.Created.Meta.Design.Images.Profile.URL, &FileBytes.PNG_FILE_1))

	tt.Describe("Overwrite profile image").Test(employee.UploadImages(200, map[string][]byte{
		"profile": FileBytes.PNG_FILE_3,
	}, nil, nil))

	tt.Describe("Get overwritten profile image").Test(employee.GetImage(200, employee.Created.Meta.Design.Images.Profile.URL, &FileBytes.PNG_FILE_3))

	// Test that GET employee returns images in meta.design.images
	tt.Describe("Get employee by ID and verify images are returned").Test(func() error {
		if err := employee.GetById(200, &c.Owner.X_Auth_Token, nil); err != nil {
			return err
		}
		if employee.Created.Meta.Design.Images.Profile.URL == "" {
			return fmt.Errorf("Expected profile image URL to be returned in GET response, but got empty string")
		}
		return nil
	}())

	tt.Describe("Get employee by email and verify images are returned").Test(func() error {
		if err := employee.GetByEmail(200, &c.Owner.X_Auth_Token, nil); err != nil {
			return err
		}
		if employee.Created.Meta.Design.Images.Profile.URL == "" {
			return fmt.Errorf("Expected profile image URL to be returned in GET response, but got empty string")
		}
		return nil
	}())

	// Upload multiple images to test all image types
	tt.Describe("Upload multiple images (logo, banner, background)").Test(employee.UploadImages(200, map[string][]byte{
		"logo":       FileBytes.PNG_FILE_1,
		"banner":     FileBytes.PNG_FILE_2,
		"background": FileBytes.PNG_FILE_3,
	}, nil, nil))

	// Verify all images are returned in GET response
	tt.Describe("Get employee and verify all uploaded images are returned").Test(func() error {
		if err := employee.GetById(200, &c.Owner.X_Auth_Token, nil); err != nil {
			return err
		}
		images := employee.Created.Meta.Design.Images
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

	tt.Describe("Employee deletion").Test(employee.Delete(200, nil, nil))

	tt.Describe("Get deleted employee by ID").Test(employee.GetById(404, &c.Owner.X_Auth_Token, nil))
}

func Test_EmployeeAppointments(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handler.NewTestErrorHandler(t)
	c := &model.Company{}

	tt.Describe("Company setup").Test(c.Set())

	employee := c.Employees[0]

	// Test employee appointments endpoint
	t.Run("Employee Appointments", func(t *testing.T) {
		// Test getting appointments when employee has no appointments
		tt.Describe("Get empty appointments list").Test(func() error {
			appointmentList, err := employee.GetAppointments(200, 1, 10, "", "", "", "UTC", nil, nil)
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

		// Test pagination with different page sizes
		tt.Describe("Get empty appointments with page size 5").Test(func() error {
			appointmentList, err := employee.GetAppointments(200, 1, 5, "", "", "", "UTC", nil, nil)
			if err != nil {
				return err
			}
			if len(appointmentList.Appointments) != 0 {
				return fmt.Errorf("expected 0 appointments, got %d", len(appointmentList.Appointments))
			}
			if appointmentList.PageSize != 5 {
				return fmt.Errorf("expected page size 5, got %d", appointmentList.PageSize)
			}
			return nil
		}())

		// Test different pagination parameters
		tt.Describe("Get appointments page 2 (empty)").Test(func() error {
			appointmentList, err := employee.GetAppointments(200, 2, 10, "", "", "", "UTC", nil, nil)
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
			_, err := employee.GetAppointments(400, 1, 10, "", "", "", "", nil, nil)
			if err != nil {
				return err
			}
			return nil
		}())

		// Test date range filtering - invalid date format
		tt.Describe("Test invalid start_date format").Test(func() error {
			_, err := employee.GetAppointments(400, 1, 10, "2025-04-21", "", "", "UTC", nil, nil)
			if err != nil {
				return err
			}
			return nil
		}())

		// Test date range filtering - invalid end_date format
		tt.Describe("Test invalid end_date format").Test(func() error {
			_, err := employee.GetAppointments(400, 1, 10, "", "04/21/2025", "", "UTC", nil, nil)
			if err != nil {
				return err
			}
			return nil
		}())

		// Test date range filtering - valid date range
		tt.Describe("Test valid date range filter").Test(func() error {
			appointmentList, err := employee.GetAppointments(200, 1, 10, "21/04/2025", "31/05/2025", "", "UTC", nil, nil)
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
			_, err := employee.GetAppointments(400, 1, 10, "01/01/2025", "15/04/2025", "", "UTC", nil, nil)
			if err != nil {
				return err
			}
			return nil
		}())

		// Test date range filtering - end date before start date
		tt.Describe("Test end_date before start_date").Test(func() error {
			_, err := employee.GetAppointments(400, 1, 10, "31/05/2025", "21/04/2025", "", "UTC", nil, nil)
			if err != nil {
				return err
			}
			return nil
		}())

		// Test date range filtering - exactly 90 days
		tt.Describe("Test exactly 90 days range").Test(func() error {
			appointmentList, err := employee.GetAppointments(200, 1, 10, "01/01/2025", "31/03/2025", "", "UTC", nil, nil)
			if err != nil {
				return err
			}
			if len(appointmentList.Appointments) != 0 {
				return fmt.Errorf("expected 0 appointments with 90-day filter, got %d", len(appointmentList.Appointments))
			}
			return nil
		}())

		// Test cancelled filter - true
		tt.Describe("Test cancelled filter true").Test(func() error {
			appointmentList, err := employee.GetAppointments(200, 1, 10, "", "", "true", "UTC", nil, nil)
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
			appointmentList, err := employee.GetAppointments(200, 1, 10, "", "", "false", "UTC", nil, nil)
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
			_, err := employee.GetAppointments(400, 1, 10, "", "", "invalid", "UTC", nil, nil)
			if err != nil {
				return err
			}
			return nil
		}())

		// Test combined filters - date range and cancelled status
		tt.Describe("Test combined date range and cancelled filters").Test(func() error {
			appointmentList, err := employee.GetAppointments(200, 1, 10, "21/04/2025", "31/05/2025", "false", "UTC", nil, nil)
			if err != nil {
				return err
			}
			if len(appointmentList.Appointments) != 0 {
				return fmt.Errorf("expected 0 appointments with combined filters, got %d", len(appointmentList.Appointments))
			}
			return nil
		}())

		// Test with only start_date (no end_date)
		tt.Describe("Test with only start_date filter").Test(func() error {
			appointmentList, err := employee.GetAppointments(200, 1, 10, "21/04/2025", "", "", "UTC", nil, nil)
			if err != nil {
				return err
			}
			if len(appointmentList.Appointments) != 0 {
				return fmt.Errorf("expected 0 appointments with start_date filter, got %d", len(appointmentList.Appointments))
			}
			return nil
		}())

		// Test with only end_date (no start_date)
		tt.Describe("Test with only end_date filter").Test(func() error {
			appointmentList, err := employee.GetAppointments(200, 1, 10, "", "31/05/2025", "", "UTC", nil, nil)
			if err != nil {
				return err
			}
			if len(appointmentList.Appointments) != 0 {
				return fmt.Errorf("expected 0 appointments with end_date filter, got %d", len(appointmentList.Appointments))
			}
			return nil
		}())

		// Test for non-existent employee UUID validation
		tt.Describe("Test non-existent employee returns 404").Test(func() error {
			nonExistentEmployee := &model.Employee{
				Created: &coreModel.Employee{},
				Company: c,
			}
			nonExistentEmployee.Created.ID = uuid.New()
			nonExistentEmployee.X_Auth_Token = c.Owner.X_Auth_Token

			_, err := nonExistentEmployee.GetAppointments(404, 1, 10, "", "", "", "UTC", &c.Owner.X_Auth_Token, nil)
			if err != nil {
				return err
			}
			return nil
		}())
	})
}
