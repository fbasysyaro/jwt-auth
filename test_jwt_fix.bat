@echo off
echo Testing JWT Fix...
echo.

echo 1. Testing Registration...
curl -s -X POST http://localhost:8080/api/v1/auth/register ^
  -H "Content-Type: application/json" ^
  -d "{\"username\": \"testuser\", \"email\": \"test@example.com\", \"password\": \"password123\"}" > register_response.json

echo Registration Response:
type register_response.json
echo.
echo.

echo 2. Testing Login...
curl -s -X POST http://localhost:8080/api/v1/auth/login ^
  -H "Content-Type: application/json" ^
  -d "{\"email\": \"test@example.com\", \"password\": \"password123\"}" > login_response.json

echo Login Response:
type login_response.json
echo.
echo.

echo 3. Check if tokens are real JWT tokens (should start with eyJ):
findstr "access_token" login_response.json
echo.

echo Test complete! 
echo - Real JWT tokens should start with "eyJ"
echo - Mock tokens looked like "access_1_email@domain.com"
echo.

del register_response.json login_response.json 2>nul