@echo off
echo Testing JWT Auth Service v2 with Refresh Tokens
echo =================================================
echo.

echo 1. Register User
curl -s -X POST http://localhost:8080/api/v1/register ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"refreshuser\",\"email\":\"refresh@example.com\",\"password\":\"password123\"}" > register.json
echo Registration Response:
type register.json
echo.
echo.

echo 2. Login User  
curl -s -X POST http://localhost:8080/api/v1/login ^
  -H "Content-Type: application/json" ^
  -d "{\"email\":\"refresh@example.com\",\"password\":\"password123\"}" > login.json
echo Login Response:
type login.json
echo.
echo.

echo 3. Extract tokens
for /f "tokens=2 delims=:," %%a in ('findstr "access_token" login.json') do set ACCESS_TOKEN=%%a
for /f "tokens=2 delims=:," %%b in ('findstr "refresh_token" login.json') do set REFRESH_TOKEN=%%b
set ACCESS_TOKEN=%ACCESS_TOKEN:"=%
set ACCESS_TOKEN=%ACCESS_TOKEN: =%
set REFRESH_TOKEN=%REFRESH_TOKEN:"=%
set REFRESH_TOKEN=%REFRESH_TOKEN: =%

echo 4. Verify Access Token
curl -s -X POST http://localhost:8080/api/v1/verify ^
  -H "Authorization: Bearer %ACCESS_TOKEN%"
echo.
echo.

echo 5. Use Refresh Token to get new tokens
curl -s -X POST http://localhost:8080/api/v1/refresh ^
  -H "Content-Type: application/json" ^
  -d "{\"refresh_token\":\"%REFRESH_TOKEN%\"}" > refresh.json
echo Refresh Response:
type refresh.json
echo.
echo.

echo 6. Logout (blacklist current access token)
curl -s -X POST http://localhost:8080/api/v1/logout ^
  -H "Authorization: Bearer %ACCESS_TOKEN%"
echo.
echo.

echo 7. Try to verify after logout (should fail)
curl -s -X POST http://localhost:8080/api/v1/verify ^
  -H "Authorization: Bearer %ACCESS_TOKEN%"
echo.

del register.json login.json refresh.json 2>nul
echo.
echo Test Complete!