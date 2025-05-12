package model

type ClientAppointmentArchive struct {
	AppointmentBase
}

func (ClientAppointmentArchive) Indexes() map[string]string {
	return map[string]string{
		"idx_client_time_active":   "CREATE INDEX IF NOT EXISTS idx_client_time_active ON public.client_appointments_archive (client_id, start_time, end_time, is_cancelled)",
		"idx_employee_time_active": "CREATE INDEX IF NOT EXISTS idx_employee_time_active ON public.client_appointments_archive (employee_id, start_time, end_time, is_cancelled)",
		"idx_branch_time_active":   "CREATE INDEX IF NOT EXISTS idx_branch_time_active ON public.client_appointments_archive (branch_id, start_time, end_time, is_cancelled)",
	}
}