.PHONY: build run stop clean test

# Docker commands
build:
	docker-compose build

run:
	docker-compose up -d

stop:
	docker-compose down

logs:
	docker-compose logs -f api

clean: stop
	docker system prune -f
	docker volume rm jwt-auth_postgres_data

# Local development commands
dev:
	go run cmd/main.go

test:
	go test -v ./...

# Docker image commands
docker-push:
	docker tag jwt-auth:latest your-registry/jwt-auth:latest
	docker push your-registry/jwt-auth:latest
