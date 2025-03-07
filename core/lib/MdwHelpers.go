package lib

import (
	"agenda-kaki-go/core/config/namespace"
	"errors"
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

// GetFromCtx retrieves an interface from Fiber context
func GetFromCtx[T any](c *fiber.Ctx, key string) (T, error) {
	interfaceData := c.Locals(key)
	var zero T

	if interfaceData == nil {
		return zero, InterfaceDataNotFound(key)
	}

	interfaceValue, ok := interfaceData.(T)

	if ok {
		return interfaceValue, nil
	}

	return zero, InvalidDataType(key)
}

// SaveBody is a fiber.Handler middleware that parses
// the request body and saves it to the Fiber context
func SaveBodyOnCtx[Body any](c *fiber.Ctx) error {
	// Check if body is a pointer. It can not be.
	if reflect.TypeOf((*Body)(nil)).Elem().Kind() == reflect.Ptr {
		return errors.New("at SaveBodyOnCtx function: body type cannot be a pointer")
	}

	var body Body
	err := BodyParser(c.Body(), &body)
	if err != nil {
		return err
	}
	c.Locals(namespace.RequestKey.Body_Parsed, &body)
	return c.Next()
}

func GetBodyFromCtx[Body any](c *fiber.Ctx) (Body, error) {
	return GetFromCtx[Body](c, namespace.RequestKey.Body_Parsed)
}

func GetClaimsFromCtx(c *fiber.Ctx) (map[string]interface{}, error) {
	return GetFromCtx[map[string]interface{}](c, namespace.RequestKey.Auth_Claims)
}

func InterfaceDataNotFound(interfaceName string) error {
	errStr := fmt.Sprintf("%s data not found in *fiber.Ctx", interfaceName)
	return errors.New(errStr)
}

func InvalidDataType(interfaceName string) error {
	errStr := fmt.Sprintf("invalid %s data type", interfaceName)
	return errors.New(errStr)
}

func MatchUserTokenWithCompanyID(c *fiber.Ctx) error {
	// Check if company_id parameter exists in request body
	body, err := GetBodyFromCtx[map[string]any](c)
	if err != nil {
		return err
	}
	companyID, ok := body["company_id"]
	if !ok {
		return MyErrors.CompanyIDNotFound.SendToClient(c)
	}
	claims, err := GetClaimsFromCtx(c)
	if err != nil {
		return err
	}
	userCompanyID, ok := claims["company_id"]
	if !ok {
		return MyErrors.CompanyIDNotFound.SendToClient(c)
	}
	if companyID != userCompanyID {
		return MyErrors.Unauthroized.SendToClient(c)
	}
	return c.Next()
}
