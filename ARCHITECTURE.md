# Clean Architecture Structure

This project follows Clean Architecture principles with clear separation of concerns:

## Directory Structure

```
internal/
├── domain/                    # Enterprise Business Rules
│   ├── entities/             # Business entities
│   │   └── user.go
│   ├── repositories/         # Repository interfaces
│   │   └── user_repository.go
│   └── services/            # Domain service interfaces
│       ├── auth_service.go
│       ├── external_services.go
│       └── token_blacklist_service.go
│
├── application/              # Application Business Rules
│   ├── dto/                 # Data Transfer Objects
│   │   └── auth_dto.go
│   └── services/            # Application service implementations
│       ├── auth_service_impl.go
│       └── *_test.go
│
├── infrastructure/           # Frameworks & Drivers
│   ├── database/            # Database implementations
│   │   └── postgres.go
│   ├── email/               # Email service implementations
│   │   └── email_service_impl.go
│   ├── jwt/                 # JWT implementations
│   │   └── jwt_manager_impl.go
│   ├── redis/               # Redis implementations
│   │   ├── redis.go
│   │   └── token_blacklist_service_impl.go
│   └── repositories/        # Repository implementations
│       └── user_repository_impl.go
│
└── interfaces/              # Interface Adapters
    ├── config/              # Configuration
    │   └── config.go
    └── http/                # HTTP layer
        ├── handlers/        # HTTP handlers
        ├── middleware/      # HTTP middleware
        ├── routes/          # Route definitions
        └── docs/            # API documentation
```

## Layer Dependencies

- **Domain**: No dependencies on other layers
- **Application**: Depends only on Domain
- **Infrastructure**: Implements Domain interfaces
- **Interfaces**: Orchestrates all layers

## Key Principles Applied

1. **Dependency Inversion**: High-level modules don't depend on low-level modules
2. **Interface Segregation**: Small, focused interfaces
3. **Single Responsibility**: Each layer has one reason to change
4. **Open/Closed**: Open for extension, closed for modification