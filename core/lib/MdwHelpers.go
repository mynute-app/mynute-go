package lib

import (
	"agenda-kaki-go/core/config/namespace"
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
)

func InterfaceDataNotFound(interfaceName string) error {
	errStr := fmt.Sprintf("%s data not found in context", interfaceName)
	return errors.New(errStr)
}

func InvalidDataType(interfaceName string) error {
	errStr := fmt.Sprintf("invalid %s data type", interfaceName)
	return errors.New(errStr)
}

// getInterface retrieves an interface from Fiber context
func GetInterface[T any](c fiber.Ctx, key namespace.ContextKey) (T, error) {
	interfaceData := c.Locals(key)
	if interfaceData == nil {
		var zero T
		return zero, InterfaceDataNotFound(string(key))
	}
	interfaceValue, ok := interfaceData.(T)
	log.Printf("model: %+v", interfaceValue)
	if !ok {
		var zero T
		return zero, InvalidDataType(string(key))
	}
	return interfaceValue, nil
}

