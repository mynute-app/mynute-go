package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v2"
)

// CreateHolidays creates a holiday
//
//	@Summary		Create holiday
//	@Description	Create a holiday
//	@Tags			Holidays
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			holiday	body		DTO.Holidays	true	"Holiday"
//	@Success		201		{object}	DTO.Holidays
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/holidays [post]
func CreateHoliday(c *fiber.Ctx) error {
	var holiday model.Holiday

	if err := Create(c, &holiday); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &holiday, &DTO.Holidays{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// GetHolidayById retrieves a holiday by ID
//
//	@Summary		Get holiday by ID
//	@Description	Retrieve a holiday by its ID
//	@Tags			Holidays
//	@Param			id	path	string	true	"Holiday ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Holidays
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/holidays/{id} [get]
func GetHolidayById(c *fiber.Ctx) error {
	var holiday model.Holiday

	if err := GetOneBy("id", c, &holiday); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &holiday, &DTO.Holidays{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// GetHolidayByName retrieves a holiday by name
//
//	@Summary		Get holiday by name
//	@Description	Retrieve a holiday by its name
//	@Tags			Holidays
//	@Param			name	path	string	true	"Holiday Name"
//	@Produce		json
//	@Success		200	{object}	DTO.Holidays
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/holidays/name/{name} [get]
func GetHolidayByName(c *fiber.Ctx) error {
	var holiday model.Holiday

	if err := GetOneBy("name", c, &holiday); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &holiday, &DTO.Holidays{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// UpdateHolidayById updates a holiday by ID
//
//	@Summary		Update holiday
//	@Description	Update a holiday
//	@Tags			Holidays
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string			true	"Holiday ID"
//	@Param			holiday	body		DTO.Holidays	true	"Holiday"
//	@Success		200		{object}	DTO.Holidays
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/holidays/{id} [patch]
func UpdateHolidayById(c *fiber.Ctx) error {
	var holiday model.Holiday

	if err := UpdateOneById(c, &holiday); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &holiday, &DTO.Holidays{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// DeleteHolidayById deletes a holiday by ID
//
//	@Summary		Delete holiday by ID
//	@Description	Delete a holiday by its ID
//	@Tags			Holidays
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Holiday ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Holidays
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/holidays/{id} [delete]
func DeleteHolidayById(c *fiber.Ctx) error {
	return DeleteOneById(c, &model.Holiday{})
}

// Holidays creates a new holidays_controller
func Holiday(Gorm *handler.Gorm) {
	endpoint := &middleware.Endpoint{DB: Gorm}
	endpoint.BulkRegisterHandler([]fiber.Handler{
		CreateHoliday,
		GetHolidayById,
		GetHolidayByName,
		UpdateHolidayById,
		DeleteHolidayById,
	})
}
