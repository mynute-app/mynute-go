package handlers

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	"log"

	"github.com/gofiber/fiber/v3"
)

type HTTP struct {
	Gorm *Gorm
}

// ActionChain holds the intermediate data for method chaining
type ActionChain struct {
	h                *HTTP
	ctx              fiber.Ctx
	sendResponse     *lib.SendResponse
	Error            error
}

// Chainable method to set the Fiber context
func (h *HTTP) FiberCtx(c fiber.Ctx) *ActionChain {
	return &ActionChain{
		h:            h,
		sendResponse: &lib.SendResponse{Ctx: c},
		ctx:          c,
	}
}

func (ac *ActionChain) RunMiddlewares(mdws []func(fiber.Ctx) (int, error)) *ActionChain {
	for _, mdw := range mdws {
		if s, err := mdw(ac.ctx); err != nil {
			ac.sendResponse.HttpError(s, err)
			ac.Error = err
			return ac
		}
	}
	return ac
}

// Final Method that executes a GET action
func (ac *ActionChain) GetBy(paramKey string) {
	if ac.Error != nil {
		return
	}

	keys := namespace.GeneralKey

	model, err := lib.GetFromCtx[interface{}](ac.ctx, keys.Model)
	if err != nil {
		ac.sendResponse.Http500(err)
		return
	}
	assocs, err := lib.GetFromCtx[[]string](ac.ctx, keys.Associations)
	if err != nil {
		ac.sendResponse.Http500(err)
		return
	}
	dto, err := lib.GetFromCtx[interface{}](ac.ctx, keys.Dto)
	if err != nil {
		ac.sendResponse.Http500(err)
		return
	}

	if paramKey == "" {
		if err := ac.h.Gorm.GetAll(model, assocs); err != nil {
			ac.sendResponse.Http500(err)
			return
		}
	} else {
		paramVal := ac.ctx.Params(paramKey)
		if err := ac.h.Gorm.GetOneBy(paramKey, paramVal, model, assocs); err != nil {
			ac.sendResponse.Http404()
			return
		}
	}

	ac.sendResponse.DTO(200, model, dto)
}

// Final Method that executes a FORCE GET action
func (ac *ActionChain) ForceGetBy(paramKey string) {
	if ac.Error != nil {
		return
	}

	keys := namespace.GeneralKey

	model, err := lib.GetFromCtx[interface{}](ac.ctx, keys.Model)
	if err != nil {
		ac.sendResponse.Http500(err)
		return
	}
	assocs, err := lib.GetFromCtx[[]string](ac.ctx, keys.Associations)
	if err != nil {
		ac.sendResponse.Http500(err)
		return
	}
	dto, err := lib.GetFromCtx[interface{}](ac.ctx, keys.Dto)
	if err != nil {
		ac.sendResponse.Http500(err)
		return
	}

	if paramKey == "" {
		if err := ac.h.Gorm.ForceGetAll(model, assocs); err != nil {
			ac.sendResponse.Http500(err)
			return
		}
	} else {
		paramVal := ac.ctx.Params(paramKey)
		if err := ac.h.Gorm.ForceGetOneBy(paramKey, paramVal, model, assocs); err != nil {
			ac.sendResponse.Http404()
			return
		}
	}

	ac.sendResponse.DTO(200, model, dto)
}

// Final Method that executes a CREATE action
func (ac *ActionChain) CreateOne() {
	if ac.Error != nil {
		return
	}
	keys := namespace.GeneralKey

	model, err := lib.GetFromCtx[interface{}](ac.ctx, keys.Model)
	if err != nil {
		ac.sendResponse.Http500(err)
		return
	}

	if err := ac.h.Gorm.Create(model); err != nil {
		ac.sendResponse.Http400(err)
		return
	}

	dto, err := lib.GetFromCtx[interface{}](ac.ctx, keys.Dto)
	if err != nil {
		ac.sendResponse.Http500(err)
		return
	}

	ac.sendResponse.DTO(201, model, dto)
}

// Final method that executes a DELETE action
func (ac *ActionChain) DeleteOneById() {
	if ac.Error != nil {
		return
	}
	keys := namespace.GeneralKey

	id := ac.ctx.Params(string(keys.QueryId))
	model, err := lib.GetFromCtx[interface{}](ac.ctx, keys.Model)
	if err != nil {
		ac.sendResponse.Http500(err)
		return
	}

	if err := ac.h.Gorm.GetOneBy("id", id, model, nil); err != nil {
		if err.Error() == "record not found" {
			ac.sendResponse.Http404()
			return
		}
		ac.sendResponse.Http500(err)
		return
	}

	if err := ac.h.Gorm.DeleteOneById(id, model); err != nil {
		if err.Error() == "record not found" {
			ac.sendResponse.Http404()
			return
		}
		ac.sendResponse.Http500(err)
		return
	}

	log.Printf("Deleted record with ID: %s", id)

	ac.sendResponse.Http204()
}

// Final method that executes a FORCE DELETE action
func (ac *ActionChain) ForceDeleteOneById() {
	if ac.Error != nil {
		return
	}
	keys := namespace.GeneralKey

	id := ac.ctx.Params(string(keys.QueryId))
	model, err := lib.GetFromCtx[interface{}](ac.ctx, keys.Model)
	if err != nil {
		ac.sendResponse.Http500(err)
		return
	}

	if err := ac.h.Gorm.ForceGetOneBy("id", id, model, nil); err != nil {
		if err.Error() == "record not found" {
			ac.sendResponse.Http404()
			return
		}
		ac.sendResponse.Http500(err)
		return
	}

	if err := ac.h.Gorm.ForceDeleteOneById(id, model); err != nil {
		if err.Error() == "record not found" {
			ac.sendResponse.Http404()
			return
		}
		ac.sendResponse.Http500(err)
		return
	}

	ac.sendResponse.Http204()
}

// Final method that executes an UPDATE action
func (ac *ActionChain) UpdateOneById() {
	if ac.Error != nil {
		return
	}
	keys := namespace.GeneralKey
	associations, err := lib.GetFromCtx[[]string](ac.ctx, keys.Associations)
	if err != nil {
		ac.sendResponse.Http500(err)
		return
	}
	changes, err := lib.GetFromCtx[map[string]interface{}](ac.ctx, keys.Changes)
	if err != nil {
		ac.sendResponse.Http500(err)
		return
	}
	model, err := lib.GetFromCtx[interface{}](ac.ctx, keys.Model)
	if err != nil {
		ac.sendResponse.Http500(err)
		return
	}

	id := ac.ctx.Params(string(keys.QueryId))

	if err := ac.h.Gorm.UpdateOneById(id, model, changes, associations); err != nil {
		ac.sendResponse.Http400(err)
		return
	}

	dto, err := lib.GetFromCtx[interface{}](ac.ctx, keys.Dto); if err != nil {
		ac.sendResponse.Http500(err)
		return
	}

	ac.sendResponse.DTO(200, model, dto)
}