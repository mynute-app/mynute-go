package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/service"

	"github.com/gofiber/fiber/v2"
)

// EmployeeController embeds service.Base in order to extend it with the functions below
type user_controller struct {
	service.Base[model.User, DTO.User]
}

// CreateUser creates an user
//
//	@Summary		Create user
//	@Description	Create an user
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			user	body		DTO.CreateUser	true	"User"
//	@Success		200		{object}	DTO.CreatedUser
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/user [post]
func (cc *user_controller) CreateUser(c *fiber.Ctx) error {
	user := &model.User{}
	if err := c.BodyParser(user); err != nil {
		return err
	}
	if err := cc.Request.Gorm.DB.Create(user).Error; err != nil {
		return err
	}
	res := &lib.SendResponse{Ctx: c}
	res.DTO(200, user, &DTO.User{})
	return nil
}

// LoginUser logs an user in
//
//	@Summary		Login
//	@Description	Log in an user
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			user	body	DTO.LoginUser	true	"User"
//	@Success		200
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/user/login [post]
func (cc *user_controller) LoginUser(c *fiber.Ctx) error {
	var body model.User
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	var user model.User
	if err := cc.Request.Gorm.GetOneBy("email", body.Email, &user, []string{}); err != nil {
		return err
	}
	if !handler.ComparePassword(user.Password, body.Password) {
		return lib.Error.Auth.InvalidLogin.SendToClient(c)
	}
	token, err := handler.JWT(c).Encode(user)
	if err != nil {
		return err
	}
	c.Response().Header.Set("Authorization", token)
	return nil
}

// VerifyUserEmail Does the email verification for an user
//
//	@Summary		Verify email
//	@Description	Verify an user's email
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			email	path	string	true	"User Email"
//	@Param			code	path	string	true	"Verification Code"
//	@Success		200		{object}	nil
//	@Failure		404		{object}	nil
//	@Router			/user/verify-email/{email}/{code} [post]
func (cc *user_controller) VerifyUserEmail(c *fiber.Ctx) error {
	res := &lib.SendResponse{Ctx: c}
	email := c.Params("email")
	var user model.User
	user.Email = email
	if err := user.ValidateEmail(); err != nil {
		return res.Send(400, err)
	}
	if err := cc.Request.Gorm.GetOneBy("email", email, &user, []string{}); err != nil {
		return err
	}
	// code := c.Params("code")
	// }
	// if user.VerificationCode != code {
	// 	return lib.Error.Auth.EmailCodeInvalid.SendToClient(c)
	// }
	user.Verified = true
	if err := cc.Request.Gorm.DB.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

// GetUserByEmail retrieves an user by email
//
//	@Summary		Get user by email
//	@Description	Retrieve an user by its email
//	@Tags			User
//	@Param			email	path	string	true	"User Email"
//	@Produce		json
//	@Success		200	{object}	DTO.User
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/user/email/{email} [get]
func (cc *user_controller) GetUserByEmail(c *fiber.Ctx) error {
	return cc.GetBy("email", c)
}

// UpdateUserById updates an user by ID
//
//	@Summary		Update user
//	@Description	Update an user
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string		true	"User ID"
//	@Param			user	body		DTO.User	true	"User"
//	@Success		200		{object}	DTO.User
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/user/{id} [patch]
func (cc *user_controller) UpdateUserById(c *fiber.Ctx) error {
	return cc.UpdateOneById(c)
}

// DeleteUserById deletes an user by ID
//
//	@Summary		Delete user
//	@Description	Delete an user
//	@Tags			User
//	@Param			id	path	string	true	"User ID"
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		404	{object}	nil
//	@Router			/user/{id} [delete]
func (cc *user_controller) DeleteUserById(c *fiber.Ctx) error {
	return cc.DeleteOneById(c)
}

// func (cc *user_controller) Login(c *fiber.Ctx) error {
// 	cc.SetAction(c)
// 	body := c.Locals(namespace.GeneralKey.Model).(*model.User)
// 	var userDatabase model.User
// 	if err := cc.Request.Gorm.GetOneBy("email", body.Email, &userDatabase, []string{}); err != nil {
// 		cc.AutoReqActions.ActionFailed(404, err)
// 		return nil
// 	}
// 	if handler.ComparePassword(userDatabase.Password, body.Password) && userDatabase.Verified {
// 		cc.AutoReqActions.Status = 401
// 		return nil
// 	}
// 	claims := handler.JWT(c).CreateClaims(userDatabase.Email)
// 	token, err := handler.JWT(c).CreateToken(claims)
// 	if err != nil {
// 		cc.AutoReqActions.ActionFailed(500, err)
// 	}
// 	log.Println("User logged in")
// 	c.Response().Header.Set("Authorization", token)

// 	return nil
// }

var UserAssociations = []string{"Appointments"}

func User(Gorm *handler.Gorm) *user_controller {
	return &user_controller{
		Base: service.Base[model.User, DTO.User]{
			Name:         namespace.UserKey.Name,
			Request:      handler.Request(Gorm),
			Associations: UserAssociations,
		},
	}
}
