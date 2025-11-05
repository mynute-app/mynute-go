package endpointSeed

import (
	"mynute-go/services/auth/config/db/model"
	resourceSeed "mynute-go/services/core/src/config/db/seed/resource"
	"mynute-go/services/core/src/config/namespace"
)

var CreateCompany = &model.EndPoint{
	Path:           "/company",
	Method:         namespace.CreateActionMethod,
	ControllerName: "CreateCompany",
	Description:    "Create a company",
}

var GetCompanyById = &model.EndPoint{
	Path:             "/company/:id",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetCompanyById",
	Description:      "View company by ID",
	DenyUnauthorized: true,
	NeedsCompanyId:   true,
	Resource:         resourceSeed.Company,
}

var GetCompanyByName = &model.EndPoint{
	Path:           "/company/name/:name",
	Method:         namespace.ViewActionMethod,
	ControllerName: "GetCompanyByName",
	Description:    "View company by name",
}

var CheckIfCompanyExistsByTaxID = &model.EndPoint{
	Path:           "/company/tax_id/:tax_id/exists",
	Method:         namespace.ViewActionMethod,
	ControllerName: "CheckIfCompanyExistsByTaxID",
	Description:    "Check if company exists by tax ID",
}

var GetCompanyByTaxId = &model.EndPoint{
	Path:             "/company/tax_id/:tax_id",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetCompanyByTaxId",
	DenyUnauthorized: true,
	NeedsCompanyId:   true,
	Description:      "View company by tax ID",
}

var GetCompanyBySubdomain = &model.EndPoint{
	Path:           "/company/subdomain/:subdomain_name",
	Method:         namespace.ViewActionMethod,
	ControllerName: "GetCompanyBySubdomain",
	Description:    "View company by subdomain",
}

var UpdateCompanyById = &model.EndPoint{
	Path:             "/company/:id",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateCompanyById",
	Description:      "Update company by ID",
	DenyUnauthorized: true,
	NeedsCompanyId:   true,
	Resource:         resourceSeed.Company,
}

var UpdateCompanyImages = &model.EndPoint{
	Path:             "/company/:id/design/images",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateCompanyImages",
	Description:      "Update company design images",
	DenyUnauthorized: true,
	NeedsCompanyId:   true,
	Resource:         resourceSeed.Company,
}

var DeleteCompanyImage = &model.EndPoint{
	Path:             "/company/:id/design/images/:image_type",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteCompanyImage",
	Description:      "Delete company design images",
	DenyUnauthorized: true,
	NeedsCompanyId:   true,
	Resource:         resourceSeed.Company,
}

var UpdateCompanyColors = &model.EndPoint{
	Path:             "/company/:id/design/colors",
	Method:           namespace.PutActionMethod,
	ControllerName:   "UpdateCompanyColors",
	Description:      "Update company design colors",
	DenyUnauthorized: true,
	NeedsCompanyId:   true,
	Resource:         resourceSeed.Company,
}

var DeleteCompanyById = &model.EndPoint{
	Path:             "/company/:id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteCompanyById",
	Description:      "Delete company by ID",
	DenyUnauthorized: true,
	NeedsCompanyId:   true,
	Resource:         resourceSeed.Company,
}

