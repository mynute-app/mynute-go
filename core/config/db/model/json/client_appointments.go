package mJSON

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ClientAppointments []*ClientAppointment

type ClientAppointment struct {
	AppointmentID    uuid.UUID `json:"appointment_id"`
	ServiceName      string    `json:"service_name"`
	ServicePrice     int64     `json:"service_price"`
	ServiceID        uuid.UUID `json:"service_id"`
	CompanyTradeName string    `json:"company_trade_name"`
	CompanyLegalName string    `json:"company_legal_name"`
	CompanyID        uuid.UUID `json:"company_id"`
	BranchAddress    string    `json:"branch_address"`
	BranchID         uuid.UUID `json:"branch_id"`
	EmployeeName     string    `json:"employee_name"`
	EmployeeID       uuid.UUID `json:"employee_id"`
	IsCancelled      bool      `json:"is_cancelled"`
	StartTime        time.Time `json:"start_time"`
	Price            *int64    `json:"price"`
	Currency         *string   `json:"currency"` // Default currency is BRL
}

func (cas *ClientAppointments) Value() (driver.Value, error) {
    // If the pointer itself is nil, or if it points to a nil or empty slice
    if cas == nil || *cas == nil || len(*cas) == 0 {
        return json.Marshal([]*ClientAppointment{}) // Explicitly marshal an empty slice to "[]"
    }
    return json.Marshal(*cas) // Marshal the actual slice
}

func (cas *ClientAppointments) Scan(value any) error {
    if value == nil {
        *cas = nil // Or *cas = ClientAppointments{} if you prefer an empty slice over nil
        return nil
    }

    bytes, ok := value.([]byte)
    if !ok {
        return fmt.Errorf("Scan: failed to assert value to []byte, got %T", value)
    }

    // Handle cases where the DB might store an empty string or "null" literally
    if len(bytes) == 0 || string(bytes) == "null" {
        *cas = nil // Or *cas = ClientAppointments{}
        return nil
    }

    // The core unmarshalling logic
    err := json.Unmarshal(bytes, cas)
    if err != nil {
        // The original error message already points to this problem.
        // You could add more context if needed.
        // e.g., return fmt.Errorf("failed to unmarshal ClientAppointments from JSON '%s': %w", string(bytes), err)
        return err
    }
    return nil
}

func (cas *ClientAppointments) Add(ca *ClientAppointment) {
	if cas == nil {
		cas = &ClientAppointments{}
	}
	*cas = append(*cas, ca)
}

func (cas *ClientAppointments) UpdateOneById(id uuid.UUID, newCa *ClientAppointment) {
	for i, ca := range *cas {
		if ca.AppointmentID == id {
			(*cas)[i] = newCa
			break
		}
	}
}
