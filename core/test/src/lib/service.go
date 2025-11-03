package lib

import (
	"fmt"
	"mynute-go/core/test/src/model"

	"github.com/google/uuid"
)

func GetServiceByID(company *model.Company, serviceIDStr string) (*model.Service, error) {
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

