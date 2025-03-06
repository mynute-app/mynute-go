package lib

import (
	"agenda-kaki-go/core/config/namespace"
	"errors"
	"fmt"

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
	var body Body
	err := BodyParser(c.Body(), &body)
	if err != nil {
		return err
	}
	c.Locals(namespace.RequestKey.Body_Parsed, &body)
	return c.Next()
}

func InterfaceDataNotFound(interfaceName string) error {
	errStr := fmt.Sprintf("%s data not found in *fiber.Ctx", interfaceName)
	return errors.New(errStr)
}

func InvalidDataType(interfaceName string) error {
	errStr := fmt.Sprintf("invalid %s data type", interfaceName)
	return errors.New(errStr)
}
