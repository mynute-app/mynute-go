package endpointSeed

import (
	"mynute-go/services/auth/config/db/model"
	resourceSeed "mynute-go/services/core/src/config/db/seed/resource"
	"mynute-go/services/core/src/config/namespace"
)

var CreateSector = &model.EndPoint{
	Path:             "/sector",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "CreateSector",
	Description:      "Creates a company sector",
	DenyUnauthorized: true,
	Resource:         resourceSeed.Company,
}

var GetSectorById = &model.EndPoint{
	Path:           "/sector/:id",
	Method:         namespace.ViewActionMethod,
	ControllerName: "GetSectorById",
	Description:    "Retrieves a company sector by ID",
}

var GetSectorByName = &model.EndPoint{
	Path:           "/sector/name/:name",
	Method:         namespace.ViewActionMethod,
	ControllerName: "GetSectorByName",
	Description:    "Retrieves a company sector by name",
}

var UpdateSectorById = &model.EndPoint{
	Path:             "/sector/:id",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateSectorById",
	Description:      "Updates a company sector by ID",
	DenyUnauthorized: true,
	Resource:         resourceSeed.Sector,
}

var DeleteSectorById = &model.EndPoint{
	Path:             "/sector/:id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteSectorById",
	Description:      "Deletes a company sector by ID",
	DenyUnauthorized: true,
	Resource:         resourceSeed.Sector,
}

