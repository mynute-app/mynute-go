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

	// If the initial type assertion succeeds, return the value
	if ok {
		return interfaceValue, nil
	}

	return zero, InvalidDataType(string(key))

	// The code below does not seem to work as expected.
	// It is supposed to handle slices, but it never gets to that point.
	// And I found out that if you pass the type to the fiber context as a pointer to the slice
	// and then pass the generic type T as an interface for this function, it just works.

	// Moral of the story, if you want this function to handle slices, pass the slice
	// as a pointer to the context and then pass the generic type T as an interface{} for this function.
	// But remember not to pass the type T as a []interface{} or *[]interface{}. 
	// Just pass the generic type T as an interface{} and it should work as expected.

	// By now I will coment out the rest below.

	// If the type assertion fails, check if T is a slice type
	// var zero T
	// zeroType := reflect.TypeOf(zero)

	// // If not a slice, then return an invalid data type error
	// if zeroType.Kind() != reflect.Slice {
	// 	return zero, InvalidDataType(string(key))
	// }

	// // Get the underlying type of the slice element
	// elemType := zeroType.Elem()

	// // Dynamically create a slice of the element type using reflection
	// slicePtr := reflect.New(reflect.MakeSlice(reflect.SliceOf(elemType), 0, 0).Type()).Interface()

	// // Validate the type assertion before returning to avoid runtime panic
	// convertedSlice, ok := slicePtr.(T)
	// if !ok {
	// 	return zero, InvalidDataType(string(key))
	// }

	// // Return the valid slice
	// return convertedSlice, nil
}

// func GetFromCtx[T any](c fiber.Ctx, key namespace.ContextKey) (T, error) {
// 	interfaceData := c.Locals(key)
// 	if interfaceData == nil {
// 		var zero T
// 		return zero, InterfaceDataNotFound(string(key))
// 	}

// 	// Type assertion for non-slice types
// 	interfaceValue, ok := interfaceData.(T)
// 	log.Printf("interfaceValue: %+v", interfaceValue)
// 	log.Printf("interfaceData: %+v", interfaceData)
// 	log.Printf("ok: %+v", ok)

// 	// If type assertion works, return the value
// 	if ok {
// 		return interfaceValue, nil
// 	}

// 	// If it fails, check if we are dealing with slices and convert
// 	interfaceDataVal := reflect.ValueOf(interfaceData)
// 	zeroType := reflect.TypeOf(*new(T))

// 	log.Printf("zeroType: %+v", zeroType)
// 	log.Printf("interfaceDataVal: %+v", interfaceDataVal)

// 	// Check if the target type is a slice
// 	if zeroType.Kind() == reflect.Slice && interfaceDataVal.Kind() == reflect.Slice {
// 		// Create a new slice of type T
// 		slice := reflect.MakeSlice(reflect.TypeOf(*new(T)), interfaceDataVal.Len(), interfaceDataVal.Cap())

// 		// Copy elements from interfaceData to the new slice
// 		for i := 0; i < interfaceDataVal.Len(); i++ {
// 			slice.Index(i).Set(interfaceDataVal.Index(i))
// 		}

// 		// Convert back to T
// 		return slice.Interface().(T), nil
// 	}

// 	// If all else fails, return an invalid data type error
// 	var zero T
// 	return zero, InvalidDataType(string(key))
// }

// func GetFromCtx[T any](c fiber.Ctx, key namespace.ContextKey) (T, error) {
// 	interfaceData := c.Locals(key)
// 	if interfaceData == nil {
// 			var zero T
// 			return zero, InterfaceDataNotFound(string(key))
// 	}

// 	// Check if the retrieved value is assignable to the target type T
// 	interfaceValue := reflect.ValueOf(interfaceData)

// 	// Log the retrieved data for debugging
// 	log.Printf("interfaceData: %+v", interfaceData)
// 	log.Printf("interfaceValue: %+v", interfaceValue.Interface())

// 	// Create a zero value for type T
// 	var zero T

// 	// Check if the type of interfaceData matches T
// 	if reflect.TypeOf(interfaceData).Kind() == reflect.TypeOf(zero).Kind() {
// 			// Safe cast to T
// 			return interfaceData.(T), nil
// 	}

// 	return zero, InvalidDataType(string(key))
// }

// GetFromCtx retrieves an interface from Fiber context
// func GetFromCtx[T any](c fiber.Ctx, key namespace.ContextKey) (T, error) {
// 	interfaceData := c.Locals(key)
// 	if interfaceData == nil {
// 		var zero T
// 		return zero, InterfaceDataNotFound(string(key))
// 	}
// 	interfaceValue, ok := interfaceData.(T)
// 	log.Printf("interfaceValue: %+v", interfaceValue)
// 	log.Printf("interfaceData: %+v", interfaceData)
// 	log.Printf("ok: %+v", ok)
// 	if !ok {
// 		var zero T
// 		return zero, InvalidDataType(string(key))
// 	}
// 	return interfaceValue, nil
// }

func InterfaceDataNotFound(interfaceName string) error {
	errStr := fmt.Sprintf("%s data not found in fiber.Ctx", interfaceName)
	return errors.New(errStr)
}

func InvalidDataType(interfaceName string) error {
	errStr := fmt.Sprintf("invalid %s data type", interfaceName)
	return errors.New(errStr)
}
