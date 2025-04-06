package model

type Permission struct {
	BaseModel
	EmployeeID   uint
	ResourceType string
	EndPointID   *uint // nil means all resources of this type
	Action       string
}
