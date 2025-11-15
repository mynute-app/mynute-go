package model

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var AllowEndpointCreation = false

// EndPoint model for Core Service
// /auth service handles endpoint-based authorization logic
type EndPoint struct {
	BaseModel
	ControllerName   string     `gorm:"type:varchar(100)" json:"controller_name"`
	Description      string     `gorm:"type:text" json:"description"`
	Method           string     `gorm:"type:varchar(6)" json:"method"`
	Path             string     `gorm:"type:text" json:"path"`
	DenyUnauthorized bool       `gorm:"default:false" json:"deny_unauthorized"`
	NeedsCompanyId   bool       `gorm:"default:false" json:"needs_company_id"`
	ResourceID       *uuid.UUID `gorm:"type:uuid" json:"resource_id"`
	Resource         *Resource  `gorm:"foreignKey:ResourceID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"resource"`
}

func (EndPoint) TableName() string  { return "public.endpoints" }
func (EndPoint) SchemaType() string { return "public" }

func (EndPoint) Indexes() map[string]string {
	return map[string]string{
		"idx_method_path": "CREATE UNIQUE INDEX idx_method_path ON routes (method, path)",
	}
}

func (r *EndPoint) BeforeCreate(tx *gorm.DB) error {
	if !AllowEndpointCreation {
		panic("EndPoint creation is not allowed")
	}
	return nil
}

type EndpointCfg struct {
	AllowCreation bool // Allow creation of endpoints
}

// EndPoints processes the given endpoints and returns them along with a cleanup function.
func EndPoints(endpoints []*EndPoint, cfg *EndpointCfg, db *gorm.DB) ([]*EndPoint, func(), error) {
	AllowEndpointCreation = cfg.AllowCreation

	// Retrieve resources from database
	resourceMap := map[string]uuid.UUID{}
	var resources []Resource
	if err := db.Find(&resources).Error; err != nil {
		return nil, nil, err
	}
	for _, r := range resources {
		resourceMap[r.Table] = r.ID
	}

	for _, edp := range endpoints {
		if edp.Resource != nil {
			if id, ok := resourceMap[edp.Resource.Table]; ok {
				edp.ResourceID = &id
			} else {
				return nil, nil, fmt.Errorf("resource not found for table: %s", edp.Resource.Table)
			}
			edp.Resource = nil
		}
	}

	deferFnc := func() {
		AllowEndpointCreation = false
	}

	return endpoints, deferFnc, nil
}

// LoadEndpointIDs loads the IDs of the endpoints from the database
// and updates the endpoint variables with their corresponding IDs.
func LoadEndpointIDs(endpoints []*EndPoint, db *gorm.DB) error {
	for _, ep := range endpoints {
		var existing EndPoint
		if err := db.
			Where("method = ? AND path = ?", ep.Method, ep.Path).
			First(&existing).Error; err != nil {
			return fmt.Errorf("failed to load endpoint ID for %s %s: %w", ep.Method, ep.Path, err)
		}
		ep.ID = existing.ID
	}
	return nil
}
