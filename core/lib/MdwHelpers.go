package lib

import (
	"agenda-kaki-go/core/config/namespace"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"
)

// GetFromCtx retrieves an interface from Fiber context
func GetFromCtx[T any](c fiber.Ctx, key namespace.ContextKey) (T, error) {
	interfaceData := c.Locals(key)
	var zero T

	if interfaceData == nil {
		return zero, InterfaceDataNotFound(string(key))
	}

	interfaceValue, ok := interfaceData.(T)

	if ok {
		return interfaceValue, nil
	}

	return zero, InvalidDataType(string(key))
}

func InterfaceDataNotFound(interfaceName string) error {
	errStr := fmt.Sprintf("%s data not found in fiber.Ctx", interfaceName)
	return errors.New(errStr)
}

func InvalidDataType(interfaceName string) error {
	errStr := fmt.Sprintf("invalid %s data type", interfaceName)
	return errors.New(errStr)
}
