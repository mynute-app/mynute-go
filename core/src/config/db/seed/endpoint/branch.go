package endpointSeed

import (
	"mynute-go/auth/model"
	resourceSeed "mynute-go/core/src/config/db/seed/resource"
	"mynute-go/core/src/config/namespace"
)

var CreateBranch = &model.EndPoint{
	Path:             "/branch",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "CreateBranch",
	Description:      "Create a branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Company,
}

var GetBranchById = &model.EndPoint{
	Path:             "/branch/:id",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetBranchById",
	Description:      "View branch by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Branch,
}

var UpdateBranchById = &model.EndPoint{
	Path:             "/branch/:id",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateBranchById",
	Description:      "Update branch by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Branch,
}

var DeleteBranchById = &model.EndPoint{
	Path:             "/branch/:id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteBranchById",
	Description:      "Delete branch by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Branch,
}

var GetEmployeeServicesByBranchId = &model.EndPoint{
	Path:             "/branch/:branch_id/employee/:employee_id/services",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetEmployeeServicesByBranchId",
	Description:      "View employee offered services at the branch by branch ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Branch,
}

var AddServiceToBranch = &model.EndPoint{
	Path:             "/branch/:branch_id/service/:service_id",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "AddServiceToBranch",
	Description:      "Add service to branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Branch,
}

var RemoveServiceFromBranch = &model.EndPoint{
	Path:             "/branch/:branch_id/service/:service_id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "RemoveServiceFromBranch",
	Description:      "Remove service from branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Branch,
}

var CreateBranchWorkSchedule = &model.EndPoint{
	Path:             "/branch/:id/work_schedule",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "CreateBranchWorkSchedule",
	Description:      "Add work schedule to branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Branch,
}

var UpdateBranchImages = &model.EndPoint{
	Path:             "/branch/:id/design/images",
	Method:           namespace.PatchActionMethod,
	ControllerName:   "UpdateBranchImages",
	Description:      "Update branch images",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Branch,
}

var DeleteBranchImage = &model.EndPoint{
	Path:             "/branch/:id/design/images/:image_type",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteBranchImage",
	Description:      "Delete branch image",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Branch,
}

var GetBranchWorkSchedule = &model.EndPoint{
	Path:           "/branch/:id/work_schedule",
	Method:         namespace.ViewActionMethod,
	ControllerName: "GetBranchWorkSchedule",
	Description:    "View work schedule for branch",
	NeedsCompanyId: true,
	Resource:       resourceSeed.Branch,
}

var GetBranchWorkRange = &model.EndPoint{
	Path:             "/branch/:id/work_range/:work_range_id",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetBranchWorkRange",
	Description:      "View work range by ID",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Branch,
}

var UpdateBranchWorkRange = &model.EndPoint{
	Path:             "/branch/:id/work_range/:work_range_id",
	Method:           namespace.PutActionMethod,
	ControllerName:   "UpdateBranchWorkRange",
	Description:      "Update work range in branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Branch,
}

var DeleteBranchWorkRange = &model.EndPoint{
	Path:             "/branch/:id/work_range/:work_range_id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteBranchWorkRange",
	Description:      "Remove work range from branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Branch,
}

var AddBranchWorkRangeServices = &model.EndPoint{
	Path:             "/branch/:id/work_range/:work_range_id/services",
	Method:           namespace.CreateActionMethod,
	ControllerName:   "AddBranchWorkRangeServices",
	Description:      "Add services to work range in branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Branch,
}

var DeleteBranchWorkRangeService = &model.EndPoint{
	Path:             "/branch/:id/work_range/:work_range_id/service/:service_id",
	Method:           namespace.DeleteActionMethod,
	ControllerName:   "DeleteBranchWorkRangeService",
	Description:      "Remove service from work range in branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Branch,
}

var GetBranchAppointmentsById = &model.EndPoint{
	Path:             "/branch/:id/appointments",
	Method:           namespace.ViewActionMethod,
	ControllerName:   "GetBranchAppointmentsById",
	Description:      "View appointments for a branch",
	NeedsCompanyId:   true,
	DenyUnauthorized: true,
	Resource:         resourceSeed.Branch,
}
