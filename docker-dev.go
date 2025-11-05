//go:build ignore
// +build ignore

package main

import (
	"os"
	"os/exec"
)

// Docker Dev Launcher
// This is a convenience wrapper to run cmd/docker-dev/main.go
// Usage: go run docker-dev.go [up|down|restart|logs]

func main() {
	// Forward all arguments to the actual docker-dev command
	args := append([]string{"run", "cmd/docker-dev/main.go"}, os.Args[1:]...)
	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
}
