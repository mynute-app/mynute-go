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

	"github.com/gofiber/fiber/v2"
)

func Create(c *fiber.Ctx, model any) error {
	var err error
	Service := service.New(c)
	defer func() { Service.DeferDB(err) }()
	if err := Service.SetModel(model).Create().Error; err != nil {
		return lib.Error.General.CreatedError.WithError(err)
	}
	return nil
}

func GetOneBy(param string, c *fiber.Ctx, model any, nested_preload *[]string, do_not_load *[]string) error {
	var err error
	Service := service.New(c)
	defer func() { Service.DeferDB(err) }()
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
	defer func() { Service.DeferDB(err) }()
	if err = Service.SetModel(model).SetNestedPreload(nested_preload).UpdateOneById().Error; err != nil {
		return err
	}
	return nil
}

func DeleteOneById(c *fiber.Ctx, model any) error {
	var err error
	Service := service.New(c)
	defer func() { Service.DeferDB(err) }()
	if err = Service.SetModel(model).DeleteOneById().Error; err != nil {
		return err
	}
	return nil
}

func LoginByPassword(user_type string, model any, c *fiber.Ctx) (string, error) {
	var err error
	Service := service.New(c)
	defer func() { Service.DeferDB(err) }()
	token, err := Service.SetModel(model).LoginByPassword(user_type)
	return token, err
}

func LoginByEmailCode(user_type string, model any, c *fiber.Ctx) (string, error) {
	var err error
	Service := service.New(c)
	defer func() { Service.DeferDB(err) }()
	token, err := Service.SetModel(model).LoginByEmailCode(user_type)
	return token, err
}

func ResetLoginvalidationCode(c *fiber.Ctx, user_email string, model any) (string, error) {
	var err error
	Service := service.New(c)
	defer func() { Service.DeferDB(err) }()
	return Service.SetModel(model).ResetLoginCodeByEmail(user_email)
}

func SendLoginValidationCodeByEmail(c *fiber.Ctx, model any) error {
	user_email := c.Params("email")
	if user_email == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("missing 'email' at params route"))
	}

	LoginValidationCode, err := ResetLoginvalidationCode(c, user_email, model)
	if err != nil {
		return err
	}

	language := c.Query("language", "en")

	// Initialize renderer
	renderer := email.NewTemplateRenderer("./static/email", "./translation/email")

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

func ResetPasswordByEmail(c *fiber.Ctx, user_email string, model any) (*DTO.PasswordReseted, error) {
	var err error
	Service := service.New(c)
	defer func() { Service.DeferDB(err) }()
	return Service.SetModel(model).ResetPasswordByEmail(user_email)
}

func SendNewPasswordByEmail(c *fiber.Ctx, user_email string, model any) error {
	password, err := ResetPasswordByEmail(c, user_email, model)
	if err != nil {
		return err
	}

	renderer := email.NewTemplateRenderer("./static/email", "./translation/email")

	language := c.Query("language", "en")

	// Render email template
	renderedEmail, err := renderer.RenderEmail("new_password", language, email.TemplateData{
		"NewPassword": password.Password,
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

func SendVerificationCodeByEmail(c *fiber.Ctx, model any) error {
	var err error
	user_email := c.Params("email")
	if user_email == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("missing 'email' at params route"))
	}

	company_id := c.Params("company_id", "")
	if company_id != "" {
		if err := lib.ChangeToCompanySchemaByContext(c); err != nil {
			return lib.Error.General.BadRequest.WithError(err)
		}
	} else {
		if err := lib.ChangeToPublicSchemaByContext(c); err != nil {
			return lib.Error.General.BadRequest.WithError(err)
		}
	}

	Service := service.New(c)
	defer func() { Service.DeferDB(err) }()

	code, err := Service.SetModel(model).GetVerificationCodeByEmail(user_email)
	if err != nil {
		return err
	}

	language := c.Query("language", "en")

	// Build verification link with code included
	protocol := c.Protocol()
	host := c.Hostname()
	var verificationLink string

	// Detect if this is an admin by checking the model type
	modelType := reflect.TypeOf(model).String()
	if modelType == "*model.Admin" {
		// Admin verification link
		verificationLink = fmt.Sprintf("%s://%s/admin/verify-email/%s/%s?lang=%s", protocol, host, user_email, code, language)
	} else if company_id != "" {
		// Employee verification link
		verificationLink = fmt.Sprintf("%s://%s/verify-email?email=%s&company_id=%s&type=employee&lang=%s&code=%s", protocol, host, user_email, company_id, language, code)
	} else {
		// Client verification link
		verificationLink = fmt.Sprintf("%s://%s/verify-email?email=%s&type=client&lang=%s&code=%s", protocol, host, user_email, language, code)
	}

	// Initialize renderer
	renderer := email.NewTemplateRenderer("./static/email", "./translation/email")

	// Render email template
	renderedEmail, err := renderer.RenderEmail("email_verification_code", language, email.TemplateData{
		"VerificationCode": code,
		"VerificationLink": verificationLink,
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
	if email == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("missing 'email' at params route"))
	}
	Service := service.New(c)
	defer func() { Service.DeferDB(err) }()
	company_id := c.Params("company_id", "")
	if company_id != "" {
		if err := lib.ChangeToCompanySchemaByContext(c); err != nil {
			return lib.Error.General.BadRequest.WithError(err)
		}
	} else {
		if err := lib.ChangeToPublicSchemaByContext(c); err != nil {
			return lib.Error.General.BadRequest.WithError(err)
		}
	}
	code := c.Params("code", "")
	if code == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("missing 'code' at params route"))
	}
	return Service.SetModel(model).VerifyEmail(email, code)
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
