# Todo API - Go FastHTTP Implementation

A high-performance RESTful API for managing todo items, built with **Go** and **FastHTTP**. This is a parallel implementation of the Rust Actix API, providing the same functionality with Go's simplicity and FastHTTP's blazing-fast performance.

## Technology Stack

**Core Framework:**
- **FastHTTP 1.67.0** - Ultra-fast HTTP server library for Go (10x faster than net/http)
- **FastHTTPRouter** - High-performance HTTP router for FastHTTP
- **pgx/v5** - PostgreSQL driver with connection pooling

**Database:**
- **PostgreSQL 18-alpine** - Relational database with UUID support

**Utilities:**
- **Google UUID** - UUID generation (v4)
- **Go 1.25** - Latest stable Go version

## Project Structure

```
api-go/
├── cmd/
│   └── main.go                        # Application entry point
├── internal/
│   ├── db/
│   │   └── db.go                      # Database connection pool management
│   ├── handlers/
│   │   └── todo.go                    # CRUD operation handlers
│   ├── models/
│   │   └── todo.go                    # Data models and DTOs
│   ├── routes/
│   │   └── routes.go                  # Route configuration and middleware
│   └── errors/
│       └── errors.go                  # Custom error handling
├── migrations/
│   └── 01_create_todos_table.sql      # Database schema migration
├── go.mod                             # Go module definition
├── go.sum                             # Dependency checksums
├── .env                               # Environment variables
├── Dockerfile                         # Multi-stage Docker build
├── docker-compose.yml                 # PostgreSQL + app orchestration
└── README.md                          # This file
```

## API Endpoints

All endpoints are prefixed with `/api/todos`:

| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| **GET** | `/api/todos` | List all todos (ordered by creation date DESC) | 200 OK |
| **POST** | `/api/todos` | Create a new todo | 201 Created |
| **GET** | `/api/todos/{id}` | Retrieve a specific todo by UUID | 200 OK |
| **PUT** | `/api/todos/{id}` | Update a todo (partial updates supported) | 200 OK |
| **DELETE** | `/api/todos/{id}` | Delete a todo | 204 No Content |

## Data Model

### Database Schema

```sql
CREATE TABLE todos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### Data Structures

**Todo Model:**
- `id: UUID` - Unique identifier (auto-generated)
- `title: string` - Todo title (required)
- `description: *string` - Optional description
- `completed: bool` - Completion status (default: false)
- `created_at: time.Time` - Creation timestamp (UTC)
- `updated_at: time.Time` - Last update timestamp (UTC)

**CreateTodoRequest:**
- `title: string` - Required title
- `description: *string` - Optional description

**UpdateTodoRequest:**
- `title: *string` - Optional title update
- `description: *string` - Optional description update
- `completed: *bool` - Optional completion status

## Installation & Setup

### Prerequisites

- Go 1.25 or higher
- PostgreSQL 18 or higher
- Docker and Docker Compose (for containerized setup)

### Local Development

1. **Clone the repository:**
```bash
cd api-go
```

2. **Install dependencies:**
```bash
go mod download
```

3. **Set up environment variables:**
```bash
cp .env .env.local
# Edit .env.local with your PostgreSQL connection details
```

4. **Start PostgreSQL:**
```bash
# Using Docker
docker run --name postgres-todo -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=todo_db -p 5432:5432 postgres:18-alpine

# Or use your local PostgreSQL installation
```

5. **Run database migrations:**
```bash
# Install migrate tool (optional, can also run SQL directly)
# Or manually execute migrations/01_create_todos_table.sql in PostgreSQL
psql -U postgres -d todo_db -f migrations/01_create_todos_table.sql
```

6. **Run the application:**
```bash
go run ./cmd/main.go
```

The API will be available at `http://127.0.0.1:8080`

### Using Docker Compose

1. **Build and start the containers:**
```bash
docker-compose up --build
```

2. **Verify the setup:**
```bash
curl http://localhost:8080/api/todos
```

3. **Stop the containers:**
```bash
docker-compose down
```

## Usage Examples

### List All Todos

```bash
curl -X GET http://localhost:8080/api/todos
```

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Learn Go",
    "description": "Study Go programming language",
    "completed": false,
    "created_at": "2025-11-05T10:30:00Z",
    "updated_at": "2025-11-05T10:30:00Z"
  }
]
```

### Create a Todo

```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Build API in Go",
    "description": "Create a FastHTTP API"
  }'
```

**Response:** (201 Created)
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "title": "Build API in Go",
  "description": "Create a FastHTTP API",
  "completed": false,
  "created_at": "2025-11-05T10:35:00Z",
  "updated_at": "2025-11-05T10:35:00Z"
}
```

### Get a Specific Todo

```bash
curl -X GET http://localhost:8080/api/todos/550e8400-e29b-41d4-a716-446655440001
```

### Update a Todo

```bash
curl -X PUT http://localhost:8080/api/todos/550e8400-e29b-41d4-a716-446655440001 \
  -H "Content-Type: application/json" \
  -d '{
    "completed": true
  }'
```

Partial updates are supported - only include the fields you want to update.

### Delete a Todo

```bash
curl -X DELETE http://localhost:8080/api/todos/550e8400-e29b-41d4-a716-446655440001
```

**Response:** (204 No Content)

## Error Handling

All errors are returned with a consistent JSON format:

```json
{
  "error": "ERROR_TYPE",
  "message": "Description of what went wrong"
}
```

### Error Types

| Error | Status | Description |
|-------|--------|-------------|
| `BAD_REQUEST` | 400 | Invalid input (e.g., empty title, invalid UUID) |
| `NOT_FOUND` | 404 | Resource doesn't exist |
| `INTERNAL_SERVER_ERROR` | 500 | Database errors or unexpected failures |
| `CONFLICT` | 409 | Reserved for future use |

### Example Error Response

```json
{
  "error": "NOT_FOUND",
  "message": "Todo not found"
}
```

## Configuration

### Environment Variables

- `DATABASE_URL` - PostgreSQL connection string (required)
  - Format: `postgres://user:password@host:port/database`
  - Example: `postgres://postgres:password@localhost:5432/todo_db`

- `PORT` - Server listening port (optional, defaults to 8080)

- `LOG_LEVEL` - Logging level (optional, for future use)

All environment variables can be set in `.env` or passed to the application.

## Key Features

- **Ultra-Fast Performance**: FastHTTP provides 10x faster performance compared to net/http
- **Connection Pooling**: Efficient PostgreSQL connection management with pgx
- **UUID Identifiers**: Type-safe UUID generation and handling
- **Timestamp Management**: Automatic created_at and updated_at with UTC timezone
- **CORS Support**: Configured to handle cross-origin requests
- **Request Logging**: Built-in request logging for debugging
- **Graceful Shutdown**: Proper signal handling for clean shutdown
- **Docker Ready**: Multi-stage Docker builds and Docker Compose orchestration
- **Error Handling**: Structured error responses with appropriate HTTP status codes

## Building & Deployment

### Build Docker Image

```bash
docker build -t todo-app-go:latest .
```

### Run with Docker

```bash
docker run -p 8080:8080 \
  -e DATABASE_URL="postgres://postgres:password@host:5432/todo_db" \
  todo-app-go:latest
```

### Build Native Binary

```bash
go build -o todo-app ./cmd/main.go
./todo-app
```

### Cross-Compilation

```bash
# Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o todo-app ./cmd/main.go

# macOS
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o todo-app ./cmd/main.go

# Windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o todo-app.exe ./cmd/main.go
```

## Performance Characteristics

FastHTTP provides significant performance benefits over standard net/http:

- **Memory Efficiency**: Zero memory allocations in hot paths
- **Request Handling**: Up to 10x faster than net/http for high-throughput scenarios
- **Connection Management**: Efficient connection pooling and reuse
- **Concurrent Requests**: Optimized for handling thousands of concurrent connections

Benchmark comparison against Rust Actix:
- Both implementations share the same database layer
- Go FastHTTP provides comparable throughput with simpler deployment

## Development

### Running Tests

```bash
go test ./...
```

### Code Organization

- `cmd/` - Application entry point and main logic
- `internal/` - Private packages
  - `db/` - Database connection management
  - `handlers/` - HTTP request handlers
  - `models/` - Data structures and DTOs
  - `routes/` - Route configuration and middleware
  - `errors/` - Error types and handling utilities

### Adding New Endpoints

1. Create handler in `internal/handlers/`
2. Register route in `internal/routes/routes.go`
3. Add corresponding models in `internal/models/` if needed

## Comparison with Rust Actix Version

| Feature | Rust Actix | Go FastHTTP |
|---------|-----------|-----------|
| Language | Rust | Go |
| Performance | Excellent | Excellent (10x net/http) |
| Type Safety | Compile-time | Runtime |
| Learning Curve | Steep | Gentle |
| Deployment | Single binary | Single binary |
| Development Speed | Moderate | Fast |
| Memory Usage | Low | Low |

Both implementations:
- Share identical API endpoints
- Use PostgreSQL with UUID identifiers
- Support CORS
- Include proper error handling
- Include request logging
- Are production-ready

## Troubleshooting

### Connection Refused
- Ensure PostgreSQL is running
- Check DATABASE_URL is correct
- Verify PostgreSQL port is accessible

### Port Already in Use
```bash
# Change PORT environment variable
export PORT=8081
go run ./cmd/main.go
```

### Database Connection Issues
```bash
# Test connection with psql
psql postgres://postgres:password@localhost:5432/todo_db
```

### Migration Issues
- Ensure migrations directory is accessible
- Check PostgreSQL permissions
- Run migrations manually if needed

## License

MIT License

## Additional Resources

- [FastHTTP Documentation](https://github.com/valyala/fasthttp)
- [pgx Documentation](https://github.com/jackc/pgx)
- [Go Official Documentation](https://golang.org/doc)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
