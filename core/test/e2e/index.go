package e2e

import (
	"fmt"
	"mynute-go/core/config/namespace"
	handlerT "mynute-go/core/test/handlers"
	modelT "mynute-go/core/test/models"
	"reflect"
)

type WorkScheduleModelTarget interface {
	GetID() string
	GetCompanyID() string
	GetAuthToken() string
	SetWorkRanges([]any)
}

func CreateWorkSchedule[T any](target WorkScheduleModelTarget, entity string, status int, schedule any, x_auth_token string, x_company_id *string) error {
	// Validar WorkRanges
	v := reflect.ValueOf(schedule)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if field := v.FieldByName("WorkRanges"); field.IsNil() {
		return fmt.Errorf("work schedule cannot be nil")
	}

	// Preparar headers
	tCompanyID := target.GetCompanyID()
	cID, err := modelT.Get_x_company_id(x_company_id, &tCompanyID)
	if err != nil {
		return err
	}

	// Requisição
	var updated T
	err = handlerT.NewHttpClient().
		Method("POST").
		URL(fmt.Sprintf("/%s/%s/work_schedule", entity, target.GetID())).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Company, cID).
		Header(namespace.HeadersKey.Auth, target.GetAuthToken()).
		Send(schedule).
		ParseResponse(&updated).
		Error
	if err != nil {
		return fmt.Errorf("failed to create work schedule: %w", err)
	}

	// Extrair WorkRanges
	val := reflect.ValueOf(updated)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	wrField := val.FieldByName("WorkRanges")
	if wrField.IsValid() && wrField.Kind() == reflect.Slice {
		wr := make([]any, wrField.Len())
		for i := 0; i < wrField.Len(); i++ {
			wr[i] = wrField.Index(i).Interface()
		}
		target.SetWorkRanges(wr)
	}

	return nil
}

func GetWorkRange() {
	// This function is intentionally left empty.
	// It serves as a placeholder for future work schedule retrieval.
}

func UpdateWorkRange() {
	// This function is intentionally left empty.
	// It serves as a placeholder for future work schedule updates.
}

func DeleteWorkRange() {
	// This function is intentionally left empty.
	// It serves as a placeholder for future work schedule deletions.
}
