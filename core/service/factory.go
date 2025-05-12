package service

import (
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"fmt"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

func Factory(c *fiber.Ctx) *service {
	tx, err := lib.Session(c)
	service := &service{
		Context: c,
		err:     err,
		MyGorm:  handler.MyGormWrapper(tx),
	}
	return service
}

type service struct {
	Context *fiber.Ctx
	MyGorm  *handler.Gorm
	err     error
	MyModel any
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

func (s *service) Model(model any) *service {
	s.MyModel = model
	return s
}

func (s *service) GetAll() error {
	if s.err != nil {
		return s.err
	}
	if err := s.MyGorm.GetAll(s.MyModel); err != nil {
		return lib.Error.General.RecordNotFound.WithError(err)
	}
	return nil
}

func (s *service) GetBy(param string) error {
	if s.err != nil {
		return s.err
	}
	if param == "" {
		return s.GetAll()
	}
	val, err := s.get_param(param)
	if err != nil {
		return err
	}
	if err := s.MyGorm.GetOneBy(param, val, s.MyModel); err != nil {
		return lib.Error.General.RecordNotFound.WithError(err)
	}
	return nil
}

func (s *service) ForceGetBy(param string) error {
	if s.err != nil {
		return s.err
	}
	if param == "" {
		return s.GetAll()
	}
	val, err := s.get_param(param)
	if err != nil {
		return err
	}
	if err := s.MyGorm.ForceGetOneBy(param, val, s.MyModel); err != nil {
		return lib.Error.General.RecordNotFound.WithError(err)
	}
	return nil
}

func (s *service) Create() error {
	if s.err != nil {
		return s.err
	}
	if err := s.MyGorm.Create(s.MyModel); err != nil {
		return err
	}
	return nil
}

func (s *service) UpdateOneById() error {
	if s.err != nil {
		return s.err
	}
	val, err := s.get_param("id")
	if err != nil {
		return err
	}
	changes := make(map[string]any)
	if err := s.Context.BodyParser(&changes); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	if err := s.MyGorm.UpdateOneById(val, s.MyModel, changes); err != nil {
		return lib.Error.General.RecordNotFound.WithError(err)
	}
	if err := s.MyGorm.GetOneBy("id", val, s.MyModel); err != nil {
		return lib.Error.General.RecordNotFound.WithError(err)
	}
	return nil
}

func (s *service) DeleteOneById() error {
	if s.err != nil {
		return s.err
	}
	val, err := s.get_param("id")
	if err != nil {
		return err
	}
	if err := s.MyGorm.DeleteOneById(val, s.MyModel); err != nil {
		return lib.Error.General.RecordNotFound.WithError(err)
	}
	return nil
}

func (s *service) ForceDeleteOneById() error {
	if s.err != nil {
		return s.err
	}
	val, err := s.get_param("id")
	if err != nil {
		return err
	}
	if err := s.MyGorm.ForceDeleteOneById(val, s.MyModel); err != nil {
		return lib.Error.General.RecordNotFound.WithError(err)
	}
	return nil
}
