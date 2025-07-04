@echo off
echo Testing Simplified JWT Auth Service
echo ====================================
echo.

echo 1. Health Check
curl -s http://localhost:8080/health
echo.
echo.

echo 2. Register User
curl -s -X POST http://localhost:8080/api/v1/register ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"testuser2\",\"email\":\"test2@example.com\",\"password\":\"password123\"}" > register.json
echo Registration Response:
type register.json
echo.
echo.

echo 3. Login User  
curl -s -X POST http://localhost:8080/api/v1/login ^
  -H "Content-Type: application/json" ^
  -d "{\"email\":\"test2@example.com\",\"password\":\"password123\"}" > login.json
echo Login Response:
type login.json
echo.
echo.

echo 4. Extract token for verification
for /f "tokens=2 delims=:," %%a in ('findstr "token" login.json') do set TOKEN=%%a
set TOKEN=%TOKEN:"=%
set TOKEN=%TOKEN: =%

echo 5. Verify Token
curl -s -X POST http://localhost:8080/api/v1/verify ^
  -H "Authorization: Bearer %TOKEN%"
echo.
echo.

echo 6. Logout
curl -s -X POST http://localhost:8080/api/v1/logout ^
  -H "Authorization: Bearer %TOKEN%"
echo.
echo.

echo 7. Try to verify after logout (should fail)
curl -s -X POST http://localhost:8080/api/v1/verify ^
  -H "Authorization: Bearer %TOKEN%"
echo.

del register.json login.json 2>nul
echo.
echo Test Complete!