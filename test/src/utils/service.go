package utilsT

import (
	"fmt"
	modelT "mynute-go/test/src/models"

	"github.com/google/uuid"
)

func GetServiceByID(company *modelT.Company, serviceIDStr string) (*modelT.Service, error) {
	serviceUUID, err := uuid.Parse(serviceIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid Service ID string from slot finder: %s, error: %v", serviceIDStr, err)
	}
	for _, serv := range company.Services { // Assuming company.Services holds all services
		if serv.Created.ID == serviceUUID {
			return serv, nil
		}
	}
	return nil, fmt.Errorf("service with ID %s (found by slot finder) not in company.Services", serviceIDStr)
}
