# Quick Start Guide - FastAPI Todo API

## Prerequisites

- Python 3.12+ (Download from https://www.python.org/)
- PostgreSQL 18+ (or Docker)
- Docker & Docker Compose (optional but recommended)

## Option 1: Quick Start with Docker Compose (Recommended)

The fastest way to get everything running:

```bash
cd api-fastapi
docker-compose up --build
```

Your API will be available at: `http://localhost:8080/api/todos`

**API Documentation:**
- Swagger UI: `http://localhost:8080/docs`
- ReDoc: `http://localhost:8080/redoc`

To stop: `Ctrl+C` or run `docker-compose down`

## Option 2: Local Development

### 1. Create Virtual Environment

```bash
cd api-fastapi

# Create virtual environment
python3 -m venv venv

# Activate it
source venv/bin/activate           # On macOS/Linux
# or
venv\Scripts\activate              # On Windows
```

### 2. Start PostgreSQL

```bash
# Using Docker (recommended)
docker run -d --name postgres-todo \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=todo_db \
  -p 5432:5432 \
  postgres:18-alpine

# Or use your local PostgreSQL installation
```

### 3. Install Dependencies

```bash
pip install -r requirements.txt
```

### 4. Run Database Migrations (Optional)

```bash
# FastAPI auto-creates tables, but you can pre-populate:
psql postgres://postgres:password@localhost:5432/todo_db < migrations/01_create_todos_table.sql
```

### 5. Run the Application

```bash
python -m uvicorn main:app --reload --host 0.0.0.0 --port 8080
```

The API will start on `http://127.0.0.1:8080`

**API Documentation:**
- Swagger UI: `http://127.0.0.1:8080/docs`
- ReDoc: `http://127.0.0.1:8080/redoc`

## Quick API Test

### List all todos
```bash
curl http://localhost:8080/api/todos
```

### Create a todo
```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Learn FastAPI",
    "description": "Study FastAPI documentation and async patterns"
  }'
```

### Get a todo (replace UUID)
```bash
curl http://localhost:8080/api/todos/550e8400-e29b-41d4-a716-446655440000
```

### Update a todo
```bash
curl -X PUT http://localhost:8080/api/todos/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{
    "completed": true,
    "title": "Mastered FastAPI"
  }'
```

### Delete a todo
```bash
curl -X DELETE http://localhost:8080/api/todos/550e8400-e29b-41d4-a716-446655440000
```

## Testing with Python requests

```python
import requests

BASE_URL = "http://localhost:8080"

# List todos
response = requests.get(f"{BASE_URL}/api/todos")
print(response.json())

# Create todo
response = requests.post(
    f"{BASE_URL}/api/todos",
    json={
        "title": "Learn FastAPI",
        "description": "Master async Python"
    }
)
print(response.json())
```

## Environment Variables

Configuration via `.env`:
```bash
# Database
DATABASE_URL=postgresql+asyncpg://postgres:password@localhost:5432/todo_db

# Server
PORT=8080
ENV=development

# Logging
LOG_LEVEL=info
```

## npm Scripts / Commands

### Development
```bash
# Run with hot reload
python -m uvicorn main:app --reload

# Run tests
pytest -v

# Check types
mypy main.py models.py schemas.py
```

### Production
```bash
# Start production server (no reload)
python -m uvicorn main:app --host 0.0.0.0 --port 8080

# With gunicorn (multiple workers)
gunicorn -w 4 -k uvicorn.workers.UvicornWorker main:app
```

### Docker
```bash
# Build image
docker build -t todo-fastapi .

# Run container
docker run -p 8080:8080 \
  -e DATABASE_URL=postgresql+asyncpg://postgres:password@postgres:5432/todo_db \
  todo-fastapi

# With docker-compose
docker-compose up --build
docker-compose down
```

## Troubleshooting

### Port Already in Use
```bash
export PORT=8081
python -m uvicorn main:app --port 8081
```

### Cannot Connect to PostgreSQL
Check the connection string in `.env`:
```bash
# Should match your PostgreSQL setup
DATABASE_URL=postgresql+asyncpg://postgres:password@localhost:5432/todo_db
```

Test the connection:
```bash
psql postgres://postgres:password@localhost:5432/todo_db -c "SELECT 1"
```

### pip install fails
```bash
# Upgrade pip
pip install --upgrade pip

# Clear cache and reinstall
pip cache purge
pip install -r requirements.txt
```

## Why FastAPI?

FastAPI is the **fastest Python web framework**:

✅ **Performance** - Comparable to Node.js and Go
✅ **Async Support** - Native async/await for concurrency
✅ **Auto Documentation** - Built-in Swagger UI and ReDoc
✅ **Type Safety** - Full Pydantic validation
✅ **Developer Experience** - Clean, intuitive API
✅ **Production Ready** - Used by major companies

## Next Steps

1. **Read the [README.md](README.md)** for comprehensive documentation
2. **Review the code** in `main.py`, `models.py`, `schemas.py`
3. **Explore the API** at `http://localhost:8080/docs`
4. **Test the endpoints** with curl or Postman
5. **Deploy to production** using Docker

## Development Workflow

```bash
# 1. Create virtual environment
python -m venv venv
source venv/bin/activate

# 2. Install dependencies
pip install -r requirements.txt

# 3. Start PostgreSQL (Docker)
docker run -d --name postgres-todo \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=todo_db \
  -p 5432:5432 \
  postgres:18-alpine

# 4. Start development server
python -m uvicorn main:app --reload

# 5. Test the API
curl http://localhost:8080/api/todos

# 6. View API docs
# Open http://localhost:8080/docs in browser

# 7. When done, stop services
Ctrl+C  # Stop FastAPI
docker stop postgres-todo
```

## Comparison with Other Todo APIs

All four implementations provide identical functionality:

| Implementation | Framework | Language | Performance |
|---|---|---|---|
| **api/** | Actix | Rust | Excellent |
| **api-go/** | FastHTTP | Go | Excellent |
| **api-fiber/** | Fiber | Go | Excellent |
| **api-fastapi/** | FastAPI | Python | Very Good |

**Choose FastAPI if:**
- You prefer Python and its ecosystem
- You want auto-generated API documentation
- You're familiar with async Python
- You need quick development iteration
- You want the easiest to learn

## Support

For more information:
- [FastAPI Documentation](https://fastapi.tiangolo.com/)
- [SQLAlchemy Async Guide](https://docs.sqlalchemy.org/en/20/orm/extensions/asyncio.html)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

---

Happy coding! 🚀
