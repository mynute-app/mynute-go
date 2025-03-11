package lib

import (
	"fmt"
	"strings"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupMetrics(app *fiber.App) {
	prometheusMiddleware := fiberprometheus.New("fiber-api") // App name for Prometheus
	app.Use(prometheusMiddleware.Middleware)                // Enable Prometheus middleware

	// Expose Prometheus metrics using promhttp
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
}


func ApiLog(c *fiber.Ctx) error {
	start := time.Now() // Start timer

	fmt.Println(">>> Request Received!")
	printRequest(c)
	fmt.Println(">>> Will process request and send answer soon...")
	fmt.Println("-------------------")

	// Call the next middleware/handler
	err := c.Next()

	// Calculate duration
	duration := time.Since(start).Milliseconds()

	// Read response body
	body := c.Response().Body()
	responseHeaders := filterHeaders(c.Response().Header.String())

	fmt.Println(">>> Response Sent!")
	fmt.Println("Response: ", c.Response().StatusCode())
	fmt.Println("Response Time:", duration, "ms")
	fmt.Println("Response Body: ", string(body))
	fmt.Println("Response Header: ", responseHeaders)
	fmt.Println(">>> Originated from Request")
	printRequest(c)
	fmt.Println("-------------------")

	return err
}

func printRequest(c *fiber.Ctx) {
	fmt.Println("Request: ", c.Method(), " ", c.OriginalURL())
	fmt.Println("From: ", c.IP())

	// Filter request headers
	requestHeaders := filterHeaders(c.Request().Header.String())
	fmt.Println("Headers: ", requestHeaders)
	fmt.Println("Body: ", string(c.Request().Body()))
}

func filterHeaders(headers string) string {
	lines := strings.Split(headers, "\n")
	var filtered []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			filtered = append(filtered, strings.TrimSpace(line))
		}
	}
	return strings.Join(filtered, "\n")
}