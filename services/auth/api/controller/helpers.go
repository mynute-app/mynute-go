package controller

import (
	"fmt"
	"mynute-go/services/auth/api/lib"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// =====================
// SHARED HELPERS
// =====================

func CreateUser(c *fiber.Ctx, modelInstance any) error {
	// Get database session
	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	// Parse the request body into the model
	if err := c.BodyParser(modelInstance); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	// Create the user record
	if err := tx.Create(modelInstance).Error; err != nil {
		// Check for unique constraint violation
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("user with this email already exists"))
		}
		return lib.Error.General.CreatedError.WithError(err)
	}

	return nil
}

func List(c *fiber.Ctx, model any, result any) (int, int, error) {

	// Parse pagination parameters
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	// Validate pagination parameters
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	tx, err := lib.Session(c)
	if err != nil {
		return 0, 0, err
	}

	if err := tx.Model(model).
		Limit(limit).
		Offset(offset).
		Find(result).Error; err != nil {
		return 0, 0, lib.Error.General.InternalError.WithError(err)
	}

	return limit, offset, nil
}

func GetOneBy(param string, c *fiber.Ctx, modelInstance any) error {
	// Get database session
	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	// Get the parameter value from context
	paramValue := c.Params(param)
	if paramValue == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("missing parameter: %s", param))
	}

	// Validate UUID if param is "id"
	if param == "id" {
		if _, err := uuid.Parse(paramValue); err != nil {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid UUID format"))
		}
	}

	// Query the database
	query := tx.Model(modelInstance).Where(param+" = ?", paramValue)

	// Execute the query
	if err := query.First(modelInstance).Error; err != nil {
		if err.Error() == "record not found" {
			return lib.Error.General.RecordNotFound
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

func UpdateOneById(c *fiber.Ctx, modelInstance any) error {
	// Get database session
	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	// Get ID from params
	id := c.Params("id")
	if id == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("missing id parameter"))
	}

	// Validate UUID
	if _, err := uuid.Parse(id); err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid UUID format"))
	}

	// Parse updates from body
	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	// Update the record
	if err := tx.Model(modelInstance).Where("id = ?", id).Updates(updates).Error; err != nil {
		// Check for unique constraint violation
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("a record with this unique value already exists"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	// Fetch the updated record
	if err := tx.Where("id = ?", id).First(modelInstance).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

func DeleteOneById(c *fiber.Ctx, modelInstance any) error {
	// Get database session
	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	// Get ID from params
	id := c.Params("id")
	if id == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("missing id parameter"))
	}

	// Validate UUID
	if _, err := uuid.Parse(id); err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid UUID format"))
	}

	// Delete the record
	if err := tx.Where("id = ?", id).Delete(modelInstance).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Return success
	return lib.ResponseFactory(c).Send(200, map[string]string{"message": "Deleted successfully"})
}

