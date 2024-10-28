package handlers

import (
	"agenda-kaki-go/core/lib"
	"log"

	"github.com/gofiber/fiber/v3"
)

func Response(c fiber.Ctx) *Res {
	return &Res{Ctx: c}
}

type Res struct {
	Ctx fiber.Ctx
}

// This function is used to send a response back to the client
// using the Data Transfer Object (DTO) pattern.
func (sr *Res) DTO(s int, source interface{}, dto interface{}) {
	if err := lib.ParseToDTO(source, dto); err != nil {
		sr.Http500(err)
	}
	sr.send(s, dto)
}

func (sr *Res) HttpError(s int, err error) {
	sr.send(s, err.Error())
}

func (sr *Res) Http400(err error) {
	sr.send(400, err.Error())
}

func (sr *Res) Http404() {
	sr.sendStatus(404)
}

func (sr *Res) Http500(err error) {
	log.Printf("An internal error occurred! \n Error: %v", err)
	sr.send(500, err.Error())
}

func (sr *Res) Http201(data any) {
	sr.send(201, data)
}

func (sr *Res) Http204() {
	sr.sendStatus(204)
}

func (sr *Res) Http200(data any) {
	sr.send(200, data)
}

func (sr *Res) send(s int, data any) {
	if err := sr.Ctx.Status(s).JSON(data); err != nil {
		sr.saveError(err)
	}
}

func (sr *Res) sendStatus(s int) {
	if err := sr.Ctx.SendStatus(s); err != nil {
		sr.saveError(err)
	}
}

func (sr *Res) saveError(err error) {
	log.Printf("An error occurred when sending the response back! \n Error: %v", err)
}
