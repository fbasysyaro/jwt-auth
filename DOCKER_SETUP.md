# Docker Setup Guide

## Quick Start

### 1. Stop any existing services
```bash
# Check what's using port 8080
check-port.bat

# Stop existing Docker containers
docker-compose down
```

### 2. Start with Docker
```bash
# Option 1: Use the script
start-docker.bat

# Option 2: Manual commands
docker-compose down
docker-compose up --build -d
```

### 3. Access Swagger
```
http://localhost:8080/api/v1/swagger/index.html
```

## Troubleshooting

### Port 8080 Already in Use
```bash
# Check what's using the port
netstat -ano | findstr :8080

# Kill the process (replace PID with actual process ID)
taskkill /PID [PID] /F

# Or use different port in docker-compose.yml
ports:
  - "8081:8080"  # Use port 8081 instead
```

### Check Container Status
```bash
# View running containers
docker ps

# View logs
docker-compose logs -f api
docker-compose logs -f postgres
docker-compose logs -f redis
```

### Reset Everything
```bash
# Stop and remove everything
docker-compose down -v

# Remove images
docker-compose down --rmi all

# Start fresh
docker-compose up --build -d
```

## Services

- **API**: http://localhost:8080
- **Swagger**: http://localhost:8080/api/v1/swagger/index.html
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

## Environment Variables

The Docker setup uses these environment variables:
- `REDIS_HOST=redis` (container name)
- `DB_HOST=postgres` (container name)
- All other configs in docker-compose.yml

## Testing with Docker

Once containers are running, use the same testing methods:

1. **Swagger UI**: http://localhost:8080/api/v1/swagger/index.html
2. **Test Scripts**: `test_api.bat` (works with Docker too)
3. **Manual curl**: Same commands as in API_TESTING.md