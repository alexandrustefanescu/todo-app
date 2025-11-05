package routes

import (
	"github.com/gofiber/fiber/v2"
	"todo-app/internal/handlers"
	"todo-app/internal/middleware"
)

// Setup configures all routes for the application
func Setup(app *fiber.App) {
	// Apply global middleware
	app.Use(middleware.CORSMiddleware())

	// API routes
	api := app.Group("/api")
	todos := api.Group("/todos")

	// List todos
	todos.Get("", handlers.ListTodos)

	// Create todo
	todos.Post("", handlers.CreateTodo)

	// Get specific todo
	todos.Get("/:id", handlers.GetTodo)

	// Update todo
	todos.Put("/:id", handlers.UpdateTodo)

	// Delete todo
	todos.Delete("/:id", handlers.DeleteTodo)
}
