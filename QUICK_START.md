# Quick Start Guide

## 1. Start the Application

### Option A: Using Docker (Recommended)
```bash
# Quick start
start-docker.bat

# Or manually
docker-compose up --build -d
```

### Option B: Local Development
```bash
# Make sure PostgreSQL and Redis are running locally
go run cmd/main.go
```

### Option C: Build and Run
```bash
go build -o jwt-auth.exe cmd/main.go
./jwt-auth.exe
```

## 2. Access Swagger Documentation

Open your browser and go to:
```
http://localhost:8080/api/v1/swagger/index.html
```

## 3. Test API Endpoints

### Using Swagger UI (Recommended)
1. Go to Swagger URL above
2. Try the endpoints in this order:
   - **POST /auth/register** - Register a new user
   - **POST /auth/login** - Login and get tokens
   - **GET /profile** - Test protected route (click 🔒 and enter: `Bearer YOUR_TOKEN`)
   - **POST /logout** - Logout user

### Using Test Scripts
```bash
# Windows
test_api.bat

# Linux/Mac
chmod +x test_api.sh
./test_api.sh
```

### Using curl (Manual)
```bash
# 1. Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"password123"}'

# 2. Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# 3. Use the access_token from login response for protected routes
curl -X GET http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 4. Available Endpoints

- `GET /health` - Health check
- `POST /api/v1/auth/register` - Register user
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/refresh` - Refresh token
- `GET /api/v1/profile` - Get user profile (protected)
- `GET /api/v1/dashboard` - Dashboard (protected)
- `POST /api/v1/logout` - Logout (protected)

## 5. Authentication

For protected routes, include the JWT token in the Authorization header:
```
Authorization: Bearer YOUR_ACCESS_TOKEN
```

## Troubleshooting

- **Port 8080 in use**: Run `check-port.bat` to see what's using it, then `docker-compose down`
- **Docker issues**: See `DOCKER_SETUP.md` for detailed troubleshooting
- **Database connection**: Ensure PostgreSQL container is running
- **Redis connection**: Ensure Redis container is running