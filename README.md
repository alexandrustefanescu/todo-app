# TODO App - Full Stack Application

A complete TODO management application with a Rust REST API backend and a separate frontend client.

## Project Structure

```
todo-app/
├── api/                          # Rust Actix Web REST API
│   ├── src/
│   │   ├── main.rs              # Server entry point
│   │   ├── db/                  # Database connection
│   │   ├── models/              # Data models
│   │   ├── handlers/            # API endpoint handlers
│   │   ├── routes/              # Route configuration
│   │   └── error/               # Error handling
│   ├── migrations/              # Database migrations
│   ├── docker-compose.yml       # PostgreSQL setup
│   ├── Dockerfile               # API container
│   ├── Cargo.toml               # Rust dependencies
│   ├── .env                     # Environment variables
│   ├── README.md                # API documentation
│   ├── QUICKSTART.md            # Quick start guide
│   └── PROJECT_STRUCTURE.md     # Architecture details
│
├── client/                       # Frontend application
│   └── (To be built)
│
├── .gitignore                   # Git ignore rules
└── README.md                    # This file
```

## Quick Start

### Running the API

Navigate to the `api` folder:

```bash
cd api

# 1. Start PostgreSQL with Docker Compose
docker-compose up

# 2. In another terminal, run migrations
sqlx migrate run

# 3. Start the Rust application
cargo run --release
```

The API will be available at `http://localhost:8080`

### Testing the API

```bash
# Create a todo
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Rust", "description": "Study Rust programming"}'

# List todos
curl http://localhost:8080/api/todos

# Get a specific todo (replace {id} with actual ID)
curl http://localhost:8080/api/todos/{id}

# Update a todo
curl -X PUT http://localhost:8080/api/todos/{id} \
  -H "Content-Type: application/json" \
  -d '{"completed": true}'

# Delete a todo
curl -X DELETE http://localhost:8080/api/todos/{id}
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/todos` | List all todos |
| POST | `/api/todos` | Create new todo |
| GET | `/api/todos/{id}` | Get specific todo |
| PUT | `/api/todos/{id}` | Update todo |
| DELETE | `/api/todos/{id}` | Delete todo |

## Technology Stack

### Backend (API)
- **Framework**: Actix Web 4.4
- **Runtime**: Tokio 1.35
- **Database**: PostgreSQL 18-alpine (Docker)
- **ORM/Query**: SQLx 0.7
- **Serialization**: Serde 1.0

### Frontend (Client)
- Coming soon...

## Prerequisites

- **Rust 1.70+** - Install from https://rustup.rs/
- **Docker & Docker Compose** - Install from https://www.docker.com/
- **sqlx-cli** - Install with: `cargo install sqlx-cli --no-default-features --features postgres`

## Documentation

- **API Documentation**: See [api/README.md](api/README.md)
- **Quick Start Guide**: See [api/QUICKSTART.md](api/QUICKSTART.md)
- **Architecture**: See [api/PROJECT_STRUCTURE.md](api/PROJECT_STRUCTURE.md)

## Development Workflow

### 1. Start the Database
```bash
cd api
docker-compose up
```

### 2. Run Migrations
```bash
cd api
sqlx migrate run
```

### 3. Start the API Server
```bash
cd api
cargo run --release
```

### 4. Test with curl or Thunder Client

## Common Commands

```bash
# Building
cargo build --release

# Running
cargo run --release

# Testing
cargo test

# Database operations
sqlx migrate run          # Apply migrations
sqlx migrate revert       # Revert last migration

# Docker operations
docker-compose up         # Start services
docker-compose down       # Stop services
docker-compose down -v    # Stop and remove volumes
```

## Folder Organization

- **api/** - Contains all backend/API code
  - Follow the structure and documentation in `api/README.md`
  - Database configurations are in `api/.env`
  - API runs on port 8080

- **client/** - Designated for frontend code
  - To be implemented (React, Vue, Angular, etc.)
  - Will consume the API from `http://localhost:8080/api`

## API Features

✅ Full CRUD operations for todos
✅ UUID identifiers
✅ Timestamp tracking (created_at, updated_at)
✅ Completion status tracking
✅ Type-safe database queries
✅ Comprehensive error handling
✅ Request logging
✅ Docker containerization

## Environment Variables

API environment variables are in `api/.env`:
```
DATABASE_URL=postgres://postgres:password@localhost:5432/todo_db
PORT=8080
RUST_LOG=debug
```

## Database

PostgreSQL 18-alpine running in Docker:
- **Host**: localhost:5432
- **Username**: postgres
- **Password**: password
- **Database**: todo_db

## Troubleshooting

### Database Connection Issues
- Ensure `docker-compose up` is running
- Check DATABASE_URL in `api/.env`
- Verify PostgreSQL container: `docker-compose ps`

### API Won't Start
- Ensure database is running and migrations applied
- Check port 8080 is not in use
- Verify DATABASE_URL is correct

### Migration Issues
- Clean and rebuild: `docker-compose down -v && docker-compose up`
- Run migrations: `sqlx migrate run`

## Next Steps

1. ✅ API is set up and ready
2. 📝 Create frontend in `client/` folder
3. 🔗 Connect frontend to API endpoints
4. 🚀 Deploy both services

## License

This project is open source and available under the MIT License.

## Support

For issues or questions about the API, see [api/README.md](api/README.md)
