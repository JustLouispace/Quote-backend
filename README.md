# Quote Backend API

A RESTful API for managing quotes, built with Go, Gin, and GORM. This project provides user authentication, CRUD operations for quotes, and a voting system. It uses a clean architecture and SQLite for storage.

## Features
- 🔒 JWT Authentication
- 📝 CRUD operations for quotes
- 🗳️ Voting system for quotes
- 🛡️ Protected routes
- 💾 SQLite database
- 📚 Clean architecture
- 🔄 CORS enabled

## Prerequisites
- Go 1.20 or higher
- Git

## Quick Start
1. Clone the repository:
   ```bash
git clone https://github.com/JustLouispace/Qoute-backend.git
cd Qoute-backend
```
2. Install dependencies:
   ```bash
go mod download
```
3. Create a `.env` file in the root directory:
   ```env
PORT=8080
JWT_SECRET=your-secret-key
GIN_MODE=debug
DATABASE_DSN=quotes.db   # Use ':memory:' to run with an in-memory database
```
4. Run the server:
   ```bash
go run main.go
```
The server will start on port 8080 by default. If you set `DATABASE_DSN=:memory:`, the database will be in-memory and reset on each run.

## API Reference
See [API.md](./API.md) for detailed documentation of all endpoints, request/response formats, authentication, and error codes.

## Project Structure
```
.
├── config/         # Configuration files
│   └── database.go # Database configuration
├── handlers/       # HTTP request handlers
│   ├── auth.go     # Authentication handlers
│   ├── quote.go    # Quote handlers
│   └── vote.go     # Voting handlers
├── middleware/     # Custom middleware
│   └── auth.go     # Authentication middleware
├── models/         # Data models
│   ├── quote.go    # Quote model
│   ├── user.go     # User model
│   └── vote.go     # Vote model
├── main.go         # Application entrypoint
├── go.mod          # Go module definition
├── go.sum          # Go module checksums
├── .env            # Environment variables (not committed)
├── README.md       # Project documentation
├── API.md          # API documentation
## API Reference
See [API.md](./API.md) for detailed documentation of all endpoints, request/response formats, authentication, and error codes.

## Quick Start
1. Register a new user:
   ```bash
   curl -X POST http://localhost:8080/register \
     -H "Content-Type: application/json" \
     -d '{"username":"testuser","password":"password123"}'
   ```
2. Login to get JWT token:
   ```bash
   curl -X POST http://localhost:8080/login \
     -H "Content-Type: application/json" \
     -d '{"username":"testuser","password":"password123"}'
   ```
3. Create a quote (with JWT token):
   ```bash
   curl -X POST http://localhost:8080/quotes \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"content":"The only way to do great work is to love what you do.","author":"Steve Jobs"}'
   ```

## API Overview

| Endpoint                   | Method | Description                 | Auth Required |
|----------------------------|--------|-----------------------------|--------------|
| `/register`                | POST   | Register a new user         | No           |
| `/login`                   | POST   | User login (get JWT)        | No           |
| `/quotes`                  | GET    | List all quotes             | Yes          |
| `/quotes`                  | POST   | Create a new quote          | Yes          |
| `/quotes/{id}`             | GET    | Get quote by ID             | Yes          |
| `/quotes/{id}`             | PUT    | Update a quote              | Yes          |
| `/quotes/{id}`             | DELETE | Delete a quote              | Yes          |
| `/quotes/{id}/vote`        | POST   | Vote for a quote            | Yes          |
| `/quotes/{id}/vote`        | DELETE | Remove vote from a quote    | Yes          |
| `/quotes/{id}/vote/count`  | GET    | Get vote count for a quote  | Yes          |
| `/quotes/{id}/vote/check`  | GET    | Check if user voted         | Yes          |
| `/health`                  | GET    | Health check                | No           |

For full details, see [API.md](./API.md).


## Development

### Adding Dependencies
```bash
go get github.com/package/name
```

### Updating Dependencies
```bash
go mod tidy
```

## Testing

Automated tests are now included! To run all tests, use:

```bash
go test ./...
```

An example test is provided for the `/health` endpoint in `main_test.go`. You can add more tests for other endpoints and handlers as needed.

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License
There is currently **no license file** in this repository. If you intend to use or distribute this project, please add an appropriate license.

Project Link: [https://github.com/JustLouispace/Qoute-backend](https://github.com/JustLouispace/Qoute-backend)