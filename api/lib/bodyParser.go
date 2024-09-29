package lib

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"errors"
)

// ParseBody function to parse the request body into the provided struct (interface{})
func ParseBody(r *fasthttp.Request, v interface{}) error {
	// Get the request body as a byte slice
	body := r.Body()

	// Check if the body is empty
	if len(body) == 0 {
		return errors.New("empty request body")
	}

	// Unmarshal the JSON body into the provided interface (v)
	if err := json.Unmarshal(body, v); err != nil {
		return err // Return the JSON unmarshalling error
	}

	return nil
}
