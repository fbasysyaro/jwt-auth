@echo off
echo Testing Simplified JWT Auth Service
echo.

echo 1. Register User
curl -X POST http://localhost:8080/api/v1/register ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"testuser\",\"email\":\"test@example.com\",\"password\":\"password123\"}" > register.json
type register.json
echo.
echo.

echo 2. Login User  
curl -X POST http://localhost:8080/api/v1/login ^
  -H "Content-Type: application/json" ^
  -d "{\"email\":\"test@example.com\",\"password\":\"password123\"}" > login.json
type login.json
echo.
echo.

echo 3. Extract token and verify
echo Please copy the token from above and test verify endpoint:
echo curl -X POST http://localhost:8080/api/v1/verify -H "Authorization: Bearer YOUR_TOKEN"
echo.

echo 4. Test logout
echo curl -X POST http://localhost:8080/api/v1/logout -H "Authorization: Bearer YOUR_TOKEN"
echo.

del register.json login.json 2>nul