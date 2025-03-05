package controllers

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"
	"agenda-kaki-go/core/service"
	"log"

	"github.com/gofiber/fiber/v2"
)

// EmployeeController embeds service.Base in order to extend it with the functions below
type userController struct {
	service.Base[models.User, DTO.User]
}

// CreateUser creates an user
//
// @Summary		Create user
// @Description	Create an user
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			user	body		DTO.CreateUser	true	"User"
// @Success		200		{object}	DTO.User
// @Failure		400		{object}	DTO.ErrorResponse
// @Router			/user [post]
func (cc *userController) CreateUser(c *fiber.Ctx) error {
	return cc.CreateOne(c)
}

// GetOneByEmail retrieves an user by email
//
// @Summary		Get user by email
// @Description	Retrieve an user by its email
// @Tags			User
// @Param			email	path	string	true	"User Email"
// @Produce		json
// @Success		200	{object}	DTO.User
// @Failure		404	{object}	DTO.ErrorResponse
// @Router			/user/email/{email} [get]
func (cc *userController) GetOneByEmail(c *fiber.Ctx) error {
	return cc.GetBy("email", c)
}

// UpdateUserById updates an user by ID
//
// @Summary		Update user
// @Description	Update an user
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			id		path		string		true	"User ID"
// @Param			user	body		DTO.User	true	"User"
// @Success		200		{object}	DTO.User
// @Failure		400		{object}	DTO.ErrorResponse
// @Router			/user/{id} [patch]
func (cc *userController) UpdateUserById(c *fiber.Ctx) error {
	return cc.UpdateOneById(c)
}

// DeleteUserById deletes an user by ID
//
// @Summary		Delete user
// @Description	Delete an user
// @Tags			User
// @Param			id	path	string	true	"User ID"
// @Produce		json
// @Success		200	{object}	DTO.User
// @Failure		404	{object}	DTO.ErrorResponse
// @Router			/user/{id} [delete]
func (cc *userController) DeleteUserById(c *fiber.Ctx) error {
	return cc.DeleteOneById(c)
}

// Login just logs an user in case the password is correct
//
// @Summary		Login
// @Description	Log in an user
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			user	body		DTO.User	true	"User"
// @Success		200		{object}	DTO.User
// @Failure		404		{object}	DTO.ErrorResponse
// @Failure		401		{object}	DTO.ErrorResponse
// @Router			/user/login [post]
func (cc *userController) Login(c *fiber.Ctx) error {
	cc.SetAction(c)
	body := c.Locals(namespace.GeneralKey.Model).(*models.User)
	var userDatabase models.User
	if err := cc.Request.Gorm.GetOneBy("email", body.Email, &userDatabase, []string{}); err != nil {
		cc.AutoReqActions.ActionFailed(404, err)
		return nil
	}
	if handlers.ComparePassword(userDatabase.Password, body.Password) && userDatabase.Verified {
		cc.AutoReqActions.Status = 401
		return nil
	}
	claims := handlers.JWT(c).CreateClaims(userDatabase.Email)
	token, err := handlers.JWT(c).CreateToken(claims)
	if err != nil {
		cc.AutoReqActions.ActionFailed(500, err)
	}
	log.Println("User logged in")
	c.Response().Header.Set("Authorization", token)

	return nil
}

func User(Gorm *handlers.Gorm) *userController {
	return &userController{
		Base: service.Base[models.User, DTO.User]{
			Name:         namespace.UserKey.Name,
			Request:      handlers.Request(Gorm),
			Middleware:   middleware.User(Gorm),
			Associations: []string{"Branches", "Services", "Appointments", "Company"},
		},
	}
}
