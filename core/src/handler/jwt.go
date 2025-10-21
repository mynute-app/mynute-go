package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	DTO "mynute-go/core/src/config/api/dto"
	"mynute-go/core/src/config/namespace"
	"mynute-go/core/src/lib"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type jsonWebToken struct {
	C   *fiber.Ctx
	Res *lib.SendResponseStruct
}

func JWT(c *fiber.Ctx) *jsonWebToken {
	return &jsonWebToken{C: c, Res: &lib.SendResponseStruct{Ctx: c}}
}

func (j *jsonWebToken) GetToken() string {
	return j.C.Get(namespace.HeadersKey.Auth)
}

func (j *jsonWebToken) Encode(data any) (string, error) {
	claims := j.CreateClaims(data)
	token, err := j.CreateToken(claims)
	if err != nil {
		return "", err
	}
	return token, nil
}

// create token
func (j *jsonWebToken) CreateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySecret)
}

func (j *jsonWebToken) CreateClaims(data any) jwt.Claims {
	return jwt.MapClaims{
		"data": data,
		"exp":  time.Now().Add(time.Hour * 2160).Unix(), // 90 days
	}
}

func (j *jsonWebToken) WhoAreYou() (*DTO.Claims, error) {
	auth_token := j.C.Get(namespace.HeadersKey.Auth)
	if auth_token == "" {
		return nil, nil
	}

	parseCallback := func(token *jwt.Token) (any, error) {
		// Validate the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return mySecret, nil
	}

	token, err := jwt.Parse(auth_token, parseCallback)
	if err != nil {
		return nil, lib.Error.Auth.InvalidToken.WithError(err)
	} else if token == nil {
		return nil, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, lib.Error.Auth.InvalidToken.WithError(errors.New("invalid jwt.MapClaims passed"))
	}

	claim_data, ok := claims["data"].(map[string]any)
	if !ok {
		return nil, lib.Error.Auth.InvalidToken.WithError(errors.New("invalid claim.data passed"))
	}

	// Parse claim_data into model.Client{} struct

	// Turn claim_data into bytes
	claim_data_bytes, err := json.Marshal(claim_data)
	if err != nil {
		return nil, lib.Error.Auth.InvalidToken.WithError(err)
	}

	var client DTO.Claims
	err = json.Unmarshal(claim_data_bytes, &client)
	if err != nil {
		return nil, lib.Error.Auth.InvalidToken.WithError(err)
	}

	return &client, nil
}

// getSecret retrieves the JWT secret from an environment variable
func getSecret() []byte {
	generateMySecret := func() []byte {
		s := fmt.Sprintf("my_secret_is_%d!", lib.GenerateRandomInt(16))
		return []byte(s)
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return generateMySecret() // Only for testing or development
	}
	return []byte(secret)
}

var mySecret = getSecret()
