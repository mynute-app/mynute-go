package endpointSeed

import (
	"mynute-go/auth/config/db/model"
	resourceSeed "mynute-go/core/src/config/db/seed/resource"
	"mynute-go/core/src/config/namespace"
)

var CreateAppointment = &model.EndPoint{
	Path:             "/appointment",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "CreateAppointment",
	Description:      "Create an appointment",
	NeedsCompanyId:   true,
	DenyUnauthorized: false,
	Resource:         resourceSeed.Branch,
}

var GetAppointmentByID = &model.EndPoint{
	Path:             "/appointment/:id",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetAppointmentByID",
	Description:      "View appointment by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Appointment,
}

var UpdateAppointmentByID = &model.EndPoint{
	Path:             "/appointment/:id",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateAppointmentByID",
	Description:      "Update appointment by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Appointment,
}

var CancelAppointmentByID = &model.EndPoint{
	Path:             "/appointment/:id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "CancelAppointmentByID",
	Description:      "This wil cancel appointment by ID. Deleting appointments is forbidden.",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Appointment,
}

