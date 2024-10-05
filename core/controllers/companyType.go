package controllers

import (
	"agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v3"
)

type CompanyType struct {
	Gorm        *handlers.Gorm
	Middleware  *middleware.CompanyType
	HttpHandler *handlers.HTTP
}

func (ctc *CompanyType) getBy(paramKey string, c fiber.Ctx) error {

	var model models.CompanyType
	var dto DTO.CompanyType
	assocs := []string{}

	err := ctc.HttpHandler.
		Model(&model).
		DTO(&dto).
		Assoc(assocs).
		FiberCtx(c).
		GetOneBy(paramKey)
	
	if err != nil {
		panic(err)
	}

	return nil
}

// func (ctc *CompanyType) getBy(paramKey string, c fiber.Ctx) error {

// 	var model models.CompanyType
// 	var dto DTO.CompanyType
// 	assocs := []string{}

// 	if err := ctc.HttpHandler.GetOneBy(paramKey, c, &model, &dto, assocs); err != nil {
// 		return err
// 	}

// 	ctc.HttpHandler.GetOneBy(paramKey).Model(&model).DTO(&dto).Assoc(assocs).FiberContext(c)

// 	return nil
// }

// var model models.CompanyType

// paramVal := c.Params(paramKey)

// if err := ctc.Gorm.GetOneBy(paramKey, paramVal, &model, nil); err != nil {
// 	return lib.Fiber404(c)
// }

// var dto DTO.CompanyType

// if err := handlers.ParseAndSendResponse(model, &dto, c); err != nil {
// 	return err
// }

func (ctc *CompanyType) updateBy(paramKey string, c fiber.Ctx) error {
	var model models.CompanyType
	var dto DTO.CompanyType
	assocs := []string{}

	updateMiddleware := func (c fiber.Ctx) error {
		ctc.Middleware.Update

	if err := ctc.HttpHandler.
	Model(&model).
	DTO(&dto).
	Assoc(assocs).
	FiberCtx(c).
	Middleware(ctc.Middleware.Update).
	UpdateOneBy(paramKey); err != nil {
		return err
	}

	return nil
}

// func (ctc *CompanyType) updateBy(paramKey string, c fiber.Ctx) error {
// 	var changes map[string]interface{}

// 	if err := lib.BodyParser(c.Body(), &changes); err != nil {
// 		return lib.FiberError(400, c, err)
// 	}

// 	if err := ctc.Middleware.Update(changes); err != nil {
// 		return lib.FiberError(400, c, err)
// 	}

// 	var model models.CompanyType

// 	var paramVal = c.Params(paramKey)

// 	if err := ctc.Gorm.UpdateOneBy(paramKey, paramVal, &model, changes, nil); err != nil {
// 		return lib.FiberError(400, c, err)
// 	}

// 	var dto DTO.CompanyType

// 	if err := lib.ParseToDTO(model, &dto); err != nil {
// 		return err
// 	}

// 	if err := c.JSON(dto); err != nil {
// 		return err
// 	}

// 	return nil
// }



// func (ctc *CompanyType) deleteBy(paramKey string, c fiber.Ctx) error {
// 	var model models.CompanyType

// 	paramVal := c.Params(paramKey)

// 	if err := ctc.Gorm.GetOneBy(paramKey, paramVal, &model, nil); err != nil {
// 		return lib.FiberError(404, c, err)
// 	}

// 	if err := ctc.Middleware.Delete(model); err != nil {
// 		return lib.FiberError(400, c, err)
// 	}

// 	if err := ctc.Gorm.DeleteOneBy(paramKey, paramVal, &model); err != nil {
// 		return lib.FiberError(400, c, err)
// 	}

// 	return nil
// }

func (ctc *CompanyType) Create(c fiber.Ctx) error {
	var model models.CompanyType

	if err := lib.BodyParser(c.Body(), &model); err != nil {
		return lib.FiberError(400, c, err)
	}

	if err := ctc.Middleware.Create(model); err != nil {
		return lib.FiberError(400, c, err)
	}

	if err := ctc.Gorm.Create(&model, nil); err != nil {
		return lib.FiberError(400, c, err)
	}

	var dto DTO.CompanyType

	if err := lib.ParseToDTO(model, &dto); err != nil {
		return lib.FiberError(500, c, err)
	}

	if err := c.JSON(dto); err != nil {
		return err
	}

	return nil
}

func (ctc *CompanyType) GetAll(c fiber.Ctx) error {
	var model []models.CompanyType

	if err := ctc.Gorm.GetAll(&model, nil); err != nil {
		return lib.FiberError(400, c, err)
	}

	var dto []DTO.CompanyType

	if err := lib.ParseToDTO(model, &dto); err != nil {
		return lib.FiberError(500, c, err)
	}

	if err := c.JSON(dto); err != nil {
		return err
	}

	return nil
}

func (ctc *CompanyType) GetOneById(c fiber.Ctx) error {
	return ctc.getBy("id", c)
}

func (ctc *CompanyType) GetOneByName(c fiber.Ctx) error {
	return ctc.getBy("name", c)
}

func (ctc *CompanyType) UpdateById(c fiber.Ctx) error {
	return ctc.updateBy("id", c)
}

func (ctc *CompanyType) UpdateByName(c fiber.Ctx) error {
	return ctc.updateBy("name", c)
}

func (ctc *CompanyType) DeleteById(c fiber.Ctx) error {
	return ctc.deleteBy("id", c)
}

func (ctc *CompanyType) DeleteByName(c fiber.Ctx) error {
	return ctc.deleteBy("name", c)
}

// func (ctc *CompanyType) Create(c fiber.Ctx) error {
// 	var model models.CompanyType
// 	var dto DTO.CompanyType
// 	assocs := []string{}

// 	CtrlService := services.Controller{Ctx: c, DB: ctc.Gorm}
// 	if err := CtrlService.Create(&model, &dto, assocs); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (ctc *CompanyType) GetAll(c fiber.Ctx) error {
// 	var model []models.CompanyType
// 	var dto []DTO.CompanyType
// 	assocs := []string{}

// 	CtrlService := services.Controller{Ctx: c, DB: ctc.DB}
// 	if err := CtrlService.GetAll(&model, &dto, assocs); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (ctc *CompanyType) getBy(param string, c fiber.Ctx) error {
// 	var model models.CompanyType
// 	var dto DTO.CompanyType
// 	assocs := []string{}

// 	CtrlService := services.Controller{Ctx: c, DB: ctc.DB}
// 	if err := CtrlService.GetOneBy(param, &model, &dto, assocs); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (ctc *CompanyType) updateBy(param string, c fiber.Ctx) error {
// 	var model models.CompanyType
// 	var dto DTO.CompanyType
// 	assocs := []string{}

// 	CtrlService := services.Controller{Ctx: c, DB: ctc.DB}
// 	if err := CtrlService.UpdateOneBy(param, &model, &dto, assocs); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (ctc *CompanyType) deleteBy(param string, c fiber.Ctx) error {
// 	var model models.CompanyType

// 	CtrlService := services.Controller{Ctx: c, DB: ctc.DB}
// 	if err := CtrlService.DeleteOneBy(param, &model); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (ctc *CompanyType) GetOneById(c fiber.Ctx) error {
// 	return ctc.getBy("id", c)
// }

// func (ctc *CompanyType) GetOneByName(c fiber.Ctx) error {
// 	return ctc.getBy("name", c)
// }

// func (ctc *CompanyType) UpdateById(c fiber.Ctx) error {
// 	return ctc.updateBy("id", c)
// }

// func (ctc *CompanyType) DeleteById(c fiber.Ctx) error {
// 	return ctc.deleteBy("id", c)
// }
