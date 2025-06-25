package e2e_test

import (
	"agenda-kaki-go/core"

	DTO "agenda-kaki-go/core/config/api/dto"
	handlerT "agenda-kaki-go/core/test/handlers"
	modelT "agenda-kaki-go/core/test/models"
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

	tt.Describe("Employee get by ID").Test(employee.GetById(200, nil, nil))
	tt.Describe("Employee get by email").Test(employee.GetByEmail(200, nil, nil))

	tt.Describe("Employee update").Test(employee.Update(200, map[string]any{
		"name": "Updated Employee Name xD",
	}, nil, nil))
	ServicesID := []DTO.ServiceID{
		{ID: c.Services[0].Created.ID},
	}
	tt.Describe("Employee update work schedule").Test(employee.CreateWorkSchedule(200, DTO.CreateEmployeeWorkSchedule{
		WorkRanges: []DTO.CreateEmployeeWorkRange{
			{
				EmployeeID: employee.Created.ID,
				BranchID:   c.Branches[0].Created.ID,
				Weekday:    1,
				StartTime:  "08:00",
				EndTime:    "12:00",
				TimeZone:   c.Branches[0].Created.TimeZone,
				Services:   ServicesID,
			},
			{
				EmployeeID: employee.Created.ID,
				BranchID:   c.Branches[0].Created.ID,
				Weekday:    1,
				StartTime:  "13:00",
				EndTime:    "17:00",
				TimeZone:   c.Branches[0].Created.TimeZone,
				Services:   ServicesID,
			},
			{
				EmployeeID: employee.Created.ID,
				BranchID:   c.Branches[0].Created.ID,
				Weekday:    2,
				StartTime:  "08:00",
				EndTime:    "12:00",
				TimeZone:   c.Branches[0].Created.TimeZone,
				Services:   ServicesID,
			},
			{
				EmployeeID: employee.Created.ID,
				BranchID:   c.Branches[0].Created.ID,
				Weekday:    2,
				StartTime:  "13:00",
				EndTime:    "17:00",
				TimeZone:   c.Branches[0].Created.TimeZone,
				Services:   ServicesID,
			},
			{
				EmployeeID: employee.Created.ID,
				BranchID:   c.Branches[0].Created.ID,
				Weekday:    3,
				StartTime:  "08:00",
				EndTime:    "12:00",
				TimeZone:   c.Branches[0].Created.TimeZone,
				Services:   ServicesID,
			},
			{
				EmployeeID: employee.Created.ID,
				BranchID:   c.Branches[0].Created.ID,
				Weekday:    3,
				StartTime:  "13:00",
				EndTime:    "17:00",
				TimeZone:   c.Branches[0].Created.TimeZone,
				Services:   ServicesID,
			},
			{
				EmployeeID: employee.Created.ID,
				BranchID:   c.Branches[0].Created.ID,
				Weekday:    4,
				StartTime:  "08:00",
				EndTime:    "12:00",
				TimeZone:   c.Branches[0].Created.TimeZone,
				Services:   ServicesID,
			},
			{
				EmployeeID: employee.Created.ID,
				BranchID:   c.Branches[0].Created.ID,
				Weekday:    4,
				StartTime:  "13:00",
				EndTime:    "17:00",
				TimeZone:   c.Branches[0].Created.TimeZone,
				Services:   ServicesID,
			},
			{
				EmployeeID: employee.Created.ID,
				BranchID:   c.Branches[0].Created.ID,
				Weekday:    5,
				StartTime:  "08:00",
				EndTime:    "12:00",
				TimeZone:   c.Branches[0].Created.TimeZone,
				Services:   ServicesID,
			},
			{
				EmployeeID: employee.Created.ID,
				BranchID:   c.Branches[0].Created.ID,
				Weekday:    5,
				StartTime:  "13:00",
				EndTime:    "17:00",
				TimeZone:   c.Branches[0].Created.TimeZone,
				Services:   ServicesID,
			},
			{
				EmployeeID: employee.Created.ID,
				BranchID:   c.Branches[0].Created.ID,
				Weekday:    6,
				StartTime:  "08:00",
				EndTime:    "12:00",
				TimeZone:   c.Branches[0].Created.TimeZone,
				Services:   ServicesID,
			},
			{
				EmployeeID: employee.Created.ID,
				BranchID:   c.Branches[0].Created.ID,
				Weekday:    6,
				StartTime:  "13:00",
				EndTime:    "17:00",
				TimeZone:   c.Branches[0].Created.TimeZone,
				Services:   ServicesID,
			},
		},
	}, nil, nil))

	tt.Describe("Changing employee company_id").Test(employee.Update(400, map[string]any{
		"company_id": uuid.New().String(),
	}, &c.Owner.X_Auth_Token, nil))

	tt.Describe("Employee deletion").Test(employee.Delete(200, nil, nil))
}
