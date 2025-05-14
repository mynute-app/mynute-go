package mJSON

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type ClientAppointments []*ClientAppointment

type ClientAppointment struct {
	AppointmentID uuid.UUID `json:"appointment_id"`
	ServiceName   string    `json:"service_name"`
	ServiceValue  uint      `json:"service_value"`
	ServiceID     uuid.UUID `json:"service_id"`
	CompanyName   string    `json:"company_name"`
	CompanyID     uuid.UUID `json:"company_id"`
	BranchAddress string    `json:"branch_address"`
	BranchID      uuid.UUID `json:"branch_id"`
	EmployeeName  string    `json:"employee_name"`
	EmployeeID    uuid.UUID `json:"employee_id"`
	IsCancelled   bool      `json:"is_cancelled"`
	StartTime     time.Time `json:"start_time"`
	ValuePaid     *uint     `json:"value_paid"`
}

func (cas *ClientAppointments) Value() (driver.Value, error) {
	if cas == nil || len(*cas) == 0 {
		return json.Marshal(&ClientAppointments{})
	}
	return json.Marshal(cas)
}

func (cas *ClientAppointments) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan ClientAppointments: expected []byte")
	}

	return json.Unmarshal(bytes, cas)
}

func (cas *ClientAppointments) Add(ca *ClientAppointment) {
	if cas == nil {
		cas = &ClientAppointments{}
	}
	*cas = append(*cas, ca)
}

func (cas *ClientAppointments) UpdateOneById(id uuid.UUID, newCa *ClientAppointment) {
	for _, ca := range *cas {
		if ca.AppointmentID != newCa.AppointmentID {
			continue
		}
		ca = newCa
	}
}