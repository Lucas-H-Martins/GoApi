# Go API Study Project

This repository contains a RESTful API built with Go, designed as a study project to demonstrate various concepts and best practices in Go development.

## Project Structure

```
.
├── cmd/                    # Application entry points
│   ├── goapi/             # Main API server
│   │   └── main.go        # API server entrypoint
│   └── migrate/           # Database migration system
│       └── main.go        # Migration system entrypoint
├── config/                # Configuration management
│   └── database.go        # Database configuration
├── docs/                  # API documentation
│   ├── docs.go           # Swagger documentation
│   ├── swagger.json      # Swagger JSON specification
│   └── swagger.yaml      # Swagger YAML specification
├── handlers/             # HTTP request handlers
│   └── user_handler.go   # User-related request handlers
├── logger/               # Logging utilities
│   └── logger.go         # Custom logger implementation
├── middleware/           # HTTP middleware components
│   ├── auth.go          # Authentication middleware
│   └── logger.go        # Request logging middleware
├── migrations/           # Database migrations
│   ├── manager.go       # Migration system core logic
│   ├── down/            # Rollback migrations
│   └── up/              # Forward migrations
├── models/              # Data models
│   └── user.go          # User model definitions
├── repository/          # Data access layer
│   ├── user_repository.go    # User repository interface
│   └── users_sql/           # SQL queries for users
├── routes/              # Route definitions
│   ├── routes.go        # Main router setup
│   └── user_routes/     # User-specific routes
├── services/            # Business logic layer
│   └── user_service.go  # User-related business logic
├── env.*               # Environment configuration files
├── .gitignore         # Git ignore file
├── Dockerfile         # Docker configuration
├── docker-compose.yml # Docker Compose configuration
├── go.mod             # Go module file
├── go.sum             # Go module checksum
└── Makefile           # Build and development commands
```

## Features

### RESTful API Design
- Follows REST principles for resource management
- Standard HTTP methods (GET, POST, PUT, DELETE)
- Consistent response patterns
- Proper status codes and error handling

### Middleware Implementation
1. **Authentication Middleware**
   - Currently uses a fixed API key for demonstration
   - Designed to be easily extensible for various auth systems:
     - Auth0
     - Keycloak
     - Azure AD
     - Custom AD systems
     - JWT-based authentication

2. **Logger Middleware**
   - Comprehensive request/response logging
   - Includes:
     - Request method and path
     - Response status
     - Processing time
     - Request headers
     - Request body (when applicable)
     - Error details

### Database Layer
- PostgreSQL database
- Repository pattern implementation
- SQL query management
- Migration system with:
  - Forward and rollback migrations
  - Transaction support
  - Migration tracking
  - Safe execution (prevents duplicates)

### API Documentation
- Swagger/OpenAPI documentation
- Interactive API testing interface
- Available at `/swagger/*`

## Makefile Commands

The project includes a Makefile with the following commands:

```bash
# Build the application
make build

# Run the application
make run

# Run tests
make test

# Run linter
make lint

# Run database migrations
make migrate-up    # Apply migrations
make migrate-down  # Rollback migrations

# Generate Swagger documentation
make swagger

# Clean build artifacts
make clean
```

## Getting Started

1. Clone the repository
2. Set up environment variables:
   - Copy `env.local` to `.env`
   - Adjust values as needed
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Run database migrations:
   ```bash
   make migrate-up
   ```
5. Start the application:
   ```bash
   make run
   ```

## API Endpoints

### Users
- `GET /users` - List users (with pagination and filtering)
- `POST /users` - Create a new user
- `GET /users/{id}` - Get user by ID
- `PUT /users/{id}` - Update user
- `DELETE /users/{id}` - Delete user

## Development

### Prerequisites
- Go 1.21 or higher
- PostgreSQL
- Make

### Environment Variables
The project uses environment-specific configuration files:
- `env.local` - Local development
- `env.dev` - Development environment
- `env.prod` - Production environment

Required variables:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=goapi_db
DB_SSL_MODE=disable
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 

