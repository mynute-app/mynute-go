package e2e_test

import (
	"agenda-kaki-go/core"
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/tests/handlers"
	"fmt"
	"testing"
)

type Company struct {
	created    model.Company
	owner      *Employee
	employees  []*Employee
	branches   []*Branch
	services   []*Service
	auth_token string
}

func Test_Company(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	company := &Company{}
	company.Create(t, 200)
	company.owner.VerifyEmail(t, 200)
	company.owner.Login(t, 200)
	company.auth_token = company.owner.auth_token
	company.Update(t, 200, map[string]any{"name": "Updated Company Name"})
	company.GetById(t, 200)
	company.GetByName(t, 200)
	company.Delete(t, 200)
}

func (c *Company) Set(t *testing.T) {
	c.Create(t, 200)
	c.owner.company = c
	c.owner.VerifyEmail(t, 200)
	c.owner.Login(t, 200)
	c.auth_token = c.owner.auth_token
	c.owner.GetById(t, 200)
	employee := &Employee{}
	employee.company = c
	employee.Create(t, 200)
	employee.VerifyEmail(t, 200)
	employee.Login(t, 200)
	c.employees = append(c.employees, employee)
	branch := &Branch{}
	branch.auth_token = c.owner.auth_token
	branch.company = c
	branch.Create(t, 200)
	c.branches = append(c.branches, branch)
	service := &Service{}
	service.auth_token = c.owner.auth_token
	service.company = c
	service.Create(t, 200)
	c.services = append(c.services, service)
	c.GetById(t, 200)
	c.employees[0].AddService(t, 200, c.services[0])
	c.employees[0].AddBranch(t, 200, c.branches[0], &c.owner.auth_token)
	c.branches[0].AddService(t, 200, c.services[0], nil)
	c.employees[0].Update(t, 200, map[string]any{"work_schedule": []model.WorkSchedule{
		{
			Monday: []model.WorkRange{
				{Start: "08:00", End: "12:00", BranchID: c.branches[0].created.ID},
				{Start: "13:00", End: "17:00", BranchID: c.branches[0].created.ID},
			},
			Tuesday: []model.WorkRange{
				{Start: "09:00", End: "12:00", BranchID: c.branches[0].created.ID},
				{Start: "13:00", End: "18:00", BranchID: c.branches[0].created.ID},
			},
			Wednesday: []model.WorkRange{
				{Start: "08:00", End: "12:00", BranchID: c.branches[0].created.ID},
				{Start: "13:00", End: "17:00", BranchID: c.branches[0].created.ID},
			},
			Thursday: []model.WorkRange{
				{Start: "08:00", End: "12:00", BranchID: c.branches[0].created.ID},
				{Start: "13:00", End: "17:00", BranchID: c.branches[0].created.ID},
			},
			Friday: []model.WorkRange{
				{Start: "08:00", End: "12:00", BranchID: c.branches[0].created.ID},
				{Start: "13:00", End: "17:00", BranchID: c.branches[0].created.ID},
			},
			Saturday: []model.WorkRange{
				{Start: "08:00", End: "12:00", BranchID: c.branches[0].created.ID},
				{Start: "13:00", End: "17:00", BranchID: c.branches[0].created.ID},
			},
			Sunday: []model.WorkRange{},
		},
	}})
	
}

func (c *Company) Create(t *testing.T, status int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/company")
	http.ExpectStatus(status)
	http.Header("Authorization", c.auth_token)
	ownerEmail := lib.GenerateRandomEmail("owner")
	ownerPswd := "1SecurePswd!"
	http.Send(DTO.CreateCompany{
		Name:          lib.GenerateRandomName("Company Name"),
		TaxID:         lib.GenerateRandomStrNumber(14),
		OwnerName:     lib.GenerateRandomName("Owner Name"),
		OwnerSurname:  lib.GenerateRandomName("Owner Surname"),
		OwnerEmail:    ownerEmail,
		OwnerPhone:    lib.GenerateRandomPhoneNumber(),
		OwnerPassword: ownerPswd,
	})
	http.ParseResponse(&c.created)
	owner := c.created.Employees[0]
	owner.Password = ownerPswd
	c.owner = &Employee{
		created: owner,
	}
}

func (c *Company) GetByName(t *testing.T, status int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL(fmt.Sprintf("/company/name/%s", c.created.Name))
	http.ExpectStatus(status)
	http.Send(nil)
	http.ParseResponse(&c.created)
}

func (c *Company) GetById(t *testing.T, status int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL(fmt.Sprintf("/company/%s", c.created.ID.String()))
	http.ExpectStatus(status)
	http.Send(nil)
	http.ParseResponse(&c.created)
}

func (c *Company) Update(t *testing.T, status int, changes map[string]any) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("PATCH")
	http.URL(fmt.Sprintf("/company/%s", c.created.ID.String()))
	http.ExpectStatus(status)
	http.Header("Authorization", c.auth_token)
	http.Send(changes)
}

func (c *Company) Delete(t *testing.T, status int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("DELETE")
	http.URL(fmt.Sprintf("/company/%s", c.created.ID.String()))
	http.ExpectStatus(status)
	http.Header("Authorization", c.auth_token)
	http.Send(nil)
}
