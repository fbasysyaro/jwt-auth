@echo off
setlocal enabledelayedexpansion

set BASE_URL=http://localhost:8080

echo === JWT Auth API Testing ===
echo.

echo 1. Testing Health Check...
curl -s -X GET %BASE_URL%/health
echo.
echo.

echo 2. Testing User Registration...
curl -s -X POST %BASE_URL%/api/v1/auth/register ^
  -H "Content-Type: application/json" ^
  -d "{\"username\": \"testuser\", \"email\": \"test@example.com\", \"password\": \"password123\"}"
echo.
echo.

echo 3. Testing User Login...
curl -s -X POST %BASE_URL%/api/v1/auth/login ^
  -H "Content-Type: application/json" ^
  -d "{\"email\": \"test@example.com\", \"password\": \"password123\"}" > login_response.json

type login_response.json
echo.

echo 4. Testing Protected Route - Profile...
echo Please copy the access_token from above and run:
echo curl -X GET %BASE_URL%/api/v1/profile -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
echo.

echo 5. Testing Protected Route - Dashboard...
echo curl -X GET %BASE_URL%/api/v1/dashboard -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
echo.

echo 6. Testing Token Refresh...
echo curl -X POST %BASE_URL%/api/v1/auth/refresh -H "Content-Type: application/json" -d "{\"refresh_token\": \"YOUR_REFRESH_TOKEN\"}"
echo.

echo 7. Testing Logout...
echo curl -X POST %BASE_URL%/api/v1/logout -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
echo.

echo === Testing Complete ===
del login_response.json 2>nul