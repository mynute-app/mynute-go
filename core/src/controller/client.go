package controller

import (
	"fmt"
	DTO "mynute-go/core/src/config/api/dto"
	dJSON "mynute-go/core/src/config/api/dto/json"
	"mynute-go/core/src/config/db/model"
	"mynute-go/core/src/config/namespace"
	"mynute-go/core/src/handler"
	"mynute-go/core/src/lib"
	"mynute-go/core/src/middleware"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateClient creates an client
//
//	@Summary		Create client
//	@Description	Create an client
//	@Tags			Client
//	@Accept			json
//	@Produce		json
//	@Param			client	body		DTO.CreateClient	true	"Client"
//	@Success		200		{object}	DTO.Client
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/client [post]
func CreateClient(c *fiber.Ctx) error {
	var err error
	var client model.Client
	if err := c.BodyParser(&client); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	tx, err := lib.Session(c)
	if err != nil {
		return err
	}
	if err := tx.Model(&model.Client{}).Create(&client).Error; err != nil {
		return err
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &client, &DTO.Client{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// LoginClientByPassword logs an client in
//
//	@Summary		Login
//	@Description	Log in an client using password
//	@Tags			Client
//	@Accept			json
//	@Produce		json
//	@Param			client	body	DTO.LoginClient	true	"Client"
//	@Success		200
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/client/login [post]
func LoginClientByPassword(c *fiber.Ctx) error {
	token, err := LoginByPassword(namespace.ClientKey.Name, &model.Client{}, c)
	if err != nil {
		return err
	}
	c.Response().Header.Set(namespace.HeadersKey.Auth, token)
	return nil
}

// LoginClientByEmailCode logs in a client using email and validation code
//
//	@Summary		Login client by email code
//	@Description	Login client using email and validation code
//	@Tags			client/auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body	DTO.LoginByEmailCode	true	"Login credentials"
//	@Success		200
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/client/login-with-code [post]
func LoginClientByEmailCode(c *fiber.Ctx) error {
	token, err := LoginByEmailCode(namespace.ClientKey.Name, &model.Client{}, c)
	if err != nil {
		return err
	}
	c.Response().Header.Set(namespace.HeadersKey.Auth, token)
	return nil
}

// SendClientLoginValidationCodeByEmail sends a login validation code to a client's email
//
//	@Summary		Send client login validation code by email
//	@Description	Send a login validation code to a client's email
//	@Tags			Client
//	@Param			email	path	string	true	"Client Email"
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/client/send-login-code/email/{email} [post]
func SendClientLoginValidationCodeByEmail(c *fiber.Ctx) error {
	if err := SendLoginValidationCodeByEmail(c, &model.Client{}); err != nil {
		return err
	}
	return nil
}

// GetClientByEmail retrieves an client by email
//
//	@Summary		Get client by email
//	@Description	Retrieve an client by its email
//	@Tags			Client
//	@Param			email	path	string	true	"Client Email"
//	@Produce		json
//	@Success		200	{object}	DTO.Client
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/client/email/{email} [get]
func GetClientByEmail(c *fiber.Ctx) error {
	email := c.Params("email")

	if email == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("missing 'email' at params route"))
	}

	// URL decode the email parameter to handle encoded characters like %40 (@)
	decodedEmail, err := url.QueryUnescape(email)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid email format: %w", err))
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var client model.ClientMeta
	if err := tx.Model(&model.Client{}).Where("email = ?", decodedEmail).First(&client).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.Client.NotFound
		}
		return lib.Error.General.InternalError.WithError(err)
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &client, &DTO.Client{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// GetClientById returns a client by ID
//
//	@Summary		Get client by ID
//	@Description	Get a client by ID
//	@Tags			Client
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			client_id		path	string	true	"Client ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Client
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/client/{client_id} [get]
func GetClientById(c *fiber.Ctx) error {
	id := c.Params("client_id")
	if id == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("missing 'id' at params route"))
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var client model.ClientMeta
	if err := tx.Model(&model.Client{}).Where("id = ?", id).First(&client).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.Client.NotFound
		}
		return lib.Error.General.InternalError.WithError(err)
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &client, &DTO.Client{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// GetClientAppointmentsById returns the appointments of a client with pagination and filters
//
//	@Summary		Get client appointments
//	@Description	Get the appointments of a client with pagination and optional filters
//	@Tags			Client
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			client_id		path	string	true	"Client ID"
//	@Param			page			query	int		false	"Page number (default: 1)"
//	@Param			page_size		query	int		false	"Page size (default: 10)"
//	@Param			start_date		query	string	false	"Start date in DD/MM/YYYY format"
//	@Param			end_date		query	string	false	"End date in DD/MM/YYYY format (max 90 days range)"
//	@Param			cancelled		query	string	false	"Filter by cancelled status: 'true' or 'false'"
//	@Param			timezone		query	string	true	"Timezone in IANA format (required)"
//	@Produce		json
//	@Success		200	{object}	DTO.AppointmentList
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/client/{client_id}/appointments [get]
func GetClientAppointmentsById(c *fiber.Ctx) error {
	client_id := c.Params("client_id")

	// Validate client_id is not empty and is a valid UUID
	if client_id == "" {
		return lib.Error.General.ResourceNotFoundError.WithError(fmt.Errorf("client_id parameter is required"))
	}

	// Validate UUID format
	if _, err := uuid.Parse(client_id); err != nil {
		return lib.Error.General.ResourceNotFoundError.WithError(fmt.Errorf("invalid client_id format"))
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	// Verify client exists
	var count int64
	if err := tx.Model(&model.Client{}).Where("id = ?", client_id).Count(&count).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	if count == 0 {
		return lib.Error.General.ResourceNotFoundError.WithError(fmt.Errorf("client not found"))
	}

	// Parse pagination parameters
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 10)
	offset := (page - 1) * pageSize

	// Parse required timezone parameter
	timezone := c.Query("timezone")
	if timezone == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("timezone parameter is required"))
	}

	// Build query with client filter
	query := tx.Model(&model.Appointment{}).Where("client_id = ?", client_id)

	// Parse and validate date range filters
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr != "" || endDateStr != "" {
		var startDate, endDate time.Time

		if startDateStr != "" {
			startDate, err = time.Parse("02/01/2006", startDateStr)
			if err != nil {
				return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid start_date format, expected DD/MM/YYYY"))
			}
		}

		if endDateStr != "" {
			endDate, err = time.Parse("02/01/2006", endDateStr)
			if err != nil {
				return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid end_date format, expected DD/MM/YYYY"))
			}
			// Set end date to end of day
			endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}

		// Validate date range
		if startDateStr != "" && endDateStr != "" {
			if endDate.Before(startDate) {
				return lib.Error.General.BadRequest.WithError(fmt.Errorf("end_date must be after start_date"))
			}

			daysDiff := endDate.Sub(startDate).Hours() / 24
			if daysDiff > 90 {
				return lib.Error.General.BadRequest.WithError(fmt.Errorf("date range cannot exceed 90 days"))
			}
		}

		// Apply date filters
		if startDateStr != "" {
			query = query.Where("start_time >= ?", startDate)
		}
		if endDateStr != "" {
			query = query.Where("start_time <= ?", endDate)
		}
	}

	// Parse cancelled filter
	cancelledStr := c.Query("cancelled")
	if cancelledStr != "" {
		if cancelledStr == "true" {
			query = query.Where("is_cancelled = ?", true)
		} else if cancelledStr == "false" {
			query = query.Where("is_cancelled = ?", false)
		} else {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("cancelled parameter must be 'true' or 'false'"))
		}
	}

	// Get total count for pagination
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Get appointments with pagination
	var appointments []model.Appointment
	if err := query.
		Offset(offset).
		Limit(pageSize).
		Find(&appointments).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	// Convert to DTO (AppointmentBasicInfo - excludes History and Comments)
	appointmentsDTO := make([]DTO.AppointmentBasicInfo, len(appointments))
	for i, apt := range appointments {
		paymentID := uuid.UUID{}
		if apt.PaymentID != nil {
			paymentID = *apt.PaymentID
		}
		cancelledEmployeeID := uuid.UUID{}
		if apt.CancelledEmployeeID != nil {
			cancelledEmployeeID = *apt.CancelledEmployeeID
		}

		appointmentsDTO[i] = DTO.AppointmentBasicInfo{
			ID:                    apt.ID,
			ServiceID:             apt.ServiceID,
			EmployeeID:            apt.EmployeeID,
			ClientID:              apt.ClientID,
			BranchID:              apt.BranchID,
			CompanyID:             apt.CompanyID,
			PaymentID:             paymentID,
			CancelledEmployeeID:   cancelledEmployeeID,
			StartTime:             apt.StartTime.Format(time.RFC3339),
			EndTime:               apt.EndTime.Format(time.RFC3339),
			TimeZone:              apt.TimeZone,
			Cancelled:             apt.IsCancelled,
			CancelTime:            apt.CancelTime.Format(time.RFC3339),
			IsFulfilled:           apt.IsFulfilled,
			IsCancelled:           apt.IsCancelled,
			IsCancelledByClient:   apt.IsCancelledByClient,
			IsCancelledByEmployee: apt.IsCancelledByEmployee,
			IsConfirmedByClient:   apt.IsConfirmedByClient,
		}
	}

	// Since this is client appointments, we don't need to fetch client info (it's the same client)
	// But we'll keep the structure consistent with the other endpoints
	var clientInfo []DTO.ClientBasicInfo

	// Get the client's basic info
	var client model.Client
	if err := tx.Select("id", "name", "surname", "email", "phone").
		Where("id = ?", client_id).
		First(&client).Error; err == nil {
		clientInfo = append(clientInfo, DTO.ClientBasicInfo{
			ID:      client.ID,
			Name:    client.Name,
			Surname: client.Surname,
			Email:   client.Email,
			Phone:   client.Phone,
		})
	}

	AppointmentList := DTO.AppointmentList{
		Appointments: appointmentsDTO,
		ClientInfo:   clientInfo,
		Page:         page,
		PageSize:     pageSize,
		TotalCount:   int(totalCount),
	}

	if err := lib.ResponseFactory(c).Send(200, &AppointmentList); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// UpdateClientById updates an client by ID
//
//	@Summary		Update client
//	@Description	Update an client
//	@Tags			Client
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			client_id	path		string		true	"Client ID"
//	@Param			client		body		DTO.Client	true	"Client"
//	@Success		200			{object}	DTO.Client
//	@Failure		400			{object}	DTO.ErrorResponse
//	@Router			/client/{client_id} [patch]
func UpdateClientById(c *fiber.Ctx) error {
	var client model.Client

	if err := UpdateOneById(c, &client, nil); err != nil {
		return err
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	if err := client.GetFullClient(tx); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &client, &DTO.Client{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// DeleteClientById deletes an client by ID
//
//	@Summary		Delete client
//	@Description	Delete an client
//	@Tags			Client
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			client_id		path		string	true	"Client ID"
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		404	{object}	nil
//	@Router			/client/{client_id} [delete]
func DeleteClientById(c *fiber.Ctx) error {
	return DeleteOneById(c, &model.Client{})
}

// UpdateClientImages updates the design images of an client
//
//	@Summary		Update client design images
//	@Description	Update the design images of an client
//	@Tags			Client
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			client_id		path		string	true	"Client ID"
//	@Accept			json
//	@Produce		json
//	@Param			profile	formData	file	false	"Profile image"
//	@Success		200		{object}	DTO.Client
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/client/{client_id}/design/images [patch]
func UpdateClientImages(c *fiber.Ctx) error {
	img_types_allowed := map[string]bool{"profile": true}

	var client model.Client
	Design, err := UpdateImagesById(c, client.TableName(), &client, img_types_allowed)
	if err != nil {
		return err
	}

	return lib.ResponseFactory(c).SendDTO(200, &Design.Images, &dJSON.Images{})
}

// DeleteClientImage deletes the design images of an client
//
//	@Summary		Delete client design images
//	@Description	Delete the design images of an client
//	@Tags			Client
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			client_id		path		string	true	"Client ID"
//	@Param			image_type		path		string	true	"Image Type"
//	@Produce		json
//	@Success		200	{object}	dJSON.Images
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/client/{client_id}/design/images/{image_type} [delete]
func DeleteClientImage(c *fiber.Ctx) error {
	img_types_allowed := map[string]bool{"profile": true}
	var client model.Client
	Design, err := DeleteImageById(c, client.TableName(), &client, img_types_allowed)
	if err != nil {
		return err
	}

	return lib.ResponseFactory(c).SendDTO(200, &Design.Images, &dJSON.Images{})
}

// ResetClientPasswordByEmail resets the password of a client by email
//
//	@Summary		Reset client password by email
//	@Description	Reset the password of a client by its email
//	@Tags			Client
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			email			path		string	true	"Client Email"
//	@Query			language																																	query		string	false	"Language code (default: en)"
//	@Produce		json
//	@Success		200	{object}	DTO.PasswordReseted
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/client/reset-password/{email} [post]
func ResetClientPasswordByEmail(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("missing 'email' at params route"))
	}
	if err := SendNewPasswordByEmail(c, email, &model.Client{}); err != nil {
		return err
	}
	return lib.ResponseFactory(c).Http200(nil)
}

// SendClientVerificationCodeByEmail sends a verification code to a client's email
//
//	@Summary		Send client verification code by email
//	@Description	Send a verification code to a client's email
//	@Tags			Client
//	@Param			email		path	string	true	"Client Email"
//	@Query			language	query	string	false	"Language for the email content"
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/client/send-verification-code/email/{email} [post]
func SendClientVerificationCodeByEmail(c *fiber.Ctx) error {
	return SendVerificationCodeByEmail(c, &model.Client{})
}

// VerifyClientEmail verifies a client's email
//
//	@Summary		Verify client email
//	@Description	Verify a client's email
//	@Tags			Client
//	@Param			email	path	string	true	"Client Email"
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/client/verify-email/{email}/{code} [get]
func VerifyClientEmail(c *fiber.Ctx) error {
	return VerifyEmail(c, &model.Client{})
}

func Client(Gorm *handler.Gorm) {
	endpoint := &middleware.Endpoint{DB: Gorm}
	endpoint.BulkRegisterHandler([]fiber.Handler{
		CreateClient,
		LoginClientByPassword,
		LoginClientByEmailCode,
		SendClientLoginValidationCodeByEmail,
		ResetClientPasswordByEmail,
		GetClientByEmail,
		GetClientById,
		GetClientAppointmentsById,
		UpdateClientById,
		DeleteClientById,
		UpdateClientImages,
		DeleteClientImage,
		SendClientVerificationCodeByEmail,
		VerifyClientEmail,
	})
}
