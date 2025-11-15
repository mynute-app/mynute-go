package policySeed

import (
	"mynute-go/services/core/config/db/model"
	endpointSeed "mynute-go/services/core/config/db/seed/endpoint"
)

// --- TENANT POLICIES: Holiday-related operations ---

var TenantAllowCreateHoliday = &model.TenantPolicy{
	Name:        "TP: CanCreateHoliday",
	Description: "Allows company managers (Owner, GM, BM) to create holidays for the company/branch.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.CreateHoliday.ID,
	Conditions:  model.JsonRawMessage(CompanyManagerCheck),
}

var TenantAllowGetHolidayById = &model.TenantPolicy{
	Name:        "TP: CanViewHolidayById",
	Description: "Allows company members to view holiday details.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetHolidayById.ID,
	Conditions:  model.JsonRawMessage(CompanyInternalUserCheck),
}

var TenantAllowUpdateHolidayById = &model.TenantPolicy{
	Name:        "TP: CanUpdateHoliday",
	Description: "Allows company managers (Owner, GM, BM) to update holidays.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateHolidayById.ID,
	Conditions:  model.JsonRawMessage(CompanyManagerCheck),
}

var TenantAllowDeleteHolidayById = &model.TenantPolicy{
	Name:        "TP: CanDeleteHoliday",
	Description: "Allows company managers (Owner, GM, BM) to delete holidays.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteHolidayById.ID,
	Conditions:  model.JsonRawMessage(CompanyManagerCheck),
}
