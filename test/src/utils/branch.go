package utilsT

import (
	"fmt"
	modelT "mynute-go/test/src/models"

	"github.com/google/uuid"
)

// Helper functions to retrieve Branch and Service objects from Company test setup data
func GetBranchByID(company *modelT.Company, branchIDStr string) (*modelT.Branch, error) {
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
