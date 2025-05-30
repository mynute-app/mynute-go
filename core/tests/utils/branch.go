package utils_test

import (
	models_test "agenda-kaki-go/core/tests/models"
	"testing"

	"github.com/google/uuid"
)

// Helper functions to retrieve Branch and Service objects from Company test setup data
func GetBranchByID(t *testing.T, company *models_test.Company, branchIDStr string) *models_test.Branch {
	t.Helper()
	branchUUID, err := uuid.Parse(branchIDStr)
	if err != nil {
		t.Fatalf("Invalid Branch ID string from slot finder: %s, error: %v", branchIDStr, err)
	}
	for _, br := range company.Branches {
		if br.Created.ID == branchUUID {
			return br
		}
	}
	t.Fatalf("Test setup error: Branch with ID %s (found by slot finder) not in company.Branches", branchIDStr)
	return nil
}