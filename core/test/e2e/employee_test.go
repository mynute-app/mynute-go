package e2e_test

import (
	"agenda-kaki-go/core"

	mJSON "agenda-kaki-go/core/config/db/model/json"
	handlerT "agenda-kaki-go/core/test/handlers"
	modelT "agenda-kaki-go/core/test/models"
	"testing"

	"github.com/google/uuid"
)

func Test_Employee(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handlerT.NewTestErrorHandler(t)
	company := &modelT.Company{}

	tt.Describe("Company setup").Test(company.Set()) // Cria company, employees, branches, services

	employee := company.Employees[0]

	tt.Describe("Employee get by ID").Test(employee.GetById(200, nil, nil))
	tt.Describe("Employee get by email").Test(employee.GetByEmail(200, nil, nil))

	tt.Describe("Employee update").Test(employee.Update(200, map[string]any{
		"name": "Updated Employee Name xD",
	}, nil, nil))

	tt.Describe("Employee update work schedule").Test(employee.UpdateWorkSchedule(200, []mJSON.WorkSchedule{
		{
			Monday: []mJSON.WorkRange{
				{Start: "08:00", End: "12:00", BranchID: company.Branches[0].Created.ID},
				{Start: "13:00", End: "17:00", BranchID: company.Branches[0].Created.ID},
			},
			Tuesday: []mJSON.WorkRange{
				{Start: "09:00", End: "12:00", BranchID: company.Branches[0].Created.ID},
				{Start: "13:00", End: "18:00", BranchID: company.Branches[0].Created.ID},
			},
			Friday: []mJSON.WorkRange{
				{Start: "08:00", End: "12:00", BranchID: company.Branches[0].Created.ID},
				{Start: "13:00", End: "17:00", BranchID: company.Branches[0].Created.ID},
			},
			Saturday: []mJSON.WorkRange{
				{Start: "08:00", End: "12:00", BranchID: company.Branches[0].Created.ID},
				{Start: "13:00", End: "17:00", BranchID: company.Branches[0].Created.ID},
			},
			Sunday: []mJSON.WorkRange{},
		},
	}, nil, nil))

	tt.Describe("Changing employee company_id").Test(employee.Update(400, map[string]any{
		"company_id": uuid.New().String(),
	}, &company.Owner.X_Auth_Token, nil))

	tt.Describe("Employee deletion").Test(employee.Delete(200, nil, nil))
}
