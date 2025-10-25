package controller

import (
	"context"
	"fmt"
	DTO "mynute-go/core/src/config/api/dto"
	mJSON "mynute-go/core/src/config/db/model/json"
	"mynute-go/core/src/lib"
	"mynute-go/core/src/lib/email"
	"mynute-go/core/src/service"
	"reflect"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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

func GetOneBy(param string, c *fiber.Ctx, model any, nested_preload *[]string, do_not_load *[]string) error {
	var err error
	Service := service.New(c)
	defer Service.DeferDB(err)
	if err = Service.
		SetModel(model).
		SetNestedPreload(nested_preload).
		SetDoNotLoad(do_not_load).
		GetBy(param).Error; err != nil {
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

func LoginByPassword(user_type string, model any, c *fiber.Ctx) (string, error) {
	var err error
	Service := service.New(c)
	defer Service.DeferDB(err)
	token, err := Service.SetModel(model).LoginByPassword(user_type)
	return token, err
}

func LoginByEmailCode(user_type string, model any, c *fiber.Ctx) (string, error) {
	var err error
	Service := service.New(c)
	defer Service.DeferDB(err)
	token, err := Service.SetModel(model).LoginByEmailCode(user_type)
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

func SendLoginValidationCodeByEmail(c *fiber.Ctx, model any) error {
	user_email := c.Query("email")
	language := c.Query("lang", "en")

	if user_email == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("missing 'email' at params route"))
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	if err := tx.Model(model).Where("email = ?", user_email).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.RecordNotFound
		}
		return err
	}

	// Initialize renderer
	renderer := email.NewTemplateRenderer("./static", "./translation")

	LoginValidationCode := lib.GenerateRandomInt(6)
	codeString := fmt.Sprintf("%d", LoginValidationCode)

	// Set code expiration to 15 minutes from now
	expiryTime := time.Now().Add(15 * time.Minute)

	// Store the code in the database using reflection
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}

	metaField := modelValue.FieldByName("Meta")
	if !metaField.IsValid() || !metaField.CanSet() {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("model does not have a Meta field"))
	}

	// Get the current Meta value or create a new one if nil
	metaValue := metaField.Interface().(mJSON.UserMeta)

	// Set the validation code and expiry
	metaValue.LoginValidationCode = &codeString
	metaValue.LoginValidationExpiry = &expiryTime

	// Set the updated Meta back to the model
	metaField.Set(reflect.ValueOf(metaValue))

	// Update the model in database
	if err := tx.Model(model).Where("email = ?", user_email).Updates(model).Error; err != nil {
		return fmt.Errorf("failed to store validation code: %w", err)
	}

	// Render email template
	renderedEmail, err := renderer.RenderEmail("login_validation_code", language, email.TemplateData{
		"LoginValidationCode": LoginValidationCode,
	})

	if err != nil {
		return fmt.Errorf("failed to render email: %w", err)
	}

	// Initialize email provider
	provider, err := email.NewProvider(nil)
	if err != nil {
		return fmt.Errorf("failed to initialize email provider: %w", err)
	}

	// Send email
	err = provider.Send(context.Background(), email.EmailData{
		To:      []string{user_email},
		Subject: renderedEmail.Subject,
		Html:    renderedEmail.HTMLBody,
	})

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
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
	designField := elem.FieldByName("Design")
	metaField := elem.FieldByName("Meta")

	// Check if model has Design field directly or through Meta field
	var hasDirectDesign bool = designField.IsValid()
	var hasMetaDesign bool = metaField.IsValid()

	if !hasDirectDesign && !hasMetaDesign {
		return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("model (%s) does not have a Design field or Meta field", model_table_name))
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

	if hasDirectDesign {
		// For models with direct Design field (Company, Branch, Service)
		if err := tx.Model(model).Where("id = ?", id).Pluck("design", &Design).Error; err != nil {
			return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("failed to fetch model (%s) with id (%s): %w", model_table_name, id, err))
		}
	} else {
		// For models with Meta field (Client, Employee)
		var meta mJSON.UserMeta
		if err := tx.Model(model).Where("id = ?", id).Pluck("meta", &meta).Error; err != nil {
			return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("failed to fetch model (%s) with id (%s): %w", model_table_name, id, err))
		}
		Design = meta.Design
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

	if hasDirectDesign {
		// For models with direct Design field (Company, Branch, Service)
		if err := tx.Model(model).Where("id = ?", id).Update("design", Design).Error; err != nil {
			return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("failed to update Design field for model (%s) with id (%s): %w", model_table_name, id, err))
		}
	} else {
		// For models with Meta field (Client, Employee)
		var meta mJSON.UserMeta
		if err := tx.Model(model).Where("id = ?", id).Pluck("meta", &meta).Error; err != nil {
			return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("failed to fetch meta for model (%s) with id (%s): %w", model_table_name, id, err))
		}
		meta.Design = Design
		if err := tx.Model(model).Where("id = ?", id).Update("meta", meta).Error; err != nil {
			return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("failed to update Meta field for model (%s) with id (%s): %w", model_table_name, id, err))
		}
	}

	return &Design, nil
}

func DeleteImageById(c *fiber.Ctx, model_table_name string, model any, img_types_allowed map[string]bool) (*mJSON.DesignConfig, error) {
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() != reflect.Ptr || modelValue.IsNil() {
		return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("model (%s) must be a non-nil pointer", model_table_name))
	}

	elem := modelValue.Elem()
	designField := elem.FieldByName("Design")
	metaField := elem.FieldByName("Meta")

	// Check if model has Design field directly or through Meta field
	var hasDirectDesign bool = designField.IsValid()
	var hasMetaDesign bool = metaField.IsValid()

	if !hasDirectDesign && !hasMetaDesign {
		return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("model (%s) does not have a Design field or Meta field", model_table_name))
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

	if hasDirectDesign {
		// For models with direct Design field (Company, Branch, Service)
		if err := tx.Model(model).Where("id = ?", id).Pluck("design", &Design).Error; err != nil {
			return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("failed to fetch model (%s) with id (%s): %w", model_table_name, id, err))
		}
	} else {
		// For models with Meta field (Client, Employee)
		var meta mJSON.UserMeta
		if err := tx.Model(model).Where("id = ?", id).Pluck("meta", &meta).Error; err != nil {
			return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("failed to fetch model (%s) with id (%s): %w", model_table_name, id, err))
		}
		Design = meta.Design
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

	if hasDirectDesign {
		// For models with direct Design field (Company, Branch, Service)
		if err := tx.Model(model).Where("id = ?", id).Update("design", Design).Error; err != nil {
			return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("failed to update Design field for model (%s) with id (%s): %w", model_table_name, id, err))
		}
	} else {
		// For models with Meta field (Client, Employee)
		var meta mJSON.UserMeta
		if err := tx.Model(model).Where("id = ?", id).Pluck("meta", &meta).Error; err != nil {
			return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("failed to fetch meta for model (%s) with id (%s): %w", model_table_name, id, err))
		}
		meta.Design = Design
		if err := tx.Model(model).Where("id = ?", id).Update("meta", meta).Error; err != nil {
			return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("failed to update Meta field for model (%s) with id (%s): %w", model_table_name, id, err))
		}
	}

	return &Design, nil
}
