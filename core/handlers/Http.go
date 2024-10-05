package handlers

import (
	"agenda-kaki-go/core/lib"
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
	middlewares  []func(fiber.Ctx) error
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
func (ac *ActionChain) Middleware(mw func(fiber.Ctx) error) *ActionChain {
	ac.middlewares = append(ac.middlewares, mw)
	return ac
}

// Chainable method to set the Fiber context
func (ac *ActionChain) FiberCtx(c fiber.Ctx) *ActionChain {
	ac.ctx = c
	ac.sendResponse = &SendResponse{c}
	return ac
}

// Final method that executes a GetOneBy action
func (ac *ActionChain) GetOneBy(paramKey string) error {
	if err := ac.executeMiddlewares(); err != nil {
		return err
	}

	paramVal := ac.ctx.Params(paramKey)

	if err := ac.h.Gorm.GetOneBy(paramKey, paramVal, ac.model, ac.assoc); err != nil {
		return lib.Fiber404(ac.ctx)
	}

	if err := ac.sendResponse.DTO(ac.model, ac.dto); err != nil {
		return lib.Fiber500(ac.ctx, err)
	}

	return nil
}

// Final method that executes a DELETE action
func (ac *ActionChain) DeleteOneBy(paramKey string) error {
	if err := ac.executeMiddlewares(); err != nil {
		return err
	}

	paramVal := ac.ctx.Params(paramKey)

	if err := ac.h.Gorm.DeleteOneBy(paramKey, paramVal, ac.model); err != nil {
		return lib.Fiber500(ac.ctx, err)
	}

	return ac.ctx.SendStatus(fiber.StatusNoContent)
}

// Final method that executes an UPDATE action
func (ac *ActionChain) UpdateOneBy(paramKey string) error {
	// Parse the request body into the model
	if err := lib.BodyParser(ac.ctx.Body(), &ac.changes); err != nil {
		return lib.Fiber500(ac.ctx, err)
	}

	if err := ac.executeMiddlewares(); err != nil {
		return err
	}

	paramVal := ac.ctx.Params(paramKey)

	if err := ac.h.Gorm.UpdateOneBy(paramKey, paramVal, ac.model, ac.changes, ac.assoc); err != nil {
		return lib.Fiber400(ac.ctx, err)
	}

	if err := ac.sendResponse.DTO(ac.model, ac.dto); err != nil {
		return err
	}

	return nil
}

// Helper method to execute middlewares and authentication
func (ac *ActionChain) executeMiddlewares() error {

	for _, mw := range ac.middlewares {
		if err := mw(ac.ctx); err != nil {
			lib.Fiber400(ac.ctx, err)
			return err
		}
	}

	return nil
}
