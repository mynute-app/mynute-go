package controllers

import "github.com/gofiber/fiber/v3"

type IController interface {
	getBy(paramKey string, c fiber.Ctx) error
	forceGetBy(paramKey string, c fiber.Ctx) error
	CreateOne(c fiber.Ctx) error
	GetAll(c fiber.Ctx) error
	GetOneById(c fiber.Ctx) error
	UpdateOneById(c fiber.Ctx) error
	DeleteOneById(c fiber.Ctx) error
	ForceDeleteOneById(c fiber.Ctx) error
	ForceGetOneById(c fiber.Ctx) error
	ForceGetAll(c fiber.Ctx) error
}

func CreateRoutes(r fiber.Router, ci IController) {
	r.Post("/", ci.CreateOne) // ok
	r.Get("/", ci.GetAll) // ok
	r.Get("/force", ci.ForceGetAll) // ok
	id := r.Group("/:id")
	id.Get("/", ci.GetOneById) // ok
	id.Patch("/", ci.UpdateOneById) // ok
	id.Delete("/", ci.DeleteOneById) // ok
	id.Delete("/force", ci.ForceDeleteOneById) // ok
	id.Get("/force", ci.ForceGetOneById) // ok
}