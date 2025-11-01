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
	return token.SignedString(getSecret())
}

func (j *jsonWebToken) CreateClaims(data any) jwt.Claims {
	JWTExpirationHours := 2160 // 90 days
	return jwt.MapClaims{
		"data": data,
		"exp":  time.Now().Add(time.Hour * time.Duration(JWTExpirationHours)).Unix(), // 90 days
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

		return getSecret(), nil
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

// WhoAreYouAdmin checks if the token belongs to an admin and returns admin claims
func (j *jsonWebToken) WhoAreYouAdmin() (*DTO.AdminClaims, error) {
	auth_token := j.C.Get(namespace.HeadersKey.Auth)
	if auth_token == "" {
		return nil, nil
	}

	parseCallback := func(token *jwt.Token) (any, error) {
		// Validate the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return getSecret(), nil
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

	// Turn claim_data into bytes
	claim_data_bytes, err := json.Marshal(claim_data)
	if err != nil {
		return nil, lib.Error.Auth.InvalidToken.WithError(err)
	}

	var adminClaim DTO.AdminClaims
	err = json.Unmarshal(claim_data_bytes, &adminClaim)
	if err != nil {
		return nil, lib.Error.Auth.InvalidToken.WithError(err)
	}

	// Verify this is an admin token
	if !adminClaim.IsAdmin || adminClaim.Type != namespace.AdminKey.Name {
		return nil, nil
	}

	return &adminClaim, nil
}

// getSecret retrieves the JWT secret from an environment variable
func getSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret != "" {
		return []byte(secret)
	}

	// For testing/development, use a deterministic fallback
	// This allows tests to work without setting JWT_SECRET
	return []byte("default-test-secret-do-not-use-in-production-12345")
}
