@echo off
echo Stopping existing containers...
docker-compose down

echo Building and starting containers...
docker-compose up --build -d

echo Waiting for services to start...
timeout /t 10

echo Services started! Access Swagger at:
echo http://localhost:8080/api/v1/swagger/index.html

echo.
echo To view logs:
echo docker-compose logs -f api

echo.
echo To stop services:
echo docker-compose down