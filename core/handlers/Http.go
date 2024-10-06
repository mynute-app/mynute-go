package handlers

import (
	"agenda-kaki-go/core/lib"
	"log"

	"github.com/gofiber/fiber/v3"
)

type HTTP struct {
	Gorm *Gorm
}

// ActionChain holds the intermediate data for method chaining
type ActionChain struct {
	h            *HTTP
	model        interface{}
	dto          interface{}
	assoc        []string
	middlewares  []func(fiber.Ctx) (int, error)
	ctx          fiber.Ctx
	sendResponse *lib.SendResponse
	changes      map[string]interface{}
	Error        error
}

// Chainable method to set the model
func (h *HTTP) Model(model interface{}) *ActionChain {
	return &ActionChain{
		h:     h,
		model: model,
	}
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

// Chainable method to set the Fiber context
func (ac *ActionChain) FiberCtx(c fiber.Ctx) *ActionChain {
	ac.ctx = c
	ac.sendResponse = &lib.SendResponse{Ctx: c}
	return ac
}

// Final method that executes a GetOneBy action
func (ac *ActionChain) GetOneBy(paramKey string) {
	if status, err := ac.executeMiddlewares(); err != nil {
		ac.sendResponse.HttpError(status, err)
		return
	}

	if paramKey == "" {
		if err := ac.h.Gorm.GetAll(ac.model, ac.assoc); err != nil {
			ac.sendResponse.Http400(err)
			return
		}
		ac.sendResponse.DTO(ac.model, ac.dto)
		return
	}

	paramVal := ac.ctx.Params(paramKey)

	if err := ac.h.Gorm.GetOneBy(paramKey, paramVal, ac.model, ac.assoc); err != nil {
		ac.sendResponse.Http404()
		return
	}

	ac.sendResponse.DTO(ac.model, ac.dto)
}

// Final method that executes a DELETE action
func (ac *ActionChain) DeleteOneBy(paramKey string) {
	ac.ctx.Locals("companyType", ac.model)

	if status, err := ac.executeMiddlewares(); err != nil {
		ac.sendResponse.HttpError(status, err)
		return
	}

	paramVal := ac.ctx.Params(paramKey)

	if err := ac.h.Gorm.DeleteOneBy(paramKey, paramVal, ac.model); err != nil {
		ac.sendResponse.Http500(err)
		return
	}

	ac.sendResponse.Http204()
}

// Final method that executes an UPDATE action
func (ac *ActionChain) UpdateOneBy(paramKey string) {
	// Parse the request body into the model
	if err := lib.BodyParser(ac.ctx.Body(), &ac.changes); err != nil {
		ac.sendResponse.Http500(err)
		return
	}

	ac.ctx.Locals("changes", ac.changes)

	if status, err := ac.executeMiddlewares(); err != nil {
		ac.sendResponse.HttpError(status, err)
		return
	}

	paramVal := ac.ctx.Params(paramKey)

	if err := ac.h.Gorm.UpdateOneBy(paramKey, paramVal, ac.model, ac.changes, ac.assoc); err != nil {
		ac.sendResponse.Http400(err)
		return
	}

	ac.sendResponse.DTO(ac.model, ac.dto)
}

// Final method that executes a CREATE action
func (ac *ActionChain) Create() {
	// Parse the request body into the model
	log.Printf("Parsing body")
	if err := lib.BodyParser(ac.ctx.Body(), &ac.model); err != nil {
		ac.sendResponse.Http500(err)
		return
	}

	log.Printf("Executing middlewares")
	if status, err := ac.executeMiddlewares(); err != nil {
		ac.sendResponse.HttpError(status, err)
		return
	}

	log.Printf("Creating model")
	log.Printf("model: %+v", ac.model)

	if err := ac.h.Gorm.Create(ac.model); err != nil {
		ac.sendResponse.Http400(err)
		return
	}

	ac.sendResponse.Http201()
}

// Helper method to execute middlewares and authentication
func (ac *ActionChain) executeMiddlewares() (int, error) {

	for _, mw := range ac.middlewares {
		if status, err := mw(ac.ctx); err != nil {
			return status, err
		}
	}

	return 0, nil
}
