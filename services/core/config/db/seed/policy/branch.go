package policySeed

import (
	"mynute-go/services/auth/config/db/model"
	endpointSeed "mynute-go/services/core/config/db/seed/endpoint"
)

var AllowCreateBranch = &model.PolicyRule{
	Name:        "SDP: CanCreateBranch",
	Description: "Allows company Owner or General Manager to create branches.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.CreateBranch.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminCheck), // Owner or GM of the target company
}

var AllowGetBranchById = &model.PolicyRule{
	Name:        "SDP: CanViewBranchById",
	Description: "Allows any user belonging to the same company to view branch details by ID.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetBranchById.ID,
	Conditions:  model.JsonRawMessage(CompanyInternalUserCheck), // Any internal user of the branch's company can view
}

var AllowUpdateBranchById = &model.PolicyRule{
	Name:        "SDP: CanUpdateBranch",
	Description: "Allows company Owner, General Manager, or assigned Branch Manager to update branches.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateBranchById.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrAssignedBranchManagerCheck),
}

var AllowDeleteBranchById = &model.PolicyRule{
	Name:        "SDP: CanDeleteBranch",
	Description: "Allows company Owner or General Manager to delete branches.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteBranchById.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminCheck), // Only Owner or GM
}

var AllowGetEmployeeServicesByBranchId = &model.PolicyRule{
	Name:        "SDP: CanViewEmployeeServicesInBranch",
	Description: "Allows company members to view employee services within a branch.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetEmployeeServicesByBranchId.ID,
	Conditions:  model.JsonRawMessage(CompanyInternalUserCheck), // Any internal user of the branch's company
}

var AllowAddServiceToBranch = &model.PolicyRule{
	Name:        "SDP: CanAddServiceToBranch",
	Description: "Allows company managers (Owner, GM, relevant BM) to add services to a branch.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.AddServiceToBranch.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrAssignedBranchManagerCheck),
}

var AllowRemoveServiceFromBranch = &model.PolicyRule{
	Name:        "SDP: CanRemoveServiceFromBranch",
	Description: "Allows company managers (Owner, GM, relevant BM) to remove services from a branch.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.RemoveServiceFromBranch.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrAssignedBranchManagerCheck),
}

var AllowCreateBranchWorkSchedule = &model.PolicyRule{
	Name:        "SDP: CanCreateBranchWorkSchedule",
	Description: "Allows company Owner, General Manager, or assigned Branch Manager to create work schedules for a branch.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.CreateBranchWorkSchedule.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrAssignedBranchManagerCheck),
}

var AllowGetBranchWorkRangeById = &model.PolicyRule{
	Name:        "SDP: CanViewBranchWorkRangeById",
	Description: "Allows company members to view branch work schedules by ID.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetBranchWorkRange.ID,
	Conditions:  model.JsonRawMessage(CompanyInternalUserCheck), // Any internal user of the branch's company can view work schedules
}

var AllowDeleteBranchWorkRangeById = &model.PolicyRule{
	Name:        "SDP: CanDeleteBranchWorkRangeById",
	Description: "Allows company Owner, General Manager, or assigned Branch Manager to delete branch work schedules.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteBranchWorkRange.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrAssignedBranchManagerCheck),
}

var AllowUpdateBranchWorkRangeById = &model.PolicyRule{
	Name:        "SDP: CanUpdateBranchWorkRangeById",
	Description: "Allows company Owner, General Manager, or assigned Branch Manager to update branch work schedules.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateBranchWorkRange.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrAssignedBranchManagerCheck),
}

var AllowAddBranchWorkRangeService = &model.PolicyRule{
	Name:        "SDP: CanAddBranchWorkRangeService",
	Description: "Allows company Owner, General Manager, or assigned Branch Manager to add services to a branch work range.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.AddBranchWorkRangeServices.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrAssignedBranchManagerCheck),
}

var AllowDeleteBranchWorkRangeService = &model.PolicyRule{
	Name:        "SDP: CanDeleteBranchWorkRangeService",
	Description: "Allows company Owner, General Manager or assigned Branch Manager to remove services from a branch work range.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteBranchWorkRangeService.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrAssignedBranchManagerCheck),
}

var AllowUpdateBranchImages = &model.PolicyRule{
	Name:        "SDP: CanUpdateBranchImages",
	Description: "Allows company Owner, General Manager, or assigned Branch Manager to update branch images.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateBranchImages.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrAssignedBranchManagerCheck),
}

var AllowDeleteBranchImage = &model.PolicyRule{
	Name:        "SDP: CanDeleteBranchImage",
	Description: "Allows company Owner, General Manager, or assigned Branch Manager to delete branch images.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteBranchImage.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrAssignedBranchManagerCheck),
}

var AllowGetBranchAppointmentsById = &model.PolicyRule{
	Name:        "SDP: CanViewBranchAppointmentsById",
	Description: "Allows company Owner, General Manager, or assigned Branch Manager to view appointments within a branch.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetBranchAppointmentsById.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrAssignedBranchManagerCheck),
}
