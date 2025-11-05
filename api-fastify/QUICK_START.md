# Quick Start Guide - Fastify Todo API

## Prerequisites

- Node.js 22+ (LTS)
- PostgreSQL 18+
- Docker & Docker Compose (optional but recommended)

## Option 1: Quick Start with Docker Compose (Recommended)

The fastest way to get everything running:

```bash
cd api-fastify
docker-compose up --build
```

Your API will be available at: `http://localhost:8080/api/todos`

To stop: `Ctrl+C` or run `docker-compose down`

## Option 2: Local Development

### 1. Start PostgreSQL

```bash
# Using Docker
docker run -d --name postgres-todo \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=todo_db \
  -p 5432:5432 \
  postgres:18-alpine

# Or use your local PostgreSQL installation
```

### 2. Set Up Environment

```bash
cd api-fastify
cp .env .env.local
# Edit .env.local if needed (default values work)
```

### 3. Install Dependencies

```bash
npm install
```

### 4. Run Database Migrations

```bash
psql postgres://postgres:password@localhost:5432/todo_db < migrations/01_create_todos_table.sql
```

### 5. Run the Application

```bash
npm start
```

The API will start on `http://127.0.0.1:8080`

## Quick API Test

```bash
# List all todos
curl http://localhost:8080/api/todos

# Create a todo
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"My First Fastify Todo","description":"Learn Fastify"}'

# Get a todo (replace UUID)
curl http://localhost:8080/api/todos/550e8400-e29b-41d4-a716-446655440000

# Update a todo
curl -X PUT http://localhost:8080/api/todos/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{"completed":true}'

# Delete a todo
curl -X DELETE http://localhost:8080/api/todos/550e8400-e29b-41d4-a716-446655440000
```

## npm Scripts

```bash
# Start production server
npm start

# Start development server with auto-reload
npm run dev

# Run tests
npm test

# Format code
npm run format

# Lint code
npm run lint
```

## Troubleshooting

### Port Already in Use
```bash
export PORT=8081
npm start
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

### npm install fails
```bash
# Clear npm cache
npm cache clean --force
rm -rf node_modules package-lock.json
npm install
```

## Why Fastify?

Fastify is the **fastest web framework for Node.js**:

✅ **Performance** - Up to 30% faster than Express
✅ **Built-in Validation** - JSON Schema validation included
✅ **Low Overhead** - Minimal dependencies
✅ **Developer Friendly** - Express-like API
✅ **Active Community** - Growing ecosystem

## Next Steps

1. **Read the [README.md](README.md)** for comprehensive documentation
2. **Review the code** in `src/` directories
3. **Explore handlers** in `src/handlers/todo.js`
4. **Test the API** with curl or Postman
5. **Deploy to production** using Docker

## Development Workflow

```bash
# 1. Install dependencies
npm install

# 2. Start development server
npm run dev

# 3. In another terminal, test endpoints
curl http://localhost:8080/api/todos

# 4. Make changes to src/ files
# 5. Server auto-reloads (with --watch)

# 6. When ready, run production build
npm start
```

## Comparison with Other Todo APIs

All four implementations provide identical functionality:

| Implementation | Framework | Language | Performance |
|---|---|---|---|
| **api/** | Actix | Rust | Excellent |
| **api-go/** | FastHTTP | Go | Excellent |
| **api-fiber/** | Fiber | Go | Excellent |
| **api-fastify/** | Fastify | Node.js | Excellent |

**Choose Fastify if:**
- You prefer Node.js and JavaScript
- You want the fastest Node.js framework
- You're familiar with Express.js
- You need a quick Node.js backend

## Support

For more information:
- [Fastify Documentation](https://www.fastify.io/)
- [Node.js Documentation](https://nodejs.org/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

---

Happy coding! 🚀
