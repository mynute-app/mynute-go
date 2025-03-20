package e2e_test

import (
	"agenda-kaki-go/core"
	"agenda-kaki-go/core/config/db/model"
	"testing"
	"time"
)

type Appointment struct {
}

func Test_Appointment(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	user := &User{}
	user.Create(t, 200)
	user.VerifyEmail(t, 200)
	user.Login(t, 200)
	user.Update(t, 200, map[string]any{"name": "Updated User Name"})
	user.GetByEmail(t, 200)
	c := &Company{}
	c.Set(t)
	b := c.branches[0]
	e := c.employees[0]
	s := c.services[0]
	for range 3 {
		user.CreateAppointment(t, 200, b, e, s, c, nil)
	}
	startTimeStr := user.created.Appointments[0].StartTime.Format(time.RFC3339)
	user.CreateAppointment(t, 400, b, c.owner, s, c, nil)
	c.owner.AddService(t, 200, s)
	user.CreateAppointment(t, 400, b, c.owner, s, c, nil)
	c.owner.AddBranch(t, 200, b)
	c.owner.Update(t, 200, map[string]any{"work_schedule": []model.WorkSchedule{
		{
			Monday: []model.WorkRange{
				{Start: "08:00", End: "12:00", BranchID: b.created.ID},
				{Start: "13:00", End: "17:00", BranchID: b.created.ID},
			},
			Tuesday: []model.WorkRange{
				{Start: "09:00", End: "12:00", BranchID: b.created.ID},
				{Start: "13:00", End: "18:00", BranchID: b.created.ID},
			},
			Wednesday: []model.WorkRange{
				{Start: "08:00", End: "12:00", BranchID: b.created.ID},
				{Start: "13:00", End: "17:00", BranchID: b.created.ID},
			},
			Thursday: []model.WorkRange{
				{Start: "08:00", End: "12:00", BranchID: b.created.ID},
				{Start: "13:00", End: "17:00", BranchID: b.created.ID},
			},
			Friday: []model.WorkRange{
				{Start: "08:00", End: "12:00", BranchID: b.created.ID},
				{Start: "13:00", End: "17:00", BranchID: b.created.ID},
			},
			Saturday: []model.WorkRange{
				{Start: "08:00", End: "12:00", BranchID: b.created.ID},
				{Start: "13:00", End: "17:00", BranchID: b.created.ID},
			},
		},
	}})
	user.CreateAppointment(t, 200, b, c.owner, s, c, nil)
	user.CreateAppointment(t, 400, b, e, s, c, &startTimeStr)
	user.CreateAppointment(t, 400, b, c.owner, s, c, &startTimeStr)
	user.Delete(t, 200)
}
