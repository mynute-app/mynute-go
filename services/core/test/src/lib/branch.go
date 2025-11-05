package lib

import (
	"fmt"
	"mynute-go/services/core/test/src/model"

	"github.com/google/uuid"
)

// Helper functions to retrieve Branch and Service objects from Company test setup data
func GetBranchByID(company *model.Company, branchIDStr string) (*model.Branch, error) {
	branchUUID, err := uuid.Parse(branchIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid Branch ID string from slot finder: %s, error: %v", branchIDStr, err)
	}
	for _, br := range company.Branches {
		if br.Created.ID == branchUUID {
			return br, nil
		}
	}
	return nil, fmt.Errorf("branch with ID %s (found by slot finder) not in company.Branches", branchIDStr)
}
