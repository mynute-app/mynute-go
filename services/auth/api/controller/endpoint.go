package controller

import (
	"fmt"
	"mynute-go/services/auth/api/lib"
	authModel "mynute-go/services/auth/config/db/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =====================
// ENDPOINT MANAGEMENT
// =====================

// ListEndpoints retrieves all endpoints
//
//	@Summary		List all endpoints
//	@Description	Retrieve all endpoints
//	@Tags			Endpoints
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Produce		json
//	@Success		200	{array}		authModel.EndPoint
//	@Failure		400	{object}	map[string]string
//	@Router			/endpoints [get]
func ListEndpoints(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var endpoints []authModel.EndPoint
	if err := tx.Preload("Resource").Find(&endpoints).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return c.JSON(endpoints)
}

// GetEndpointById retrieves an endpoint by ID
//
//	@Summary		Get endpoint by ID
//	@Description	Retrieve an endpoint by its ID
//	@Tags			Endpoints
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			id				path		string	true	"Endpoint ID"
//	@Produce		json
//	@Success		200	{object}	authModel.EndPoint
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/endpoints/{id} [get]
func GetEndpointById(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	endpointID := c.Params("id")
	if endpointID == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("endpoint ID is required"))
	}

	id, err := uuid.Parse(endpointID)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid endpoint ID format"))
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var endpoint authModel.EndPoint
	if err := tx.Preload("Resource").Where("id = ?", id).First(&endpoint).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.ResourceNotFoundError.WithError(fmt.Errorf("endpoint not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	return c.JSON(endpoint)
}

// CreateEndpoint creates a new endpoint
//
//	@Summary		Create endpoint
//	@Description	Create a new endpoint
//	@Tags			Endpoints
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string					true	"X-Auth-Token"
//	@Accept			json
//	@Produce		json
//	@Param			endpoint		body		EndpointCreateRequest	true	"Endpoint data"
//	@Success		201				{object}	authModel.EndPoint
//	@Failure		400				{object}	map[string]string
//	@Router			/endpoints [post]
func CreateEndpoint(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	var req EndpointCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	// Validate method
	validMethods := map[string]bool{
		"GET": true, "POST": true, "PUT": true, "PATCH": true, "DELETE": true,
	}
	if !validMethods[req.Method] {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid HTTP method"))
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	// Check if endpoint with same method and path already exists
	var existing authModel.EndPoint
	if err := tx.Where("method = ? AND path = ?", req.Method, req.Path).First(&existing).Error; err == nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("endpoint with method %s and path %s already exists", req.Method, req.Path))
	} else if err != gorm.ErrRecordNotFound {
		return lib.Error.General.InternalError.WithError(err)
	}

	endpoint := authModel.EndPoint{
		BaseModel:        authModel.BaseModel{ID: uuid.New()},
		ControllerName:   req.ControllerName,
		Description:      req.Description,
		Method:           req.Method,
		Path:             req.Path,
		DenyUnauthorized: req.DenyUnauthorized,
		NeedsCompanyId:   req.NeedsCompanyId,
	}

	// Handle resource ID if provided
	if req.ResourceID != nil {
		resourceID, err := uuid.Parse(*req.ResourceID)
		if err != nil {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid resource ID"))
		}

		// Check if resource exists
		var resource authModel.Resource
		if err := tx.Where("id = ?", resourceID).First(&resource).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return lib.Error.General.BadRequest.WithError(fmt.Errorf("resource not found"))
			}
			return lib.Error.General.InternalError.WithError(err)
		}

		endpoint.ResourceID = &resourceID
	}

	// Temporarily allow endpoint creation
	authModel.AllowEndpointCreation = true
	defer func() { authModel.AllowEndpointCreation = false }()

	if err := tx.Create(&endpoint).Error; err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}

	// Load the resource relation
	if err := tx.Preload("Resource").Where("id = ?", endpoint.ID).First(&endpoint).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return c.Status(201).JSON(endpoint)
}

// UpdateEndpointById updates an endpoint by ID
//
//	@Summary		Update endpoint
//	@Description	Update an endpoint
//	@Tags			Endpoints
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string					true	"X-Auth-Token"
//	@Accept			json
//	@Produce		json
//	@Param			id				path		string					true	"Endpoint ID"
//	@Param			endpoint		body		EndpointUpdateRequest	true	"Endpoint data"
//	@Success		200				{object}	authModel.EndPoint
//	@Failure		400				{object}	map[string]string
//	@Failure		404				{object}	map[string]string
//	@Router			/endpoints/{id} [patch]
func UpdateEndpointById(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	endpointID := c.Params("id")
	if endpointID == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("endpoint ID is required"))
	}

	id, err := uuid.Parse(endpointID)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid endpoint ID format"))
	}

	var req EndpointUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var endpoint authModel.EndPoint
	if err := tx.Where("id = ?", id).First(&endpoint).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.ResourceNotFoundError.WithError(fmt.Errorf("endpoint not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	// Update fields if provided
	if req.ControllerName != nil {
		endpoint.ControllerName = *req.ControllerName
	}
	if req.Description != nil {
		endpoint.Description = *req.Description
	}
	if req.Method != nil {
		validMethods := map[string]bool{
			"GET": true, "POST": true, "PUT": true, "PATCH": true, "DELETE": true,
		}
		if !validMethods[*req.Method] {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid HTTP method"))
		}
		endpoint.Method = *req.Method
	}
	if req.Path != nil {
		endpoint.Path = *req.Path
	}
	if req.DenyUnauthorized != nil {
		endpoint.DenyUnauthorized = *req.DenyUnauthorized
	}
	if req.NeedsCompanyId != nil {
		endpoint.NeedsCompanyId = *req.NeedsCompanyId
	}
	if req.ResourceID != nil {
		if *req.ResourceID == "" {
			endpoint.ResourceID = nil
		} else {
			resourceID, err := uuid.Parse(*req.ResourceID)
			if err != nil {
				return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid resource ID"))
			}

			// Check if resource exists
			var resource authModel.Resource
			if err := tx.Where("id = ?", resourceID).First(&resource).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return lib.Error.General.BadRequest.WithError(fmt.Errorf("resource not found"))
				}
				return lib.Error.General.InternalError.WithError(err)
			}

			endpoint.ResourceID = &resourceID
		}
	}

	if err := tx.Save(&endpoint).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Load the resource relation
	if err := tx.Preload("Resource").Where("id = ?", endpoint.ID).First(&endpoint).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return c.JSON(endpoint)
}

// DeleteEndpointById deletes an endpoint by ID
//
//	@Summary		Delete endpoint
//	@Description	Delete an endpoint by its ID
//	@Tags			Endpoints
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			id				path		string	true	"Endpoint ID"
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/endpoints/{id} [delete]
func DeleteEndpointById(c *fiber.Ctx) error {
	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	endpointID := c.Params("id")
	if endpointID == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("endpoint ID is required"))
	}

	id, err := uuid.Parse(endpointID)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid endpoint ID format"))
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var endpoint authModel.EndPoint
	if err := tx.Where("id = ?", id).First(&endpoint).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.ResourceNotFoundError.WithError(fmt.Errorf("endpoint not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	// Check if there are policies referencing this endpoint
	var policyCount int64
	if err := tx.Model(&authModel.PolicyRule{}).Where("end_point_id = ?", id).Count(&policyCount).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	if policyCount > 0 {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("cannot delete endpoint: %d policies are referencing it", policyCount))
	}

	if err := tx.Delete(&endpoint).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return c.JSON(fiber.Map{
		"message": "Endpoint deleted successfully",
	})
}

// =====================
// REQUEST TYPES
// =====================

type EndpointCreateRequest struct {
	ControllerName   string  `json:"controller_name" validate:"required,min=3,max=100"`
	Description      string  `json:"description"`
	Method           string  `json:"method" validate:"required,oneof=GET POST PUT PATCH DELETE"`
	Path             string  `json:"path" validate:"required"`
	DenyUnauthorized bool    `json:"deny_unauthorized"`
	NeedsCompanyId   bool    `json:"needs_company_id"`
	ResourceID       *string `json:"resource_id,omitempty" validate:"omitempty,uuid"`
}

type EndpointUpdateRequest struct {
	ControllerName   *string `json:"controller_name,omitempty" validate:"omitempty,min=3,max=100"`
	Description      *string `json:"description,omitempty"`
	Method           *string `json:"method,omitempty" validate:"omitempty,oneof=GET POST PUT PATCH DELETE"`
	Path             *string `json:"path,omitempty"`
	DenyUnauthorized *bool   `json:"deny_unauthorized,omitempty"`
	NeedsCompanyId   *bool   `json:"needs_company_id,omitempty"`
	ResourceID       *string `json:"resource_id,omitempty" validate:"omitempty,uuid"`
}
