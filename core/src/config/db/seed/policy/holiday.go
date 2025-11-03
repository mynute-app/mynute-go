package policySeed

import (
	"mynute-go/auth/model"
	endpointSeed "mynute-go/core/src/config/db/seed/endpoint"
)

var AllowCreateHoliday = &model.PolicyRule{
	Name:        "SDP: CanCreateHoliday",
	Description: "Allows company managers (Owner, GM, BM) to create holidays for the company/branch.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.CreateHoliday.ID,
	Conditions:  model.JsonRawMessage(CompanyManagerCheck), // Any manager in the relevant company context. Add branch check if needed.
}

var AllowGetHolidayById = &model.PolicyRule{
	Name:        "SDP: CanViewHolidayById",
	Description: "Allows company members to view holiday details.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.GetHolidayById.ID,
	Conditions:  model.JsonRawMessage(CompanyInternalUserCheck), // Any internal user of the holiday's company
}

var AllowUpdateHolidayById = &model.PolicyRule{
	Name:        "SDP: CanUpdateHoliday",
	Description: "Allows company managers (Owner, GM, BM) to update holidays.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.UpdateHolidayById.ID,
	Conditions:  model.JsonRawMessage(CompanyManagerCheck), // Any manager of the holiday's company
}

var AllowDeleteHolidayById = &model.PolicyRule{
	Name:        "SDP: CanDeleteHoliday",
	Description: "Allows company managers (Owner, GM, BM) to delete holidays.",
	Effect:      "Allow",
	EndPointID:  endpointSeed.DeleteHolidayById.ID,
	Conditions:  model.JsonRawMessage(CompanyManagerCheck), // Any manager of the holiday's company
}
