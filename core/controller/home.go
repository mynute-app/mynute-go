package controller

import (
	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	return c.SendString(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>MyNute API</title>
			<style>
				body {
					font-family: sans-serif;
					background-color: #f9fafb;
					display: flex;
					flex-direction: column;
					justify-content: center;
					align-items: center;
					height: 100vh;
					margin: 0;
				}
				h1 {
					color: #1f2937;
					font-size: 2rem;
				}
				p {
					color: #4b5563;
				}
				code {
					background: #e5e7eb;
					padding: 4px 6px;
					border-radius: 5px;
					font-size: 0.9rem;
				}
			</style>
		</head>
		<body>
			<h1>🚀 Up and running! 🚀</h1>
			<p>Welcome to the <code>Mynute API</code> written in Go.</p>
		</body>
		</html>
	`)
}
