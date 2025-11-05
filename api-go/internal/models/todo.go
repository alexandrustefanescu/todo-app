package models

import (
	"time"

	"github.com/google/uuid"
)

// Todo represents a todo item in the database
type Todo struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CreateTodoRequest is the request payload for creating a new todo
type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description *string `json:"description"`
}

// UpdateTodoRequest is the request payload for updating a todo
type UpdateTodoRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}

// TodoResponse is the response payload for todo operations
type TodoResponse struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// FromTodo converts a Todo to a TodoResponse
func (tr *TodoResponse) FromTodo(t *Todo) {
	tr.ID = t.ID
	tr.Title = t.Title
	tr.Description = t.Description
	tr.Completed = t.Completed
	tr.CreatedAt = t.CreatedAt
	tr.UpdatedAt = t.UpdatedAt
}

// ErrorResponse is the error response payload
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
