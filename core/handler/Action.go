package handler

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

type actions struct {
	res           *lib.SendResponseStruct
	ctx           *fiber.Ctx
	myGormWrapper *Gorm
}

func Actions(c *fiber.Ctx) *actions {
	tx, err := lib.Session(c)
	if err != nil {
		return nil
	}
	return &actions{
		res:           &lib.SendResponseStruct{Ctx: c},
		ctx:           c,
		myGormWrapper: &Gorm{DB: tx},
	}
}

func (a *actions) GetAll() error {
	model_array, err := lib.GetFromCtx[any](a.ctx, namespace.GeneralKey.ModelArr)
	if err != nil {
		return err
	}
	if err := a.myGormWrapper.GetAll(model_array); err != nil {
		return lib.Error.General.RecordNotFound.WithError(err)
	}
	dto_array, err := lib.GetFromCtx[any](a.ctx, namespace.GeneralKey.DtoArr)
	if err != nil {
		return err
	}
	return a.res.SendDTO(200, model_array, dto_array)
}

func (a *actions) GetBy(paramKey string) error {
	if paramKey == "" {
		return a.GetAll()
	} // Get the parameter value from the context
	paramVal := a.ctx.Params(paramKey)
	// Decode URL-encoded characters
	cleanedParamVal, err := url.QueryUnescape(paramVal)
	if err != nil {
		return err
	}
	// Get the model from the context
	model, err := lib.GetFromCtx[any](a.ctx, namespace.GeneralKey.Model)
	if err != nil {
		return err
	}
	if err := a.myGormWrapper.GetOneBy(paramKey, cleanedParamVal, model); err != nil {
		return lib.Error.General.RecordNotFound.WithError(err)
	}
	dto, err := lib.GetFromCtx[any](a.ctx, namespace.GeneralKey.Dto)
	if err != nil {
		return err
	}
	return a.res.SendDTO(200, model, dto)
}

func (a *actions) ForceGetBy(paramKey string) error {
	if paramKey == "" {
		return a.GetAll()
	} // Get the parameter value from the context
	paramVal := a.ctx.Params(paramKey)
	// Decode URL-encoded characters
	cleanedParamVal, err := url.QueryUnescape(paramVal)
	if err != nil {
		return err
	}
	// Get the model from the context
	model, err := lib.GetFromCtx[any](a.ctx, namespace.GeneralKey.Model)
	if err != nil {
		return err
	}
	if err := a.myGormWrapper.ForceGetOneBy(paramKey, cleanedParamVal, model); err != nil {
		return lib.Error.General.RecordNotFound.WithError(err)
	}
	dto, err := lib.GetFromCtx[any](a.ctx, namespace.GeneralKey.Dto)
	if err != nil {
		return err
	}
	return a.res.SendDTO(200, model, dto)
}
