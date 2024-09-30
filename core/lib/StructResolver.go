package lib

import (
	"errors"
	"reflect"
)

// Function to check the struct pointer and resolve the struct
func ResolvePointerStruct(v interface{}) (reflect.Value, error) {
	val := reflect.ValueOf(v)
	// Ensure the value is a pointer to a struct
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return val, errors.New("interface is not a pointer to a struct")
	}
	return val.Elem(), nil
}

// Function to check and resolve a struct
func ResolveStruct(v interface{}) (reflect.Value, error) {
	val := reflect.ValueOf(v)
	// Ensure the value is a struct
	if val.Kind() != reflect.Struct {
		return val, errors.New("interface is not a struct")
	}
	return val, nil
}

