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
	if metaValue.Login.ValidationCode == nil {
		return "", lib.Error.Auth.InvalidLogin.WithError(fmt.Errorf("no validation code found"))
	}

	// Check if expiry field is nil
	if metaValue.Login.ValidationExpiry == nil {
		return "", lib.Error.Auth.InvalidLogin.WithError(fmt.Errorf("validation code has expired"))
	}

	storedCode := *metaValue.Login.ValidationCode
	expiryTime := *metaValue.Login.ValidationExpiry

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
	metaValue.Login.ValidationCode = nil
	metaValue.Login.ValidationExpiry = nil
	metaValue.Login.ValidationRequestedAt = nil
	metaField.Set(reflect.ValueOf(metaValue))

	// Update the model to clear the validation code
	// Only update Meta and Verified fields, not Password
	if err := s.MyGorm.DB.Model(s.Model).
		Where("email = ?", body.Email).
		Select("Meta", "Verified").
		Updates(s.Model).Error; err != nil {
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

func (s *service) prepare_email(user_email string) (string, error) {
	if user_email == "" {
		return "", lib.Error.General.BadRequest.WithError(fmt.Errorf("email parameter is empty"))
	}
	cleanedEmail, err := url.QueryUnescape(user_email)
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

func (s *service) ResetLoginCodeByEmail(user_email string) (string, error) {
	user_email, err := s.prepare_email(user_email)
	if err != nil {
		return "", lib.Error.General.BadRequest.WithError(err)
	}
	if err := s.MyGorm.DB.Model(s.Model).Where("email = ?", user_email).First(&s.Model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", lib.Error.General.RecordNotFound
		}
		return "", err
	}

	// Store the code in the database using reflection
	modelValue := reflect.ValueOf(s.Model)
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}

	metaField := modelValue.FieldByName("Meta")
	if !metaField.IsValid() || !metaField.CanSet() {
		return "", lib.Error.General.BadRequest.WithError(fmt.Errorf("model does not have a Meta field"))
	}

	metaValue := metaField.Interface().(mJSON.UserMeta)

	if metaValue.Login.ValidationRequestedAt != nil && metaValue.Login.ValidationRequestsCount != nil {
		waitTime := time.Minute
		if *metaValue.Login.ValidationRequestsCount > 3 {
			waitTime = 5 * time.Minute
		}
		if time.Since(*metaValue.Login.ValidationRequestedAt) < waitTime {
			return "", lib.Error.General.TooManyRequests.WithError(fmt.Errorf("a validation code was sent recently; please wait before requesting another"))
		}
	}

	LoginValidationCode := lib.GenerateRandomInt(6)
	codeString := fmt.Sprintf("%d", LoginValidationCode)

	// Set code expiration to 15 minutes from now
	expiryTime := time.Now().Add(15 * time.Minute)

	// Set the validation code and expiry
	metaValue.Login.ValidationCode = &codeString
	metaValue.Login.ValidationExpiry = &expiryTime
	now := time.Now()
	metaValue.Login.ValidationRequestedAt = &now

	// Increment request count
	if metaValue.Login.ValidationRequestsCount == nil {
		count := 1
		metaValue.Login.ValidationRequestsCount = &count
	} else {
		*metaValue.Login.ValidationRequestsCount++
	}

	// Set the updated Meta back to the model
	metaField.Set(reflect.ValueOf(metaValue))

	// Update the model in database
	if err := s.MyGorm.DB.
		Model(s.Model).
		Where("email = ?", user_email).
		Updates(s.Model).Error; err != nil {
		return "", fmt.Errorf("failed to store validation code: %w", err)
	}

	return codeString, nil
}

func (s *service) ResetPasswordByEmail(user_email string) (*DTO.PasswordReseted, error) {
	if s.Error != nil {
		return nil, s.Error
	}
	user_email, err := s.prepare_email(user_email)
	if err != nil {
		return nil, lib.Error.General.BadRequest.WithError(err)
	}
	// Get the user by email
	if err := s.MyGorm.DB.
		Model(s.Model).
		Where("email = ?", user_email).
		First(s.Model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, lib.Error.General.RecordNotFound
		}
		return nil, lib.Error.General.InternalError.WithError(err)
	}

	modelValue := reflect.ValueOf(s.Model)
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}

	metaField := modelValue.FieldByName("Meta")
	if !metaField.IsValid() || !metaField.CanSet() {
		return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("model does not have a Meta field"))
	}

	// Get the current Meta value or create a new one if nil
	metaValue := metaField.Interface().(mJSON.UserMeta)

	if metaValue.Login.NewPasswordRequestedAt != nil && metaValue.Login.NewPasswordRequestsCount != nil {
		waitTime := time.Minute
		if *metaValue.Login.NewPasswordRequestsCount > 3 {
			waitTime = 5 * time.Minute
		}
		if time.Since(*metaValue.Login.NewPasswordRequestedAt) < waitTime {
			return nil, lib.Error.General.TooManyRequests.WithError(fmt.Errorf("a validation code was sent recently; please wait before requesting another"))
		}
	}

	passwordField := modelValue.FieldByName("Password")
	if !passwordField.IsValid() || !passwordField.CanSet() {
		return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("model does not have a Password field"))
	}
	newPassword := lib.GenerateValidPassword()
	passwordField.SetString(newPassword)

	// Update request tracking
	now := time.Now()
	metaValue.Login.NewPasswordRequestedAt = &now

	// Increment request count
	if metaValue.Login.NewPasswordRequestsCount == nil {
		count := 1
		metaValue.Login.NewPasswordRequestsCount = &count
	} else {
		*metaValue.Login.NewPasswordRequestsCount++
	}

	// Set the updated Meta back to the model
	metaField.Set(reflect.ValueOf(metaValue))

	if err := s.MyGorm.DB.
		Model(s.Model).
		Where("email = ?", user_email).
		Updates(s.Model).Error; err != nil {
		return nil, lib.Error.General.InternalError.WithError(err)
	}

	return &DTO.PasswordReseted{Password: newPassword}, nil
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
