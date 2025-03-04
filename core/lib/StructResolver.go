package lib

import (
	"errors"
	"reflect"
)

// Function to check the struct pointer and resolve the struct or slice
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

	return val, errors.New("interface is not a pointer to a struct or slice")
}

// Function to check and resolve a struct or slice
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
