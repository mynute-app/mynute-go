package endpointSeed

import (
	"mynute-go/auth/model"
	resourceSeed "mynute-go/core/src/config/db/seed/resource"
	"mynute-go/core/src/config/namespace"
)

var CreateService = &model.EndPoint{
	Path:             "/service",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "CreateService",
	Description:      "Create a service",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Company,
}

var GetServiceById = &model.EndPoint{
	Path:             "/service/:id",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetServiceById",
	Description:      "View service by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Service,
}

var GetServiceByName = &model.EndPoint{
	Path:           "/service/name/:name",
	Method:         namespace.ViewActionMethod,
	ControllerName: "GetServiceByName",
	Description:    "View service by name",
	NeedsCompanyId: true,
}

var UpdateServiceById = &model.EndPoint{
	Path:             "/service/:id",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateServiceById",
	Description:      "Update service by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Service,
}

var DeleteServiceById = &model.EndPoint{
	Path:             "/service/:id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteServiceById",
	Description:      "Delete service by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Service,
}

var UpdateServiceImages = &model.EndPoint{
	Path:             "/service/:id/design/images",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateServiceImages",
	Description:      "Update images of a service",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Service,
}

var DeleteServiceImage = &model.EndPoint{
	Path:             "/service/:id/design/images/:image_type",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteServiceImage",
	Description:      "Delete an image of a service",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Service,
}

var GetServiceAvailability = &model.EndPoint{
	Path:           "/service/:id/availability",
	Method:         namespace.ViewActionMethod,
	ControllerName: "GetServiceAvailability",
	Description:    "Get availability of a service",
	NeedsCompanyId: true,
	Resource:       resourceSeed.Service,
}
