package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/valyala/fasthttp"
	"todo-app/internal/db"
	"todo-app/internal/routes"
)

func init() {
	// Load environment variables
	loadEnv()
}

// loadEnv loads environment variables from the system
func loadEnv() {
	// This is a simple version - in production you'd use godotenv package
	// For now, we rely on system environment variables or docker env
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Println("Warning: DATABASE_URL not set, this will be required at runtime")
	}
}

func main() {
	// Initialize database connection
	err := db.Init()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create router
	router := routes.NewRouter()

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Add logging middleware
	handler := routes.CORSMiddleware(routes.RequestLogger(router.Handler))

	// Start server
	server := &fasthttp.Server{
		Handler: handler,
		Name:    "todo-app",
	}

	// Handle graceful shutdown
	go func() {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigch
		log.Printf("Received signal: %v\n", sig)
		server.Shutdown()
	}()

	addr := fmt.Sprintf("127.0.0.1:%s", port)
	log.Printf("Starting server on %s\n", addr)
	if err := server.ListenAndServe(addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
