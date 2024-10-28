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

// HttpActionChain holds the intermediate data for method chaining
type HttpActionChain struct {
	h            *HTTP
	ctx          fiber.Ctx
	sendResponse *lib.SendResponse
	Error        error
	Status       int
}

// Chainable method to set the Fiber context
func (h *HTTP) FiberCtx(c fiber.Ctx) *HttpActionChain {
	return &HttpActionChain{
		h:            h,
		sendResponse: &lib.SendResponse{Ctx: c},
		ctx:          c,
	}
}

// Centralized error handling
func (ac *HttpActionChain) SendError(status int, err error) {
	ac.Error = err
	ac.Status = status
	ac.sendResponse.HttpError(status, err)
}

// Final Method that executes a GET action
func (ac *HttpActionChain) GetBy(paramKey string) {
	if ac.Error != nil {
		return
	}

	keys := namespace.GeneralKey

	modelArr, err := lib.GetFromCtx[interface{}](ac.ctx, keys.ModelArr)
	if err != nil {
		ac.sendResponse.Http500(err)
		return
	}
	assocs, err := lib.GetFromCtx[[]string](ac.ctx, keys.Associations)
	if err != nil {
		ac.sendResponse.Http500(err)
		return
	}
	dtoArr, err := lib.GetFromCtx[interface{}](ac.ctx, keys.DtoArr)
	if err != nil {
		ac.sendResponse.Http500(err)
		return
	}

	if paramKey == "" {
		if err := ac.h.Gorm.GetAll(modelArr, assocs); err != nil {
			ac.sendResponse.Http500(err)
			return
		}
	} else {
		paramVal := ac.ctx.Params(paramKey)
		if err := ac.h.Gorm.GetOneBy(paramKey, paramVal, modelArr, assocs); err != nil {
			ac.sendResponse.Http404()
			return
		}
	}

	ac.sendResponse.DTO(200, modelArr, dtoArr)
}

// Final Method that executes a FORCE GET action
func (ac *HttpActionChain) ForceGetBy(paramKey string) {
	if ac.Error != nil {
		return
	}

	keys := namespace.GeneralKey

	modelArr, err := lib.GetFromCtx[interface{}](ac.ctx, keys.ModelArr)
	if err != nil {
		ac.sendResponse.Http500(err)
		return
	}
	assocs, err := lib.GetFromCtx[[]string](ac.ctx, keys.Associations)
	if err != nil {
		ac.sendResponse.Http500(err)
		return
	}
	dtoArr, err := lib.GetFromCtx[interface{}](ac.ctx, keys.DtoArr)
	if err != nil {
		ac.sendResponse.Http500(err)
		return
	}

	if paramKey == "" {
		if err := ac.h.Gorm.ForceGetAll(modelArr, assocs); err != nil {
			ac.sendResponse.Http500(err)
			return
		}
	} else {
		paramVal := ac.ctx.Params(paramKey)
		if err := ac.h.Gorm.ForceGetOneBy(paramKey, paramVal, modelArr, assocs); err != nil {
			ac.sendResponse.Http404()
			return
		}
	}

	ac.sendResponse.DTO(200, modelArr, dtoArr)
}

// Final Method that executes a CREATE action
func (ac *HttpActionChain) CreateOne() {
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
func (ac *HttpActionChain) DeleteOneById() {
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
func (ac *HttpActionChain) ForceDeleteOneById() {
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
func (ac *HttpActionChain) UpdateOneById() {
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

	dto, err := lib.GetFromCtx[interface{}](ac.ctx, keys.Dto)
	if err != nil {
		ac.sendResponse.Http500(err)
		return
	}

	ac.sendResponse.DTO(200, model, dto)
}
