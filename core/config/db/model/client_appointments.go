package model

type ClientAppointment struct {
	AppointmentBase
}

func (ClientAppointment) Indexes() map[string]string {
	return map[string]string{
		"idx_client_time_active":   "CREATE INDEX IF NOT EXISTS idx_client_time_active ON public.client_appointments_index (client_id, start_time, end_time, is_cancelled)",
		"idx_employee_time_active": "CREATE INDEX IF NOT EXISTS idx_employee_time_active ON public.client_appointments_index (employee_id, start_time, end_time, is_cancelled)",
		"idx_branch_time_active":   "CREATE INDEX IF NOT EXISTS idx_branch_time_active ON public.client_appointments_index (branch_id, start_time, end_time, is_cancelled)",
	}
}

