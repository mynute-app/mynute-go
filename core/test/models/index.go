package modelT

import (
	"agenda-kaki-go/core/config/namespace"
	handlerT "agenda-kaki-go/core/test/handlers"
	"fmt"
	"reflect"
)

type WorkScheduleModelTarget interface {
	GetID() string
	GetCompanyID() string
	GetAuthToken() string
	SetWorkRanges([]any) error
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

	cID, err := Get_x_company_id(x_company_id, &tCompanyID)
	if err != nil {
		return err
	}

	tAuthToken := target.GetAuthToken()

	token, err := Get_x_auth_token(&x_auth_token, &tAuthToken)
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
		Header(namespace.HeadersKey.Auth, token).
		Send(schedule).
		ParseResponse(&updated).
		Error
	if err != nil {
		return fmt.Errorf("failed to create work schedule: %w", err)
	}

	if err := put_work_schedule(updated, target); err != nil {
		return err
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

func put_work_schedule[T any](updated T, target WorkScheduleModelTarget) error {
	val := reflect.ValueOf(updated)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	wrField := val.FieldByName("WorkRanges")
	if wrField.IsValid() && wrField.Kind() == reflect.Slice {
		wr := make([]any, wrField.Len())
		for i := range wrField.Len() {
			wr[i] = wrField.Index(i).Interface()
		}
		err := target.SetWorkRanges(wr)
		if err != nil {
			return fmt.Errorf("failed to set work ranges: %w", err)
		}
	}
	return nil
}
