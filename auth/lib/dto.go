package lib

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

// ResolvePointerStruct checks the struct pointer and resolves the struct or slice
func ResolvePointerStruct(v any) (reflect.Value, error) {
	val := reflect.ValueOf(v)

	// Ensure the value is a pointer
	if val.Kind() != reflect.Ptr {
		return val, errors.New("interface is not a pointer")
	}

	// If it's a pointer to a slice, return the slice
	if val.Elem().Kind() == reflect.Slice {
		return val.Elem(), nil
	}

	// If it's a pointer to a struct, return the struct
	if val.Elem().Kind() == reflect.Struct {
		return val.Elem(), nil
	}

	fmt.Printf("Resolved DTO Struct: %+v\n", val.Elem().Interface())

	return val, errors.New("interface is not a pointer to a struct or slice")
}

// ResolveStruct checks and resolves a struct or slice
func ResolveStruct(v any) (reflect.Value, error) {
	val := reflect.ValueOf(v)

	// Dereference the pointer if it's a pointer
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Ensure the value is a struct or a slice
	if val.Kind() == reflect.Struct || val.Kind() == reflect.Slice {
		return val, nil
	}

	return val, errors.New("interface is not a struct or slice")
}

// ConvertToDTO recursively converts a source struct or slice of structs to a destination DTO struct or slice of DTOs
func ParseToDTO(source any, dto any) error {
	sourceVal := reflect.ValueOf(source)
	dtoVal := reflect.ValueOf(dto)

	// Handle if the source is a slice
	if sourceVal.Kind() == reflect.Slice && dtoVal.Kind() == reflect.Ptr && dtoVal.Elem().Kind() == reflect.Slice {
		// It's a slice, handle accordingly
		return copySlice(sourceVal, dtoVal.Elem())
	}

	// Resolve and validate the destination (DTO) as a pointer to a struct or slice
	dtoVal, err := ResolvePointerStruct(dto)
	if err != nil {
		return err
	}

	// Resolve and validate the source as a struct or slice
	sourceVal, err = ResolveStruct(source)
	if err != nil {
		return err
	}

	// Call the recursive function to copy the fields
	if err := copyMatchingFields(sourceVal, dtoVal); err != nil {
		return err
	}

	return nil
}

// Helper function to copy slice elements
func copySlice(sourceVal, dtoVal reflect.Value) error {
	if sourceVal.IsNil() {
		dtoVal.Set(reflect.MakeSlice(dtoVal.Type(), 0, 0))
		return nil
	}

	// Initialize the destination slice with the same length
	dtoVal.Set(reflect.MakeSlice(dtoVal.Type(), sourceVal.Len(), sourceVal.Len()))

	for i := 0; i < sourceVal.Len(); i++ {
		sourceElem := sourceVal.Index(i)
		destElem := dtoVal.Index(i)

		// If the elements are structs, recursively copy their fields
		if sourceElem.Kind() == reflect.Struct && destElem.Kind() == reflect.Struct {
			if err := copyMatchingFields(sourceElem, destElem); err != nil {
				return err
			}
		} else if sourceElem.Type() == destElem.Type() {
			// Directly copy if the types match
			destElem.Set(sourceElem)
		}
	}

	return nil
}

// copyMatchingFields is a helper function that recursively copies fields from the source to the destination
func copyMatchingFields(sourceVal, dtoVal reflect.Value) error {
	// Handle the case where either source or destination is a slice
	if sourceVal.Kind() == reflect.Slice && dtoVal.Kind() == reflect.Slice {
		return copySlice(sourceVal, dtoVal)
	}

	// Ensure both are structs before calling NumField
	if sourceVal.Kind() != reflect.Struct || dtoVal.Kind() != reflect.Struct {
		return errors.New("both source and destination must be structs or slices")
	}

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
				// Special handling for time.Time
				if sourceField.Type() == reflect.TypeOf(time.Time{}) && dtoField.Type() == reflect.TypeOf(time.Time{}) {
					if dtoField.CanSet() {
						dtoField.Set(sourceField)
					} else {
						return errors.New("cannot set field: " + dtoFieldType.Name)
					}
				} else {
					if err := copyMatchingFields(sourceField, dtoField); err != nil {
						return err
					}
				}
			} else if sourceField.Kind() == reflect.Slice && dtoField.Kind() == reflect.Slice {
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

