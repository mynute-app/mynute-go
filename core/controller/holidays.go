package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/middleware"
	"agenda-kaki-go/core/service"

	"github.com/gofiber/fiber/v2"
)

type HolidaysController struct {
	service.Base[model.Holidays, DTO.Holidays]
}

// CreateHolidays creates a holiday
//
//	@Summary		Create holiday
//	@Description	Create a holiday
//	@Tags			Holidays
//	@Accept			json
//	@Produce		json
//	@Param			holiday	body		DTO.Holidays	true	"Holiday"
//	@Success		200		{object}	DTO.Holidays
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/holidays [post]
func (cc *HolidaysController) CreateHoliday(c *fiber.Ctx) error {
	return cc.CreateOne(c)
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
func (cc *HolidaysController) GetHolidayById(c *fiber.Ctx) error {
	return cc.GetBy("id", c)
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
func (cc *HolidaysController) GetHolidayByName(c *fiber.Ctx) error {
	return cc.GetBy("name", c)
}

// UpdateHolidayById updates a holiday by ID
//
//	@Summary		Update holiday
//	@Description	Update a holiday
//	@Tags			Holidays
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string			true	"Holiday ID"
//	@Param			holiday	body		DTO.Holidays	true	"Holiday"
//	@Success		200		{object}	DTO.Holidays
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/holidays/{id} [patch]
func (cc *HolidaysController) UpdateHolidayById(c *fiber.Ctx) error {
	return cc.UpdateOneById(c)
}

// DeleteHolidayById deletes a holiday by ID
//
//	@Summary		Delete holiday by ID
//	@Description	Delete a holiday by its ID
//	@Tags			Holidays
//	@Param			id	path	string	true	"Holiday ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Holidays
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/holidays/{id} [delete]
func (cc *HolidaysController) DeleteHolidayById(c *fiber.Ctx) error {
	return cc.DeleteOneById(c)
}

// Holidays creates a new HolidaysController
func Holidays(Gorm *handler.Gorm) *HolidaysController {
	return &HolidaysController{
		Base: service.Base[model.Holidays, DTO.Holidays]{
			Name:         namespace.HolidaysKey.Name,
			Request:      handler.Request(Gorm),
			Middleware:   middleware.Holidays(Gorm),
			Associations: []string{},
		},
	}
}
