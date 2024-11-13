package handlers

import (
	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
)

type Authentication struct {
	C       fiber.Ctx
	Res     *Res
	Request *Req
}

func HashPassword(password string) (string, error) {
	// Generate a hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword, password string) bool {
	// Compare the password with the hashed one
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}
