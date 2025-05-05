package mJSON

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type AppointmentHistory struct {
	FieldChanges []FieldChange `json:"field_changes"`
}

type FieldChange struct {
	CreatedAt time.Time `json:"created_at"`
	Field     string    `json:"field"`
	OldValue  string    `json:"old_value"`
	NewValue  string    `json:"new_value"`
	Reason    string    `json:"reason"`
}

func (ah *AppointmentHistory) Value() (driver.Value, error) {
	return json.Marshal(ah)
}

func (ah *AppointmentHistory) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan WorkSchedule: expected []byte")
	}

	return json.Unmarshal(bytes, ah)
}

func (ah *AppointmentHistory) IsEmpty() bool {
	return ah == nil || len(ah.FieldChanges) == 0
}

func (ah *AppointmentHistory) FilterByField(field string) []FieldChange {
	if ah == nil {
		return nil
	}
	var filteredChanges []FieldChange
	for _, change := range ah.FieldChanges {
		if change.Field == field {
			filteredChanges = append(filteredChanges, change)
		}
	}
	return filteredChanges
}
