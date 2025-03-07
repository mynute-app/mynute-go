package handler

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type jsonWebToken struct {
	C   *fiber.Ctx
	Res *lib.SendResponse
}

func JWT(c *fiber.Ctx) *jsonWebToken {
	return &jsonWebToken{C: c, Res: &lib.SendResponse{Ctx: c}}
}

func (j *jsonWebToken) GetToken() string {
	return j.C.Get("Authorization")
}

// create token
func (j *jsonWebToken) CreateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySecret)
}

func (j *jsonWebToken) CreateClaims(data any) jwt.Claims {
	return jwt.MapClaims{
		"data": data,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	}
}

// WhoAreYou decrypts and validates the JWT token, saving user data in context if valid
func (j *jsonWebToken) WhoAreYou() error {
	saveUserData := func(value any) {
		j.C.Locals(namespace.GeneralKey.UserData, value)
	}

	// Retrieve the token from the Authorization header
	tokenString := j.GetToken()
	if tokenString == "" {
		saveUserData(nil)
		return errors.New("missing auth token")
	}

	keyFunc := func(token *jwt.Token) (any, error) {
		// Validate the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return mySecret, nil
	}

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, keyFunc)

	if err != nil {
		return err
	} else if token == nil {
		return errors.New("invalid token")
	}

	// Check token validity and extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return errors.New("invalid token")
	}

	// Store claims (user data) in Fiber's Locals
	saveUserData(claims)
	return nil
}

// getSecret retrieves the JWT secret from an environment variable
func getSecret() []byte {
	generateMySecret := func() []byte {
		s := fmt.Sprintf("my_secret_is_%d!", lib.GenerateRandomIntOfExactly(16))
		return []byte(s)
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return generateMySecret() // Only for testing or development
	}
	return []byte(secret)
}

var mySecret = getSecret()
