package lib

import (
	"log"

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
func (sr *SendResponse) DTO(s int, source any, dto any) *SendResponse {
	if err := ParseToDTO(source, dto); err != nil {
		sr.Http500(err)
	}
	sr.send(s, dto)
	return sr
}

func (sr *SendResponse) HttpError(s int, err error) *SendResponse {
	sr.send(s, err.Error())
	return sr
}

func (sr *SendResponse) Http400(err error) *SendResponse {
	sr.send(400, err.Error())
	return sr
}

func (sr *SendResponse) Http404() *SendResponse {
	sr.sendStatus(404)
	return sr
}

func (sr *SendResponse) Http401(err error) *SendResponse {
	sr.send(401, err.Error())
	return sr
}

func (sr *SendResponse) Http500(err error) *SendResponse {
	log.Printf("An internal error occurred! \n Error: %v", err)
	sr.send(500, err.Error())
	return sr
}

func (sr *SendResponse) Http201(data any) *SendResponse {
	sr.send(201, data)
	return sr
}

func (sr *SendResponse) Http204() *SendResponse {
	sr.sendStatus(204)
	return sr
}

func (sr *SendResponse) Http200(data any) *SendResponse {
	sr.send(200, data)
	return sr
}

func (sr *SendResponse) send(s int, data any) *SendResponse {
	if data == nil {
		sr.sendStatus(s)
		return sr
	}
	if err := sr.Ctx.Status(s).JSON(data); err != nil {
		sr.saveError(err)
	}
	return sr
}

func (sr *SendResponse) sendStatus(s int) *SendResponse {
	if err := sr.Ctx.SendStatus(s); err != nil {
		sr.saveError(err)
	}
	return sr
}

func (sr *SendResponse) saveError(err error) *SendResponse {
	log.Printf("An error occurred when sending the response back! \n Error: %v", err)
	return sr
}
