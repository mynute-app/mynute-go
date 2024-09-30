package services

import (
	"agenda-kaki-go/core/lib"
)

// ConvertToDTO converts a source struct to a destination struct
func ConvertToDTO(source interface{}, dto interface{}) error {
	// Resolve and validate the destination (DTO) as a pointer to a struct
	dtoVal, err := lib.ResolvePointerStruct(dto)
	if err != nil {
		return err
	}

	// Resolve and validate the source as a struct
	sourceVal, err := lib.ResolveStruct(source)
	if err != nil {
		return err
	}

	// Iterate through the fields of the destination (DTO) struct
	for i := 0; i < dtoVal.NumField(); i++ {
		dtoField := dtoVal.Field(i)
		dtoFieldType := dtoVal.Type().Field(i)

		// Look for a matching field in the source by name
		sourceField := sourceVal.FieldByName(dtoFieldType.Name)

		// Check if the source field exists, is valid, and has the same type as the destination field
		if sourceField.IsValid() && sourceField.Type() == dtoField.Type() {
			// Set the destination field with the value from the source field
			dtoField.Set(sourceField)
		}
	}

	return nil
}
