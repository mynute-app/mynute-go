package controller

import (
	"fmt"
	DTO "mynute-go/core/config/api/dto"
	mJSON "mynute-go/core/config/db/model/json"
	"mynute-go/core/lib"
	"mynute-go/core/service"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

func Create(c *fiber.Ctx, model any) error {
	var err error
	Service := service.New(c)
	defer Service.DeferDB(err)
	if err := Service.SetModel(model).Create().Error; err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}
	return nil
}

func GetOneBy(param string, c *fiber.Ctx, model any, nested_preload *[]string) error {
	var err error
	Service := service.New(c)
	defer Service.DeferDB(err)
	if err = Service.SetModel(model).SetNestedPreload(nested_preload).GetBy(param).Error; err != nil {
		return err
	}
	return nil
}

func UpdateOneById(c *fiber.Ctx, model any, nested_preload *[]string) error {
	var err error
	Service := service.New(c)
	defer Service.DeferDB(err)
	if err = Service.SetModel(model).SetNestedPreload(nested_preload).UpdateOneById().Error; err != nil {
		return err
	}
	return nil
}

func DeleteOneById(c *fiber.Ctx, model any) error {
	var err error
	Service := service.New(c)
	defer Service.DeferDB(err)
	if err = Service.SetModel(model).DeleteOneById().Error; err != nil {
		return err
	}
	return nil
}

func Login(user_type string, model any, c *fiber.Ctx) (string, error) {
	var err error
	Service := service.New(c)
	defer Service.DeferDB(err)
	token, err := Service.SetModel(model).Login(user_type)
	return token, err
}

func ResetPasswordByEmail(c *fiber.Ctx, model any) (DTO.PasswordReseted, error) {
	var err error
	email := c.Params("email")
	Service := service.New(c)
	defer Service.DeferDB(err)
	new_pass, err := Service.SetModel(model).ResetPasswordByEmail(email)
	return new_pass, err
}

func VerifyEmail(c *fiber.Ctx, model any) error {
	var err error
	email := c.Params("email")
	Service := service.New(c)
	defer Service.DeferDB(err)
	return Service.SetModel(model).VerifyEmail(email)
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

	var Design mJSON.DesignConfig
	if err := tx.Model(model).Where("id = ?", id).Pluck("design", &Design).Error; err != nil {
		return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("failed to fetch model (%s) with id (%s): %w", model_table_name, id, err))
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

	var Design mJSON.DesignConfig
	if err := tx.Model(model).Where("id = ?", id).Pluck("design", &Design).Error; err != nil {
		return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("failed to fetch model (%s) with id (%s): %w", model_table_name, id, err))
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
