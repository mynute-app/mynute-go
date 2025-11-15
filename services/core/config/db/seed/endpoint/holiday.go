package endpointSeed

import (
	"mynute-go/services/core/config/db/model"
	resourceSeed "mynute-go/services/core/config/db/seed/resource"
	"mynute-go/services/core/config/namespace"
)

var CreateHoliday = &model.EndPoint{
	Path:             "/holiday",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "CreateHoliday",
	Description:      "Create a holiday",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Company,
}

var GetHolidayById = &model.EndPoint{
	Path:             "/holiday/:id",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetHolidayById",
	Description:      "View holiday by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Holiday,
}

var GetHolidayByName = &model.EndPoint{
	Path:           "/holiday/name/:name",
	Method:         namespace.ViewActionMethod,
	ControllerName: "GetHolidayByName",
	Description:    "View holiday by name",
	NeedsCompanyId: true,
}

var UpdateHolidayById = &model.EndPoint{
	Path:             "/holiday/:id",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateHolidayById",
	Description:      "Update holiday by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Holiday,
}

var DeleteHolidayById = &model.EndPoint{
	Path:             "/holiday/:id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteHolidayById",
	Description:      "Delete holiday by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Holiday,
}

