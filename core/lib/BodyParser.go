package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// ParseBody function to parse the request body into the provided struct (any)
func BodyParser(b []byte, v any) error {
	// Ensure v is a pointer
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("bodyParser: expected a non-nil pointer but received something else")
	}

	// Check if the body is empty
	if len(b) == 0 {
		return errors.New("empty request body")
	}

	// Unmarshal the JSON body into the provided pointer (v)
	if err := json.Unmarshal(b, v); err != nil {
		fmt.Printf("Error when unmarshalling JSON: %v\n", err.Error())
		return err
	}

	return nil
}
