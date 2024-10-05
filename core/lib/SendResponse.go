package lib

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

type SendResponse struct {
	c fiber.Ctx
}

func (sr *SendResponse) DTO(source interface{}, dto interface{}) error {
	if err := ParseToDTO(source, dto); err != nil {
		return sr.Http500(err)
	}

	sr.Http200(dto)

	return nil
}

func (sr *SendResponse) HttpError(s int, err error) error {
	return sr.send(s, err.Error())
}

func (sr *SendResponse) Http400(err error) error {
	return sr.send(400, err.Error())
}

func (sr *SendResponse) Http404() error {
	return sr.send(404, nil)
}

func (sr *SendResponse) Http500(err error) error {
	log.Printf("An internal error occurred! \n Error: %v", err)
	return sr.send(500, err.Error())
}

func (sr *SendResponse) Http201() error {
	return sr.send(201, nil)
}

func (sr *SendResponse) Http204() error {
	return sr.send(204, nil)
}

func (sr *SendResponse) Http200(data any) error {
	return sr.send(200, data)
}

func (sr *SendResponse) send(s int, data any) error {
	if data != nil {
		if err := sr.c.Status(s).JSON(data); err != nil {
			sr.saveError(err)
		}
		return nil
	}
	if err := sr.c.SendStatus(s); err != nil {
		sr.saveError(err)
		return err
	}
	return nil
}

func (sr *SendResponse) saveError(err error) {
	log.Printf("An error occurred when sending the response back! \n Error: %v", err)
}
