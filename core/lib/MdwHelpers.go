package lib

import (
	"agenda-kaki-go/core/config/namespace"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/gofiber/fiber/v3"
)

func GetFromCtx[T any](c fiber.Ctx, key namespace.ContextKey) (T, error) {
	interfaceData := c.Locals(key)
	if interfaceData == nil {
			var zero T
			return zero, InterfaceDataNotFound(string(key))
	}

	// Check if the retrieved value is assignable to the target type T
	interfaceValue := reflect.ValueOf(interfaceData)

	// Log the retrieved data for debugging
	log.Printf("interfaceData: %+v", interfaceData)
	log.Printf("interfaceValue: %+v", interfaceValue.Interface())

	// Create a zero value for type T
	var zero T

	// Check if the type of interfaceData matches T
	if reflect.TypeOf(interfaceData).Kind() == reflect.TypeOf(zero).Kind() {
			// Safe cast to T
			return interfaceData.(T), nil
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
