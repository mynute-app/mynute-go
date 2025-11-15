package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Resource struct {
	BaseModel
	Name        string             `json:"name" gorm:"unique;not null"`
	Description string             `json:"description"`
	Table       string             `json:"table"`
	References  ResourceReferences `gorm:"type:jsonb" json:"references"`
}

func (Resource) TableName() string  { return "public.resources" }
func (Resource) SchemaType() string { return "public" }

// --- Define ResourceReference first ---
type ResourceReference struct {
	DatabaseKey string `json:"database_key"` // The key in the database, e.g. "id", "tax_id".
	RequestKey  string `json:"request_key"`  // The key in the request body, query, path, header...
	RequestRef  string `json:"request_ref"`  // "query", "body", "header", "path".
}

// --- Define the custom slice type ---
type ResourceReferences []ResourceReference

// --- Implement the Valuer interface for ResourceReferences ---
func (r ResourceReferences) Value() (driver.Value, error) {
	if len(r) == 0 {
		return nil, nil // Store empty slice as NULL in DB
	}
	// Marshal the slice into JSON bytes
	return json.Marshal(r)
}

// --- Implement the Scanner interface for ResourceReferences ---
func (r *ResourceReferences) Scan(value any) error {
	// Get bytes from the database value
	bytes, ok := value.([]byte)
	if !ok {
		// Handle the case where the database value might be nil
		if value == nil {
			*r = nil // Set the slice to nil if DB value is NULL
			return nil
		}
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	// Handle empty byte slice if necessary (e.g., if DB stores '' instead of NULL/[])
	if len(bytes) == 0 {
		*r = nil
		return nil
	}

	// Unmarshal the JSON bytes into the slice (use pointer *r)
	return json.Unmarshal(bytes, r)
}

func SingleQueryRef() ResourceReference {
	return ResourceReference{
		DatabaseKey: "id",
		RequestKey:  "id",
		RequestRef:  "query",
	}
}

func SinglePathRef() ResourceReference {
	return ResourceReference{
		DatabaseKey: "id",
		RequestKey:  "id",
		RequestRef:  "path",
	}
}

func MultipleQueryRef(name, dbKey string) ResourceReference {
	return ResourceReference{
		DatabaseKey: dbKey,
		RequestKey:  name,
		RequestRef:  "query",
	}
}

func MultiplePathRef(name, dbKey string) ResourceReference {
	return ResourceReference{
		DatabaseKey: dbKey,
		RequestKey:  name,
		RequestRef:  "path",
	}
}

func MultipleBodyRef(name, dbKey string) ResourceReference {
	return ResourceReference{
		DatabaseKey: dbKey,
		RequestKey:  name,
		RequestRef:  "body",
	}
}

var SingleIdQueryRef = ResourceReference{
	DatabaseKey: "id",
	RequestKey:  "id",
	RequestRef:  "query",
}
