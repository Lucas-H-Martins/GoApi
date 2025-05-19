# GoAPI

A simple Go API using MVC pattern with authentication and PostgreSQL database.

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Bash (for environment setup script)

## Environment Setup

The application supports three environments:
- `local`: Local development
- `dev`: Development server
- `prod`: Production server

To set up an environment:

1. Run the setup script:
```bash
chmod +x scripts/setup_env.sh
./scripts/setup_env.sh <environment>
```

For `dev` and `prod` environments, the script will prompt for database credentials and store them securely.

### Environment Files

Environment configurations are stored in `config/env.<environment>` files:
- `config/env.local`: Local development settings
- `config/env.dev`: Development server settings
- `config/env.prod`: Production server settings

### Secret Management

Values prefixed with `!` in environment files are treated as paths to secret files. For example:
```
DB_PASSWORD=!secrets/dev/db_password
```
This will look for the password in the file `secrets/dev/db_password`.

## Database Setup

1. Create a PostgreSQL database:
```bash
createdb goapi_db
```

2. Run the migrations:
```bash
psql -d goapi_db -f migrations/001_create_users_table.sql
```

## Getting Started

1. Install dependencies:
```bash
go mod tidy
```

2. Set up your environment:
```bash
./scripts/setup_env.sh local  # or dev/prod
```

3. Run the application:
```bash
go run main.go
```

The server will start on the configured port (default: 8080 for local/dev, 80 for prod).

## Available Endpoints

### Public Endpoints
- GET `/hello` - Returns a hello world message (no authentication required)

### Protected Endpoints (require authentication)
#### Users
- POST `/users` - Create a new user
- GET `/users` - List all users
- GET `/users/{id}` - Get a specific user
- PUT `/users/{id}` - Update a user
- DELETE `/users/{id}` - Delete a user

## API Documentation

The API documentation is available through Swagger UI. To access it:

1. Start the server:
```bash
go run ./cmd/goapi/main.go
```

2. Open your browser and navigate to:
```
http://localhost:8080/swagger/index.html
```

The Swagger UI provides:
- Interactive documentation of all API endpoints
- Request/response schemas
- Ability to test endpoints directly from the browser
- Detailed parameter descriptions and requirements

## Authentication

Protected endpoints require a Bearer token in the Authorization header:

```
Authorization: Bearer your-secret-token-123
```

### Example cURL Requests

Public endpoint:
```bash
curl http://localhost:8080/hello
```

Protected endpoints:
```bash
# List users
curl -H "Authorization: Bearer your-secret-token-123" http://localhost:8080/users

# Create user
curl -X POST -H "Authorization: Bearer your-secret-token-123" \
     -H "Content-Type: application/json" \
     -d '{"name":"John Doe","email":"john@example.com"}' \
     http://localhost:8080/users

# Get user by ID
curl -H "Authorization: Bearer your-secret-token-123" http://localhost:8080/users/1

# Update user
curl -X PUT -H "Authorization: Bearer your-secret-token-123" \
     -H "Content-Type: application/json" \
     -d '{"name":"John Updated","email":"john.updated@example.com"}' \
     http://localhost:8080/users/1

# Delete user
curl -X DELETE -H "Authorization: Bearer your-secret-token-123" http://localhost:8080/users/1
```

## Example Responses

Success Response:
```json
{
    "message": "Hello, World!"
}
```

User Response:
```json
{
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2024-03-21T10:00:00Z",
    "updated_at": "2024-03-21T10:00:00Z"
}
```

Unauthorized Response:
```json
{
    "error": "No authorization header provided"
}
``` 

