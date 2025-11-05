package errors

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"todo-app/internal/models"
)

// ErrorType represents the type of error
type ErrorType string

const (
	NotFound            ErrorType = "NOT_FOUND"
	BadRequest          ErrorType = "BAD_REQUEST"
	InternalServerError ErrorType = "INTERNAL_SERVER_ERROR"
	Conflict            ErrorType = "CONFLICT"
)

// APIError represents an API error with type and message
type APIError struct {
	Type    ErrorType
	Message string
	Status  int
}

// NewNotFound creates a NOT_FOUND error
func NewNotFound(message string) *APIError {
	return &APIError{
		Type:    NotFound,
		Message: message,
		Status:  fiber.StatusNotFound,
	}
}

// NewBadRequest creates a BAD_REQUEST error
func NewBadRequest(message string) *APIError {
	return &APIError{
		Type:    BadRequest,
		Message: message,
		Status:  fiber.StatusBadRequest,
	}
}

// NewInternalServerError creates an INTERNAL_SERVER_ERROR error
func NewInternalServerError(message string) *APIError {
	return &APIError{
		Type:    InternalServerError,
		Message: message,
		Status:  fiber.StatusInternalServerError,
	}
}

// NewConflict creates a CONFLICT error
func NewConflict(message string) *APIError {
	return &APIError{
		Type:    Conflict,
		Message: message,
		Status:  fiber.StatusConflict,
	}
}

// HandleError sends an error response
func HandleError(c *fiber.Ctx, err *APIError) error {
	response := models.ErrorResponse{
		Error:   string(err.Type),
		Message: err.Message,
	}
	return c.Status(err.Status).JSON(response)
}

// HandleInternalError logs and returns an internal server error
func HandleInternalError(c *fiber.Ctx, err error) error {
	log.Printf("Internal error: %v\n", err)
	return HandleError(c, NewInternalServerError("An internal error occurred"))
}
