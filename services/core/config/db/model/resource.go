package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// Resource model for Core Service
// /auth service handles resource-based authorization logic
type Resource struct {
	BaseModel
	Name        string             `json:"name" gorm:"unique;not null"`
	Description string             `json:"description"`
	Table       string             `json:"table"`
	References  ResourceReferences `gorm:"type:jsonb" json:"references"`
}

func (Resource) TableName() string  { return "public.resources" }
func (Resource) SchemaType() string { return "public" }

// ResourceReference defines how to find a resource from request data
type ResourceReference struct {
	DatabaseKey string `json:"database_key"` // The key in the database, e.g. "id", "tax_id"
	RequestKey  string `json:"request_key"`  // The key in the request body, query, path, header
	RequestRef  string `json:"request_ref"`  // "query", "body", "header", "path"
}

// ResourceReferences is a slice of ResourceReference
type ResourceReferences []ResourceReference

// Value implements the driver.Valuer interface for JSONB storage
func (r ResourceReferences) Value() (driver.Value, error) {
	if len(r) == 0 {
		return nil, nil
	}
	return json.Marshal(r)
}

// Scan implements the sql.Scanner interface for JSONB retrieval
func (r *ResourceReferences) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		if value == nil {
			*r = nil
			return nil
		}
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	if len(bytes) == 0 {
		*r = nil
		return nil
	}

	return json.Unmarshal(bytes, r)
}

// Helper functions to create common resource reference patterns
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
