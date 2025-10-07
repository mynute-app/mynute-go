package controller

import (
	"mynute-go/core/src/lib"

	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	root, err := lib.FindProjectRoot()
	if err != nil {
		return err
	}
	return c.SendFile(root + "/static/index.html")
}
