# JWT Authentication Service

A complete JWT authentication service written in Go using the Gin framework. This service is designed to be used as a standalone authentication service that can be integrated with other applications through Docker.

## Features

- User registration and login
- JWT token-based authentication
- Access and refresh tokens
- Token validation and refresh
- Protected routes
- User profile
- Logout functionality
- CORS support
- PostgreSQL database

## Prerequisites

### For Docker Deployment (Recommended)
- Docker
- Docker Compose

### For Local Development
- Go 1.20 or later
- PostgreSQL 12 or later
- Make (optional, for using Makefile commands)

## Setup

### Using Docker (Recommended)

1. Clone the repository
```bash
git clone <repository-url>
cd jwt-auth
```

2. Build and run the containers:
```bash
docker-compose up -d
```

The service will be available at `http://localhost:8080`

To stop the service:
```bash
docker-compose down
```

To view logs:
```bash
docker-compose logs -f api
```

### Integration with Other Applications

To integrate this authentication service with your application:

1. Add this service to your application's docker-compose.yml:
```yaml
services:
  auth_service:
    image: your-registry/jwt-auth:latest
    environment:
      - JWT_SECRET=your_jwt_secret_key
      - JWT_ACCESS_EXPIRY=1h
      - JWT_REFRESH_EXPIRY=24h
    ports:
      - "8080:8080"
    networks:
      - your_network
```

2. Use the authentication endpoints from your application:
```
http://auth_service:8080/api/v1/auth/*
```

### Local Development

1. Clone the repository
```bash
git clone <repository-url>
cd jwt-auth
```

2. Create a `.env` file in the root directory:
```env
# Server Configuration
SERVER_HOST=localhost
SERVER_PORT=8080

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=jwt_auth
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your_jwt_secret_key
JWT_ACCESS_EXPIRY=1h
JWT_REFRESH_EXPIRY=24h
```

3. Create the database:
```bash
createdb jwt_auth
```

4. Run the migrations:
```bash
psql -d jwt_auth -f migrations/000001_init.sql
```

5. Install dependencies:
```bash
go mod download
```

## Running the Service

Start the server:
```bash
go run cmd/main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Public Routes

- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/refresh` - Refresh access token
- `GET /health` - Health check endpoint

### Protected Routes (Requires Authentication)

- `GET /api/v1/profile` - Get user profile
- `POST /api/v1/logout` - Logout user
- `GET /api/v1/dashboard` - Example protected route

## Authentication

Include the JWT token in the Authorization header:
```
Authorization: Bearer <your-token>
```

## Project Structure

```
.
├── cmd/
│   └── main.go
├── internal/
│   ├── application/
│   │   ├── dto/
│   │   └── services/
│   ├── domain/
│   │   ├── entities/
│   │   ├── repositories/
│   │   └── services/
│   ├── infrastructure/
│   │   ├── database/
│   │   ├── jwt/
│   │   └── repositories/
│   └── interfaces/
│       ├── config/
│       └── http/
│           ├── handlers/
│           ├── middleware/
│           └── routes/
├── migrations/
└── .env
```

## Error Handling

The service returns appropriate HTTP status codes and error messages in a consistent format:

```json
{
  "error": "error_code",
  "message": "Descriptive error message"
}
```

## Success Response Format

Successful responses follow this format:

```json
{
  "message": "Success message",
  "data": {} // Optional data payload
}
```

## License

MIT License
