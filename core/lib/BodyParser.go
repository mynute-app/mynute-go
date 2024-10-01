package lib

import (
	"encoding/json"
	"errors"
)

// ParseBody function to parse the request body into the provided struct (interface{})
func BodyParser(b []byte, v interface{}) error {
	// Check if the body is empty
	if len(b) == 0 {
		return errors.New("empty request body")
	}

	// Unmarshal the JSON body into the provided interface (v)
	if err := json.Unmarshal(b, v); err != nil {
		return err // Return the JSON unmarshalling error
	}

	return nil
}