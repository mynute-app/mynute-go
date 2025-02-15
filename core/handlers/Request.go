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
func (r *Req) FiberCtx(c *fiber.Ctx) *ReqActions {
	return &ReqActions{
		req: r,
		res: Response(c),
		ctx: c,
	}
}

// ReqActions holds the intermediate data for method chaining
type ReqActions struct {
	req    *Req
	res    *Res
	ctx    *fiber.Ctx
	Error  error
	Status int
}

// Centralized error handling
func (ac *ReqActions) SendError(status int, err error) {
	ac.Error = err
	ac.Status = status
	ac.res.HttpError(status, err)
}

// Final Method that executes a GET action
func (ac *ReqActions) GetBy(paramKey string) {
	if ac.Error != nil {
		return
	}

	keys := namespace.GeneralKey

	modelArr, err := lib.GetFromCtx[interface{}](ac.ctx, keys.ModelArr)
	if err != nil {
		ac.res.Http500(err)
		return
	}
	model, err := lib.GetFromCtx[interface{}](ac.ctx, keys.Model)
	if err != nil {
		ac.res.Http500(err)
		return
	}
	assocs, err := lib.GetFromCtx[[]string](ac.ctx, keys.Associations)
	if err != nil {
		ac.res.Http500(err)
		return
	}
	dtoArr, err := lib.GetFromCtx[interface{}](ac.ctx, keys.DtoArr)
	if err != nil {
		ac.res.Http500(err)
		return
	}
	dto, err := lib.GetFromCtx[interface{}](ac.ctx, keys.Dto)
	if err != nil {
		ac.res.Http500(err)
		return
	}
	if paramKey == "" {
		if err := ac.req.Gorm.GetAll(modelArr, assocs); err != nil { // 🚨 Aqui pode estar o erro
			ac.res.Http500(err)
			return
		}
		ac.res.DTO(200, modelArr, dtoArr)
	} else {
		paramVal := ac.ctx.Params(paramKey)
		if err := ac.req.Gorm.GetOneBy(paramKey, paramVal, &model, assocs); err != nil {
			ac.res.Http404()
			return
		}
		ac.res.DTO(200, model, dto)
	}
}

// Final Method that executes a FORCE GET action
func (ac *ReqActions) ForceGetBy(paramKey string) {
	if ac.Error != nil {
		return
	}

	keys := namespace.GeneralKey

	modelArr, err := lib.GetFromCtx[interface{}](ac.ctx, keys.ModelArr)
	if err != nil {
		ac.res.Http500(err)
		return
	}
	model, err := lib.GetFromCtx[interface{}](ac.ctx, keys.Model)
	if err != nil {
		ac.res.Http500(err)
		return
	}
	assocs, err := lib.GetFromCtx[[]string](ac.ctx, keys.Associations)
	if err != nil {
		ac.res.Http500(err)
		return
	}
	dtoArr, err := lib.GetFromCtx[interface{}](ac.ctx, keys.DtoArr)
	if err != nil {
		ac.res.Http500(err)
		return
	}
	dto, err := lib.GetFromCtx[interface{}](ac.ctx, keys.Dto)
	if err != nil {
		ac.res.Http500(err)
		return
	}

	if paramKey == "" {
		if err := ac.req.Gorm.ForceGetAll(modelArr, assocs); err != nil {
			ac.res.Http500(err)
			return
		}
		ac.res.DTO(200, modelArr, dtoArr)
		return
	} else {
		paramVal := ac.ctx.Params(paramKey)
		if err := ac.req.Gorm.ForceGetOneBy(paramKey, paramVal, model, assocs); err != nil {
			ac.res.Http404()
			return
		}
		ac.res.DTO(200, model, dto)
		return
	}
}

// Final Method that executes a CREATE action
func (ac *ReqActions) CreateOne() {
	if ac.Error != nil {
		return
	}

	keys := namespace.GeneralKey

	model, err := lib.GetFromCtx[interface{}](ac.ctx, keys.Model)
	if err != nil {
		ac.res.Http500(err)
		return
	}

	if err := ac.req.Gorm.Create(model); err != nil {
		ac.res.Http400(err)
		return
	}

	dto, err := lib.GetFromCtx[interface{}](ac.ctx, keys.Dto)
	if err != nil {
		ac.res.Http500(err)
		return
	}

	ac.res.DTO(201, model, dto)
}

// Final method that executes a DELETE action
func (ac *ReqActions) DeleteOneById() {
	if ac.Error != nil {
		return
	}
	GeneralKey := namespace.GeneralKey
	QueryKey := namespace.QueryKey

	id := ac.ctx.Params(QueryKey.Id)
	model, err := lib.GetFromCtx[interface{}](ac.ctx, GeneralKey.Model)
	if err != nil {
		ac.res.Http500(err)
		return
	}

	if err := ac.req.Gorm.GetOneBy("id", id, model, nil); err != nil {
		if err.Error() == "record not found" {
			ac.res.Http404()
			return
		}
		ac.res.Http500(err)
		return
	}

	if err := ac.req.Gorm.DeleteOneById(id, model); err != nil {
		if err.Error() == "record not found" {
			ac.res.Http404()
			return
		}
		ac.res.Http500(err)
		return
	}

	log.Printf("Deleted record with ID: %s", id)

	ac.res.Http204()
}

// Final method that executes a FORCE DELETE action
func (ac *ReqActions) ForceDeleteOneById() {
	if ac.Error != nil {
		return
	}
	GeneralKey := namespace.GeneralKey
	QueryKey := namespace.QueryKey

	id := ac.ctx.Params(string(QueryKey.Id))
	model, err := lib.GetFromCtx[interface{}](ac.ctx, GeneralKey.Model)
	if err != nil {
		ac.res.Http500(err)
		return
	}

	if err := ac.req.Gorm.ForceGetOneBy("id", id, model, nil); err != nil {
		if err.Error() == "record not found" {
			ac.res.Http404()
			return
		}
		ac.res.Http500(err)
		return
	}

	if err := ac.req.Gorm.ForceDeleteOneById(id, model); err != nil {
		if err.Error() == "record not found" {
			ac.res.Http404()
			return
		}
		ac.res.Http500(err)
		return
	}

	ac.res.Http204()
}

// Final method that executes an UPDATE action
func (ac *ReqActions) UpdateOneById() {
	if ac.Error != nil {
		return
	}
	keys := namespace.GeneralKey
	associations, err := lib.GetFromCtx[[]string](ac.ctx, keys.Associations)
	if err != nil {
		ac.res.Http500(err)
		return
	}
	// changes, err := lib.GetFromCtx[map[string]interface{}](ac.ctx, keys.Changes)
	// if err != nil {
	// 	ac.res.Http500(err)
	// 	return
	// }
	model, err := lib.GetFromCtx[interface{}](ac.ctx, keys.Model)
	if err != nil {
		ac.res.Http500(err)
		return
	}

	id := ac.ctx.Params(namespace.QueryKey.Id)

	if err := ac.req.Gorm.UpdateOneById(id, model, model, associations); err != nil {
		ac.res.Http400(err)
		return
	}

	dto, err := lib.GetFromCtx[interface{}](ac.ctx, keys.Dto)
	if err != nil {
		ac.res.Http500(err)
		return
	}

	ac.res.DTO(200, model, dto)
}
