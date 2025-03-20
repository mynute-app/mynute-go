package e2e_test

import (
	"agenda-kaki-go/core"
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/tests/handlers"
	"fmt"
	"testing"
	"time"
)

type User struct {
	created    model.User
	auth_token string
}

func Test_User(t *testing.T) {
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
	user.CreateAppointment(t, 200, b, e, s, c, nil)
	startTimeStr := user.created.Appointments[0].StartTime.Format(time.RFC3339)
	user.CreateAppointment(t, 400, b, e, s, c, &startTimeStr)
	user.Delete(t, 200)
}

func (u *User) Set(t *testing.T) {
	u.Create(t, 200)
	u.VerifyEmail(t, 200)
	u.Login(t, 200)
}

func (u *User) Create(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/user")
	http.ExpectStatus(s)
	email := lib.GenerateRandomEmail("user")
	pswd := "1SecurePswd!"
	http.Send(DTO.CreateUser{
		Email:    email,
		Name:     lib.GenerateRandomName("User Name"),
		Surname:  lib.GenerateRandomName("User Surname"),
		Password: pswd,
		Phone:    lib.GenerateRandomPhoneNumber(),
	})
	http.ParseResponse(&u.created)
	u.created.Password = pswd
}

func (u *User) Update(t *testing.T, s int, changes map[string]any) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("PATCH")
	http.URL("/user/" + fmt.Sprintf("%v", u.created.ID))
	http.ExpectStatus(s)
	http.Header("Authorization", u.auth_token)
	http.Send(changes)
}

func (u *User) GetByEmail(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL("/user/email/" + u.created.Email)
	http.ExpectStatus(s)
	http.Header("Authorization", u.auth_token)
	http.Send(nil)
}

func (u *User) Delete(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("DELETE")
	http.URL(fmt.Sprintf("/user/%v", u.created.ID))
	http.ExpectStatus(s)
	http.Header("Authorization", u.auth_token)
	http.Send(nil)
}

func (u *User) VerifyEmail(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL(fmt.Sprintf("/user/verify-email/%v/%s", u.created.Email, "12345"))
	http.ExpectStatus(s)
	http.Header("Authorization", u.auth_token)
	http.Send(nil)
}

func (u *User) Login(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/user/login")
	http.ExpectStatus(s)
	http.Send(map[string]any{
		"email":    u.created.Email,
		"password": "1SecurePswd!",
	})
	auth := http.ResHeaders["Authorization"]
	if len(auth) == 0 {
		t.Errorf("Authorization header not found")
		return
	}
	u.auth_token = auth[0]
}

func (u *User) CreateAppointment(t *testing.T, s int, b *Branch, e *Employee, srvc *Service, c *Company, startTime *string) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/appointment")
	http.ExpectStatus(s)
	http.Header("Authorization", u.auth_token)
	if startTime == nil {
		tempStartTime := lib.GenerateDateRFC3339(2027, 10, 29)
		startTime = &tempStartTime
	}
	http.Send(DTO.Appointment{
		BranchID:   b.created.ID,
		ServiceID:  srvc.created.ID,
		EmployeeID: e.created.ID,
		UserID:     u.created.ID,
		CompanyID:  c.created.ID,
		StartTime:  *startTime,
	})
	var newAppointment model.Appointment
	http.ParseResponse(&newAppointment)
	u.created.Appointments = append(u.created.Appointments, newAppointment)
	e.GetById(t, 200)
	b.GetById(t, 200)
	u.GetByEmail(t, 200)
}

func Test_User_Create_Success(t *testing.T) {
	server := core.NewServer().Run("test")
	user := &User{}
	user.Create(t, 200)
	server.Shutdown()
}

func Test_Login_Success(t *testing.T) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/login")
	http.ExpectStatus(200)
	http.Send(map[string]any{
		"email":    "test@email.com",
		"password": "1SecurePswd!",
	})
}
