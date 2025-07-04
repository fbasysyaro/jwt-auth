# Simple JWT Authentication Service

A minimal JWT authentication service for other applications.

## Features

1. **User Registration** - Register new users from other systems
2. **User Login** - Authenticate users and get JWT tokens
3. **Token Verification** - Validate JWT tokens for other systems
4. **User Logout** - Blacklist tokens for secure logout

## Quick Start

```bash
# With Docker
docker-compose up --build -d

# Local development
go run main.go
```

## API Endpoints

### 1. Register User
```bash
POST /api/v1/register
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com", 
  "password": "password123"
}

Response:
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 2. Login User
```bash
POST /api/v1/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}

Response: Same as register
```

### 3. Verify Token
```bash
POST /api/v1/verify
Authorization: Bearer <token>

# OR

POST /api/v1/verify
Content-Type: application/json

{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}

Response:
{
  "valid": true,
  "user": {
    "id": 1,
    "username": "john_doe", 
    "email": "john@example.com"
  }
}
```

### 4. Logout User
```bash
POST /api/v1/logout
Authorization: Bearer <token>

Response:
{
  "message": "Logged out successfully"
}
```

## Integration with Other Systems

### Example: Verify user in another application
```go
func verifyUser(token string) (*User, error) {
    resp, err := http.Post("http://jwt-auth:8080/api/v1/verify", 
        "application/json", 
        strings.NewReader(`{"token":"`+token+`"}`))
    
    if err != nil || resp.StatusCode != 200 {
        return nil, errors.New("invalid token")
    }
    
    // Parse response and return user
    return user, nil
}
```

## Environment Variables

- `PORT` - Server port (default: 8080)
- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database user (default: postgres)
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name (default: jwt_auth)
- `JWT_SECRET` - JWT signing secret

## Project Structure

```
/
├── main.go           # Application entry point
├── config/
│   └── config.go     # Configuration and DB setup
├── handler/
│   └── auth.go       # HTTP handlers (register, login, verify, logout)
└── repository/
    ├── user.go       # User database operations
    └── token.go      # JWT token operations
```