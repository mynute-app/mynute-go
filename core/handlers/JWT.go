package handlers

import (
	"agenda-kaki-go/tests/lib"
	"errors"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

type jsonWebToken struct {
	C   fiber.Ctx
	Res *Res
}

func JWT(c fiber.Ctx) *jsonWebToken {
	return &jsonWebToken{C: c, Res: Response(c)}
}

func (j *jsonWebToken) GetToken() string {
	return j.C.Get("Authorization")
}

// WhoAreYou decrypts and validates the JWT token, saving user data in context if valid
func (j *jsonWebToken) WhoAreYou() error {
	// Retrieve the token from the Authorization header
	tokenString := j.GetToken()
	if tokenString == "" {
		return j.C.Next()
	}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return mySecret, nil
	}

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, keyFunc)

	if err != nil {
		return j.Res.Http400(err).Next()
	}

	// Check token validity and extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Store claims (user data) in Fiber's Locals
		j.C.Locals("userData", claims)
		return nil
	}

	return j.Res.Http400(errors.New("invalid token")).Next()
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
