package e2e_test

import (
	"agenda-kaki-go/core"
	"agenda-kaki-go/core/config/db/model"
	"testing"
	"time"
)

type Appointment struct {
}

// As the appointment dates are being generated randomly,
// this tests can fail sometimes. When it fails, take a look
// closely to check if the error is related to the appointment date
// being in conflict with another appointment.
// If so, just run the test again and it should pass.
// If the error is not related to the appointment date, then fix it.
func Test_Appointment(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	client := &Client{}
	client.Create(t, 200)
	client.VerifyEmail(t, 200)
	client.Login(t, 200)
	client.Update(t, 200, map[string]any{"name": "Updated Client Name"})
	client.GetByEmail(t, 200)
	c := &Company{}
	c.Set(t)
	b := c.branches[0]
	e := c.employees[0]
	s := c.services[0]
	client.CreateAppointment(t, 200, b, e, s, c, nil)
	startTimeStr := client.created.Appointments[0].StartTime.Format(time.RFC3339)
	client.CreateAppointment(t, 400, b, c.owner, s, c, nil)
	c.owner.AddService(t, 200, s)
	client.CreateAppointment(t, 400, b, c.owner, s, c, nil)
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
	client.CreateAppointment(t, 200, b, c.owner, s, c, nil)
	client.CreateAppointment(t, 400, b, e, s, c, &startTimeStr)
	client.CreateAppointment(t, 400, b, c.owner, s, c, &startTimeStr)
	client.Delete(t, 200)
}
