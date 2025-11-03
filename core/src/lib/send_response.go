package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

type SendResponseStruct struct {
	Ctx *fiber.Ctx
}

func ResponseFactory(c *fiber.Ctx) *SendResponseStruct {
	return &SendResponseStruct{
		Ctx: c,
	}
}

func (sr *SendResponseStruct) Next() error {
	return sr.Ctx.Next()
}

func (sr *SendResponseStruct) SendDTO(s int, source any, dto any) error {
	if source == nil || dto == nil {
		return sr.sendStatus(s)
	}
	IsPointerToStruct := func(v any) bool {
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

func (sr *SendResponseStruct) HttpError(s int, err error) error {
	if err2 := sr.send(s, err); err2 != nil {
		return nil
	}
	return nil
}

func (sr *SendResponseStruct) Http400(err error) error {
	if err2 := sr.send(400, err); err2 != nil {
		return nil
	}
	return nil
}

func (sr *SendResponseStruct) Http404() error {
	if err := sr.sendStatus(404); err != nil {
		return nil
	}
	return nil
}

func (sr *SendResponseStruct) Http401(err error) error {
	if err2 := sr.send(401, err); err2 != nil {
		return nil
	}
	return nil
}

func (sr *SendResponseStruct) Http500(err error) error {
	log.Printf("An internal error occurred! \n Error: %v", err)
	if err2 := sr.send(500, err); err2 != nil {
		return nil
	}
	return nil
}

func (sr *SendResponseStruct) Http201(data any) error {
	if err := sr.send(201, data); err != nil {
		return err
	}
	return nil
}

func (sr *SendResponseStruct) Http204() error {
	if err := sr.sendStatus(204); err != nil {
		return err
	}
	return nil
}

func (sr *SendResponseStruct) Http200(data any) error {
	if err := sr.send(200, data); err != nil {
		return err
	}
	return nil
}

func (sr *SendResponseStruct) Send(s int, data any) error {
	if err := sr.send(s, data); err != nil {
		return err
	}
	return nil
}

func (sr *SendResponseStruct) send(s int, data any) error {
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

func (sr *SendResponseStruct) sendError(s int, err error) error {
	sr.saveError(err)

	// Se for ErrorStruct, retorna direto
	if errStruct, ok := err.(ErrorStruct); ok {
		return sr.Ctx.Status(s).JSON(errStruct)
	}

	// Fallback para erro genérico
	return sr.Ctx.Status(s).JSON(map[string]any{
		"error": err.Error(),
	})
}

func (sr *SendResponseStruct) sendStatus(s int) error {
	if err := sr.Ctx.SendStatus(s); err != nil {
		sr.saveError(err)
		fmt.Printf("Failed to send status: %d\n", s)
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}
	return nil
}

func (sr *SendResponseStruct) saveError(err error) *SendResponseStruct {
	fmt.Printf("An error occurred!\n>>> %v\n", err.Error())
	return sr
}

