package lib

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type SendResponse struct {
	Ctx *fiber.Ctx
}

// This function is used to send a response back to the client
// using the Data Transfer Object (DTO) pattern.
func (sr *SendResponse) DTO(s int, source interface{}, dto interface{}) {
	if err := ParseToDTO(source, dto); err != nil {
		sr.Http500(err)
	}
	sr.send(s, dto)
}

func (sr *SendResponse) HttpError(s int, err error) {
	sr.send(s, err.Error())
}

func (sr *SendResponse) Http400(err error) {
	sr.send(400, err.Error())
}

func (sr *SendResponse) Http404() {
	sr.sendStatus(404)
}

func (sr *SendResponse) Http500(err error) {
	log.Printf("An internal error occurred! \n Error: %v", err)
	sr.send(500, err.Error())
}

func (sr *SendResponse) Http201(data any) {
	sr.send(201, data)
}

func (sr *SendResponse) Http204() {
	sr.sendStatus(204)
}

func (sr *SendResponse) Http200(data any) {
	sr.send(200, data)
}

func (sr *SendResponse) send(s int, data any) {
	if err := sr.Ctx.Status(s).JSON(data); err != nil {
		sr.saveError(err)
	}
}

func (sr *SendResponse) sendStatus(s int) {
	if err := sr.Ctx.SendStatus(s); err != nil {
		sr.saveError(err)
	}
}

func (sr *SendResponse) saveError(err error) {
	log.Printf("An error occurred when sending the response back! \n Error: %v", err)
}
