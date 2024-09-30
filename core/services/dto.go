package services

import (
	"agenda-kaki-go/core/lib"
	"errors"
	"reflect"
)

// ConvertToDTO recursively converts a source struct to a destination DTO struct
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

	// Call the recursive function to copy the fields
	if err := copyMatchingFields(sourceVal, dtoVal); err != nil {
		return err
	}

	return nil
}

// copyMatchingFields is a helper function that recursively copies fields from the source to the destination
func copyMatchingFields(sourceVal, dtoVal reflect.Value) error {
	// Iterate through the fields of the destination (DTO) struct
	for i := 0; i < dtoVal.NumField(); i++ {
		dtoField := dtoVal.Field(i)
		dtoFieldType := dtoVal.Type().Field(i)

		// Look for a matching field in the source by name
		sourceField := sourceVal.FieldByName(dtoFieldType.Name)

		// If the source field is valid
		if sourceField.IsValid() {
			// Check if the field is a struct, and if so, recursively copy its fields
			if sourceField.Kind() == reflect.Struct && dtoField.Kind() == reflect.Struct {
				// Recursively copy nested struct fields
				if err := copyMatchingFields(sourceField, dtoField); err != nil {
					return err
				}
			} else if sourceField.Kind() == reflect.Slice && dtoField.Kind() == reflect.Slice {
				// Handle slice fields
				if err := copySlice(sourceField, dtoField); err != nil {
					return err
				}
			} else if sourceField.Type() == dtoField.Type() {
				// If the types match and it's not a struct or slice, directly copy the value
				if dtoField.CanSet() {
					dtoField.Set(sourceField)
				} else {
					return errors.New("cannot set field: " + dtoFieldType.Name)
				}
			}
		}
	}

	return nil
}

// copySlice handles copying slices from the source to the destination
func copySlice(sourceField, dtoField reflect.Value) error {
	// Initialize the destination slice
	if sourceField.IsNil() {
		dtoField.Set(reflect.MakeSlice(dtoField.Type(), 0, 0))
	} else {
		// Create a new slice with the same type as the destination
		newSlice := reflect.MakeSlice(dtoField.Type(), sourceField.Len(), sourceField.Cap())

		// Copy each element from the source slice to the destination slice
		for i := 0; i < sourceField.Len(); i++ {
			sourceElem := sourceField.Index(i)
			destElem := newSlice.Index(i)

			// Handle structs in the slice (recursively copy fields if they are structs)
			if sourceElem.Kind() == reflect.Struct && destElem.Kind() == reflect.Struct {
				if err := copyMatchingFields(sourceElem, destElem); err != nil {
					return err
				}
			} else if sourceElem.Type() == destElem.Type() {
				destElem.Set(sourceElem)
			}
		}

		// Set the destination field with the newly populated slice
		dtoField.Set(newSlice)
	}

	return nil
}
