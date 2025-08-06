package e2e_test

import (
	"fmt"
	"mynute-go/core"

	DTO "mynute-go/core/config/api/dto"
	"mynute-go/core/config/db/model"
	"mynute-go/core/lib"
	FileBytes "mynute-go/core/lib/file_bytes"
	handlerT "mynute-go/core/test/handlers"
	modelT "mynute-go/core/test/models"
	"testing"

	"github.com/google/uuid"
)

func Test_Employee(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handlerT.NewTestErrorHandler(t)
	c := &modelT.Company{}

	tt.Describe("Company setup").Test(c.Set()) // Cria company, employees, branches, services

	employee := c.Employees[0]
	branch := c.Branches[0]
	service := c.Services[0]

	tt.Describe("Employee get by ID").Test(employee.GetById(200, nil, nil))
	tt.Describe("Employee get by email").Test(employee.GetByEmail(200, nil, nil))

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

	tt.Describe("Employee update").Test(employee.Update(200, map[string]any{
		"name":     "Should Succeed Update on Employee Name",
		"password": new_password,
	}, nil, nil))

	tt.Describe("Employee update").Test(employee.Update(401, map[string]any{
		"password": "NewPswrd1@!",
	}, nil, nil))

	employee.Created.Password = new_password // Update the password in the employee model
	tt.Describe("Employee get by email").Test(employee.GetByEmail(401, nil, nil))

	tt.Describe("Employee login").Test(employee.Login(200, nil))

	ServicesID := []DTO.ServiceBase{
		{ID: service.Created.ID},
	}

	EmployeeWorkSchedule := modelT.GetExampleEmployeeWorkSchedule(employee.Created.ID, branch.Created.ID, ServicesID)

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

	removeAllServicesFromWorkRange := func(work_range model.EmployeeWorkRange) error {
		for _, service := range work_range.Services {
			if err := employee.RemoveServiceFromWorkRange(200, work_range.ID.String(), service.ID.String(), nil, nil); err != nil {
				return err
			}
		}
		return nil
	}

	tt.Describe("Removing all services from employee work range").Test(removeAllServicesFromWorkRange(wr))

	checkIfAllServicesRemoved := func(work_range model.EmployeeWorkRange) error {
		for _, ewr := range employee.Created.WorkSchedule {
			if ewr.ID == work_range.ID && len(ewr.Services) > 0 {
				return fmt.Errorf("Employee work range %s still has services associated: %v", work_range.ID, ewr.Services)
			}
		}
		return nil
	}

	tt.Describe("Checking if all services were removed from employee work range").Test(checkIfAllServicesRemoved(wr))

	AddAllServicesBackToWorkRange := func(work_range model.EmployeeWorkRange) error {
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

	tt.Describe("Get profile image").Test(employee.GetImage(200, employee.Created.Design.Images.Profile.URL, &FileBytes.PNG_FILE_1))

	tt.Describe("Overwrite profile image").Test(employee.UploadImages(200, map[string][]byte{
		"profile": FileBytes.PNG_FILE_3,
	}, nil, nil))

	tt.Describe("Get overwritten profile image").Test(employee.GetImage(200, employee.Created.Design.Images.Profile.URL, &FileBytes.PNG_FILE_3))

	tt.Describe("Employee deletion").Test(employee.Delete(200, nil, nil))

	tt.Describe("Get deleted employee by ID").Test(employee.GetById(404, &c.Owner.X_Auth_Token, nil))
}
