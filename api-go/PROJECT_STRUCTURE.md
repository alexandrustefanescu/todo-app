# Go FastHTTP Todo API - Project Structure & Architecture

## Overview

This document describes the architecture and organization of the Go FastHTTP Todo API implementation.

## Directory Structure

```
api-go/
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
│   │   └── routes.go                      # Route definitions and middleware
│   └── errors/
│       └── errors.go                      # Custom error types and response handling
├── migrations/
│   └── 01_create_todos_table.sql          # Initial database schema
├── go.mod                                  # Go module definition
├── go.sum                                  # Dependency version checksums
├── .env                                    # Environment variables (local development)
├── .gitignore                             # Git ignore rules
├── Dockerfile                              # Multi-stage container build
├── docker-compose.yml                      # Service orchestration configuration
├── Makefile                               # Build and development tasks
├── README.md                              # User documentation
└── PROJECT_STRUCTURE.md                   # This file
```

## Core Components

### 1. Main Application (cmd/main.go)

**Purpose:** Application entry point and server initialization

**Key Functions:**
- `init()` - Environment setup
- `loadEnv()` - Load environment variables
- `main()` - Initialize database, create router, start server

**Responsibilities:**
- Initialize PostgreSQL connection pool
- Create and configure FastHTTP router
- Set up middleware (CORS, request logging)
- Handle graceful shutdown with signal handling
- Bind to configured port and start listening

**Configuration:**
- `DATABASE_URL` - PostgreSQL connection string (required)
- `PORT` - Server port (default: 8080)

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
- Default max connections: pool default (typically 25)

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

// Request DTO
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
    // Identical to Todo
}

// Error response
type ErrorResponse struct {
    Error   string
    Message string
}
```

**JSON Serialization:**
- All types use `encoding/json` struct tags
- Supports both marshaling (Go → JSON) and unmarshaling (JSON → Go)
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

**Handler Pattern:**
```go
func Handler(ctx *fasthttp.RequestCtx) {
    // 1. Parse request (query params, body, path params)
    // 2. Validate input
    // 3. Execute business logic (database operations)
    // 4. Write response (success or error)
}
```

**Key Implementation Details:**

- **Context Usage:** Each handler receives `*fasthttp.RequestCtx` for request/response handling
- **Path Parameters:** Accessed via `ctx.UserValue("id")`
- **Request Body:** Parsed with `json.Unmarshal(ctx.PostBody(), &req)`
- **Response Writing:** Centralized via error helpers
- **Error Handling:** Returns structured error responses with HTTP status codes
- **Validation:** Empty title check, UUID format validation

### 5. Routing (internal/routes/routes.go)

**Purpose:** Define API routes and apply middleware

**Router Configuration:**
- Uses `fasthttprouter.Router` for high-performance routing
- Supports path parameters (`:id`)
- Automatic 404 handling

**Routes:**
```
GET    /api/todos      → ListTodos
POST   /api/todos      → CreateTodo
GET    /api/todos/:id  → GetTodo
PUT    /api/todos/:id  → UpdateTodo
DELETE /api/todos/:id  → DeleteTodo
```

**Middleware:**

1. **CORSMiddleware** - Adds CORS headers to responses
   - Allows all origins (`*`)
   - Allows common HTTP methods
   - Allows Content-Type header

2. **RequestLogger** - Logs all incoming requests
   - Format: `METHOD PATH REMOTE_ADDR`
   - Useful for debugging and monitoring

**Middleware Composition:**
```go
handler := CORSMiddleware(RequestLogger(router.Handler))
```

### 6. Error Handling (internal/errors/errors.go)

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
- `NewBadRequest()` - 400 errors (validation failures)
- `NewInternalServerError()` - 500 errors (database/server issues)
- `NewConflict()` - 409 errors (reserved for future use)

**Response Functions:**
- `WriteErrorResponse()` - Marshal error to JSON and write to response
- `WriteJSONResponse()` - Marshal data to JSON with status code

**Error Response Format:**
```json
{
  "error": "ERROR_TYPE",
  "message": "Description of what went wrong"
}
```

### 7. Database Migrations (migrations/01_create_todos_table.sql)

**Purpose:** Initialize database schema

**Contents:**
- Creates `todos` table with UUID primary key
- Defines columns: id, title, description, completed, created_at, updated_at
- Creates indexes on common query fields

**Execution:**
```bash
psql -U postgres -d todo_db -f migrations/01_create_todos_table.sql
# Or via docker-compose (automatic on PostgreSQL startup)
```

## Data Flow

### GET /api/todos (List)

```
Request
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
Response (200 OK)
```

### POST /api/todos (Create)

```
Request with JSON body
    ↓
CreateTodo handler
    ↓
Unmarshal JSON → CreateTodoRequest
    ↓
Validate (title not empty)
    ↓
Generate UUID and timestamps
    ↓
INSERT INTO todos
    ↓
Scan result into Todo struct
    ↓
Marshal to JSON
    ↓
Response (201 Created)
```

### PUT /api/todos/:id (Update)

```
Request with JSON body and ID
    ↓
UpdateTodo handler
    ↓
Parse UUID from path param
    ↓
Unmarshal JSON → UpdateTodoRequest
    ↓
Check todo exists (SELECT id FROM todos WHERE id = ?)
    ↓
UPDATE todos (COALESCE for partial updates)
    ↓
Scan result into Todo struct
    ↓
Marshal to JSON
    ↓
Response (200 OK)
```

### DELETE /api/todos/:id

```
Request with ID
    ↓
DeleteTodo handler
    ↓
Parse UUID from path param
    ↓
Check todo exists
    ↓
DELETE FROM todos WHERE id = ?
    ↓
Response (204 No Content)
```

## Technology Details

### FastHTTP vs net/http

**FastHTTP Advantages:**
- 10x faster throughput
- Zero allocations in hot paths
- Efficient memory pooling
- Better for high-concurrency scenarios
- Simpler API for request/response handling

**Trade-offs:**
- Smaller ecosystem (fewer libraries)
- Less standard middleware available
- Requires custom middleware implementation

### pgx Driver Benefits

- **Performance:** Built for speed with prepared statement caching
- **Safety:** Prevents SQL injection with parameterized queries
- **Features:** Supports UUIDs, JSON, arrays natively
- **Connection Pooling:** Built-in connection pool management
- **Modern API:** Context-based cancellation and timeout support

### UUID Generation

- Uses `github.com/google/uuid` package
- Generates v4 (random) UUIDs
- Database generates UUIDs as default (gen_random_uuid())
- Application can also generate UUIDs

### Time Handling

- All timestamps use `time.Time` with UTC timezone
- Database stores as TIMESTAMPTZ (timezone-aware)
- JSON marshaling uses RFC3339 format
- Automatic created_at/updated_at management

## Deployment Considerations

### Local Development
```bash
# With Go installed
go mod download
go run ./cmd/main.go

# With make
make run
```

### Docker Deployment
```bash
# Single container
docker build -t todo-app-go .
docker run -p 8080:8080 -e DATABASE_URL=... todo-app-go

# With docker-compose
docker-compose up
```

### Performance Tuning

1. **Connection Pool Size:** Adjust pgx pool max connections
2. **Timeouts:** Add context timeouts to database operations
3. **Indexes:** Consider additional indexes for large datasets
4. **Caching:** Implement response caching for frequently accessed todos
5. **Monitoring:** Add metrics collection (request duration, errors)

## Security Considerations

1. **SQL Injection:** Prevented by parameterized queries in pgx
2. **Input Validation:** Title requirement, UUID format validation
3. **CORS:** Currently allows all origins (configure for production)
4. **Rate Limiting:** Not implemented (add if needed)
5. **Authentication:** Not implemented (add for production)
6. **HTTPS:** Implement via reverse proxy (nginx, etc.)

## Future Enhancements

1. **Testing:** Add unit and integration tests
2. **Logging:** Structured logging with log levels
3. **Validation:** More comprehensive input validation
4. **Filtering:** Add query parameters for filtering/sorting
5. **Pagination:** Implement limit/offset pagination
6. **Caching:** Add Redis caching layer
7. **Metrics:** Prometheus/Grafana monitoring
8. **Authentication:** JWT or OAuth2 integration
9. **Documentation:** OpenAPI/Swagger specs
10. **Graceful Degradation:** Circuit breakers, retries

## Comparison with Rust Actix

| Aspect | Go FastHTTP | Rust Actix |
|--------|-----------|-----------|
| Compile Time | Seconds | Minutes |
| Binary Size | ~15 MB | ~20 MB |
| Deployment | Simple | Simple |
| Code Verbosity | Low | Moderate |
| Type Safety | Runtime | Compile-time |
| Testing Framework | Standard | Multiple options |
| Learning Curve | Gentle | Steep |
| Performance | Excellent | Excellent |
| Community Size | Large | Growing |

Both are production-ready and provide excellent performance for this use case.

## Summary

This Go FastHTTP implementation mirrors the Rust Actix API with:

- **Same Endpoints:** Identical REST interface
- **Same Database:** PostgreSQL with same schema
- **Better Performance:** 10x faster than standard net/http
- **Simpler Deployment:** Single binary, minimal dependencies
- **Easier Learning:** Go's simplicity vs Rust's complexity

The architecture is clean, maintainable, and follows Go idioms for package organization and error handling.
