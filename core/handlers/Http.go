package handlers

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"

	"github.com/gofiber/fiber/v3"
)

type HTTP struct {
	Gorm *Gorm
}

// ActionChain holds the intermediate data for method chaining
type ActionChain struct {
	h                *HTTP
	model            interface{}
	dto              interface{}
	assoc            []string
	middlewares      []func(fiber.Ctx) (int, error)
	ctx              fiber.Ctx
	sendResponse     *lib.SendResponse
	interfaceNameKey namespace.ContextKey
	Error            error
}

// Chainable method to set the Fiber context
func (h *HTTP) FiberCtx(c fiber.Ctx) *ActionChain {
	return &ActionChain{
		h: h,
		sendResponse: &lib.SendResponse{Ctx: c},
		ctx: c,
	}
}

// Chainable method to set the model
func (ac *ActionChain) Model(model interface{}) *ActionChain {
	ac.model = model
	return ac
}

// Chainable method to set the DTO
func (ac *ActionChain) DTO(dto interface{}) *ActionChain {
	ac.dto = dto
	return ac
}

// Chainable method to set the associations
func (ac *ActionChain) Assoc(assoc []string) *ActionChain {
	ac.assoc = assoc
	return ac
}

// Chainable method to add a middleware
func (ac *ActionChain) Middleware(mw func(fiber.Ctx) (int, error)) *ActionChain {
	ac.middlewares = append(ac.middlewares, mw)
	return ac
}

func (ac *ActionChain) InterfaceKey(key namespace.ContextKey) *ActionChain {
	ac.interfaceNameKey = key
	return ac
}

// // Final method that executes a GetOneBy action
// func (ac *ActionChain) GetOneBy(paramKey string) {
// 	if status, err := ac.executeMiddlewares(); err != nil {
// 		ac.sendResponse.HttpError(status, err)
// 		return
// 	}

// 	if paramKey == "" {
// 		if err := ac.h.Gorm.GetAll(ac.model, ac.assoc); err != nil {
// 			ac.sendResponse.Http400(err)
// 			return
// 		}
// 		ac.sendResponse.DTO(200, ac.model, ac.dto)
// 		return
// 	}

// 	paramVal := ac.ctx.Params(paramKey)

// 	if err := ac.h.Gorm.GetOneBy(paramKey, paramVal, ac.model, ac.assoc); err != nil {
// 		ac.sendResponse.Http404()
// 		return
// 	}

// 	ac.sendResponse.DTO(200, ac.model, ac.dto)
// }

// Final Method that executes a GET action
func (ac *ActionChain) GetBy(paramKey string) {
	if status, err := ac.executeMiddlewares(); err != nil {
		ac.sendResponse.HttpError(status, err)
		return
	}

	keys := namespace.GeneralKey

	model, err := lib.GetFromCtx[*interface{}](ac.ctx, keys.Model)
	if err != nil {
		ac.sendResponse.Http500(err)
		return
	}
	assocs, err := lib.GetFromCtx[[]string](ac.ctx, keys.Associations)
	if err != nil {
		ac.sendResponse.Http500(err)
		return
	}

	if paramKey == "" {
		if err := ac.h.Gorm.GetAll(model, assocs); err != nil {
			ac.sendResponse.Http400(err)
			return
		}
	} else {
		paramVal := ac.ctx.Params(paramKey)
		if err := ac.h.Gorm.GetOneBy(paramKey, paramVal, model, assocs); err != nil {
			ac.sendResponse.Http404()
			return
		}
	}

	ac.sendResponse.DTO(200, ac.model, ac.dto)
}

// Final Method that executes a CREATE action
func (ac *ActionChain) Create() {
	if status, err := ac.executeMiddlewares(); err != nil {
		ac.sendResponse.HttpError(status, err)
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

	if status, err := ac.executeMiddlewares(); err != nil {
		ac.sendResponse.HttpError(status, err)
		return
	}

	keys := namespace.GeneralKey

	id := ac.ctx.Params(string(keys.QueryId))
	model, err := lib.GetFromCtx[interface{}](ac.ctx, keys.Model)
	if err != nil {
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

	ac.sendResponse.Http204()
}

// Final method that executes an UPDATE action
func (ac *ActionChain) UpdateOneById() {
	if status, err := ac.executeMiddlewares(); err != nil {
		ac.sendResponse.HttpError(status, err)
		return
	}
	keys := namespace.GeneralKey
	associations, err := lib.GetFromCtx[[]string](ac.ctx, keys.Associations); if err != nil {
		ac.sendResponse.Http500(err)
		return
	}
	changes, err := lib.GetFromCtx[map[string]interface{}](ac.ctx, keys.Changes); if err != nil {
		ac.sendResponse.Http500(err)
		return
	}
	model, err := lib.GetFromCtx[interface{}](ac.ctx, keys.Model); if err != nil {
		ac.sendResponse.Http500(err)
		return
	}

	id := ac.ctx.Params(string(keys.QueryId))

	if err := ac.h.Gorm.UpdateOneById(id, model, changes, associations); err != nil {
		ac.sendResponse.Http400(err)
		return
	}

	ac.sendResponse.DTO(200, ac.model, ac.dto)
}

// // Final method that executes a CREATE action
// func (ac *ActionChain) Create() {
// 	// Parse the request body into the model
// 	log.Printf("Parsing body")
// 	if err := lib.BodyParser(ac.ctx.Body(), &ac.model); err != nil {
// 		ac.sendResponse.Http500(err)
// 		return
// 	}

// 	ac.ctx.Locals(ac.interfaceNameKey, &ac.model)

// 	log.Printf("Executing middlewares")
// 	if status, err := ac.executeMiddlewares(); err != nil {
// 		ac.sendResponse.HttpError(status, err)
// 		return
// 	}

// 	log.Printf("Creating model")
// 	log.Printf("model: %+v", ac.model)

// 	if err := ac.h.Gorm.Create(ac.model); err != nil {
// 		ac.sendResponse.Http400(err)
// 		return
// 	}

// 	ac.sendResponse.Http201(ac.model)
// }

// Helper method to execute middlewares and authentication
func (ac *ActionChain) executeMiddlewares() (int, error) {

	for _, mw := range ac.middlewares {
		if status, err := mw(ac.ctx); err != nil {
			return status, err
		}
	}

	return 0, nil
}
