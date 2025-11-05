# TODO App - Rust + Actix Web API

A RESTful API for managing TODO items built with Rust, Actix Web framework, and PostgreSQL.

## Features

- **CRUD Operations**: Create, read, update, and delete todos
- **Async/Await**: Full async implementation with Tokio runtime
- **Type Safety**: Leveraging Rust's type system for safety
- **Error Handling**: Comprehensive error handling with custom error types
- **Logging**: Built-in logging with env_logger
- **Database**: PostgreSQL with SQLx for compile-time query checking

## Prerequisites

### Option 1: Local Setup
- Rust 1.70+ (Install from https://rustup.rs/)
- PostgreSQL 12+ (https://www.postgresql.org/download/)
- Cargo (comes with Rust)

### Option 2: Docker Setup (Recommended)
- Docker (https://www.docker.com/products/docker-desktop)
- Docker Compose

## Setup

### Option 1: Using Docker Compose for PostgreSQL (Recommended)

The easiest way to get started is using Docker Compose to run PostgreSQL, then run the Rust application locally.

#### Prerequisites:
- Docker and Docker Compose installed
- Rust 1.70+ installed

#### Steps:
```bash
# Start PostgreSQL with Docker Compose
docker-compose up

# In another terminal, navigate to the project and run:
cd /Users/alexandrustefanescu/Desktop/todo-app

# Install sqlx-cli if not already installed
cargo install sqlx-cli --no-default-features --features postgres

# Run migrations (creates the todos table)
sqlx migrate run

# Run the Rust application
cargo run --release
```

PostgreSQL will be available at `localhost:5432`
The application will be available at `http://localhost:8080`

To stop the PostgreSQL container:
```bash
docker-compose down
```

To stop and remove all data:
```bash
docker-compose down -v
```

### Option 2: Local Setup

#### 1. Clone and Navigate to Project
```bash
cd /Users/alexandrustefanescu/Desktop/todo-app
```

#### 2. Create PostgreSQL Database
```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE todo_db;

# Exit psql
\q
```

#### 3. Configure Environment Variables
Update the `.env` file with your PostgreSQL credentials:
```
DATABASE_URL=postgres://postgres:password@localhost:5432/todo_db
PORT=8080
RUST_LOG=debug
```

#### 4. Run Database Migrations
```bash
# Install sqlx-cli if not already installed
cargo install sqlx-cli --no-default-features --features postgres

# Run migrations
sqlx migrate run
```

#### 5. Build and Run
```bash
# Build the project
cargo build --release

# Run the application
cargo run
```

The server will start at `http://127.0.0.1:8080`

## API Endpoints

### List All Todos
```
GET /api/todos
```

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Learn Rust",
    "description": "Study Rust programming language",
    "completed": false,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
]
```

### Get Single Todo
```
GET /api/todos/{id}
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Learn Rust",
  "description": "Study Rust programming language",
  "completed": false,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Create Todo
```
POST /api/todos
Content-Type: application/json

{
  "title": "Learn Rust",
  "description": "Study Rust programming language"
}
```

**Response:** `201 Created`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Learn Rust",
  "description": "Study Rust programming language",
  "completed": false,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Update Todo
```
PUT /api/todos/{id}
Content-Type: application/json

{
  "title": "Learn Rust Advanced",
  "completed": true
}
```

**Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Learn Rust Advanced",
  "description": "Study Rust programming language",
  "completed": true,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T11:45:00Z"
}
```

### Delete Todo
```
DELETE /api/todos/{id}
```

**Response:** `204 No Content`

## Error Responses

### Bad Request (400)
```json
{
  "error": "BAD_REQUEST",
  "message": "Title cannot be empty"
}
```

### Not Found (404)
```json
{
  "error": "NOT_FOUND",
  "message": "Todo with id {id} not found"
}
```

### Internal Server Error (500)
```json
{
  "error": "INTERNAL_SERVER_ERROR",
  "message": "Database error: ..."
}
```

## Testing with curl

### Create a todo
```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Buy groceries", "description": "Milk, eggs, bread"}'
```

### List todos
```bash
curl http://localhost:8080/api/todos
```

### Get specific todo
```bash
curl http://localhost:8080/api/todos/{id}
```

### Update todo
```bash
curl -X PUT http://localhost:8080/api/todos/{id} \
  -H "Content-Type: application/json" \
  -d '{"completed": true}'
```

### Delete todo
```bash
curl -X DELETE http://localhost:8080/api/todos/{id}
```

## Project Structure

```
todo-app/
├── src/
│   ├── main.rs           # Application entry point
│   ├── db/
│   │   └── mod.rs        # Database connection setup
│   ├── models/
│   │   ├── mod.rs        # Models module
│   │   └── todo.rs       # Todo model and DTOs
│   ├── handlers/
│   │   ├── mod.rs        # Handlers module
│   │   └── todo.rs       # Todo CRUD handlers
│   ├── routes/
│   │   └── mod.rs        # Route configuration
│   └── error/
│       └── mod.rs        # Error handling
├── migrations/
│   └── 01_create_todos_table.sql  # Database schema
├── Cargo.toml            # Rust dependencies
├── .env                  # Environment configuration
├── .gitignore            # Git ignore rules
└── README.md             # This file
```

## Development

### Building
```bash
cargo build
```

### Running Tests
```bash
cargo test
```

### Code Quality
```bash
# Format code
cargo fmt

# Check for common mistakes
cargo clippy
```

### Checking Dependencies
```bash
cargo outdated
```

## Database Schema

The `todos` table has the following structure:

```sql
CREATE TABLE todos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_completed ON todos(completed);
CREATE INDEX idx_created_at ON todos(created_at DESC);
```

## Dependencies

- **actix-web**: Web framework
- **tokio**: Async runtime
- **serde**: Serialization/deserialization
- **sqlx**: SQL toolkit with compile-time query checking
- **uuid**: UUID generation
- **chrono**: Date and time handling
- **dotenv**: Environment variable loading
- **log/env_logger**: Logging

## Docker

### Run PostgreSQL with Docker Compose
```bash
# Start PostgreSQL container
docker-compose up

# Run in background
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the container
docker-compose down

# Stop and remove all data
docker-compose down -v
```

### Run Individual PostgreSQL Container
```bash
docker run -d --name postgres-todo \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=todo_db \
  -p 5432:5432 \
  postgres:18-alpine
```

### Docker Compose Configuration

The `docker-compose.yml` includes:
- **postgres**: PostgreSQL 18-alpine with persistent volume and health checks
- Automatic initialization of database schema from `./migrations/`

Database credentials:
- Username: `postgres`
- Password: `password`
- Database: `todo_db`
- Port: `5432`

### Building and Running the Application

To build the Docker image for the Rust application (optional):
```bash
docker build -t todo-app .
```

To run the application in Docker (after building):
```bash
docker run -d --name todo-app \
  -e DATABASE_URL=postgres://postgres:password@host.docker.internal:5432/todo_db \
  -e RUST_LOG=debug \
  -p 8080:8080 \
  todo-app
```

## Future Enhancements

- [x] Docker configuration
- [ ] Authentication and authorization
- [ ] Pagination and filtering
- [ ] Request validation with validators
- [ ] Integration tests
- [ ] Deployment documentation
- [ ] WebSocket support for real-time updates
- [ ] User-based todo separation

## License

This project is open source and available under the MIT License.

## Troubleshooting

### Connection refused
- Ensure PostgreSQL is running
- Check DATABASE_URL in .env file

### Database error
- Verify the database exists
- Run migrations: `sqlx migrate run`

### Port already in use
- Change PORT in .env file
- Or kill the process using port 8080

### Docker Issues

#### Docker Compose won't start
```bash
# Clean up any existing containers and volumes
docker-compose down -v

# Rebuild and start fresh
docker-compose up --build
```

#### PostgreSQL container is not healthy
```bash
# Check container logs
docker-compose logs postgres

# Rebuild the compose setup
docker-compose restart postgres
```

#### Application can't connect to database
- Ensure PostgreSQL container is running: `docker-compose ps`
- Check the DATABASE_URL matches the service name: `postgres` (not localhost)
- Wait for PostgreSQL to be ready (health check may take 10-15 seconds)

#### Clean up Docker resources
```bash
# Stop and remove all containers
docker-compose down

# Remove all containers and volumes
docker-compose down -v

# Remove the Docker image
docker rmi todo-app
```

## Support

For issues or questions, please open an issue on the project repository.
