package handler

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	"net/url"

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
		res: &lib.SendResponse{Ctx: c},
		ctx: c,
	}
}

// AutoReqActions holds the intermediate data for method chaining
type AutoReqActions struct {
	req      *Req
	res      *lib.SendResponse
	ctx      *fiber.Ctx
	ctxVal   *ContextValues
	mute_res bool
	Error    error
	Status   int
}

// Struct to hold fetched values from context
type ContextValues struct {
	ModelArr any
	Model    any
	Assocs   []string
	DtoArr   any
	Dto      any
	Changes  map[string]any
}

// Centralized method to fetch values from context
func (ac *AutoReqActions) fetchContextValues() error {
	keys := namespace.GeneralKey

	modelArr, err := lib.GetFromCtx[any](ac.ctx, keys.ModelArr)
	if err != nil {
		return err
	}
	model, err := lib.GetFromCtx[any](ac.ctx, keys.Model)
	if err != nil {
		return err
	}
	assocs, err := lib.GetFromCtx[[]string](ac.ctx, keys.Associations)
	if err != nil {
		return err
	}
	dtoArr, err := lib.GetFromCtx[any](ac.ctx, keys.DtoArr)
	if err != nil {
		return err
	}
	dto, err := lib.GetFromCtx[any](ac.ctx, keys.Dto)
	if err != nil {
		return err
	}
	changes, err := lib.GetFromCtx[map[string]any](ac.ctx, keys.Changes)
	if err != nil {
		return err
	}

	ac.ctxVal = &ContextValues{
		ModelArr: modelArr,
		Model:    model,
		Assocs:   assocs,
		DtoArr:   dtoArr,
		Dto:      dto,
		Changes:  changes,
	}

	return nil
}

// Standardized success response
func (ac *AutoReqActions) ActionSuccess(status int, data any, dto any) error {
	ac.Error = nil
	ac.Status = status
	if !ac.mute_res {
		if err := ac.res.SendDTO(status, data, dto); err != nil {
			return err
		}
	}
	return nil
}

// Standardized failure response
func (ac *AutoReqActions) ActionFailed(status int, err error) error {
	ac.Error = err
	ac.Status = status
	if !ac.mute_res {
		if err := ac.res.HttpError(status, err); err != nil {
			return err
		}
	}
	return nil
}

func (ac *AutoReqActions) MuteResponse(mute bool) {
	ac.mute_res = mute
}

// Final method that executes a GET action
func (ac *AutoReqActions) GetBy(paramKey string) error {
	if ac.Error != nil {
		return ac.Error
	}

	err := ac.fetchContextValues()
	if err != nil {
		return err
	}

	if paramKey == "" {
		if err := ac.req.Gorm.GetAll(ac.ctxVal.ModelArr, ac.ctxVal.Assocs); err != nil {
			return err
		}
		return ac.ActionSuccess(200, ac.ctxVal.ModelArr, ac.ctxVal.DtoArr)
	} // Get the parameter value from the context
	paramVal := ac.ctx.Params(paramKey)
	// Decode URL-encoded characters
	cleanedParamVal, err := url.QueryUnescape(paramVal)
	if err != nil {
		return err
	}
	if err := ac.req.Gorm.GetOneBy(paramKey, cleanedParamVal, ac.ctxVal.Model, ac.ctxVal.Assocs); err != nil {
		return err
	}
	return ac.ActionSuccess(200, ac.ctxVal.Model, ac.ctxVal.Dto)
}

// Final method that executes a FORCE GET action
func (ac *AutoReqActions) ForceGetBy(paramKey string) error {
	if ac.Error != nil {
		return ac.Error
	}

	err := ac.fetchContextValues()
	if err != nil {
		return err
	}

	if paramKey == "" {
		if err := ac.req.Gorm.ForceGetAll(ac.ctxVal.ModelArr, ac.ctxVal.Assocs); err != nil {
			return err
		}
		return ac.ActionSuccess(200, ac.ctxVal.ModelArr, ac.ctxVal.DtoArr)
	}
	paramVal := ac.ctx.Params(paramKey)
	if err := ac.req.Gorm.ForceGetOneBy(paramKey, paramVal, ac.ctxVal.Model, ac.ctxVal.Assocs); err != nil {
		if err := ac.ActionFailed(404, err); err != nil {
			return err
		}
	}
	return ac.ActionSuccess(200, ac.ctxVal.Model, ac.ctxVal.Dto)
}

// Final method that executes a CREATE action
func (ac *AutoReqActions) CreateOne() error {
	if ac.Error != nil {
		return ac.Error
	}

	err := ac.fetchContextValues()
	if err != nil {
		return err
	}

	if err := ac.req.Gorm.Create(ac.ctxVal.Model, ac.ctxVal.Assocs); err != nil {
		return ac.ActionFailed(400, err)
	}

	return ac.ActionSuccess(200, ac.ctxVal.Model, ac.ctxVal.Dto)
}

// Final method that executes a DELETE action
func (ac *AutoReqActions) DeleteOneById() error {
	if ac.Error != nil {
		return ac.Error
	}

	id := ac.ctx.Params(namespace.QueryKey.Id)
	err := ac.fetchContextValues()
	if err != nil {
		return err
	}

	if err := ac.req.Gorm.GetOneBy("id", id, ac.ctxVal.Model, nil); err != nil {
		return ac.ActionFailed(404, err)
	}

	if err := ac.req.Gorm.DeleteOneById(id, ac.ctxVal.Model); err != nil {
		return err
	}

	return ac.ActionSuccess(200, nil, nil)
}

// Final method that executes a FORCE DELETE action
func (ac *AutoReqActions) ForceDeleteOneById() error {
	if ac.Error != nil {
		return ac.Error
	}

	id := ac.ctx.Params(namespace.QueryKey.Id)
	err := ac.fetchContextValues()
	if err != nil {
		return err
	}

	if err := ac.req.Gorm.ForceGetOneBy("id", id, ac.ctxVal.Model, nil); err != nil {
		return ac.ActionFailed(404, err)
	}

	if err := ac.req.Gorm.ForceDeleteOneById(id, ac.ctxVal.Model); err != nil {
		return err
	}

	return ac.ActionSuccess(204, nil, nil)
}

// Final method that executes an UPDATE action
func (ac *AutoReqActions) UpdateOneById() error {
	if ac.Error != nil {
		return ac.Error
	}

	id := ac.ctx.Params(namespace.QueryKey.Id)
	err := ac.fetchContextValues()
	if err != nil {
		return err
	}

	if err := ac.req.Gorm.UpdateOneById(id, ac.ctxVal.Model, ac.ctxVal.Changes, ac.ctxVal.Assocs); err != nil {
		return ac.ActionFailed(400, err)
	}

	return ac.ActionSuccess(200, ac.ctxVal.Model, ac.ctxVal.Dto)
}
