# TODO App - Full Stack Application

A complete TODO management application with **five production-ready API implementations** (Rust Actix, Go FastHTTP, Go Fiber, Node.js Fastify, Python FastAPI) and a Vanilla JavaScript frontend client.

## Project Structure

```
todo-app/
├── api/                          # Rust Actix Web REST API
│   ├── src/                      # Rust source code
│   ├── migrations/               # Database migrations
│   ├── docker-compose.yml        # PostgreSQL setup
│   ├── Dockerfile                # API container
│   ├── Cargo.toml                # Rust dependencies (Updated Nov 2025)
│   ├── README.md                 # API documentation
│   ├── QUICK_START.md            # Quick start guide
│   └── PROJECT_STRUCTURE.md      # Architecture details
│
├── api-go/                       # Go FastHTTP Router REST API
│   ├── cmd/main.go               # Application entry point
│   ├── internal/                 # Database, handlers, routes, models
│   ├── migrations/               # Database migrations
│   ├── docker-compose.yml        # PostgreSQL setup
│   ├── Dockerfile                # API container
│   ├── go.mod                    # Go dependencies
│   ├── README.md                 # API documentation
│   ├── QUICK_START.md            # Quick start guide
│   └── IMPLEMENTATION_SUMMARY.md # Architecture details
│
├── api-fiber/                    # Go Fiber Web Framework REST API
│   ├── cmd/main.go               # Application entry point
│   ├── internal/                 # Database, handlers, routes, models
│   ├── migrations/               # Database migrations
│   ├── docker-compose.yml        # PostgreSQL setup
│   ├── Dockerfile                # API container
│   ├── go.mod                    # Go dependencies
│   ├── README.md                 # API documentation
│   ├── QUICK_START.md            # Quick start guide
│   └── IMPLEMENTATION_SUMMARY.md # Architecture details
│
├── api-fastify/                  # Node.js Fastify REST API
│   ├── src/                      # JavaScript source code
│   ├── migrations/               # Database migrations
│   ├── docker-compose.yml        # PostgreSQL setup
│   ├── Dockerfile                # API container
│   ├── package.json              # npm dependencies (Latest Nov 2025)
│   ├── README.md                 # API documentation
│   ├── QUICK_START.md            # Quick start guide
│   └── IMPLEMENTATION_SUMMARY.md # Architecture details
│
├── api-fastapi/                  # Python FastAPI REST API
│   ├── main.py                   # FastAPI application & endpoints
│   ├── models.py                 # SQLAlchemy ORM models
│   ├── schemas.py                # Pydantic validation schemas
│   ├── migrations/               # Database migrations
│   ├── docker-compose.yml        # PostgreSQL setup
│   ├── Dockerfile                # API container
│   ├── requirements.txt           # Python dependencies (Latest Nov 2025)
│   ├── README.md                 # API documentation
│   ├── QUICK_START.md            # Quick start guide
│   └── IMPLEMENTATION_SUMMARY.md # Architecture details
│
├── client/                       # Frontend application (Vanilla JavaScript + Bun)
│   ├── index.html                # Main HTML file
│   ├── app.js                    # Application logic
│   ├── styles.css                # Styling
│   ├── README.md                 # Client documentation
│   ├── BUN_SETUP.md              # Bun runtime setup
│   └── package.json              # Dependencies (Bun)
│
├── API_IMPLEMENTATIONS_GUIDE.md  # Comprehensive comparison of all 4 APIs
├── COMPARISON.md                 # Rust vs Go detailed comparison
├── GETTING_STARTED.md            # Selection guide for choosing an API
├── RUST_VERSION_UPDATE.md        # Details on latest Rust dependency updates
├── .gitignore                    # Git ignore rules
└── README.md                     # This file
```

## Quick Start

### Fastest Way: Using Docker Compose

Choose your preferred API implementation:

**Rust Actix:**
```bash
cd api && docker-compose up --build
```

**Go FastHTTP:**
```bash
cd api-go && docker-compose up --build
```

**Go Fiber:**
```bash
cd api-fiber && docker-compose up --build
```

**Node.js Fastify:**
```bash
cd api-fastify && docker-compose up --build
```

**Python FastAPI:**
```bash
cd api-fastapi && docker-compose up --build
```

The API will be available at `http://localhost:8080/api/todos`

### Local Development Setup

For detailed setup instructions for each implementation, see:
- [api/QUICK_START.md](api/QUICK_START.md) - Rust Actix
- [api-go/QUICK_START.md](api-go/QUICK_START.md) - Go FastHTTP
- [api-fiber/QUICK_START.md](api-fiber/QUICK_START.md) - Go Fiber
- [api-fastify/QUICK_START.md](api-fastify/QUICK_START.md) - Node.js Fastify
- [api-fastapi/QUICK_START.md](api-fastapi/QUICK_START.md) - Python FastAPI

### Running the Client

Navigate to the `client` folder:

```bash
cd client

# 1. Install dependencies with Bun
bun install

# 2. Start the development server
bun run dev
```

The client will be available at `http://localhost:3000` (or the port configured in your setup)

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

## Backend API Implementations

We provide **five production-ready implementations** of the same Todo API with identical endpoints:

| Implementation | Framework | Language | Status | Latest Versions |
|---|---|---|---|---|
| **api/** | Actix-web | Rust 2021 | ✅ Complete | 4.11 actix-web, 1.41 tokio, 0.8 sqlx |
| **api-go/** | FastHTTP Router | Go 1.25 | ✅ Complete | 1.67.0 fasthttp, 5.7.1 pgx |
| **api-fiber/** | Fiber (Express-like) | Go 1.25 | ✅ Complete | 2.52.5 fiber, 5.7.1 pgx |
| **api-fastify/** | Fastify (Express-like) | Node.js 22 LTS | ✅ Complete | 5.6.1 fastify, 8.12.0 pg |
| **api-fastapi/** | FastAPI | Python 3.12 | ✅ Complete | 0.115.4 fastapi, 2.0.36 sqlalchemy, 0.31.0 asyncpg |

### Choosing Your API Implementation

- **Rust Actix** (`api/`) - Best for type safety and compile-time guarantees
- **Go FastHTTP** (`api-go/`) - Best for raw performance and low-level control
- **Go Fiber** (`api-fiber/`) - Best for Express.js developers and balance of DX + performance
- **Node.js Fastify** (`api-fastify/`) - Best for Node.js teams and rapid development
- **Python FastAPI** (`api-fastapi/`) - Best for Python teams and auto API documentation

See [GETTING_STARTED.md](GETTING_STARTED.md) for detailed comparison and selection guide.

### Frontend (Client)
- **Language**: Vanilla JavaScript
- **Runtime**: Bun
- **Styling**: CSS3
- **Build Tool**: Bun

## Prerequisites

- **Rust 1.70+** - Install from https://rustup.rs/
- **Docker & Docker Compose** - Install from https://www.docker.com/
- **sqlx-cli** - Install with: `cargo install sqlx-cli --no-default-features --features postgres`
- **Bun** - Install from https://bun.sh/

## Documentation

### Main Documentation
- **[API_IMPLEMENTATIONS_GUIDE.md](API_IMPLEMENTATIONS_GUIDE.md)** - Comprehensive comparison of all 5 API implementations
- **[GETTING_STARTED.md](GETTING_STARTED.md)** - Selection guide to choose the right implementation for your needs
- **[COMPARISON.md](COMPARISON.md)** - Detailed Rust vs Go comparison
- **[RUST_VERSION_UPDATE.md](RUST_VERSION_UPDATE.md)** - Latest dependency versions for Rust Actix (Nov 2025)

### Per-Implementation Documentation
**Rust Actix (api/):**
- [README.md](api/README.md) - Full API documentation
- [QUICK_START.md](api/QUICK_START.md) - Quick start guide
- [PROJECT_STRUCTURE.md](api/PROJECT_STRUCTURE.md) - Architecture details

**Go FastHTTP (api-go/):**
- [README.md](api-go/README.md) - Full API documentation
- [QUICK_START.md](api-go/QUICK_START.md) - Quick start guide
- [IMPLEMENTATION_SUMMARY.md](api-go/IMPLEMENTATION_SUMMARY.md) - Architecture details

**Go Fiber (api-fiber/):**
- [README.md](api-fiber/README.md) - Full API documentation
- [QUICK_START.md](api-fiber/QUICK_START.md) - Quick start guide
- [IMPLEMENTATION_SUMMARY.md](api-fiber/IMPLEMENTATION_SUMMARY.md) - Architecture details

**Node.js Fastify (api-fastify/):**
- [README.md](api-fastify/README.md) - Full API documentation
- [QUICK_START.md](api-fastify/QUICK_START.md) - Quick start guide
- [IMPLEMENTATION_SUMMARY.md](api-fastify/IMPLEMENTATION_SUMMARY.md) - Architecture details

**Python FastAPI (api-fastapi/):**
- [README.md](api-fastapi/README.md) - Full API documentation
- [QUICK_START.md](api-fastapi/QUICK_START.md) - Quick start guide
- [IMPLEMENTATION_SUMMARY.md](api-fastapi/IMPLEMENTATION_SUMMARY.md) - Architecture details

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

### API Implementations

- **api/** - Rust Actix Web REST API
  - Latest versions: Actix-web 4.11, Tokio 1.41, SQLx 0.8
  - See `api/README.md` for full documentation
  - Database: PostgreSQL 18

- **api-go/** - Go FastHTTP Router REST API
  - Latest versions: FastHTTP 1.67.0, pgx 5.7.1, Go 1.25
  - See `api-go/README.md` for full documentation
  - Database: PostgreSQL 18

- **api-fiber/** - Go Fiber Web Framework REST API
  - Latest versions: Fiber 2.52.5, pgx 5.7.1, Go 1.25
  - See `api-fiber/README.md` for full documentation
  - Database: PostgreSQL 18

- **api-fastify/** - Node.js Fastify REST API
  - Latest versions: Fastify 5.6.1, pg 8.12.0, Node.js 22 LTS
  - See `api-fastify/README.md` for full documentation
  - Database: PostgreSQL 18

- **api-fastapi/** - Python FastAPI REST API
  - Latest versions: FastAPI 0.115.4, SQLAlchemy 2.0.36, asyncpg 0.31.0, Python 3.12
  - See `api-fastapi/README.md` for full documentation
  - Database: PostgreSQL 18

### Frontend

- **client/** - Vanilla JavaScript frontend
  - Built with Vanilla JavaScript and Bun runtime
  - See `client/README.md` for client documentation
  - See `client/BUN_SETUP.md` for Bun setup instructions
  - Consumes the API from `http://localhost:8080/api`

### Documentation

- **API_IMPLEMENTATIONS_GUIDE.md** - Comprehensive 15KB+ guide comparing all 4 implementations
- **GETTING_STARTED.md** - Selection guide for choosing an API implementation
- **COMPARISON.md** - Detailed Rust vs Go comparison
- **RUST_VERSION_UPDATE.md** - Latest dependency versions documentation

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

## Project Status

### APIs (All Complete & Production-Ready)
1. ✅ **Rust Actix** - Full CRUD operations, latest dependencies (4.11 actix-web, 1.41 tokio, 0.8 sqlx)
2. ✅ **Go FastHTTP** - Full CRUD operations, latest dependencies (1.67.0 fasthttp, 5.7.1 pgx)
3. ✅ **Go Fiber** - Full CRUD operations, latest dependencies (2.52.5 fiber, 5.7.1 pgx)
4. ✅ **Node.js Fastify** - Full CRUD operations, latest dependencies (5.6.1 fastify, 8.12.0 pg)
5. ✅ **Python FastAPI** - Full CRUD operations, latest dependencies (0.115.4 fastapi, 2.0.36 sqlalchemy, 0.31.0 asyncpg)

### Frontend
1. ✅ Frontend created with Vanilla JavaScript + Bun
2. ✅ Frontend connected to API endpoints
3. ✅ Fully functional todo management interface

### Documentation
1. ✅ Comprehensive API Implementation Guide (15KB+)
2. ✅ Selection guide for choosing an API
3. ✅ Individual documentation for each implementation
4. ✅ Latest dependency version documentation

### Overall Status
🚀 **Ready for production deployment with 4 independent implementation choices**

## Next Steps

1. 🔧 Deploy frontend and API to production
2. 📊 Add additional features (filters, search, etc.)
3. 🔐 Implement authentication if needed
4. 📱 Add responsive mobile design enhancements

## License

This project is open source and available under the MIT License.

## Support

For issues or questions about the API, see [api/README.md](api/README.md)
