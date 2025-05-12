package model

type ClientAppointment struct {
	AppointmentBase
}

var ClientAppointmentTableName = "public.client_appointments"
func (ClientAppointment) TableName() string {
	return ClientAppointmentTableName
}

func (ClientAppointment) Indexes() map[string]string {
	return AppointmentIndexes(ClientAppointmentTableName)
}