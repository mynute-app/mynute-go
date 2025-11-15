package policySeed

import (
	"mynute-go/services/core/config/db/model"
	endpointSeed "mynute-go/services/core/config/db/seed/endpoint"
)

// --- TENANT POLICIES: Branch-related operations ---

var TenantAllowCreateBranch = &model.TenantPolicy{
	Name:        "TP: CanCreateBranch",
	Description: "Allows company Owner or General Manager to create branches.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.CreateBranch.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminCheck),
}

var TenantAllowGetBranchById = &model.TenantPolicy{
	Name:        "TP: CanViewBranchById",
	Description: "Allows any user belonging to the same company to view branch details by ID.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetBranchById.ID,
	Conditions:  model.JsonRawMessage(CompanyInternalUserCheck),
}

var TenantAllowUpdateBranchById = &model.TenantPolicy{
	Name:        "TP: CanUpdateBranch",
	Description: "Allows company Owner, General Manager, or assigned Branch Manager to update branches.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateBranchById.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrAssignedBranchManagerCheck),
}

var TenantAllowDeleteBranchById = &model.TenantPolicy{
	Name:        "TP: CanDeleteBranch",
	Description: "Allows company Owner or General Manager to delete branches.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteBranchById.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminCheck),
}

var TenantAllowGetEmployeeServicesByBranchId = &model.TenantPolicy{
	Name:        "TP: CanViewEmployeeServicesInBranch",
	Description: "Allows company members to view employee services within a branch.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetEmployeeServicesByBranchId.ID,
	Conditions:  model.JsonRawMessage(CompanyInternalUserCheck),
}

var TenantAllowAddServiceToBranch = &model.TenantPolicy{
	Name:        "TP: CanAddServiceToBranch",
	Description: "Allows company managers to add services to a branch.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.AddServiceToBranch.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrAssignedBranchManagerCheck),
}

var TenantAllowRemoveServiceFromBranch = &model.TenantPolicy{
	Name:        "TP: CanRemoveServiceFromBranch",
	Description: "Allows company managers to remove services from a branch.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.RemoveServiceFromBranch.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrAssignedBranchManagerCheck),
}

var TenantAllowCreateBranchWorkSchedule = &model.TenantPolicy{
	Name:        "TP: CanCreateBranchWorkSchedule",
	Description: "Allows company managers to create work schedules for a branch.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.CreateBranchWorkSchedule.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrAssignedBranchManagerCheck),
}

var TenantAllowGetBranchWorkRangeById = &model.TenantPolicy{
	Name:        "TP: CanViewBranchWorkRangeById",
	Description: "Allows company members to view branch work schedules by ID.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetBranchWorkRange.ID,
	Conditions:  model.JsonRawMessage(CompanyInternalUserCheck),
}

var TenantAllowDeleteBranchWorkRangeById = &model.TenantPolicy{
	Name:        "TP: CanDeleteBranchWorkRangeById",
	Description: "Allows company managers to delete branch work schedules.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteBranchWorkRange.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrAssignedBranchManagerCheck),
}

var TenantAllowUpdateBranchWorkRangeById = &model.TenantPolicy{
	Name:        "TP: CanUpdateBranchWorkRangeById",
	Description: "Allows company managers to update branch work schedules.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateBranchWorkRange.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrAssignedBranchManagerCheck),
}

var TenantAllowAddBranchWorkRangeService = &model.TenantPolicy{
	Name:        "TP: CanAddBranchWorkRangeService",
	Description: "Allows company managers to add services to a branch work range.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.AddBranchWorkRangeServices.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrAssignedBranchManagerCheck),
}

var TenantAllowDeleteBranchWorkRangeService = &model.TenantPolicy{
	Name:        "TP: CanDeleteBranchWorkRangeService",
	Description: "Allows company managers to remove services from a branch work range.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteBranchWorkRangeService.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrAssignedBranchManagerCheck),
}

var TenantAllowUpdateBranchImages = &model.TenantPolicy{
	Name:        "TP: CanUpdateBranchImages",
	Description: "Allows company managers to update branch images.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateBranchImages.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrAssignedBranchManagerCheck),
}

var TenantAllowDeleteBranchImage = &model.TenantPolicy{
	Name:        "TP: CanDeleteBranchImage",
	Description: "Allows company managers to delete branch images.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteBranchImage.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrAssignedBranchManagerCheck),
}

var TenantAllowGetBranchAppointmentsById = &model.TenantPolicy{
	Name:        "TP: CanViewBranchAppointmentsById",
	Description: "Allows company managers to view appointments within a branch.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetBranchAppointmentsById.ID,
	Conditions:  model.JsonRawMessage(CompanyAdminOrAssignedBranchManagerCheck),
}
