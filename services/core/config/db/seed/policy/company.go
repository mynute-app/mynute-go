package policySeed

import (
	"mynute-go/services/auth/config/db/model"
	endpointSeed "mynute-go/services/core/config/db/seed/endpoint"
)

var AllowGetCompanyById = &model.PolicyRule{
	Name:        "SDP: CanViewCompanyById",
	Description: "Allows any member (employee/manager) of the company to view its details.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetCompanyById.ID,
	Conditions:  model.JsonRawMessage(CompanyMembershipAccessCheck),
}

var AllowUpdateCompanyById = &model.PolicyRule{
	Name:        "SDP: CanUpdateCompany",
	Description: "Allows the company Owner or General Manager to update company details.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateCompanyById.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminCheck), // Only Owner or GM of this company
}

var AllowDeleteCompanyById = &model.PolicyRule{
	Name:        "SDP: CanDeleteCompany",
	Description: "Allows ONLY the company Owner to delete the company.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteCompanyById.ID,
	Conditions:  model.JsonRawMessage(CompanyOwnerCheck), // Only Owner of this company
}

var AllowUpdateCompanyImages = &model.PolicyRule{
	Name:        "SDP: CanUpdateCompanyImages",
	Description: "Allows company Owner or General Manager to update company images.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateCompanyImages.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminCheck), // Only Owner or GM of this company
}

var AllowDeleteCompanyImage = &model.PolicyRule{
	Name:        "SDP: CanDeleteCompanyImage",
	Description: "Allows company Owner or General Manager to delete company images.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteCompanyImage.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminCheck), // Only Owner or GM of this company
}

var AllowUpdateCompanyColors = &model.PolicyRule{
	Name:        "SDP: CanUpdateCompanyColors",
	Description: "Allows company Owner or General Manager to update company colors.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateCompanyColors.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminCheck), // Only Owner or GM of this company
}
