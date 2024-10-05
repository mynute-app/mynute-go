package controllers

import (
	"agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v3"
)

type Company struct {
	Gorm        *handlers.Gorm
	Middleware  *middleware.Company
	HttpHandler *handlers.HTTP
}

// func (cc *Company) getBy(paramKey string, c fiber.Ctx) error {
// 	var model models.Company

// 	assocs := []string{"CompanyTypes"}

// 	paramVal := c.Params(paramKey)

// 	if err := cc.Gorm.GetOneBy(paramKey, paramVal, &model, assocs); err != nil {
// 		return lib.Fiber404(c)
// 	}

// 	var dto DTO.Company

// 	if err := lib.ParseToDTO(model, &dto); err != nil {
// 		return lib.Fiber500(c, err)
// 	}

// 	if err := c.JSON(dto); err != nil {
// 		return lib.Fiber500(c, err)
// 	}

// 	return nil
// }

// func (cc *Company) updateBy(paramKey string, c fiber.Ctx) error {
// 	var changes map[string]interface{}

// 	if err := lib.BodyParser(c.Body(), &changes); err != nil {
// 		return lib.FiberError(400, c, err)
// 	}

// 	if err := cc.Middleware.Update(changes); err != nil {
// 		return lib.FiberError(400, c, err)
// 	}

// 	var model models.Company

// 	assocs := []string{"CompanyTypes"}
// 	paramVal := c.Params(paramKey)

// 	if err := cc.Gorm.UpdateOneBy(paramKey, paramVal, &model, changes, assocs); err != nil {
// 		return lib.FiberError(400, c, err)
// 	}

// 	var dto DTO.Company

// 	if err := lib.ParseToDTO(model, &dto); err != nil {
// 		return lib.FiberError(500, c, err)
// 	}

// 	if err := c.JSON(dto); err != nil {
// 		log.Printf("An internal error occurred! %v", err)
// 		return err
// 	}

// 	return nil
// }

// func (cc *Company) deleteBy(paramKey string, c fiber.Ctx) error {
// 	var model models.Company

// 	paramVal := c.Params(paramKey)

// 	if err := cc.Gorm.DeleteOneBy(paramKey, paramVal, &model); err != nil {
// 		return lib.FiberError(400, c, err)
// 	}

// 	return nil
// }

// func (cc *Company) Create(c fiber.Ctx) error {
// 	var model models.Company

// 	log.Printf("lib.BodyParser")
// 	if err := lib.BodyParser(c.Body(), &model); err != nil {
// 		return lib.FiberError(500, c, err)
// 	}
// 	log.Printf("Middleware.Create")
// 	if err := cc.Middleware.Create(model); err != nil {
// 		return lib.FiberError(400, c, err)
// 	}

// 	assocs := []string{"CompanyTypes"}

// 	log.Printf("DB.Create")
// 	log.Printf("model: %+v", model)
// 	if err := cc.Gorm.Create(&model, assocs); err != nil {
// 		return lib.FiberError(400, c, err)
// 	}

// 	var dto DTO.Company

// 	log.Printf("lib.ParseToDTO")
// 	if err := lib.ParseToDTO(model, &dto); err != nil {
// 		return lib.FiberError(500, c, err)
// 	}

// 	log.Printf("c.JSON")
// 	if err := c.JSON(dto); err != nil {
// 		log.Printf("An internal error occurred! %v", err)
// 		return err
// 	}

// 	return nil
// }

// func (cc *Company) GetAll(c fiber.Ctx) error {
// 	var model []models.Company

// 	assocs := []string{"CompanyTypes"}

// 	if err := cc.Gorm.GetAll(&model, assocs); err != nil {
// 		return lib.FiberError(404, c, err)
// 	}

// 	var dto []DTO.Company

// 	if err := lib.ParseToDTO(model, &dto); err != nil {
// 		return lib.FiberError(500, c, err)
// 	}

// 	if err := c.JSON(dto); err != nil {
// 		log.Printf("An internal error occurred! %v", err)
// 		return err
// 	}

// 	return nil
// }

func (cc *Company) getBy(paramKey string, c fiber.Ctx) error {
	var model models.Company
	var dto DTO.Company
	var assocs = []string{"CompanyTypes"}

	cc.HttpHandler.
		Model(&model).
		DTO(&dto).
		Assoc(assocs).
		FiberCtx(c).
		GetOneBy(paramKey)

	return nil
}

func (cc *Company) updateBy(paramKey string, c fiber.Ctx) error {
	var model models.Company
	var dto DTO.Company
	var assocs = []string{"CompanyTypes"}

	cc.HttpHandler.
		Model(&model).
		DTO(&dto).
		Assoc(assocs).
		FiberCtx(c).
		UpdateOneBy(paramKey)

	return nil
}

func (cc *Company) deleteBy(paramKey string, c fiber.Ctx) error {
	var model models.Company
	var dto DTO.Company
	var assocs = []string{"CompanyTypes"}

	cc.HttpHandler.
		Model(&model).
		DTO(&dto).
		Assoc(assocs).
		FiberCtx(c).
		DeleteOneBy(paramKey)

	return nil
}

func (cc *Company) Create(c fiber.Ctx) error {
	var model models.Company
	var dto DTO.Company
	var assocs = []string{"CompanyTypes"}

	cc.HttpHandler.
		Model(&model).
		DTO(&dto).
		Assoc(assocs).
		FiberCtx(c).
		Create()

	return nil
}

func (cc *Company) GetAll(c fiber.Ctx) error {
	return cc.getBy("", c)
}


func (cc *Company) GetOneById(c fiber.Ctx) error {
	return cc.getBy("id", c)
}

func (cc *Company) GetOneByName(c fiber.Ctx) error {
	return cc.getBy("name", c)
}

func (cc *Company) GetOneByTaxId(c fiber.Ctx) error {
	return cc.getBy("tax_id", c)
}

func (cr *Company) UpdateById(c fiber.Ctx) error {
	return cr.updateBy("id", c)
}

func (cr *Company) UpdateByName(c fiber.Ctx) error {
	return cr.updateBy("name", c)
}

func (cr *Company) UpdateByTaxId(c fiber.Ctx) error {
	return cr.updateBy("tax_id", c)
}

func (cr *Company) DeleteById(c fiber.Ctx) error {
	return cr.deleteBy("id", c)
}

func (cr *Company) DeleteByName(c fiber.Ctx) error {
	return cr.deleteBy("name", c)
}

func (cr *Company) DeleteByTaxId(c fiber.Ctx) error {
	return cr.deleteBy("tax_id", c)
}
