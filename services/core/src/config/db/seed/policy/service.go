package policySeed

import (
	"mynute-go/services/auth/config/db/model"
	endpointSeed "mynute-go/services/core/src/config/db/seed/endpoint"
)

var AllowCreateService = &model.PolicyRule{
	Name:        "SDP: CanCreateService",
	Description: "Allows company managers (Owner, GM, BM) to create services.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.CreateService.ID,
	Conditions:  model.JsonRawMessage(CompanyManagerCheck), // Any manager of the company context
}

var AllowGetServiceById = &model.PolicyRule{
	Name:        "SDP: CanViewServiceById",
	Description: "Allows company members to view service details.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetServiceById.ID,
	Conditions:  model.JsonRawMessage(CompanyInternalUserCheck), // Any internal user of the service's company
}

var AllowUpdateServiceById = &model.PolicyRule{
	Name:        "SDP: CanUpdateService",
	Description: "Allows company managers (Owner, GM, BM) to update services.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateServiceById.ID,
	Conditions:  model.JsonRawMessage(CompanyManagerCheck), // Any manager of the service's company
}

var AllowDeleteServiceById = &model.PolicyRule{
	Name:        "SDP: CanDeleteService",
	Description: "Allows company managers (Owner, GM, BM) to delete services.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteServiceById.ID,
	Conditions:  model.JsonRawMessage(CompanyManagerCheck), // Any manager of the service's company
}

var AllowUpdateServiceImages = &model.PolicyRule{
	Name:        "SDP: CanUpdateServiceImages",
	Description: "Allows company managers (Owner, GM) to update service images.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateServiceImages.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminCheck), // Any manager of the service's company
}

var AllowDeleteServiceImage = &model.PolicyRule{
	Name:        "SDP: CanDeleteServiceImage",
	Description: "Allows company managers (Owner, GM) to delete service images.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteServiceImage.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminCheck), // Any manager of the service's company
}

