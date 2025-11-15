package policySeed

import (
	"mynute-go/services/core/config/db/model"
	endpointSeed "mynute-go/services/core/config/db/seed/endpoint"
)

// --- TENANT POLICIES: Service-related operations ---

var TenantAllowCreateService = &model.TenantPolicy{
	Name:        "TP: CanCreateService",
	Description: "Allows company managers (Owner, GM, BM) to create services.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.CreateService.ID,
	Conditions:  model.JsonRawMessage(CompanyManagerCheck),
}

var TenantAllowGetServiceById = &model.TenantPolicy{
	Name:        "TP: CanViewServiceById",
	Description: "Allows company members to view service details.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetServiceById.ID,
	Conditions:  model.JsonRawMessage(CompanyInternalUserCheck),
}

var TenantAllowUpdateServiceById = &model.TenantPolicy{
	Name:        "TP: CanUpdateService",
	Description: "Allows company managers (Owner, GM, BM) to update services.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateServiceById.ID,
	Conditions:  model.JsonRawMessage(CompanyManagerCheck),
}

var TenantAllowDeleteServiceById = &model.TenantPolicy{
	Name:        "TP: CanDeleteService",
	Description: "Allows company managers (Owner, GM, BM) to delete services.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteServiceById.ID,
	Conditions:  model.JsonRawMessage(CompanyManagerCheck),
}

var TenantAllowUpdateServiceImages = &model.TenantPolicy{
	Name:        "TP: CanUpdateServiceImages",
	Description: "Allows company managers (Owner, GM) to update service images.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateServiceImages.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminCheck),
}

var TenantAllowDeleteServiceImage = &model.TenantPolicy{
	Name:        "TP: CanDeleteServiceImage",
	Description: "Allows company managers (Owner, GM) to delete service images.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteServiceImage.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminCheck),
}
