package errors

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"log"
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
		Status:  fasthttp.StatusNotFound,
	}
}

// NewBadRequest creates a BAD_REQUEST error
func NewBadRequest(message string) *APIError {
	return &APIError{
		Type:    BadRequest,
		Message: message,
		Status:  fasthttp.StatusBadRequest,
	}
}

// NewInternalServerError creates an INTERNAL_SERVER_ERROR error
func NewInternalServerError(message string) *APIError {
	return &APIError{
		Type:    InternalServerError,
		Message: message,
		Status:  fasthttp.StatusInternalServerError,
	}
}

// NewConflict creates a CONFLICT error
func NewConflict(message string) *APIError {
	return &APIError{
		Type:    Conflict,
		Message: message,
		Status:  fasthttp.StatusConflict,
	}
}

// WriteErrorResponse writes an error response to the HTTP context
func WriteErrorResponse(ctx *fasthttp.RequestCtx, err *APIError) {
	ctx.SetStatusCode(err.Status)
	ctx.SetContentType("application/json; charset=utf-8")

	response := models.ErrorResponse{
		Error:   string(err.Type),
		Message: err.Message,
	}

	body, marshallErr := json.Marshal(response)
	if marshallErr != nil {
		log.Printf("Error marshalling error response: %v\n", marshallErr)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(`{"error":"INTERNAL_SERVER_ERROR","message":"Failed to marshal error response"}`)
		return
	}

	ctx.Write(body)
}

// WriteJSONResponse writes a successful JSON response to the HTTP context
func WriteJSONResponse(ctx *fasthttp.RequestCtx, statusCode int, data interface{}) {
	ctx.SetStatusCode(statusCode)
	ctx.SetContentType("application/json; charset=utf-8")

	body, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling response: %v\n", err)
		WriteErrorResponse(ctx, NewInternalServerError("Failed to marshal response"))
		return
	}

	ctx.Write(body)
}
