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
	AppointmentID uuid.UUID `json:"appointment_id"`
	CompanySchema string    `json:"company_schema"`
	IsCancelled   bool      `json:"is_cancelled"`
	StartTime     time.Time `json:"start_time"`
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
	if *cas == nil {
		*cas = make(ClientAppointments, 0, 1)
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
