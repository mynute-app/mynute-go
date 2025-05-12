package model

type ClientAppointmentArchive struct {
	BaseModel
	AppointmentBase
}

var ClientAppointmentArchiveTableName = "public.client_appointments_archive"

func (ClientAppointmentArchive) TableName() string {
	return ClientAppointmentArchiveTableName
}

func (ClientAppointmentArchive) Indexes() map[string]string {
	return AppointmentIndexes(ClientAppointmentArchiveTableName)
}