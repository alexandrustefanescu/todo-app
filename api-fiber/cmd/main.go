package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"todo-app/internal/db"
	"todo-app/internal/routes"
)

func main() {
	// Initialize database connection
	err := db.Init()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Todo API",
	})

	// Middleware
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} ${latency}\n",
	}))

	// Setup routes
	routes.Setup(app)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Handle graceful shutdown
	go func() {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigch
		log.Printf("Received signal: %v\n", sig)
		app.Shutdown()
	}()

	// Start server
	addr := fmt.Sprintf("127.0.0.1:%s", port)
	log.Printf("Starting server on %s\n", addr)
	if err := app.Listen(addr); err != nil && err != fiber.ErrShutdown {
		log.Fatalf("Server error: %v", err)
	}
	log.Println("Server shut down gracefully")
}
