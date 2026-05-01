# Go-User-Loader

A simple HTTP service written in Go that loads users from CSV or JSON files and exposes them via REST endpoints

## Features

- Load users from CSV or JSON files
- Expose users via HTTP API
- Health check endpoint
- Proper HTTP status codes and error handling

## Requirements

- Go 1.22+

## Running the service

```
go run ./cmd/main.go -file users.csv -port 8080



```

### CLI Flags

- `-file` (required): Path to users file (`.csv` or `.json`)
- `-port` (optional): Port the server listens on (default: 8080)


Response example:
```json
{
  "count": 2,
  "items": [
    {
      "id": 1,
      "first_name": "John",
      "last_name": "Doe",
      "email": "john@example.com"
    }
  ]
}
## API Endpoints

### GET /users
Returns the list of users.