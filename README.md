# Qoute Backend

A Go-based backend service built with Gin framework.

## Prerequisites

- Go 1.21 or higher
- Git

## Project Structure

```
.
├── main.go           # Application entry point
├── go.mod           # Go module file
├── .env             # Environment variables
├── handlers/        # HTTP request handlers
├── models/          # Data models
├── middleware/      # Custom middleware
└── config/          # Configuration files
```

## Setup

1. Clone the repository:
```bash
git clone https://github.com/yourusername/Qoute-backend.git
cd Qoute-backend
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file in the root directory with the following content:
```
PORT=8080
GIN_MODE=debug
```

## Running the Application

To start the server:
```bash
go run main.go
```

The server will start on port 8080 by default. You can change this by modifying the PORT in the .env file.

## API Endpoints

- `GET /health` - Health check endpoint

## Development

To add new dependencies:
```bash
go get github.com/package/name
```

To tidy up dependencies:
```bash
go mod tidy
``` 