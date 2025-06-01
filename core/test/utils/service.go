package utils_test

import (
	models_test "agenda-kaki-go/core/test/models"
	"testing"

	"github.com/google/uuid"
)

func GetServiceByID(t *testing.T, company *models_test.Company, serviceIDStr string) *models_test.Service {
	t.Helper()
	serviceUUID, err := uuid.Parse(serviceIDStr)
	if err != nil {
		t.Fatalf("Invalid Service ID string from slot finder: %s, error: %v", serviceIDStr, err)
	}
	for _, serv := range company.Services { // Assuming company.Services holds all services
		if serv.Created.ID == serviceUUID {
			return serv
		}
	}
	t.Fatalf("Test setup error: Service with ID %s (found by slot finder) not in company.Services", serviceIDStr)
	return nil
}
