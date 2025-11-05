package handlers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"todo-app/internal/db"
	"todo-app/internal/errors"
	"todo-app/internal/models"
)

// ListTodos retrieves all todos from the database
func ListTodos(ctx *fasthttp.RequestCtx) {
	query := `
		SELECT id, title, description, completed, created_at, updated_at
		FROM todos
		ORDER BY created_at DESC
	`

	rows, err := db.Pool.Query(context.Background(), query)
	if err != nil {
		log.Printf("Error querying todos: %v\n", err)
		errors.WriteErrorResponse(ctx, errors.NewInternalServerError("Failed to fetch todos"))
		return
	}
	defer rows.Close()

	var todos []models.TodoResponse
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			log.Printf("Error scanning todo: %v\n", err)
			errors.WriteErrorResponse(ctx, errors.NewInternalServerError("Failed to process todos"))
			return
		}

		var response models.TodoResponse
		response.FromTodo(&todo)
		todos = append(todos, response)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating todos: %v\n", err)
		errors.WriteErrorResponse(ctx, errors.NewInternalServerError("Failed to fetch todos"))
		return
	}

	if todos == nil {
		todos = []models.TodoResponse{}
	}

	errors.WriteJSONResponse(ctx, fasthttp.StatusOK, todos)
}

// GetTodo retrieves a single todo by ID
func GetTodo(ctx *fasthttp.RequestCtx) {
	idStr := ctx.UserValue("id").(string)
	id, err := uuid.Parse(idStr)
	if err != nil {
		errors.WriteErrorResponse(ctx, errors.NewBadRequest("Invalid todo ID format"))
		return
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
		errors.WriteErrorResponse(ctx, errors.NewNotFound("Todo not found"))
		return
	}

	var response models.TodoResponse
	response.FromTodo(&todo)

	errors.WriteJSONResponse(ctx, fasthttp.StatusOK, response)
}

// CreateTodo creates a new todo
func CreateTodo(ctx *fasthttp.RequestCtx) {
	var req models.CreateTodoRequest
	err := json.Unmarshal(ctx.PostBody(), &req)
	if err != nil {
		errors.WriteErrorResponse(ctx, errors.NewBadRequest("Invalid request body"))
		return
	}

	// Validate title is not empty
	if req.Title == "" {
		errors.WriteErrorResponse(ctx, errors.NewBadRequest("Title is required and cannot be empty"))
		return
	}

	id := uuid.New()
	now := time.Now().UTC()

	query := `
		INSERT INTO todos (id, title, description, completed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, description, completed, created_at, updated_at
	`

	var todo models.Todo
	err = db.Pool.QueryRow(
		context.Background(),
		query,
		id, req.Title, req.Description, false, now, now,
	).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		log.Printf("Error creating todo: %v\n", err)
		errors.WriteErrorResponse(ctx, errors.NewInternalServerError("Failed to create todo"))
		return
	}

	var response models.TodoResponse
	response.FromTodo(&todo)

	errors.WriteJSONResponse(ctx, fasthttp.StatusCreated, response)
}

// UpdateTodo updates an existing todo
func UpdateTodo(ctx *fasthttp.RequestCtx) {
	idStr := ctx.UserValue("id").(string)
	id, err := uuid.Parse(idStr)
	if err != nil {
		errors.WriteErrorResponse(ctx, errors.NewBadRequest("Invalid todo ID format"))
		return
	}

	var req models.UpdateTodoRequest
	err = json.Unmarshal(ctx.PostBody(), &req)
	if err != nil {
		errors.WriteErrorResponse(ctx, errors.NewBadRequest("Invalid request body"))
		return
	}

	// Check if todo exists
	checkQuery := `SELECT id FROM todos WHERE id = $1`
	var existingID uuid.UUID
	err = db.Pool.QueryRow(context.Background(), checkQuery, id).Scan(&existingID)
	if err != nil {
		log.Printf("Error checking todo existence: %v\n", err)
		errors.WriteErrorResponse(ctx, errors.NewNotFound("Todo not found"))
		return
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
		errors.WriteErrorResponse(ctx, errors.NewInternalServerError("Failed to update todo"))
		return
	}

	var response models.TodoResponse
	response.FromTodo(&todo)

	errors.WriteJSONResponse(ctx, fasthttp.StatusOK, response)
}

// DeleteTodo deletes a todo by ID
func DeleteTodo(ctx *fasthttp.RequestCtx) {
	idStr := ctx.UserValue("id").(string)
	id, err := uuid.Parse(idStr)
	if err != nil {
		errors.WriteErrorResponse(ctx, errors.NewBadRequest("Invalid todo ID format"))
		return
	}

	// Check if todo exists
	checkQuery := `SELECT id FROM todos WHERE id = $1`
	var existingID uuid.UUID
	err = db.Pool.QueryRow(context.Background(), checkQuery, id).Scan(&existingID)
	if err != nil {
		log.Printf("Error checking todo existence: %v\n", err)
		errors.WriteErrorResponse(ctx, errors.NewNotFound("Todo not found"))
		return
	}

	deleteQuery := `DELETE FROM todos WHERE id = $1`
	_, err = db.Pool.Exec(context.Background(), deleteQuery, id)
	if err != nil {
		log.Printf("Error deleting todo: %v\n", err)
		errors.WriteErrorResponse(ctx, errors.NewInternalServerError("Failed to delete todo"))
		return
	}

	ctx.SetStatusCode(fasthttp.StatusNoContent)
}
