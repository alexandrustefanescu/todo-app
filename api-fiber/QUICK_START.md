# Quick Start Guide - Fiber Todo API

## Prerequisites

- Go 1.25 or higher
- PostgreSQL 18 or higher
- Docker & Docker Compose (for containerized setup)

## Option 1: Quick Start with Docker Compose (Recommended)

The fastest way to get everything running:

```bash
cd api-fiber
docker-compose up --build
```

Your API will be available at: `http://localhost:8080/api/todos`

To stop: `Ctrl+C` or run `docker-compose down`

## Option 2: Local Development

### 1. Start PostgreSQL

```bash
# Using Docker
docker run -d \
  --name postgres-todo \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=todo_db \
  -p 5432:5432 \
  postgres:18-alpine

# Or use your local PostgreSQL installation
```

### 2. Set Up Environment

```bash
cd api-fiber
cp .env .env.local
# Edit .env.local if needed (default values work with the above setup)
```

### 3. Install Dependencies

```bash
go mod download
go mod tidy
```

### 4. Run Migrations

```bash
psql postgres://postgres:password@localhost:5432/todo_db < migrations/01_create_todos_table.sql
```

### 5. Run the Application

```bash
go run ./cmd/main.go
```

The API will start on `http://127.0.0.1:8080`

## Quick API Test

```bash
# List all todos
curl http://localhost:8080/api/todos

# Create a todo
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"My First Fiber Todo","description":"Learn Fiber"}'

# Get a todo (replace UUID with real ID)
curl http://localhost:8080/api/todos/550e8400-e29b-41d4-a716-446655440000

# Update a todo
curl -X PUT http://localhost:8080/api/todos/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{"completed":true}'

# Delete a todo
curl -X DELETE http://localhost:8080/api/todos/550e8400-e29b-41d4-a716-446655440000
```

## Using Make

If you have `make` installed, useful commands:

```bash
make help              # Show all available targets
make build             # Build the binary
make run               # Build and run
make docker-up         # Start with docker-compose
make docker-down       # Stop containers
make test              # Run tests (when added)
make clean             # Clean build artifacts
make fmt               # Format code
make lint              # Lint code
```

## Troubleshooting

### Port 8080 Already in Use
```bash
export PORT=8081
go run ./cmd/main.go
```

### Cannot Connect to PostgreSQL
Check the connection string in `.env`:
```bash
# Should match your PostgreSQL setup
DATABASE_URL=postgres://postgres:password@localhost:5432/todo_db
```

Test the connection:
```bash
psql postgres://postgres:password@localhost:5432/todo_db -c "SELECT 1"
```

### Go Module Errors
```bash
# Clear cache and re-download
go clean -modcache
go mod download
```

## Why Fiber?

Fiber is an Express.js-inspired framework built on FastHTTP:

✅ **Express-like API** - Familiar if you know Express.js or Node.js
✅ **High Performance** - 10x faster than Go's standard net/http
✅ **Rich Middleware** - Built-in logging, compression, etc.
✅ **Developer-friendly** - Minimal boilerplate, clean syntax
✅ **Active Community** - Growing ecosystem with more resources

## Next Steps

1. **Read the [README.md](README.md)** for comprehensive documentation
2. **Review [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md)** for architecture details
3. **Explore the code** in `internal/` directories
4. **Add tests** to `*_test.go` files
5. **Implement features** like filtering, pagination, authentication

## Project Structure

```
api-fiber/
├── cmd/main.go                    # Server startup
├── internal/
│   ├── db/db.go                   # Database connection
│   ├── handlers/todo.go           # HTTP handlers
│   ├── models/todo.go             # Data types
│   ├── routes/routes.go           # Route definitions
│   ├── middleware/cors.go         # CORS middleware
│   └── errors/errors.go           # Error handling
├── migrations/                    # Database schemas
├── Dockerfile                     # Container build
├── docker-compose.yml             # Service orchestration
├── Makefile                       # Build tasks
└── README.md                      # Full documentation
```

## API Comparison

All three implementations provide identical functionality:

| Implementation | Framework | Language | Performance |
|---|---|---|---|
| **api/** | Actix | Rust | Excellent |
| **api-go/** | FastHTTP Router | Go | Excellent |
| **api-fiber/** | Fiber | Go | Excellent |

**Choose Fiber if:**
- You're familiar with Express.js
- You want the easiest API to learn
- You prefer a rich middleware ecosystem
- You want the best of both worlds (fasthttp + nice API)

## Support

For more information:
- [Fiber Documentation](https://docs.gofiber.io/)
- [Fiber GitHub](https://github.com/gofiber/fiber)
- [pgx Documentation](https://github.com/jackc/pgx)
- [Go Official Docs](https://golang.org/doc)
- [PostgreSQL Docs](https://www.postgresql.org/docs/)

---

Happy coding! 🚀
