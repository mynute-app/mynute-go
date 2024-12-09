package handlers

import (
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/tests/lib"
	"errors"
	"fmt"
	"os"
	"time"

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

// create token
func (j *jsonWebToken) CreateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySecret)
}

func (j *jsonWebToken) CreateClaims(user string) jwt.Claims {
	return jwt.MapClaims{
		"userId": user,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	}
}

// WhoAreYou decrypts and validates the JWT token, saving user data in context if valid
func (j *jsonWebToken) WhoAreYou() error {
	saveUserData := func(value interface{}) {
		j.C.Locals(namespace.GeneralKey.UserData, value)
	}

	// Retrieve the token from the Authorization header
	tokenString := j.GetToken()
	if tokenString == "" {
		saveUserData(nil)
		
		return j.Res.Http401(errors.New("missing token"))
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
		return j.Res.Http401(err)
	}

	// Check token validity and extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return j.Res.Http401(errors.New("invalid token"))
	}

	// Store claims (user data) in Fiber's Locals
	saveUserData(claims)
	return j.C.Next()
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
