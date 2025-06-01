package e2e_test

import (
	"agenda-kaki-go/core"

	mJSON "agenda-kaki-go/core/config/db/model/json"
	handlerT "agenda-kaki-go/core/test/handlers"
	modelT "agenda-kaki-go/core/test/models"
	"testing"
)


func Test_Employee(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	company := &modelT.Company{}
	tt := handlerT.NewTestErrorHandler(t)
	tt.Test(company.Set())
	employee := company.Employees[0]
	tt.Test(employee.GetById(200))
	tt.Test(employee.GetByEmail(200))
	tt.Test(employee.Update(200, map[string]any{"name": "Updated Employee Name xD"}))
	tt.Test(employee.Update(200, map[string]any{"work_schedule": []mJSON.WorkSchedule{
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
	}}))
	employee.Delete(200)
}
