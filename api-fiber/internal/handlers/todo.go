package handlers

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"todo-app/internal/db"
	apperrors "todo-app/internal/errors"
	"todo-app/internal/models"
)

// ListTodos retrieves all todos from the database
func ListTodos(c *fiber.Ctx) error {
	query := `
		SELECT id, title, description, completed, created_at, updated_at
		FROM todos
		ORDER BY created_at DESC
	`

	rows, err := db.Pool.Query(context.Background(), query)
	if err != nil {
		log.Printf("Error querying todos: %v\n", err)
		return apperrors.HandleError(c, apperrors.NewInternalServerError("Failed to fetch todos"))
	}
	defer rows.Close()

	var todos []models.TodoResponse
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			log.Printf("Error scanning todo: %v\n", err)
			return apperrors.HandleError(c, apperrors.NewInternalServerError("Failed to process todos"))
		}

		var response models.TodoResponse
		response.FromTodo(&todo)
		todos = append(todos, response)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating todos: %v\n", err)
		return apperrors.HandleError(c, apperrors.NewInternalServerError("Failed to fetch todos"))
	}

	if todos == nil {
		todos = []models.TodoResponse{}
	}

	return c.Status(fiber.StatusOK).JSON(todos)
}

// GetTodo retrieves a single todo by ID
func GetTodo(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return apperrors.HandleError(c, apperrors.NewBadRequest("Invalid todo ID format"))
	}

	query := `
		SELECT id, title, description, completed, created_at, updated_at
		FROM todos
		WHERE id = $1
	`

	var todo models.Todo
	err = db.Pool.QueryRow(context.Background(), query, id).Scan(
		&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt,
	)
	if err != nil {
		log.Printf("Error querying todo: %v\n", err)
		return apperrors.HandleError(c, apperrors.NewNotFound("Todo not found"))
	}

	var response models.TodoResponse
	response.FromTodo(&todo)

	return c.Status(fiber.StatusOK).JSON(response)
}

// CreateTodo creates a new todo
func CreateTodo(c *fiber.Ctx) error {
	var req models.CreateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.HandleError(c, apperrors.NewBadRequest("Invalid request body"))
	}

	// Validate title is not empty
	if req.Title == "" {
		return apperrors.HandleError(c, apperrors.NewBadRequest("Title is required and cannot be empty"))
	}

	id := uuid.New()
	now := time.Now().UTC()

	query := `
		INSERT INTO todos (id, title, description, completed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, description, completed, created_at, updated_at
	`

	var todo models.Todo
	err := db.Pool.QueryRow(
		context.Background(),
		query,
		id, req.Title, req.Description, false, now, now,
	).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		log.Printf("Error creating todo: %v\n", err)
		return apperrors.HandleError(c, apperrors.NewInternalServerError("Failed to create todo"))
	}

	var response models.TodoResponse
	response.FromTodo(&todo)

	return c.Status(fiber.StatusCreated).JSON(response)
}

// UpdateTodo updates an existing todo
func UpdateTodo(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return apperrors.HandleError(c, apperrors.NewBadRequest("Invalid todo ID format"))
	}

	var req models.UpdateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.HandleError(c, apperrors.NewBadRequest("Invalid request body"))
	}

	// Check if todo exists
	checkQuery := `SELECT id FROM todos WHERE id = $1`
	var existingID uuid.UUID
	err = db.Pool.QueryRow(context.Background(), checkQuery, id).Scan(&existingID)
	if err != nil {
		log.Printf("Error checking todo existence: %v\n", err)
		return apperrors.HandleError(c, apperrors.NewNotFound("Todo not found"))
	}

	// Build dynamic update query
	updateQuery := `
		UPDATE todos
		SET
			title = COALESCE($2, title),
			description = COALESCE($3, description),
			completed = COALESCE($4, completed),
			updated_at = $5
		WHERE id = $1
		RETURNING id, title, description, completed, created_at, updated_at
	`

	var todo models.Todo
	err = db.Pool.QueryRow(
		context.Background(),
		updateQuery,
		id,
		req.Title,
		req.Description,
		req.Completed,
		time.Now().UTC(),
	).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		log.Printf("Error updating todo: %v\n", err)
		return apperrors.HandleError(c, apperrors.NewInternalServerError("Failed to update todo"))
	}

	var response models.TodoResponse
	response.FromTodo(&todo)

	return c.Status(fiber.StatusOK).JSON(response)
}

// DeleteTodo deletes a todo by ID
func DeleteTodo(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return apperrors.HandleError(c, apperrors.NewBadRequest("Invalid todo ID format"))
	}

	// Check if todo exists
	checkQuery := `SELECT id FROM todos WHERE id = $1`
	var existingID uuid.UUID
	err = db.Pool.QueryRow(context.Background(), checkQuery, id).Scan(&existingID)
	if err != nil {
		log.Printf("Error checking todo existence: %v\n", err)
		return apperrors.HandleError(c, apperrors.NewNotFound("Todo not found"))
	}

	deleteQuery := `DELETE FROM todos WHERE id = $1`
	_, err = db.Pool.Exec(context.Background(), deleteQuery, id)
	if err != nil {
		log.Printf("Error deleting todo: %v\n", err)
		return apperrors.HandleError(c, apperrors.NewInternalServerError("Failed to delete todo"))
	}

	return c.SendStatus(fiber.StatusNoContent)
}
