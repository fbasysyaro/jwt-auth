@echo off
echo Checking what is using port 8080...
netstat -ano | findstr :8080

echo.
echo Docker containers status:
docker ps

echo.
echo To kill process using port 8080, find the PID above and run:
echo taskkill /PID [PID_NUMBER] /F