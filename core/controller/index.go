package controller

import (
	mJSON "agenda-kaki-go/core/config/db/model/json"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/service"
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Create(c *fiber.Ctx, model any) error {
	var err error
	if err := c.BodyParser(model); err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}
	Service := service.Factory(c)
	defer Service.DeferDB(err)
	if err := Service.Create(model).Error; err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}
	HasID := func(model any) (uuid.UUID, bool) {
		val := reflect.ValueOf(model)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		field := val.FieldByName("ID")
		if !field.IsValid() || field.Type() != reflect.TypeOf(uuid.UUID{}) {
			return uuid.Nil, false
		}
		return field.Interface().(uuid.UUID), true
	}
	id, ok := HasID(model)
	if !ok {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("model does have ID field"))
	}
	if err := Service.MyGorm.GetOneBy("id", id.String(), model); err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}
	return nil
}

func GetOneBy(param string, c *fiber.Ctx, model any) error {
	var err error
	Service := service.Factory(c)
	defer Service.DeferDB(err)
	if err := Service.GetBy(param, model).Error; err != nil {
		return err
	}
	return nil
}

func UpdateOneById(c *fiber.Ctx, model any) error {
	var err error
	Service := service.Factory(c)
	defer Service.DeferDB(err)
	return Service.UpdateOneById(model).Error
}

func DeleteOneById(c *fiber.Ctx, model any) error {
	var err error
	Service := service.Factory(c)
	defer Service.DeferDB(err)
	return Service.DeleteOneById(model).Error
}

func UpdateImagesById(c *fiber.Ctx, model_table_name string, model any, img_types_allowed map[string]bool) (*mJSON.DesignConfig, error) {
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() != reflect.Ptr || modelValue.IsNil() {
		return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("model (%s) must be a non-nil pointer", model_table_name))
	}

	elem := modelValue.Elem()
	field := elem.FieldByName("Design")
	if !field.IsValid() {
		return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("model (%s) does not have a Design field", model_table_name))
	}
	id := c.Params("id")
	if id == "" {
		return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("id parameter is required at path"))
	}
	tx, err := lib.Session(c)
	if err != nil {
		return nil, err
	}
	if err := tx.First(model, "id = ?", id).Error; err != nil {
		return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("failed to fetch model (%s) with id (%s): %w", model_table_name, id, err))
	}

	modelElementValue := reflect.ValueOf(model).Elem()
	Design, ok := modelElementValue.FieldByName("Design").Interface().(mJSON.DesignConfig)
	if !ok {
		return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("model (%s) does not have a valid Design field", model_table_name))
	}

	var uploaded_img_types = make([]string, 0)

	defer func() {
		r := recover()
		if r != nil || err != nil {
			for _, img_type := range uploaded_img_types {
				_ = Design.Images.Delete(img_type, model_table_name, id)
			}
		}
	}()

	for img_type := range img_types_allowed {
		file, err := c.FormFile(img_type)
		if err != nil {
			continue
		}
		_, err = Design.Images.Save(img_type, model_table_name, id, file)
		if err != nil {
			return nil, err
		}
		uploaded_img_types = append(uploaded_img_types, img_type)
	}

	if err := tx.Model(model).Where("id = ?", id).Update("design", Design).Error; err != nil {
		return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("failed to update Design field for model (%s) with id (%s): %w", model_table_name, id, err))
	}

	return &Design, nil
}

func DeleteImageById(c *fiber.Ctx, model_table_name string, model any, img_types_allowed map[string]bool) (*mJSON.DesignConfig, error) {
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() != reflect.Ptr || modelValue.IsNil() {
		return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("model (%s) must be a non-nil pointer", model_table_name))
	}

	elem := modelValue.Elem()
	field := elem.FieldByName("Design")
	if !field.IsValid() {
		return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("model (%s) does not have a Design field", model_table_name))
	}
	id := c.Params("id")
	if id == "" {
		return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("id parameter is required at path"))
	}
	tx, err := lib.Session(c)
	if err != nil {
		return nil, err
	}
	if err := tx.First(model, "id = ?", id).Error; err != nil {
		return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("failed to fetch model (%s) with id (%s): %w", model_table_name, id, err))
	}

	modelElementValue := reflect.ValueOf(model).Elem()
	Design, ok := modelElementValue.FieldByName("Design").Interface().(mJSON.DesignConfig)
	if !ok {
		return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("model (%s) does not have a valid Design field", model_table_name))
	}

	var uploaded_img_types = make([]string, 0)

	defer func() {
		r := recover()
		if r != nil || err != nil {
			for _, img_type := range uploaded_img_types {
				_ = Design.Images.Delete(img_type, model_table_name, id)
			}
		}
	}()

	image_type := c.Params("image_type")
	if image_type == "" {
		return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("image_type parameter is required at path"))
	}

	if _, ok := img_types_allowed[image_type]; !ok {
		return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("image_type (%s) is not allowed for model (%s)", image_type, model_table_name))
	}

	if err := Design.Images.Delete(image_type, model_table_name, id); err != nil {
		return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("failed to delete image of type (%s) for model (%s) with id (%s): %w", image_type, model_table_name, id, err))
	}

	if err := tx.Model(model).Where("id = ?", id).Update("design", Design).Error; err != nil {
		return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("failed to update Design field for model (%s) with id (%s): %w", model_table_name, id, err))
	}

	return &Design, nil
}