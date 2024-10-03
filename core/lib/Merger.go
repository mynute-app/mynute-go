package lib

import "encoding/json"

func MergeMapIntoInterface(model interface{}, changes map[string]interface{}) error {
	// Marshal the changes map into JSON
	changesJSON, err := json.Marshal(changes)
	if err != nil {
		return err
	}

	// Unmarshal the JSON into the model struct
	if err := json.Unmarshal(changesJSON, model); err != nil {
		return err
	}

	return nil
}