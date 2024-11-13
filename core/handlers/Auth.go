package handlers

import "github.com/gofiber/fiber/v3"

type Authentication struct {
	C   fiber.Ctx
	Res *Res
	Request *Req
}

func Auth(c fiber.Ctx) *Authentication {
	return &Authentication{C: c, Res: Response(c)}
}
