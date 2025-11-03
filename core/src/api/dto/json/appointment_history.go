package dJSON

type AppointmentHistory struct {
	FieldChanges []FieldChange `json:"field_changes"`
}

type FieldChange struct {
	CreatedAt string `json:"created_at" example:"2021-01-01T09:00:00Z"`
	Field     string `json:"field" example:"field_name"`
	OldValue  string `json:"old_value" example:"old_value"`
	NewValue  string `json:"new_value" example:"new_value"`
	Reason    string `json:"reason" example:"Some reason."`
}
