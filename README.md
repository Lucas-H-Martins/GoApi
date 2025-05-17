# GoAPI

A simple Go API using MVC pattern with authentication.

## Getting Started

1. Install dependencies:
```bash
go mod tidy
```

2. Run the application:
```bash
go run main.go
```

The server will start on port 8080.

## Available Endpoints

### Public Endpoints
- GET `/hello` - Returns a hello world message (no authentication required)

### Protected Endpoints
- GET `/protected/hello` - Returns a hello world message (requires authentication)

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

Protected endpoint:
```bash
curl -H "Authorization: Bearer your-secret-token-123" http://localhost:8080/protected/hello
```

## Example Responses

Success Response:
```json
{
    "message": "Hello, World!"
}
```

Unauthorized Response:
```json
{
    "error": "No authorization header provided"
}
``` 