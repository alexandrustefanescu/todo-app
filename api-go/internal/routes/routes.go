package routes

import (
	"log"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttprouter"
	"todo-app/internal/handlers"
)

// Router represents the HTTP router
type Router struct {
	*fasthttprouter.Router
}

// NewRouter creates a new router with all routes configured
func NewRouter() *Router {
	router := fasthttprouter.New()

	// Configure routes
	router.GET("/api/todos", handlers.ListTodos)
	router.POST("/api/todos", handlers.CreateTodo)
	router.GET("/api/todos/:id", handlers.GetTodo)
	router.PUT("/api/todos/:id", handlers.UpdateTodo)
	router.DELETE("/api/todos/:id", handlers.DeleteTodo)

	// Handle 404
	router.NotFound = func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.WriteString(`{"error":"NOT_FOUND","message":"Endpoint not found"}`)
	}

	return &Router{router}
}

// RequestLogger logs all incoming HTTP requests
func RequestLogger(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		log.Printf("%s %s %s\n", ctx.Method(), ctx.RequestURI(), ctx.RemoteAddr())
		next(ctx)
	}
}

// CORSHandler handles CORS headers for cross-origin requests
func CORSHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if string(ctx.Method()) == "OPTIONS" {
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}
}

// CORSMiddleware returns a middleware that applies CORS headers
func CORSMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		CORSHandler(ctx)
		next(ctx)
	}
}
