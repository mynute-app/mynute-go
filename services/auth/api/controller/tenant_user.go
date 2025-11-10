package controller

import (
	"mynute-go/services/auth/api/lib"
	"mynute-go/services/auth/config/db/model"
	DTO "mynute-go/services/auth/config/dto"

	"github.com/gofiber/fiber/v2"
)

// =====================
// TENANT USER MANAGEMENT
// =====================

// CreateTenantUser creates a tenant user
//
//	@Summary		Create tenant user
//	@Description	Create a tenant user
//	@Tags			Tenant User
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			tenant	body		DTO.CreateTenant	true	"Tenant User"
//	@Success		200			{object}	DTO.TenantBase
//	@Failure		400			{object}	DTO.ErrorResponse
//	@Router			/tenant/users [post]
func CreateTenantUser(c *fiber.Ctx) error {
	var user model.TenantUser
	if err := CreateUser(c, &user); err != nil {
		return err
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &user, &DTO.TenantUserBase{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// GetTenantUserById retrieves a tenant user by ID
//
//	@Summary		Get tenant user by ID
//	@Description	Retrieve a tenant user by its ID
//	@Tags			Tenant User
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Tenant User ID"
//	@Produce		json
//	@Success		200	{object}	DTO.TenantBase
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/tenant/users/{id} [get]
func GetTenantUserById(c *fiber.Ctx) error {
	var user model.TenantUser
	if err := GetOneBy("id", c, &user); err != nil {
		return err
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &user, &DTO.TenantUserBase{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// GetTenantUserByEmail retrieves a tenant user by email
//
//	@Summary		Get tenant user by email
//	@Description	Retrieve a tenant user by its email
//	@Tags			Tenant User
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Failure		401				{object}	nil
//	@Param			email			path		string	true	"Tenant User Email"
//	@Produce		json
//	@Success		200	{object}	DTO.TenantBase
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/tenant/users/email/{email} [get]
func GetTenantUserByEmail(c *fiber.Ctx) error {
	var user model.TenantUser
	if err := GetOneBy("email", c, &user); err != nil {
		return err
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &user, &DTO.TenantUserBase{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// UpdateTenantUserById updates a tenant user by ID
//
//	@Summary		Update tenant user
//	@Description	Update a tenant user
//	@Tags			Tenant User
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string						true	"Tenant User ID"
//	@Param			tenant	body		DTO.UpdateTenantSwagger	true	"Tenant User"
//	@Success		200			{object}	DTO.TenantBase
//	@Failure		400			{object}	DTO.ErrorResponse
//	@Router			/tenant/users/{id} [patch]
func UpdateTenantUserById(c *fiber.Ctx) error {
	var user model.TenantUser
	if err := UpdateOneById(c, &user); err != nil {
		return err
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &user, &DTO.TenantUserBase{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// DeleteTenantUserById deletes a tenant user by ID
//
//	@Summary		Delete tenant user by ID
//	@Description	Delete a tenant user by its ID
//	@Tags			Tenant User
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Tenant User ID"
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/tenant/users/{id} [delete]
func DeleteTenantUserById(c *fiber.Ctx) error {
	return DeleteOneById(c, &model.TenantUser{})
}
