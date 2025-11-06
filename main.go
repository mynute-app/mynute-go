package main

import (
	"fmt"
	"log"
	"mynute-go/services/auth"
	"mynute-go/services/core"
	"mynute-go/services/email"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

//	@title			Mynute Microservices
//	@version		1.0
//	@description	Mynute microservices architecture
//	@termsOfService	http://swagger.io/terms/
//	@contact.name	API Support
//	@contact.email	support@mynute.com
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

func main() {
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║          Mynute Microservices Architecture                    ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("This project uses separate microservices.")
	fmt.Println()
	fmt.Println("Available services:")
	fmt.Println("  • Business Service (Core) - Port 4000")
	fmt.Println("  • Auth Service - Port 4001")
	fmt.Println("  • Email Service - Port 4002")
	fmt.Println()
	fmt.Println("Running Options:")
	fmt.Println("  1. Run all services:     go run .")
	fmt.Println("  2. Run individually:")
	fmt.Println("      Core:   go run ./cmd/business-service")
	fmt.Println("      Auth:   go run ./cmd/auth-service")
	fmt.Println("      Email:  go run ./cmd/email-service")
	fmt.Println()
	fmt.Println("Docker Compose:")
	fmt.Println("  • Development: docker-compose -f docker-compose.dev.yml up")
	fmt.Println("  • Production:  docker-compose -f docker-compose.prod.yml up")
	fmt.Println()
	fmt.Println("Go Workspace:")
	fmt.Println("  This project uses Go modules in a monorepo structure.")
	fmt.Println("  All services share the same go.mod at the root level.")
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
