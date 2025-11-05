# Project Structure Overview

## Directory Layout

```
todo-app/
├── src/                          # Rust source code
│   ├── main.rs                  # Application entry point & server setup
│   ├── db/
│   │   └── mod.rs               # Database connection pool management
│   ├── models/
│   │   ├── mod.rs               # Models module exports
│   │   └── todo.rs              # Todo struct, DTOs, and serialization
│   ├── handlers/
│   │   ├── mod.rs               # Handlers module exports
│   │   └── todo.rs              # CRUD operation handlers (list, get, create, update, delete)
│   ├── routes/
│   │   └── mod.rs               # Route configuration and setup
│   └── error/
│       └── mod.rs               # Custom error types and responses
│
├── migrations/                    # Database schema migrations
│   └── 01_create_todos_table.sql # Initial schema creation
│
├── Dockerfile                     # Multi-stage Docker build configuration
├── docker-compose.yml             # Docker Compose setup for PostgreSQL + App
├── .dockerignore                  # Docker build ignore patterns
│
├── Cargo.toml                     # Rust project manifest and dependencies
├── Cargo.lock                     # Locked dependency versions
│
├── .env                           # Environment variables (local setup)
├── .gitignore                     # Git ignore patterns
│
├── README.md                      # Full project documentation
├── QUICKSTART.md                  # Quick start guide
└── PROJECT_STRUCTURE.md           # This file
```

## File Descriptions

### Core Application Files

#### `src/main.rs`
- Application entry point
- Sets up Actix Web HTTP server
- Initializes database connection pool
- Configures middleware (logging)
- Registers routes
- Loads environment variables

#### `src/db/mod.rs`
- Database connection pool creation
- PostgreSQL connection management
- Connection pool configuration (max connections, etc.)

#### `src/models/todo.rs`
- `Todo`: Main database model with UUID, title, description, status, timestamps
- `TodoResponse`: DTO for API responses (same as Todo but organized separately)
- `CreateTodoRequest`: DTO for creating new todos
- `UpdateTodoRequest`: DTO for partial updates
- Serialization/deserialization with serde
- SQL row mapping with sqlx

#### `src/models/mod.rs`
- Module exports for models

#### `src/handlers/todo.rs`
- `list_todos()`: GET /api/todos - List all todos
- `get_todo()`: GET /api/todos/{id} - Get single todo
- `create_todo()`: POST /api/todos - Create new todo
- `update_todo()`: PUT /api/todos/{id} - Update todo
- `delete_todo()`: DELETE /api/todos/{id} - Delete todo
- Input validation
- Database query execution
- Error handling

#### `src/handlers/mod.rs`
- Module exports for handlers

#### `src/routes/mod.rs`
- Route configuration function
- Endpoint definitions
- Path parameter mapping
- HTTP method routing

#### `src/error/mod.rs`
- `ApiError` enum: Custom error types
  - NotFound (404)
  - BadRequest (400)
  - InternalServerError (500)
  - Conflict (409)
- `ErrorResponse` struct: JSON error response format
- `ResponseError` trait implementation
- Error to HTTP response conversion
- Database error mapping

### Configuration Files

#### `Cargo.toml`
Core dependencies:
- `actix-web`: Web framework
- `tokio`: Async runtime
- `serde`: JSON serialization
- `sqlx`: SQL database toolkit
- `uuid`: UUID generation
- `chrono`: DateTime handling
- `dotenv`: Environment loading
- `log/env_logger`: Logging

#### `.env`
- `DATABASE_URL`: PostgreSQL connection string
- `PORT`: Server port (default: 8080)
- `RUST_LOG`: Logging level (default: debug)

#### `Dockerfile`
- Two-stage build for minimal image size
- Builder stage: Compiles Rust application
- Runtime stage: Minimal Debian image with runtime dependencies
- Exposes port 8080

#### `docker-compose.yml`
- PostgreSQL 18-alpine service with:
  - Persistent volume for data
  - Health check configuration
  - Auto-initialization with migration SQL
- Rust application service with:
  - Build configuration
  - Environment variables
  - Dependency on PostgreSQL health
  - Network connectivity

### Documentation Files

#### `README.md`
- Project overview and features
- Local setup instructions (with Rust)
- Docker setup instructions (with Docker Compose)
- Complete API endpoint documentation with examples
- Database schema explanation
- Dependency list
- Troubleshooting guide
- Future enhancements

#### `QUICKSTART.md`
- Fast setup guide (Docker vs Local)
- Quick testing instructions with curl
- Common operations
- Environment variables overview

#### `PROJECT_STRUCTURE.md`
- This file
- Directory layout
- File descriptions
- Technology stack overview

### Database Files

#### `migrations/01_create_todos_table.sql`
SQL schema creation:
```sql
CREATE TABLE todos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Performance indexes
CREATE INDEX idx_completed ON todos(completed);
CREATE INDEX idx_created_at ON todos(created_at DESC);
```

## Technology Stack

### Backend Framework
- **Actix Web 4.x**: High-performance async web framework
- **Tokio**: Asynchronous runtime for concurrent operations

### Database
- **PostgreSQL 18-alpine**: Relational database with UUID support
- **SQLx 0.7**: Async SQL toolkit with compile-time checking

### Serialization
- **Serde**: Serialization/deserialization framework
- **Serde JSON**: JSON support

### Utilities
- **UUID**: Unique identifier generation
- **Chrono**: DateTime handling with timezone support
- **DotEnv**: Environment variable management
- **Log/Env Logger**: Structured logging

### Containerization
- **Docker**: Container runtime
- **Docker Compose**: Multi-container orchestration

## API Endpoints Summary

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/todos` | List all todos |
| POST | `/api/todos` | Create new todo |
| GET | `/api/todos/{id}` | Get single todo |
| PUT | `/api/todos/{id}` | Update todo |
| DELETE | `/api/todos/{id}` | Delete todo |

## Data Model

### Todo
```
{
  id: UUID (auto-generated),
  title: String (required),
  description: String? (optional),
  completed: Boolean (default: false),
  created_at: DateTime (auto-set),
  updated_at: DateTime (auto-updated)
}
```

## Development Workflow

### Local Development
1. Install Rust and PostgreSQL
2. Set up `.env` with database credentials
3. Run migrations
4. `cargo run` for development (auto-reload not included)
5. `cargo test` for testing

### Docker Development
1. `docker-compose up --build`
2. Changes require rebuild: `docker-compose up --build`

### Production Build
1. `cargo build --release` (outputs to `target/release/todo-app`)
2. Docker: `docker build -t todo-app .` then `docker run`

## Key Features

✅ Type-safe Rust implementation
✅ Async/await concurrency
✅ SQL compile-time checking
✅ Custom error handling
✅ Request logging
✅ UUID identifiers
✅ Timestamp support
✅ Docker containerization
✅ PostgreSQL persistence
✅ Health checks
✅ Database migrations
✅ Environment configuration

## Next Steps for Enhancement

- [ ] Add authentication (JWT/OAuth)
- [ ] Implement pagination
- [ ] Add request validation
- [ ] Write integration tests
- [ ] Add CORS middleware
- [ ] Implement WebSocket support
- [ ] Add user-based filtering
- [ ] Create API documentation (Swagger/OpenAPI)
- [ ] Add database query caching
- [ ] Implement rate limiting
