package main

import (
	"fmt"
	"log"
	"mynute-go/services/auth"
	"mynute-go/services/core"
	_ "mynute-go/docs"
	"mynute-go/services/email"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// @title						Mynute Microservices
// @version					1.0
// @description				Mynute microservices architecture
// @termsOfService				http://swagger.io/terms/
// @contact.name				API Support
// @contact.email				support@mynute.com
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html

func main() {
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║          Mynute Microservices Architecture                    ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("This project has been refactored into separate microservices.")
	fmt.Println()
	fmt.Println("Available services:")
	fmt.Println("  • Business Service (Core) - Port 4000")
	fmt.Println("      Run: go run cmd/business-service/main.go (It will run the microservice alone)")
	fmt.Println("  • Auth Service - Port 4001")
	fmt.Println("      Run: go run cmd/auth-service/main.go (It will run the microservice alone)")
	fmt.Println("  • Email Service - Port 4002")
	fmt.Println("      Run: go run cmd/email-service/main.go (It will run the microservice alone)")
	fmt.Println()
	fmt.Println("Running Options:")
	fmt.Println("  1. Run all services:     go run main.go")
	fmt.Println("  2. Run individually:     See commands above")
	fmt.Println()
	fmt.Println("Docker Compose:")
	fmt.Println("  • Core/Business: docker-compose -f core/docker-compose.dev.yml up")
	fmt.Println("  • Auth:          docker-compose -f auth/docker-compose.dev.yml up")
	fmt.Println("  • Email:         docker-compose -f email/docker-compose.dev.yml up")
	fmt.Println()
	fmt.Println("For more information, see the documentation in /docs")
	fmt.Println()
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println()
	fmt.Println("🚀 Starting all services...")
	fmt.Println()

	var wg sync.WaitGroup
	wg.Add(3)

	// Channel to handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Start Business Service
	var coreServer *core.Server
	go func() {
		defer wg.Done()
		log.Println("🚀 Starting Business Service on port 4000...")
		coreServer = core.NewServer()
		coreServer.Run("listen")
	}()

	// Start Auth Service
	var authServer *auth.Server
	go func() {
		defer wg.Done()
		log.Println("🚀 Starting Auth Service on port 4001...")
		authServer = auth.NewServer()
		authServer.Run("listen")
	}()

	// Start Email Service
	var emailServer *email.Server
	go func() {
		defer wg.Done()
		log.Println("🚀 Starting Email Service on port 4002...")
		emailServer = email.NewServer()
		emailServer.Run("listen")
	}()

	// Wait for interrupt signal
	<-quit
	fmt.Println()
	log.Println("⏳ Shutting down services gracefully...")

	// Shutdown both servers
	if coreServer != nil {
		coreServer.Shutdown()
	}
	if authServer != nil {
		authServer.Shutdown()
	}
	if emailServer != nil {
		emailServer.Shutdown()
	}

	wg.Wait()
	log.Println("✅ All services stopped successfully")
}
