package lib

import (
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

func InterfaceDataNotFound(interfaceName string) error {
	errStr := fmt.Sprintf("%s data not found in *fiber.Ctx", interfaceName)
	return errors.New(errStr)
}

func InvalidDataType(interfaceName string) error {
	errStr := fmt.Sprintf("invalid %s data type", interfaceName)
	return errors.New(errStr)
}
