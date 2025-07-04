@echo off
echo Testing Corrected Login Flow
echo ============================
echo.

echo 1. Register User (to get access token)
curl -s -X POST http://localhost:8080/api/v1/register ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"flowtest\",\"email\":\"flowtest@example.com\",\"password\":\"password123\"}" > register.json
echo Registration Response:
type register.json
echo.
echo.

echo 2. Extract access token from registration
for /f "tokens=2 delims=:," %%a in ('findstr "access_token" register.json') do set ACCESS_TOKEN=%%a
set ACCESS_TOKEN=%ACCESS_TOKEN:"=%
set ACCESS_TOKEN=%ACCESS_TOKEN: =%

echo 3. Login with username, password, and access token
curl -s -X POST http://localhost:8080/api/v1/login ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"flowtest\",\"password\":\"password123\",\"access_token\":\"%ACCESS_TOKEN%\"}"
echo.
echo.

echo 4. Test login with wrong password (should fail)
curl -s -X POST http://localhost:8080/api/v1/login ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"flowtest\",\"password\":\"wrongpassword\",\"access_token\":\"%ACCESS_TOKEN%\"}"
echo.
echo.

echo 5. Test login without access token (should fail)
curl -s -X POST http://localhost:8080/api/v1/login ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"flowtest\",\"password\":\"password123\"}"
echo.

del register.json 2>nul
echo.
echo Test Complete!