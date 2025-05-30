package e2e_test

import (
	"agenda-kaki-go/core"

	mJSON "agenda-kaki-go/core/config/db/model/json"
	models_test "agenda-kaki-go/core/tests/models"
	"testing"
)

// import (
// 	"agenda-kaki-go/core"
// 	"agenda-kaki-go/core/config/db/model"
// 	handler "agenda-kaki-go/core/tests/handlers"
// 	"fmt"
// 	"testing"

// 	"github.com/prometheus/common/server"
// )



func Test_Employee(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	company := &models_test.Company{}
	company.Set(t)
	employee := company.Employees[0]
	employee.GetById(t, 200)
	employee.GetByEmail(t, 200)
	employee.Update(t, 200, map[string]any{"name": "Updated Employee Name xD"})
	employee.Update(t, 200, map[string]any{"work_schedule": []mJSON.WorkSchedule{
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
	}})
	employee.Delete(t, 200)
}