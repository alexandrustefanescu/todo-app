# Quick Start Guide - Go FastHTTP Todo API

## Prerequisites

- Go 1.25 or higher
- PostgreSQL 18 or higher
- Docker & Docker Compose (for containerized setup)

## Option 1: Quick Start with Docker Compose (Recommended)

The fastest way to get everything running:

```bash
cd api-go
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
cd api-go
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
  -d '{"title":"My First Todo","description":"Learn Go"}'

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

## API Comparison with Rust Actix Version

Both APIs have:
- ✅ Identical endpoints (`/api/todos`, `/api/todos/:id`)
- ✅ Same database schema
- ✅ Same request/response formats
- ✅ Full CRUD operations
- ✅ CORS support
- ✅ Proper error handling

**Main Differences:**
- **Go FastHTTP**: Simpler, faster to develop, 10x faster than net/http
- **Rust Actix**: Type-safe at compile-time, steeper learning curve

## Next Steps

1. **Read the [README.md](README.md)** for comprehensive documentation
2. **Review [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md)** for architecture details
3. **Explore the code** in `internal/` directories
4. **Add tests** to `*_test.go` files
5. **Implement features** like filtering, pagination, authentication

## Project Structure

```
api-go/
├── cmd/main.go                    # Server startup
├── internal/
│   ├── db/db.go                   # Database connection
│   ├── handlers/todo.go           # HTTP handlers
│   ├── models/todo.go             # Data types
│   ├── routes/routes.go           # Route definitions
│   └── errors/errors.go           # Error handling
├── migrations/                    # Database schemas
├── Dockerfile                     # Container build
├── docker-compose.yml             # Service orchestration
├── Makefile                       # Build tasks
└── README.md                      # Full documentation
```

## Performance Notes

FastHTTP provides excellent performance:
- **10x faster** than standard Go net/http
- **Zero allocations** in hot paths
- **Efficient pooling** for connections and memory

This makes it ideal for high-throughput APIs like this todo manager.

## Support

For more information:
- [FastHTTP GitHub](https://github.com/valyala/fasthttp)
- [pgx Documentation](https://github.com/jackc/pgx)
- [Go Official Docs](https://golang.org/doc)
- [PostgreSQL Docs](https://www.postgresql.org/docs/)

---

Happy coding! 🚀
