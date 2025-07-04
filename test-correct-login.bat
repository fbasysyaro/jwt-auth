@echo off
echo Testing Correct Login Flow
echo ==========================
echo.

echo 1. Register User (to get access token)
curl -s -X POST http://localhost:8080/api/v1/register ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"correcttest\",\"email\":\"correcttest@example.com\",\"password\":\"password123\"}" > register.json
echo Registration Response:
type register.json
echo.
echo.

echo 2. Extract access token from registration
for /f "tokens=2 delims=:," %%a in ('findstr "access_token" register.json') do set ACCESS_TOKEN=%%a
set ACCESS_TOKEN=%ACCESS_TOKEN:"=%
set ACCESS_TOKEN=%ACCESS_TOKEN: =%

echo 3. Login with username/password in body + access token in Authorization header
curl -s -X POST http://localhost:8080/api/v1/login ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer %ACCESS_TOKEN%" ^
  -d "{\"username\":\"correcttest\",\"password\":\"password123\"}"
echo.
echo.

echo 4. Test login without Authorization header (should fail)
curl -s -X POST http://localhost:8080/api/v1/login ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"correcttest\",\"password\":\"password123\"}"
echo.

del register.json 2>nul
echo.
echo Test Complete!