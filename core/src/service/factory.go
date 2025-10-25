package service

import (
	"encoding/json"
	"fmt"
	DTO "mynute-go/core/src/config/api/dto"
	database "mynute-go/core/src/config/db"
	mJSON "mynute-go/core/src/config/db/model/json"
	"mynute-go/core/src/handler"
	"mynute-go/core/src/lib"
	"net/url"
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func New(c *fiber.Ctx) *service {
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
	Model         any
	Context       *fiber.Ctx
	MyGorm        *handler.Gorm
	DeferDB       func(err error)
	NestedPreload *[]string
	DoNotLoad     *[]string
	Error         error
}

func (s *service) SetNestedPreload(preloads *[]string) *service {
	s.NestedPreload = preloads
	return s
}

func (s *service) SetDoNotLoad(do_not_load *[]string) *service {
	if s.Error != nil {
		return s
	}
	s.MyGorm.SetDoNotLoad(do_not_load)
	return s
}

func (s *service) SetModel(model any) *service {
	if s.Error != nil {
		return s
	}
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() != reflect.Ptr || modelValue.IsNil() {
		s.Error = lib.Error.General.InternalError.WithError(fmt.Errorf("model must be a non-nil pointer"))
		return s
	}
	if modelValue.Elem().Kind() != reflect.Struct {
		s.Error = lib.Error.General.InternalError.WithError(fmt.Errorf("model must point to a struct"))
		return s
	}
	s.Model = model
	return s
}

func (s *service) get_param(param string) (string, error) {
	paramVal := s.Context.Params(param)
	if paramVal == "" {
		return "", lib.Error.General.BadRequest.WithError(fmt.Errorf("parameter %s not found on route parameters", param))
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
	if err := s.MyGorm.SetNestedPreload(s.NestedPreload).GetAll(model); err != nil {
		s.Error = lib.Error.General.RecordNotFound.WithError(err)
	}
	return s
}

func (s *service) GetBy(param string) *service {
	if s.Error != nil {
		return s
	}
	if param == "" {
		return s.GetAll(s.Model)
	}
	val, err := s.get_param(param)
	if err != nil {
		s.Error = err
		return s
	}
	if err := s.MyGorm.SetNestedPreload(s.NestedPreload).GetOneBy(param, val, s.Model); err != nil {
		s.Error = lib.Error.General.RecordNotFound.WithError(err)
	}
	return s
}

func (s *service) ForceGetBy(param string) *service {
	if s.Error != nil {
		return s
	}
	if param == "" {
		return s.GetAll(s.Model)
	}
	val, err := s.get_param(param)
	if err != nil {
		s.Error = err
		return s
	}
	if err := s.MyGorm.SetNestedPreload(s.NestedPreload).ForceGetOneBy(param, val, s.Model); err != nil {
		s.Error = lib.Error.General.RecordNotFound.WithError(err)
	}
	return s
}

func (s *service) Create() *service {
	if s.Error != nil {
		return s
	}
	if err := s.Context.BodyParser(s.Model); err != nil {
		s.Error = lib.Error.General.InternalError.WithError(err)
		return s
	}
	if err := s.MyGorm.Create(s.Model); err != nil {
		s.Error = err
		return s
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
	id, ok := HasID(s.Model)
	if !ok {
		s.Error = lib.Error.General.InternalError.WithError(fmt.Errorf("model does not have ID field"))
		return s
	}
	if err := s.MyGorm.GetOneBy("id", id.String(), s.Model); err != nil {
		s.Error = lib.Error.General.CreatedError.WithError(err)
		return s
	}
	return s
}

func (s *service) UpdateOneById() *service {
	if s.Error != nil {
		return s
	}
	val, err := s.get_param("id")
	if err != nil {
		s.Error = err
		return s
	}
	if err := s.Context.BodyParser(s.Model); err != nil {
		s.Error = lib.Error.General.InternalError.WithError(err)
		return s
	}
	if err := s.MyGorm.SetNestedPreload(s.NestedPreload).UpdateOneById(val, s.Model); err != nil {
		s.Error = lib.Error.General.UpdatedError.WithError(err)
		return s
	}
	return s
}

func (s *service) DeleteOneById() *service {
	if s.Error != nil {
		return s
	}
	val, err := s.get_param("id")
	if err != nil {
		return s
	}
	if err := s.MyGorm.DeleteOneById(val, s.Model); err != nil {
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

func (s *service) LoginByPassword(user_type string) (string, error) {
	if s.Error != nil {
		return "", s.Error
	}
	var err error
	var body DTO.LoginClient
	if err := s.Context.BodyParser(&body); err != nil {
		return "", err
	}
	if err := s.MyGorm.DB.
		Model(s.Model).
		Where("email = ?", body.Email).
		First(s.Model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", lib.Error.Client.NotFound
		}
		return "", lib.Error.General.InternalError.WithError(err)
	}

	val := reflect.ValueOf(s.Model)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	verifiedField := val.FieldByName("Verified")
	passwordField := val.FieldByName("Password")

	if !verifiedField.IsValid() || !passwordField.IsValid() {
		return "", lib.Error.General.InternalError.WithError(fmt.Errorf("model must have Verified and Password fields"))
	}

	verified := verifiedField.Bool()
	password := passwordField.String()

	if !verified {
		return "", lib.Error.Client.NotVerified
	}
	if !handler.ComparePassword(password, body.Password) {
		return "", lib.Error.Auth.InvalidLogin
	}

	userBytes, err := json.Marshal(s.Model)
	if err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}
	var claims DTO.Claims
	if err := json.Unmarshal(userBytes, &claims); err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}
	claims.Type = user_type
	token, err := handler.JWT(s.Context).Encode(&claims)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *service) LoginByEmailCode(user_type string) (string, error) {
	if s.Error != nil {
		return "", s.Error
	}
	var err error
	var body DTO.LoginByEmailCode
	if err := s.Context.BodyParser(&body); err != nil {
		return "", err
	}

	// Find user by email
	if err := s.MyGorm.DB.
		Model(s.Model).
		Where("email = ?", body.Email).
		First(s.Model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", lib.Error.Client.NotFound
		}
		return "", lib.Error.General.InternalError.WithError(err)
	}

	val := reflect.ValueOf(s.Model)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	verifiedField := val.FieldByName("Verified")
	metaField := val.FieldByName("Meta")

	if !verifiedField.IsValid() || !metaField.IsValid() {
		return "", lib.Error.General.InternalError.WithError(fmt.Errorf("model must have Verified and Meta fields"))
	}

	// Get the Meta field and extract validation code info
	metaValue := metaField.Interface().(mJSON.UserMeta)

	// Check if code field is nil or empty
	if metaValue.LoginValidationCode == nil {
		return "", lib.Error.Auth.InvalidLogin.WithError(fmt.Errorf("no validation code found"))
	}

	// Check if expiry field is nil
	if metaValue.LoginValidationExpiry == nil {
		return "", lib.Error.Auth.InvalidLogin.WithError(fmt.Errorf("validation code has expired"))
	}

	storedCode := *metaValue.LoginValidationCode
	expiryTime := *metaValue.LoginValidationExpiry

	// Check if code has expired
	if time.Now().After(expiryTime) {
		return "", lib.Error.Auth.InvalidLogin.WithError(fmt.Errorf("validation code has expired"))
	}

	// Validate the code
	if storedCode != body.Code {
		return "", lib.Error.Auth.InvalidLogin.WithError(fmt.Errorf("invalid validation code"))
	}

	verified := verifiedField.Bool()

	// Verifies the user in case not verified
	if !verified {
		verifiedField.SetBool(true)
	}

	// Clear the validation code after successful login
	metaValue.LoginValidationCode = nil
	metaValue.LoginValidationExpiry = nil
	metaField.Set(reflect.ValueOf(metaValue))

	// Update the model to clear the validation code
	if err := s.MyGorm.DB.Model(s.Model).Where("email = ?", body.Email).Updates(s.Model).Error; err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}

	// Generate JWT token
	userBytes, err := json.Marshal(s.Model)
	if err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}
	var claims DTO.Claims
	if err := json.Unmarshal(userBytes, &claims); err != nil {
		return "", lib.Error.General.InternalError.WithError(err)
	}
	claims.Type = user_type
	token, err := handler.JWT(s.Context).Encode(&claims)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *service) prepare_email(email string) (string, error) {
	if email == "" {
		return "", lib.Error.General.BadRequest.WithError(fmt.Errorf("email parameter is empty"))
	}
	cleanedEmail, err := url.QueryUnescape(email)
	if err != nil {
		return "", lib.Error.General.BadRequest.WithError(err)
	}
	if err := lib.ValidatorV10.Var(cleanedEmail, "email"); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return "", lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid email: %w", err))
		} else {
			return "", lib.Error.General.InternalError.WithError(err)
		}
	}
	return cleanedEmail, nil
}

func (s *service) ResetPasswordByEmail(email string) (DTO.PasswordReseted, error) {
	if s.Error != nil {
		return DTO.PasswordReseted{}, s.Error
	}
	email, err := s.prepare_email(email)
	if err != nil {
		return DTO.PasswordReseted{}, lib.Error.General.BadRequest.WithError(err)
	}
	newPassword := lib.GenerateValidPassword()
	if err := s.MyGorm.DB.
		Model(s.Model).
		Where("email = ?", email).
		Update("password", newPassword).Error; err != nil {
		return DTO.PasswordReseted{}, lib.Error.General.InternalError.WithError(err)
	}
	return DTO.PasswordReseted{Password: newPassword}, nil
}

func (s *service) VerifyEmail(email string) error {
	if s.Error != nil {
		return s.Error
	}
	email, err := s.prepare_email(email)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}
	if err := s.MyGorm.DB.
		Model(s.Model).
		Where("email = ?", email).
		Update("verified", true).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}
