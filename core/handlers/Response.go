package handlers

import (
	"agenda-kaki-go/core/lib"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func Response(c *fiber.Ctx) *Res {
	return &Res{Ctx: c}
}

type Res struct {
	Ctx *fiber.Ctx
}

// This function is used to send a response back to the client
// using the Data Transfer Object (DTO) pattern.
func (sr *Res) DTO(s int, source any, dto any) {
	fmt.Printf("Source: %+v\n", source)
	fmt.Printf("DTO: %+v\n", dto)
	if err := lib.ParseToDTO(source, dto); err != nil {
		sr.Http500(err)
	}
	sr.send(s, dto)
}

func (sr *Res) Next() error {
	return sr.Ctx.Next()
}

func (sr *Res) HttpError(s int, err error) *Res {
	sr.send(s, err.Error())
	return sr
}

func (sr *Res) Http400(err error) *Res {
	sr.send(400, err.Error())
	return sr
}
func (sr *Res) Http401(err error) error {
	sr.send(401, err.Error())
	return nil
}
func (sr *Res) Http404() *Res {
	sr.sendStatus(404)
	return sr
}

func (sr *Res) Http500(err error) *Res {
	log.Printf("An internal error occurred! \n Error: %v", err)
	sr.send(500, err.Error())
	return sr
}

func (sr *Res) Http201(data any) *Res {
	sr.send(201, data)
	return sr
}

func (sr *Res) Http204() *Res {
	sr.sendStatus(204)
	return sr
}

func (sr *Res) Http200(data any) *Res {
	sr.send(200, data)
	return sr
}

func (sr *Res) send(s int, data any) {
	if data == nil {
		sr.sendStatus(s)
		return
	}
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
