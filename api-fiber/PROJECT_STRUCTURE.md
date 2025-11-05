# Fiber Todo API - Project Structure & Architecture

## Overview

This document describes the architecture and organization of the Fiber Todo API implementation. Fiber is an Express.js-inspired framework built on FastHTTP, providing a clean and intuitive API while maintaining excellent performance.

## Directory Structure

```
api-fiber/
├── cmd/                                    # Application entry point
│   └── main.go                             # Server initialization and startup logic
├── internal/                               # Private packages (not exposed outside module)
│   ├── db/
│   │   └── db.go                          # Database connection pool initialization
│   ├── handlers/
│   │   └── todo.go                        # HTTP request handlers for CRUD operations
│   ├── models/
│   │   └── todo.go                        # Data structures, DTOs, and type definitions
│   ├── routes/
│   │   └── routes.go                      # Route definitions and setup
│   ├── middleware/
│   │   └── cors.go                        # CORS middleware configuration
│   └── errors/
│       └── errors.go                      # Custom error types and response handling
├── migrations/
│   └── 01_create_todos_table.sql          # Initial database schema
├── go.mod                                  # Go module definition with latest versions
├── go.sum                                  # Dependency version checksums
├── .env                                    # Environment variables (local development)
├── .gitignore                             # Git ignore rules
├── Dockerfile                              # Multi-stage container build
├── docker-compose.yml                      # Service orchestration configuration
├── Makefile                               # Build and development tasks
├── README.md                              # User documentation
├── QUICK_START.md                         # Quick setup guide
└── PROJECT_STRUCTURE.md                   # This file
```

## Core Components

### 1. Main Application (cmd/main.go)

**Purpose:** Application entry point and Fiber app initialization

**Key Functions:**
- `main()` - Initialize database, create Fiber app, start server

**Responsibilities:**
- Initialize PostgreSQL connection pool
- Create Fiber application with configuration
- Register middleware (logger, CORS)
- Setup routes
- Handle graceful shutdown with signal handling
- Bind to configured port and start listening

**Key Features:**
- Built-in Fiber logger middleware
- Custom app name configuration
- Graceful shutdown support
- Signal handling (SIGINT, SIGTERM)

### 2. Database Layer (internal/db/db.go)

**Purpose:** Manage PostgreSQL connections and pool

**Key Functions:**
- `Init()` - Initialize connection pool
- `Close()` - Close connection pool gracefully

**Features:**
- Creates pgx connection pool
- Connection reuse and pooling
- Health checks via Ping()
- Graceful shutdown

**Connection Pool Details:**
- Driver: pgx/v5 (high-performance PostgreSQL driver)
- Uses standard Go database/sql interface
- Automatic connection management
- Default behavior handles thousands of concurrent connections

### 3. Models (internal/models/todo.go)

**Purpose:** Define data structures and DTOs

**Type Definitions:**

```go
// Database model
type Todo struct {
    ID          uuid.UUID
    Title       string
    Description *string
    Completed   bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// Request DTOs
type CreateTodoRequest struct {
    Title       string
    Description *string
}

type UpdateTodoRequest struct {
    Title       *string
    Description *string
    Completed   *bool
}

// Response DTO
type TodoResponse struct {
    // Same as Todo
}

// Error response
type ErrorResponse struct {
    Error   string
    Message string
}
```

**JSON Serialization:**
- All types use `encoding/json` struct tags
- Supports both marshaling and unmarshaling
- Pointers used for optional fields in requests

### 4. Handlers (internal/handlers/todo.go)

**Purpose:** HTTP request handling for CRUD operations

**Handler Functions:**

| Function | HTTP Method | Route | Description |
|----------|-------------|-------|-------------|
| `ListTodos()` | GET | `/api/todos` | Fetch all todos (DESC by creation date) |
| `GetTodo()` | GET | `/api/todos/:id` | Fetch single todo by UUID |
| `CreateTodo()` | POST | `/api/todos` | Create new todo with validation |
| `UpdateTodo()` | PUT | `/api/todos/:id` | Partial update with null coalescing |
| `DeleteTodo()` | DELETE | `/api/todos/:id` | Delete todo by UUID |

**Handler Pattern with Fiber:**
```go
func Handler(c *fiber.Ctx) error {
    // 1. Parse request (query params, body, path params)
    // 2. Validate input
    // 3. Execute business logic (database operations)
    // 4. Return response (c.JSON, c.SendStatus, etc.)
}
```

**Key Implementation Details:**

- **Context:** Each handler receives `*fiber.Ctx` for request/response handling
- **Path Parameters:** Accessed via `c.Params("id")`
- **Request Body:** Parsed with `c.BodyParser(&req)`
- **Response Writing:** Clean methods like `c.Status().JSON()` and `c.SendStatus()`
- **Error Handling:** Returns error from handler (Fiber handles it)
- **Validation:** Title required, UUID format validation

### 5. Routing (internal/routes/routes.go)

**Purpose:** Define API routes and endpoint setup

**Router Configuration:**
- Uses Fiber's built-in routing
- Supports path parameters (`:id`)
- Automatic 404 handling
- Route grouping for organization

**Route Structure:**
```
GET    /api/todos      → ListTodos
POST   /api/todos      → CreateTodo
GET    /api/todos/:id  → GetTodo
PUT    /api/todos/:id  → UpdateTodo
DELETE /api/todos/:id  → DeleteTodo
```

**Setup Pattern:**
- Uses route groups for clean organization
- Global middleware applied at app level
- Per-group middleware support available

### 6. Middleware (internal/middleware/cors.go)

**Purpose:** Handle cross-origin requests

**CORS Configuration:**
- Allows all origins (`*`)
- Allows common HTTP methods (GET, POST, PUT, DELETE, OPTIONS)
- Allows Content-Type and Authorization headers

**Fiber Middleware Approach:**
- Uses Fiber's built-in CORS middleware
- Easy to configure and extend
- Can be applied globally or per-route

### 7. Error Handling (internal/errors/errors.go)

**Purpose:** Centralized error handling and response formatting

**Error Types:**
```go
const (
    NotFound            = "NOT_FOUND"
    BadRequest          = "BAD_REQUEST"
    InternalServerError = "INTERNAL_SERVER_ERROR"
    Conflict            = "CONFLICT"
)
```

**Error Structure:**
```go
type APIError struct {
    Type    ErrorType
    Message string
    Status  int
}
```

**Helper Functions:**
- `NewNotFound()` - 404 errors
- `NewBadRequest()` - 400 errors
- `NewInternalServerError()` - 500 errors
- `NewConflict()` - 409 errors
- `HandleError()` - Marshal and send error response
- `HandleInternalError()` - Log and return internal error

**Error Response Format:**
```json
{
  "error": "ERROR_TYPE",
  "message": "Description of what went wrong"
}
```

### 8. Database Migrations (migrations/01_create_todos_table.sql)

**Purpose:** Initialize database schema

**Contents:**
- Creates `todos` table with UUID primary key
- Defines all columns with proper types
- Creates indexes on frequently queried fields

**Execution:**
```bash
psql -U postgres -d todo_db -f migrations/01_create_todos_table.sql
# Or via docker-compose (automatic)
```

## Data Flow

### GET /api/todos (List)

```
HTTP Request
    ↓
Fiber Router matches route
    ↓
ListTodos handler
    ↓
Query database (SELECT * FROM todos ORDER BY created_at DESC)
    ↓
Scan rows into Todo structs
    ↓
Convert to TodoResponse structs
    ↓
Marshal to JSON
    ↓
c.Status(200).JSON(todos)
    ↓
HTTP Response
```

### POST /api/todos (Create)

```
HTTP Request with JSON body
    ↓
Fiber Router matches route
    ↓
CreateTodo handler
    ↓
c.BodyParser(&req) - Parse JSON
    ↓
Validate (title not empty)
    ↓
Generate UUID and timestamps
    ↓
INSERT INTO todos
    ↓
Scan result into Todo struct
    ↓
Convert to TodoResponse
    ↓
c.Status(201).JSON(response)
    ↓
HTTP Response
```

### PUT /api/todos/:id (Update)

```
HTTP Request with JSON body
    ↓
Fiber Router matches route with :id parameter
    ↓
UpdateTodo handler
    ↓
c.Params("id") - Get path parameter
    ↓
c.BodyParser(&req) - Parse JSON
    ↓
Check todo exists
    ↓
UPDATE todos (COALESCE for partial updates)
    ↓
Scan result
    ↓
c.Status(200).JSON(response)
    ↓
HTTP Response
```

### DELETE /api/todos/:id

```
HTTP Request
    ↓
Fiber Router matches route
    ↓
DeleteTodo handler
    ↓
Parse UUID from path param
    ↓
Check todo exists
    ↓
DELETE FROM todos WHERE id = ?
    ↓
c.SendStatus(204)
    ↓
HTTP Response (No Content)
```

## Technology Details

### Fiber Framework

**Advantages:**
- Express.js-like API (familiar syntax)
- Built on FastHTTP (10x faster than net/http)
- Rich middleware ecosystem
- Better developer experience than raw FastHTTP
- Active community and development
- Excellent documentation

**Key Features:**
- Route grouping and organization
- Built-in middleware (logger, compression, recovery)
- Clean error handling integration
- Body parser with multiple formats
- Static file serving
- Middleware chaining

### pgx Driver Benefits

- **Performance:** Built for speed with prepared statement caching
- **Safety:** Prevents SQL injection with parameterized queries
- **Features:** Supports UUIDs, JSON, arrays natively
- **Connection Pooling:** Built-in pool management
- **Modern API:** Context-based cancellation and timeout support

### UUID Generation

- Uses `github.com/google/uuid` package
- Generates v4 (random) UUIDs
- Database generates UUIDs as default
- Application can also generate UUIDs

### Time Handling

- All timestamps use `time.Time` with UTC timezone
- Database stores as TIMESTAMPTZ (timezone-aware)
- JSON marshaling uses RFC3339 format
- Automatic created_at/updated_at management

## Fiber vs Other Go Web Frameworks

| Aspect | Fiber | FastHTTP Router | Standard net/http |
|--------|-------|-----------------|------------------|
| **Learning Curve** | Gentle | Moderate | Easy |
| **Performance** | Excellent | Excellent | Good |
| **API Clarity** | Excellent | Good | Good |
| **Middleware** | Rich | Manual | Excellent |
| **Community** | Growing | Small | Large |
| **Ecosystem** | Good | Minimal | Excellent |

## Development Workflow

### Local Development
1. Run PostgreSQL in Docker or locally
2. Set DATABASE_URL in .env
3. Run migrations
4. Use `go run ./cmd/main.go` for development
5. API auto-reloads with tools like `air`

### Testing
- Create `*_test.go` files in same package
- Use standard `testing` package
- Run with `go test ./...`

### Deployment
- Build native binary: `go build -o todo-app ./cmd/main.go`
- Or use Docker: `docker build -t todo-app-fiber .`
- Deploy as single binary or container

## Security Considerations

1. **SQL Injection:** Prevented by parameterized queries
2. **Input Validation:** Title requirement, UUID format checks
3. **CORS:** Currently allows all origins (configure for production)
4. **Rate Limiting:** Not implemented (add if needed)
5. **Authentication:** Not implemented (add for production)
6. **HTTPS:** Implement via reverse proxy (nginx, etc.)

## Performance Tuning

1. **Connection Pool Size:** pgx manages automatically
2. **Timeouts:** Add context timeouts to database operations
3. **Indexes:** Already created on common query fields
4. **Caching:** Consider Redis for frequently accessed todos
5. **Compression:** Fiber has built-in compression middleware

## Future Enhancements

1. **Testing:** Add unit and integration tests
2. **Validation:** More comprehensive input validation
3. **Filtering:** Add query parameters for filtering/sorting
4. **Pagination:** Implement limit/offset pagination
5. **Caching:** Add Redis caching layer
6. **Metrics:** Prometheus/Grafana monitoring
7. **Authentication:** JWT or OAuth2 integration
8. **Documentation:** OpenAPI/Swagger specs
9. **Rate Limiting:** Prevent API abuse
10. **Graceful Degradation:** Circuit breakers, retries

## Comparison with Other Implementations

| Aspect | Fiber | FastHTTP Router | Rust Actix |
|--------|-------|-----------------|-----------|
| **API Style** | Express-like | Low-level | Actor-based |
| **Code Size** | Small | Small | Small |
| **Learning Curve** | Gentle | Moderate | Steep |
| **Performance** | Excellent | Excellent | Excellent |
| **Ecosystem** | Growing | Minimal | Growing |
| **Type Safety** | Runtime | Runtime | Compile-time |
| **Development Speed** | Fast | Medium | Slow |

**Best For:**
- **Fiber:** Teams with Node.js/Express background, fast development
- **FastHTTP Router:** Minimal dependencies, low-level control
- **Rust Actix:** Maximum compile-time safety, performance tuning

## Summary

This Fiber implementation provides:

- **Clean API:** Express.js-like syntax
- **High Performance:** Built on FastHTTP
- **Full Functionality:** Complete CRUD operations
- **Production Ready:** Error handling, logging, graceful shutdown
- **Easy Maintenance:** Clear code organization
- **Rich Ecosystem:** Access to Fiber middleware

The architecture follows Go best practices with clear separation of concerns, proper error handling, and efficient database connection management.
