package policySeed

import (
	"mynute-go/services/core/config/db/model"
	endpointSeed "mynute-go/services/core/config/db/seed/endpoint"
)

// --- TENANT POLICIES: Company-related operations ---

// ===== COMPANY POLICIES =====

var TenantAllowGetCompanyById = &model.TenantPolicy{
	Name:        "TP: CanViewCompanyById",
	Description: "Allows any member (employee/manager) of the company to view its details.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetCompanyById.ID,
	Conditions:  model.JsonRawMessage(CompanyMembershipAccessCheck),
}

var TenantAllowUpdateCompanyById = &model.TenantPolicy{
	Name:        "TP: CanUpdateCompany",
	Description: "Allows the company Owner or General Manager to update company details.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateCompanyById.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminCheck),
}

var TenantAllowDeleteCompanyById = &model.TenantPolicy{
	Name:        "TP: CanDeleteCompany",
	Description: "Allows ONLY the company Owner to delete the company.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteCompanyById.ID,
	Conditions:  model.JsonRawMessage(CompanyOwnerCheck),
}

var TenantAllowUpdateCompanyImages = &model.TenantPolicy{
	Name:        "TP: CanUpdateCompanyImages",
	Description: "Allows company Owner or General Manager to update company images.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateCompanyImages.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminCheck),
}

var TenantAllowDeleteCompanyImage = &model.TenantPolicy{
	Name:        "TP: CanDeleteCompanyImage",
	Description: "Allows company Owner or General Manager to delete company images.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteCompanyImage.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminCheck),
}

var TenantAllowUpdateCompanyColors = &model.TenantPolicy{
	Name:        "TP: CanUpdateCompanyColors",
	Description: "Allows company Owner or General Manager to update company colors.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateCompanyColors.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminCheck),
}
