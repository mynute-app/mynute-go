package service

import (
	database "agenda-kaki-go/core/config/db"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"fmt"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

func Factory(c *fiber.Ctx) *service {
	var err error
	tx, end, err := database.ContextTransaction(c)
	service := &service{
		Context: c,
		Error:   err,
		MyGorm:  handler.MyGormWrapper(tx),
		DeferDB: end,
	}
	return service
}

type service struct {
	Context *fiber.Ctx
	MyGorm  *handler.Gorm
	DeferDB func(err error)
	Error   error
}

func (s *service) get_param(param string) (string, error) {
	paramVal := s.Context.Params(param)
	if paramVal == "" {
		return "", lib.Error.General.NotFoundError.WithError(fmt.Errorf("parameter %s not found on route parameters", param))
	}
	cleanedParamVal, err := url.QueryUnescape(paramVal)
	if err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}
	return cleanedParamVal, nil
}

func (s *service) GetAll(model any) *service {
	if s.Error != nil {
		return s
	}
	if err := s.MyGorm.GetAll(model); err != nil {
		s.Error = lib.Error.General.RecordNotFound.WithError(err)
	}
	return s
}

func (s *service) GetBy(param string, model any) *service {
	if s.Error != nil {
		return s
	}
	if param == "" {
		return s.GetAll(model)
	}
	val, err := s.get_param(param)
	if err != nil {
		s.Error = err
		return s
	}
	if err := s.MyGorm.GetOneBy(param, val, model); err != nil {
		s.Error = lib.Error.General.RecordNotFound.WithError(err)
	}
	return s
}

func (s *service) ForceGetBy(param string, model any) *service {
	if s.Error != nil {
		return s
	}
	if param == "" {
		return s.GetAll(model)
	}
	val, err := s.get_param(param)
	if err != nil {
		s.Error = err
		return s
	}
	if err := s.MyGorm.ForceGetOneBy(param, val, model); err != nil {
		s.Error = lib.Error.General.RecordNotFound.WithError(err)
	}
	return s
}

func (s *service) Create(model any) *service {
	if s.Error != nil {
		return s
	}
	if err := s.MyGorm.Create(model); err != nil {
		s.Error = err
	}
	return s
}

func (s *service) UpdateOneById(model any) *service {
	if s.Error != nil {
		return s
	}
	val, err := s.get_param("id")
	if err != nil {
		s.Error = err
		return s
	}
	changes := make(map[string]any)
	if err := s.Context.BodyParser(&changes); err != nil {
		s.Error = lib.Error.General.InternalError.WithError(err)
		return s
	}
	if err := s.MyGorm.UpdateOneById(val, model, changes); err != nil {
		s.Error = lib.Error.General.UpdatedError.WithError(err)
		return s
	}
	if err := s.MyGorm.GetOneBy("id", val, model); err != nil {
		s.Error = lib.Error.General.RecordNotFound.WithError(err)
		return s
	}
	return s
}

func (s *service) DeleteOneById(model any) *service {
	if s.Error != nil {
		return s
	}
	val, err := s.get_param("id")
	if err != nil {
		return s
	}
	if err := s.MyGorm.DeleteOneById(val, model); err != nil {
		s.Error = lib.Error.General.RecordNotFound.WithError(err)
	}
	return s
}

func (s *service) ForceDeleteOneById(model any) *service {
	if s.Error != nil {
		return s
	}
	val, err := s.get_param("id")
	if err != nil {
		return s
	}
	if err := s.MyGorm.ForceDeleteOneById(val, model); err != nil {
		s.Error = lib.Error.General.RecordNotFound.WithError(err)
	}
	return s
}
