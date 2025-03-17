package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

type SendResponse struct {
	Ctx *fiber.Ctx
}

func (sr *SendResponse) Next() error {
	return sr.Ctx.Next()
}

// This function is used to send a response back to the client
// using the Data Transfer Object (DTO) pattern.
// func (sr *SendResponse) DTO(s int, source any, dto any) *SendResponse {
// 	if source == nil || dto == nil {
// 		sr.sendStatus(s)
// 		return sr
// 	}
// 	if err := ParseToDTO(source, dto); err != nil {
// 		sr.Http500(err)
// 	}
// 	sr.send(s, dto)
// 	return sr
// }

func (sr *SendResponse) SendDTO(s int, source any, dto any) error {
	if source == nil || dto == nil {
		return sr.sendStatus(s)
	}
	IsPointerToStruct := func (v any) bool {
		val := reflect.ValueOf(v)
		return val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct
	}
	// Return error if source or dto are not pointer to struct
	if !IsPointerToStruct(source) || !IsPointerToStruct(dto) {
		return sr.Http500(errors.New("source and dto must be pointer to struct"))
	}
	source_bytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(source_bytes, dto); err != nil {
		return err
	}
	return sr.send(s, dto)
}

func (sr *SendResponse) HttpError(s int, err error) error {
	sr.send(s, err)
	return nil
}

func (sr *SendResponse) Http400(err error) error {
	sr.send(400, err)
	return nil
}

func (sr *SendResponse) Http404() error {
	sr.sendStatus(404)
	return nil
}

func (sr *SendResponse) Http401(err error) error {
	sr.send(401, err)
	return nil
}

func (sr *SendResponse) Http500(err error) error {
	log.Printf("An internal error occurred! \n Error: %v", err)
	sr.send(500, err)
	return nil
}

func (sr *SendResponse) Http201(data any) error {
	sr.send(201, data)
	return nil
}

func (sr *SendResponse) Http204() error {
	sr.sendStatus(204)
	return nil
}

func (sr *SendResponse) Http200(data any) error {
	sr.send(200, data)
	return nil
}

func (sr *SendResponse) Send(s int, data any) error {
	sr.send(s, data)
	return nil
}

func (sr *SendResponse) send(s int, data any) error {
	if data == nil {
		return sr.sendStatus(s)
	}
	// Check if data is of error
	if error_passed, ok := data.(error); ok {
		return sr.sendError(s, error_passed)
	}
	if err := sr.Ctx.Status(s).JSON(data); err != nil {
		sr.saveError(err)
	}
	return nil
}

func (sr *SendResponse) sendError(s int, err error) error {
	sr.saveError(err)
	if resErr := sr.Ctx.Status(s).JSON(map[string]any{"Error": err.Error()}); resErr != nil {
		sr.saveError(resErr)
		return resErr
	}
	return nil
}

func (sr *SendResponse) sendStatus(s int) error {
	if err := sr.Ctx.SendStatus(s); err != nil {
		sr.saveError(err)
		fmt.Printf("Failed to send status: %d\n", s)
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}
	return nil
}

func (sr *SendResponse) saveError(err error) *SendResponse {
	fmt.Printf("An error occurred!\n>>> %v\n", err.Error())
	return sr
}
