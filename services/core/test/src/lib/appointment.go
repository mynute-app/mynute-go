package lib

import (
	"fmt"
	"mynute-go/services/core/src/config/namespace"
	"mynute-go/services/core/test/src/handler"
	"mynute-go/services/core/test/src/model"
)

func GetAppointment(s int, appointment_id string, company_id, token string, a *model.Appointment) error {
	http := handler.NewHttpClient()
	if err := http.
		Method("GET").
		URL("/appointment/"+appointment_id).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, token).
		Header(namespace.HeadersKey.Company, company_id).
		Send(nil).Error; err != nil {
		return fmt.Errorf("failed to get appointment by ID: %w", err)
	}
	if a != nil {
		if err := http.ParseResponse(&a.Created).Error; err != nil {
			return fmt.Errorf("failed to parse appointment response: %w", err)
		}
	}
	return nil
}
