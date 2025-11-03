package routes

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
)

func Prometheus(App *fiber.App) {
	prometheus := fiberprometheus.New("fiber_app")
	prometheus.RegisterAt(App, "/metrics")
	App.Use(prometheus.Middleware)
}
