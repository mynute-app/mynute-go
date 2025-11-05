package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
)

// Docker Compose Dev Runner
// Starts all microservices docker-compose.dev.yml files
// Usage: go run cmd/docker-dev/main.go

type Service struct {
	Name    string
	Project string
	File    string
}

var services = []Service{
	{Name: "Core Service", Project: "mynute-go-core", File: "services/core/docker-compose.dev.yml"},
	{Name: "Auth Service", Project: "mynute-go-auth", File: "services/auth/docker-compose.dev.yml"},
	{Name: "Email Service", Project: "mynute-go-email", File: "services/email/docker-compose.dev.yml"},
}

func main() {
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║          Docker Compose Development Environment              ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	action := "up"
	if len(os.Args) > 1 {
		action = os.Args[1]
	}

	switch action {
	case "up":
		startAllServices()
	case "down":
		stopAllServices()
	case "restart":
		stopAllServices()
		startAllServices()
	case "logs":
		showLogs()
	default:
		showHelp()
	}
}

func startAllServices() {
	fmt.Println("🚀 Starting all services...")
	fmt.Println()

	var wg sync.WaitGroup
	errors := make(chan error, len(services))

	for _, svc := range services {
		wg.Add(1)
		go func(s Service) {
			defer wg.Done()
			fmt.Printf("📦 Starting %s...\n", s.Name)
			cmd := exec.Command("docker-compose", "-p", s.Project, "-f", s.File, "up", "-d", "--force-recreate")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				errors <- fmt.Errorf("%s failed: %v", s.Name, err)
			} else {
				fmt.Printf("✅ %s started successfully\n", s.Name)
			}
		}(svc)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	hasErrors := false
	for err := range errors {
		log.Println("❌", err)
		hasErrors = true
	}

	if !hasErrors {
		fmt.Println()
		fmt.Println("✅ All services started successfully!")
		fmt.Println()
		fmt.Println("To stop all services: go run cmd/docker-dev/main.go down")
		fmt.Println("To view logs: go run cmd/docker-dev/main.go logs")
	}
}

func stopAllServices() {
	fmt.Println("🛑 Stopping all services...")
	fmt.Println()

	var wg sync.WaitGroup
	for _, svc := range services {
		wg.Add(1)
		go func(s Service) {
			defer wg.Done()
			fmt.Printf("📦 Stopping %s...\n", s.Name)
			cmd := exec.Command("docker-compose", "-p", s.Project, "-f", s.File, "down")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Printf("❌ %s failed to stop: %v\n", s.Name, err)
			} else {
				fmt.Printf("✅ %s stopped successfully\n", s.Name)
			}
		}(svc)
	}

	wg.Wait()
	fmt.Println()
	fmt.Println("✅ All services stopped!")
}

func showLogs() {
	fmt.Println("📋 Showing logs for all services...")
	fmt.Println("Press Ctrl+C to exit")
	fmt.Println()

	// Channel to handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	for _, svc := range services {
		wg.Add(1)
		go func(s Service) {
			defer wg.Done()
			cmd := exec.Command("docker-compose", "-p", s.Project, "-f", s.File, "logs", "-f")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		}(svc)
	}

	// Wait for interrupt signal
	<-quit
	fmt.Println("\n\n🛑 Stopping log stream...")
}

func showHelp() {
	fmt.Println("Usage: go run cmd/docker-dev/main.go [command]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  up       - Start all services (default)")
	fmt.Println("  down     - Stop all services")
	fmt.Println("  restart  - Restart all services")
	fmt.Println("  logs     - Show logs from all services")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/docker-dev/main.go")
	fmt.Println("  go run cmd/docker-dev/main.go up")
	fmt.Println("  go run cmd/docker-dev/main.go down")
	fmt.Println("  go run cmd/docker-dev/main.go logs")
}
