package controller

import (
	DTO "mynute-go/core/config/api/dto"
	"mynute-go/core/config/db/model"
	"mynute-go/core/handler"
	"mynute-go/core/lib"
	"mynute-go/core/middleware"

	"github.com/gofiber/fiber/v2"
)

// CreateSector creates a sector
//
//	@Summary		Create sector
//	@Description	Create a sector
//	@Tags			Sector
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			sector	body		DTO.Sector	true	"sector"
//	@Success		200		{object}	DTO.Sector
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/sector [post]
func CreateSector(c *fiber.Ctx) error {
	var sector model.Sector

	if err := Create(c, &sector); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &sector, &DTO.Sector{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// GetSectorByName retrieves a sector by ID
//
//	@Summary		Get sector by ID
//	@Description	Retrieve a sector by its ID
//	@Tags			Sector
//	@Param			id	path	string	true	"sector ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Sector
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/sector/name/{name} [get]
func GetSectorByName(c *fiber.Ctx) error {
	var sector model.Sector

	if err := GetOneBy("name", c, &sector, nil, nil); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &sector, &DTO.Sector{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// GetSectorById retrieves a sector by ID
//
//	@Summary		Get sector by ID
//	@Description	Retrieve a sector by its ID
//	@Tags			Sector
//	@Param			id	path	string	true	"sector ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Sector
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/sector/{id} [get]
func GetSectorById(c *fiber.Ctx) error {
	var sector model.Sector

	if err := GetOneBy("id", c, &sector, nil, nil); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &sector, &DTO.Sector{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil

}

// UpdateSectorById updates a sector by ID
//
//	@Summary		Update sector by ID
//	@Description	Update a sector by its ID
//	@Tags			Sector
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"sector ID"
//	@Accept			json
//	@Produce		json
//	@Param			sector	body		DTO.Sector	true	"sector"
//	@Success		200		{object}	DTO.Sector
//	@Failure		404		{object}	DTO.ErrorResponse
//	@Router			/sector/{id} [patch]
func UpdateSectorById(c *fiber.Ctx) error {
	var sector model.Sector

	if err := UpdateOneById(c, &sector, nil); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &sector, &DTO.Sector{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// DeleteSectorById deletes a sector by ID
//
//	@Summary		Delete sector by ID
//	@Description	Delete a sector by its ID
//	@Tags			Sector
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"sector ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Sector
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/sector/{id} [delete]
func DeleteSectorById(c *fiber.Ctx) error {
	return DeleteOneById(c, &model.Sector{})
}

func Sector(Gorm *handler.Gorm) {
	endpoint := &middleware.Endpoint{DB: Gorm}
	endpoint.BulkRegisterHandler([]fiber.Handler{
		CreateSector,
		GetSectorByName,
		GetSectorById,
		UpdateSectorById,
		DeleteSectorById,
	})
}
