package handlers

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	"log"

	"github.com/gofiber/fiber/v2"
)

func Request(Gorm *Gorm) *Req {
	return &Req{Gorm: Gorm}
}

type Req struct {
	Gorm *Gorm
}

// Chainable method to set the Fiber context
func (r *Req) SetAutomatedActions(c *fiber.Ctx) *AutoReqActions {
	return &AutoReqActions{
		req: r,
		res: Response(c),
		ctx: c,
	}
}

// AutoReqActions holds the intermediate data for method chaining
type AutoReqActions struct {
	req    *Req
	res    *Res
	ctx    *fiber.Ctx
	Error  error
	Status int
}

// Struct to hold fetched values from context
type ContextValues struct {
	ModelArr any
	Model    any
	Assocs   []string
	DtoArr   any
	Dto      any
}

// Centralized method to fetch values from context
func (ac *AutoReqActions) fetchContextValues() (*ContextValues, error) {
	keys := namespace.GeneralKey

	modelArr, err := lib.GetFromCtx[any](ac.ctx, keys.ModelArr)
	if err != nil {
		return nil, err
	}
	model, err := lib.GetFromCtx[any](ac.ctx, keys.Model)
	if err != nil {
		return nil, err
	}
	assocs, err := lib.GetFromCtx[[]string](ac.ctx, keys.Associations)
	if err != nil {
		return nil, err
	}
	dtoArr, err := lib.GetFromCtx[any](ac.ctx, keys.DtoArr)
	if err != nil {
		return nil, err
	}
	dto, err := lib.GetFromCtx[any](ac.ctx, keys.Dto)
	if err != nil {
		return nil, err
	}

	return &ContextValues{
		ModelArr: modelArr,
		Model:    model,
		Assocs:   assocs,
		DtoArr:   dtoArr,
		Dto:      dto,
	}, nil
}

// Standardized success response
func (ac *AutoReqActions) ActionSuccess(status int, data any, dto any) {
	ac.Error = nil
	ac.Status = status
	ac.res.DTO(status, data, dto)
}

// Standardized failure response
func (ac *AutoReqActions) ActionFailed(status int, err error) {
	ac.Error = err
	ac.Status = status
	ac.res.HttpError(status, err)
}

// Final method that executes a GET action
func (ac *AutoReqActions) GetBy(paramKey string) {
	if ac.Error != nil {
		return
	}

	ctxValues, err := ac.fetchContextValues()
	if err != nil {
		ac.ActionFailed(500, err)
		return
	}

	if paramKey == "" {
		if err := ac.req.Gorm.GetAll(ctxValues.ModelArr, ctxValues.Assocs); err != nil {
			ac.ActionFailed(500, err)
			return
		}
		ac.ActionSuccess(200, ctxValues.ModelArr, ctxValues.DtoArr)
	} else {
		paramVal := ac.ctx.Params(paramKey)
		if err := ac.req.Gorm.GetOneBy(paramKey, paramVal, ctxValues.Model, ctxValues.Assocs); err != nil {
			ac.ActionFailed(404, nil)
			return
		}
		ac.ActionSuccess(200, ctxValues.Model, ctxValues.Dto)
	}
}

// Final method that executes a FORCE GET action
func (ac *AutoReqActions) ForceGetBy(paramKey string) {
	if ac.Error != nil {
		return
	}

	ctxValues, err := ac.fetchContextValues()
	if err != nil {
		ac.ActionFailed(500, err)
		return
	}

	if paramKey == "" {
		if err := ac.req.Gorm.ForceGetAll(ctxValues.ModelArr, ctxValues.Assocs); err != nil {
			ac.ActionFailed(500, err)
			return
		}
		ac.ActionSuccess(200, ctxValues.ModelArr, ctxValues.DtoArr)
	} else {
		paramVal := ac.ctx.Params(paramKey)
		if err := ac.req.Gorm.ForceGetOneBy(paramKey, paramVal, ctxValues.Model, ctxValues.Assocs); err != nil {
			ac.ActionFailed(404, nil)
			return
		}
		ac.ActionSuccess(200, ctxValues.Model, ctxValues.Dto)
	}
}

// Final method that executes a CREATE action
func (ac *AutoReqActions) CreateOne() {
	if ac.Error != nil {
		return
	}

	ctxValues, err := ac.fetchContextValues()
	if err != nil {
		ac.ActionFailed(500, err)
		return
	}

	if err := ac.req.Gorm.Create(ctxValues.Model); err != nil {
		ac.ActionFailed(400, nil)
		return
	}

	ac.ActionSuccess(201, ctxValues.Model, ctxValues.Dto)
}

// Final method that executes a DELETE action
func (ac *AutoReqActions) DeleteOneById() {
	if ac.Error != nil {
		return
	}

	id := ac.ctx.Params(namespace.QueryKey.Id)
	ctxValues, err := ac.fetchContextValues()
	if err != nil {
		ac.ActionFailed(500, err)
		return
	}

	if err := ac.req.Gorm.GetOneBy("id", id, ctxValues.Model, nil); err != nil {
		ac.ActionFailed(404, nil)
		return
	}

	if err := ac.req.Gorm.DeleteOneById(id, ctxValues.Model); err != nil {
		ac.ActionFailed(500, err)
		return
	}

	log.Printf("Deleted record with ID: %s", id)
	ac.ActionSuccess(204, nil, nil)
}

// Final method that executes a FORCE DELETE action
func (ac *AutoReqActions) ForceDeleteOneById() {
	if ac.Error != nil {
		return
	}

	id := ac.ctx.Params(namespace.QueryKey.Id)
	ctxValues, err := ac.fetchContextValues()
	if err != nil {
		ac.ActionFailed(500, err)
		return
	}

	if err := ac.req.Gorm.ForceGetOneBy("id", id, ctxValues.Model, nil); err != nil {
		ac.ActionFailed(404, nil)
		return
	}

	if err := ac.req.Gorm.ForceDeleteOneById(id, ctxValues.Model); err != nil {
		ac.ActionFailed(500, err)
		return
	}

	ac.ActionSuccess(204, nil, nil)
}

// Final method that executes an UPDATE action
func (ac *AutoReqActions) UpdateOneById() {
	if ac.Error != nil {
		return
	}

	id := ac.ctx.Params(namespace.QueryKey.Id)
	ctxValues, err := ac.fetchContextValues()
	if err != nil {
		ac.ActionFailed(500, err)
		return
	}

	if err := ac.req.Gorm.UpdateOneById(id, ctxValues.Model, ctxValues.Model, ctxValues.Assocs); err != nil {
		ac.ActionFailed(400, nil)
		return
	}

	ac.ActionSuccess(200, ctxValues.Model, ctxValues.Dto)
}
