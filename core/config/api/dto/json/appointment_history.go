package dJSON

type AppointmentHistory struct {
	FieldChanges []FieldChange `json:"field_changes"`
}

type FieldChange struct {
	CreatedAt string `json:"created_at" example:"00000000-0000-0000-0000-000000000000"`
	Field     string    `json:"field"`
	OldValue  string    `json:"old_value"`
	NewValue  string    `json:"new_value"`
	Reason    string    `json:"reason"`
}