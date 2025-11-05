# Todo API - Fastify Implementation

A high-performance RESTful API for managing todo items, built with **Fastify** - the fastest web framework for Node.js. This is the fourth and final implementation of the todo API, demonstrating the same functionality as Rust Actix, Go FastHTTP, and Go Fiber versions.

## Technology Stack

**Core Framework:**
- **Fastify 5.6.1** - The fastest web framework for Node.js
- **Node.js 22 (LTS)** - Latest long-term support version
- **PostgreSQL 18** - Robust relational database
- **pg 8.12.0** - Native PostgreSQL driver

**Utilities:**
- **UUID 10.0.0** - UUID generation
- **Joi 17.14.0** - Schema validation
- **Pino 9.6.0** - High-performance logging
- **@fastify/cors 9.0.1** - CORS support
- **@fastify/helmet 11.2.1** - Security headers

## Project Structure

```
api-fastify/
├── src/
│   ├── index.js                    # Application entry point
│   ├── handlers/
│   │   └── todo.js                 # CRUD operation handlers
│   ├── models/
│   │   └── todo.js                 # Data models and schemas
│   ├── routes/
│   │   └── todo.js                 # Route definitions
│   ├── db/
│   │   └── connection.js           # Database connection pool
│   └── utils/
│       └── logger.js               # Pino logger setup
├── migrations/
│   └── 01_create_todos_table.sql   # Database schema
├── package.json                    # npm dependencies
├── .env                            # Environment variables
├── Dockerfile                      # Docker container build
├── docker-compose.yml              # Service orchestration
└── README.md                       # This file
```

## API Endpoints

All endpoints are under `/api/todos`:

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
- `description: string | null` - Optional description
- `completed: boolean` - Completion status (default: false)
- `created_at: Date` - Creation timestamp (UTC)
- `updated_at: Date` - Last update timestamp (UTC)

## Installation & Setup

### Prerequisites

- Node.js 22+ (LTS)
- PostgreSQL 18+
- Docker and Docker Compose (for containerized setup)

### Local Development

1. **Clone the repository:**
```bash
cd api-fastify
```

2. **Install dependencies:**
```bash
npm install
```

3. **Set up environment variables:**
```bash
cp .env .env.local
# Edit .env.local with your database credentials
```

4. **Start PostgreSQL:**
```bash
# Using Docker
docker run -d --name postgres-todo \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=todo_db \
  -p 5432:5432 \
  postgres:18-alpine

# Or use your local PostgreSQL installation
```

5. **Run database migrations:**
```bash
psql -U postgres -d todo_db -f migrations/01_create_todos_table.sql
```

6. **Run the application:**
```bash
npm start
```

The API will be available at `http://127.0.0.1:8080`

### Using Docker Compose

1. **Build and start all services:**
```bash
docker-compose up --build
```

2. **Verify the setup:**
```bash
curl http://localhost:8080/api/todos
```

3. **Stop the services:**
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
    "title": "Learn Fastify",
    "description": "Study the fastest Node.js framework",
    "completed": false,
    "created_at": "2025-11-05T10:30:00.000Z",
    "updated_at": "2025-11-05T10:30:00.000Z"
  }
]
```

### Create a Todo

```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Build API with Fastify",
    "description": "Create a high-performance Node.js API"
  }'
```

**Response:** (201 Created)
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "title": "Build API with Fastify",
  "description": "Create a high-performance Node.js API",
  "completed": false,
  "created_at": "2025-11-05T10:35:00.000Z",
  "updated_at": "2025-11-05T10:35:00.000Z"
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

Partial updates are supported - only include fields you want to update.

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

- `HOST` - Server listening host (optional, defaults to 127.0.0.1)

- `NODE_ENV` - Environment mode (development or production)

- `LOG_LEVEL` - Logging level (debug, info, warn, error)

All environment variables can be set in `.env` or as system variables.

## Key Features

- **Fastest Node.js Framework** - Built for maximum performance
- **Express.js-like API** - Familiar syntax for Node developers
- **Schema Validation** - Automatic request validation with Fastify schemas and Joi
- **Error Handling** - Structured JSON error responses
- **CORS Support** - Configured for cross-origin requests
- **Security Headers** - Helmet.js for security
- **Request Logging** - Pino logger with pretty-printing in development
- **Connection Pooling** - Efficient PostgreSQL connection management
- **Docker Ready** - Multi-stage Docker build
- **Graceful Shutdown** - Proper signal handling
- **Partial Updates** - Update only the fields you need

## Building & Deployment

### Build Docker Image

```bash
docker build -t todo-app-fastify:latest .
```

### Run with Docker

```bash
docker run -p 8080:8080 \
  -e DATABASE_URL="postgres://postgres:password@host:5432/todo_db" \
  todo-app-fastify:latest
```

### Run Locally

```bash
npm install
npm start
```

### Development Mode (with auto-reload)

```bash
npm run dev
```

## Performance Characteristics

Fastify provides excellent performance benefits:

- **Throughput** - 50,000+ requests/second (typical)
- **Latency** - <5ms average response time
- **Memory** - 50-100MB runtime footprint
- **Startup** - <100ms server startup time
- **Concurrency** - Handles thousands of concurrent connections

Fastify is the fastest web framework for Node.js, consistently outperforming Express and Hapi in benchmarks.

## Development

### Running Tests

```bash
npm test
```

### Code Formatting

```bash
npm run format
```

### Linting

```bash
npm run lint
```

### Code Organization

- `src/index.js` - Application entry point and middleware setup
- `src/handlers/` - HTTP request handlers
- `src/models/` - Data structures and validation schemas
- `src/routes/` - Route definitions
- `src/db/` - Database connection and pooling
- `src/utils/` - Utility functions (logging, etc.)

### Adding New Endpoints

1. Create handler in `src/handlers/`
2. Define schemas in `src/models/`
3. Register route in `src/routes/`
4. Test with curl or REST client

## Framework Comparison

### Fastify vs Express vs Hapi

| Feature | Fastify | Express | Hapi |
|---------|---------|---------|------|
| **Performance** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐ |
| **Request Speed** | Fastest | Good | Good |
| **Built-in Validation** | ✅ | ❌ | ✅ |
| **API Style** | Clean | Minimal | Verbose |
| **Learning Curve** | Easy | Very Easy | Moderate |
| **Ecosystem** | Growing | Largest | Moderate |

## Troubleshooting

### Connection Refused
- Ensure PostgreSQL is running
- Check DATABASE_URL is correct
- Verify PostgreSQL port is accessible

### Port Already in Use
```bash
# Change PORT environment variable
export PORT=8081
npm start
```

### Database Connection Issues
```bash
# Test connection with psql
psql postgres://postgres:password@localhost:5432/todo_db
```

### Module Not Found
```bash
# Clear npm cache and reinstall
rm -rf node_modules package-lock.json
npm install
```

## Next Steps

1. **Read the code** - Explore `src/` directory
2. **Run locally** - Start the development server
3. **Test endpoints** - Use curl or Postman
4. **Deploy** - Push Docker image to registry
5. **Monitor** - Set up logging and monitoring

## License

MIT License

## Additional Resources

- [Fastify Documentation](https://www.fastify.io/)
- [Fastify GitHub](https://github.com/fastify/fastify)
- [Node.js Documentation](https://nodejs.org/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Pino Logger](https://getpino.io/)
- [Joi Validation](https://joi.dev/)
