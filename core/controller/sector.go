package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/service"

	"github.com/gofiber/fiber/v2"
)

type sector_controller struct {
	service.Base[model.Sector, DTO.Sector]
}

// CreateSector creates a sector
//
//	@Summary		Create sector
//	@Description	Create a sector
//	@Tags			Sector
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Authorization"
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			sector	body		DTO.Sector	true	"sector"
//	@Success		201		{object}	DTO.Sector
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/sector [post]
func (cc *sector_controller) CreateSector(c *fiber.Ctx) error {
	return cc.CreateOne(c)
}

// GetSectorByName retrieves a sector by ID
//
//	@Summary		Get sector by ID
//	@Description	Retrieve a sector by its ID
//	@Tags			Sector
//	@Param			id	path	string	true	"sector ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Sector
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/sector/name/{name} [get]
func (cc *sector_controller) GetSectorByName(c *fiber.Ctx) error {
	return cc.GetBy("name", c)
}

// GetSectorById retrieves a sector by ID
//
//	@Summary		Get sector by ID
//	@Description	Retrieve a sector by its ID
//	@Tags			Sector
//	@Param			id	path	string	true	"sector ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Sector
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/sector/{id} [get]
func (cc *sector_controller) GetSectorById(c *fiber.Ctx) error {
	return cc.GetBy("id", c)
}

// UpdateSectorById updates a sector by ID
//
//	@Summary		Update sector by ID
//	@Description	Update a sector by its ID
//	@Tags			Sector
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Authorization"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"sector ID"
//	@Accept			json
//	@Produce		json
//	@Param			sector	body		DTO.Sector	true	"sector"
//	@Success		200		{object}	DTO.Sector
//	@Failure		404		{object}	DTO.ErrorResponse
//	@Router			/sector/{id} [patch]
func (cc *sector_controller) UpdateSectorById(c *fiber.Ctx) error {
	return cc.UpdateOneById(c)
}

// DeleteSectorById deletes a sector by ID
//
//	@Summary		Delete sector by ID
//	@Description	Delete a sector by its ID
//	@Tags			Sector
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Authorization"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"sector ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Sector
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/sector/{id} [delete]
func (cc *sector_controller) DeleteSectorById(c *fiber.Ctx) error {
	return cc.DeleteOneById(c)
}

func Sector(Gorm *handler.Gorm) *sector_controller {
	sc := &sector_controller{
		Base: service.Base[model.Sector, DTO.Sector]{
			Name:         namespace.SectorKey.Name,
			Request:      handler.Request(Gorm),
			Associations: []string{},
		},
	}
	route := &handler.Route{DB: Gorm.DB}
	route.Register("/sector", "POST", "private", sc.CreateSector, "Creates a company sector").Save()
	route.Register("/sector/:id", "GET", "public", sc.GetSectorById, "Retrieves a company sector by ID").Save()
	route.Register("/sector/name/:name", "GET", "public", sc.GetSectorByName, "Retrieves a company sector by name").Save()
	route.Register("/sector/:id", "PATCH", "private", sc.UpdateSectorById, "Updates a company sector by ID").Save()
	route.Register("/sector/:id", "DELETE", "private", sc.DeleteSectorById, "Deletes a company sector by ID").Save()
	return sc
}
