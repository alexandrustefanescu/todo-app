# Quick Start Guide

Get your TODO API running in minutes!

## Option 1: Docker Compose + Local Rust (Fastest - Recommended)

The fastest way to get started:

```bash
# 1. Navigate to the project
cd /Users/alexandrustefanescu/Desktop/todo-app

# 2. Start PostgreSQL with Docker Compose
docker-compose up

# 3. In another terminal, install sqlx-cli (one time only)
cargo install sqlx-cli --no-default-features --features postgres

# 4. Run migrations to create the database schema
sqlx migrate run

# 5. Start the Rust application
cargo run --release
```

Your API is now running at `http://localhost:8080`
PostgreSQL is running at `localhost:5432`

### Test the API:
```bash
# Create a todo
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Buy groceries", "description": "Milk, eggs, bread"}'

# List todos
curl http://localhost:8080/api/todos
```

### Stop the services:
```bash
# Press Ctrl+C in the terminal running cargo run (stops the API)
# Press Ctrl+C in the terminal running docker-compose up (stops PostgreSQL)
```

## Option 2: Local Setup

If you prefer to run locally without Docker:

```bash
# 1. Prerequisites
# - Install Rust: https://rustup.rs/
# - Install PostgreSQL: https://www.postgresql.org/download/

# 2. Create database
psql -U postgres
CREATE DATABASE todo_db;
\q

# 3. Navigate to project
cd /Users/alexandrustefanescu/Desktop/todo-app

# 4. Install sqlx-cli
cargo install sqlx-cli --no-default-features --features postgres

# 5. Run migrations
sqlx migrate run

# 6. Start the server
cargo run --release
```

## Testing the API

### Using curl

#### Create a Todo
```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Learn Rust",
    "description": "Study the Rust programming language"
  }'
```

#### List All Todos
```bash
curl http://localhost:8080/api/todos
```

#### Get a Specific Todo
Replace `{id}` with the todo ID from the create response:
```bash
curl http://localhost:8080/api/todos/{id}
```

#### Update a Todo
```bash
curl -X PUT http://localhost:8080/api/todos/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "completed": true
  }'
```

#### Delete a Todo
```bash
curl -X DELETE http://localhost:8080/api/todos/{id}
```

### Using Postman

1. Import the following collection into Postman
2. Set base URL to `http://localhost:8080`
3. Create requests for each endpoint shown above

## Project Structure

```
todo-app/
├── src/
│   ├── main.rs          # Application entry point
│   ├── db/mod.rs        # Database connection
│   ├── models/          # Data structures
│   ├── handlers/        # API endpoints
│   ├── routes/mod.rs    # Route configuration
│   └── error/mod.rs     # Error handling
├── migrations/          # Database migrations (SQL schema)
├── Dockerfile           # Docker build configuration (optional)
├── docker-compose.yml   # PostgreSQL 18-alpine setup
├── Cargo.toml           # Rust dependencies
├── .env                 # Environment variables
└── README.md            # Full documentation
```

## What's Included

- ✅ Full REST API with CRUD operations
- ✅ PostgreSQL database (18-alpine in Docker)
- ✅ Docker & Docker Compose setup
- ✅ Type-safe Rust with Actix Web
- ✅ Async/Await runtime with Tokio
- ✅ Comprehensive error handling
- ✅ Request logging
- ✅ UUID support
- ✅ Timestamps with timezone support

## Next Steps

1. **Start the API** using either Docker or local setup above
2. **Test some endpoints** using curl examples above
3. **Read [README.md](README.md)** for full documentation
4. **Explore the code** in the `src/` directory
5. **Customize** the API for your needs

## Environment Variables

Default configuration (works with Docker Compose):
```
DATABASE_URL=postgres://postgres:password@postgres:5432/todo_db
PORT=8080
RUST_LOG=debug
```

For local setup, update these in `.env`:
```
DATABASE_URL=postgres://postgres:password@localhost:5432/todo_db
PORT=8080
RUST_LOG=debug
```

## Troubleshooting

### Docker Issues
See the [Docker Issues](README.md#docker-issues) section in README.md

### Connection Issues
- Ensure the database is running
- Check DATABASE_URL in .env
- Verify port 8080 is available

### Build Issues
- Make sure you have Rust 1.70+ installed
- Run `cargo clean` then `cargo build`

## API Response Examples

### Successful Create (201 Created)
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Buy groceries",
  "description": "Milk, eggs, bread",
  "completed": false,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Error Response (404 Not Found)
```json
{
  "error": "NOT_FOUND",
  "message": "Todo with id {id} not found"
}
```

## More Information

- Read the full [README.md](README.md) for detailed documentation
- Check out [Actix Web documentation](https://actix.rs/)
- Learn more about [Rust](https://www.rust-lang.org/)

Happy coding! 🦀
