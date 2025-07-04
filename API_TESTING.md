# API Testing Guide

## Swagger Documentation

### Access Swagger UI
1. Start the application: `go run cmd/main.go`
2. Open browser: `http://localhost:8080/api/v1/swagger/index.html`

## Manual Testing with curl

### 1. Health Check
```bash
curl -X GET http://localhost:8080/health
```

### 2. Register User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 3. Login User
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

**Save the access_token from response for next requests**

### 4. Get Profile (Protected)
```bash
curl -X GET http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 5. Access Dashboard (Protected)
```bash
curl -X GET http://localhost:8080/api/v1/dashboard \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 6. Refresh Token
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

### 7. Logout (Protected)
```bash
curl -X POST http://localhost:8080/api/v1/logout \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Testing with Postman

### Import Collection
Create a Postman collection with these requests:

1. **Register** - POST `/api/v1/auth/register`
2. **Login** - POST `/api/v1/auth/login`
3. **Profile** - GET `/api/v1/profile`
4. **Dashboard** - GET `/api/v1/dashboard`
5. **Refresh** - POST `/api/v1/auth/refresh`
6. **Logout** - POST `/api/v1/logout`

### Environment Variables
Set these in Postman environment:
- `base_url`: `http://localhost:8080`
- `access_token`: (set from login response)
- `refresh_token`: (set from login response)

## Expected Responses

### Success Response Format
```json
{
  "message": "Success message",
  "data": {}
}
```

### Error Response Format
```json
{
  "error": "error_code",
  "message": "Descriptive error message"
}
```

### Auth Response Format
```json
{
  "access_token": "token_string",
  "refresh_token": "refresh_token_string",
  "token_type": "Bearer",
  "expires_in": 3600,
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

## Testing Flow

1. **Register** → Get tokens
2. **Login** → Get fresh tokens
3. **Profile** → Verify authentication works
4. **Dashboard** → Test another protected route
5. **Refresh** → Get new tokens
6. **Logout** → Blacklist current token
7. **Try Profile again** → Should fail (401)