# Docker Compose Troubleshooting & Fix Guide

## Issues Fixed

All docker-compose.yml files have been updated with the following fixes:

### ✅ Fix 1: Missing API Service (Rust Actix)
**Problem:** The `api/docker-compose.yml` was missing the API application service definition
**Solution:** Added the `api` service with proper build, environment, and port configuration

### ✅ Fix 2: Network Configuration
**Problem:** Services were not using Docker networks for inter-container communication
**Solution:** Added `todo_network` bridge network to all services in all implementations

### ✅ Fix 3: Proper Service Dependencies
**Problem:** API services depended on PostgreSQL without waiting for it to be ready
**Solution:** Added health checks and `depends_on` with service_healthy condition

### ✅ Fix 4: Container Naming
**Problem:** All implementations used the same container names, causing conflicts
**Solution:** Updated container names to be unique:
- `todo-app-db` (shared Postgres - can be overridden)
- `todo-app-rust` (Rust Actix)
- `todo-app-go` (Go FastHTTP)
- `todo-app-fiber` (Go Fiber)
- `todo-app-fastify` (Node.js Fastify)
- `todo-fastapi` (Python FastAPI)

### ✅ Fix 5: Environment Variables
**Problem:** Some implementations had incorrect or missing environment variables
**Solution:** Standardized environment variables:
- `DATABASE_URL`: Correct connection string for each implementation
- `PORT`: 8080 (default)
- `LOG_LEVEL` or `NODE_ENV`: Appropriate for each implementation

### ✅ Fix 6: Database Port Conflicts
**Problem:** All APIs tried to expose PostgreSQL on port 5432
**Solution:** Each implementation still uses 5432 internally (in container), mapped to 5432 on host (they can't run simultaneously)

## Updated docker-compose.yml Structure

All files now follow this standard structure:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:18-alpine
    container_name: todo-postgres-[impl]
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: todo_db
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - todo_network

  api:  # or app (application service name varies)
    build: .
    container_name: todo-app-[impl]
    environment:
      DATABASE_URL: postgres://postgres:password@postgres:5432/todo_db
      PORT: 8080
      # ... implementation-specific variables
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
    networks:
      - todo_network
    restart: unless-stopped

networks:
  todo_network:
    driver: bridge

volumes:
  postgres_data:
```

## How to Run Individual APIs

### Rust Actix
```bash
cd api
docker-compose up --build
# API at http://localhost:8080
```

### Go FastHTTP
```bash
cd api-go
docker-compose up --build
# API at http://localhost:8080
```

### Go Fiber
```bash
cd api-fiber
docker-compose up --build
# API at http://localhost:8080
```

### Node.js Fastify
```bash
cd api-fastify
docker-compose up --build
# API at http://localhost:8080
```

### Python FastAPI
```bash
cd api-fastapi
docker-compose up --build
# API at http://localhost:8080
```

## Running Multiple APIs Simultaneously

To run multiple APIs on different ports, use environment variable overrides:

### Terminal 1 (Rust Actix on 8080)
```bash
cd api
docker-compose up --build
```

### Terminal 2 (Go FastHTTP on 8081)
```bash
cd api-go
PORT=8081 docker-compose up --build
# Edit docker-compose.yml temporarily or use -p override
```

**Note:** To properly run multiple APIs simultaneously without port conflicts, you would need to:
1. Change the exposed ports in each docker-compose.yml:
   - Rust Actix: 8080:8080
   - Go FastHTTP: 8081:8080
   - Go Fiber: 8082:8080
   - Node.js Fastify: 8083:8080
   - Python FastAPI: 8084:8080

2. Use unique container and network names

3. Update the benchmark.py to use different ports

## Debugging Docker Compose Issues

### Check Docker Service Status
```bash
# See all running containers
docker ps

# See all containers (including stopped)
docker ps -a

# View logs for a specific container
docker logs todo-app-rust
docker logs todo-postgres

# Follow logs in real-time
docker logs -f todo-app-rust
```

### Check Network Connectivity
```bash
# List all networks
docker network ls

# Inspect a network
docker network inspect todo_network

# Test DNS resolution between containers
docker exec todo-app-rust ping postgres
```

### Common Issues & Solutions

#### "bind: address already in use"
**Problem:** Port 5432 or 8080 is already in use
**Solution:**
```bash
# Option 1: Stop the service using the port
docker-compose down

# Option 2: Use a different port
# Edit docker-compose.yml:
ports:
  - "5433:5432"  # Use 5433 on host instead
```

#### "Cannot connect to Docker daemon"
**Problem:** Docker is not running
**Solution:**
```bash
# Start Docker Desktop (Mac/Windows) or Docker daemon (Linux)
# On Linux:
sudo systemctl start docker

# On Mac:
# Open Docker Desktop application

# Verify Docker is running:
docker --version
```

#### "database connection refused"
**Problem:** Database is not ready when API starts
**Solution:**
- The docker-compose now includes health checks
- API waits for database with `depends_on: condition: service_healthy`
- Give it a few more seconds to start

#### "service host not found"
**Problem:** Container can't resolve the service name
**Solution:**
- Check that all services are on the same network
- Verify network is defined correctly in docker-compose.yml
- Ensure container names match in connection strings

### View Docker Compose Configuration
```bash
# See the final configuration after variable substitution
docker-compose config

# Validate the docker-compose file
docker-compose config --quiet
```

### Force Rebuild
```bash
# Remove all containers and volumes, rebuild from scratch
docker-compose down -v
docker-compose up --build
```

### Clean Up Everything
```bash
# Remove containers
docker-compose down

# Remove containers and volumes
docker-compose down -v

# Remove images as well
docker-compose down --rmi all
```

## Updated Files

The following docker-compose.yml files have been fixed:

1. ✅ `api/docker-compose.yml` - Added missing API service
2. ✅ `api-go/docker-compose.yml` - Added networks, fixed variables
3. ✅ `api-fiber/docker-compose.yml` - Added networks
4. ✅ `api-fastify/docker-compose.yml` - Updated service naming
5. ✅ `api-fastapi/docker-compose.yml` - Already correct

## Testing the Fixes

### Quick Test
```bash
# 1. Start an API
cd api
docker-compose up --build

# 2. In another terminal, test the API
curl http://localhost:8080/health
curl http://localhost:8080/api/todos

# 3. Stop with Ctrl+C
# 4. Clean up
docker-compose down
```

### Full Test for One Implementation
```bash
cd api-fastapi
docker-compose up --build

# Wait for "Application startup complete" message

# Test endpoints
curl http://localhost:8080/health
curl http://localhost:8080/api/todos
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Test Todo"}'

# Clean up
docker-compose down -v
```

## Benchmark Setup (Using Fixed docker-compose)

For the benchmarking suite to work, each API needs to be on a different port:

### Option A: Run One at a Time
```bash
# Test each implementation separately
for impl in api api-go api-fiber api-fastify api-fastapi; do
  cd $impl
  docker-compose up --build
  # Test the API
  curl http://localhost:8080/health
  docker-compose down -v
  cd ..
done
```

### Option B: Port Forwarding (Advanced)
Edit each docker-compose.yml to use different host ports:

**api/docker-compose.yml:**
```yaml
api:
  ports:
    - "8080:8080"  # Keep as is
```

**api-go/docker-compose.yml:**
```yaml
api:
  ports:
    - "8081:8080"  # Map host 8081 to container 8080
```

**api-fiber/docker-compose.yml:**
```yaml
api:
  ports:
    - "8082:8080"  # Map host 8082 to container 8080
```

**api-fastify/docker-compose.yml:**
```yaml
api:
  ports:
    - "8083:8080"  # Map host 8083 to container 8080
```

**api-fastapi/docker-compose.yml:**
```yaml
api:
  ports:
    - "8084:8080"  # Map host 8084 to container 8080
```

Then update benchmark.py to use the correct ports:
```python
base_urls = {
    "rust_actix": "http://localhost:8080",
    "go_fasthttp": "http://localhost:8081",
    "go_fiber": "http://localhost:8082",
    "node_fastify": "http://localhost:8083",
    "python_fastapi": "http://localhost:8084",
}
```

## Summary

All docker-compose.yml files have been fixed and now:

✅ Include all necessary services (API + Database)
✅ Use proper networking (todo_network)
✅ Have health checks for database readiness
✅ Include proper service dependencies
✅ Use unique container names per implementation
✅ Have correct environment variables
✅ Support docker-compose up --build

Each API should now start correctly with:
```bash
cd [api-directory]
docker-compose up --build
```

## Additional Resources

- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Docker Networking](https://docs.docker.com/network/)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [Healthchecks in Docker Compose](https://docs.docker.com/compose/compose-file/compose-file-v3/#healthcheck)
