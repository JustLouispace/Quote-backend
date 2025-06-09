# Quote Backend API

A RESTful API for managing quotes, built with Go, Gin, and GORM.

## Features

- 🔐 JWT Authentication
- 📝 CRUD operations for quotes
- 🛡️ Protected routes
- 💾 SQLite database
- 📚 Clean architecture
- 🔄 CORS enabled

## Prerequisites

- Go 1.20 or higher
- Git

## Installation

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
```

## Running the Application

Start the server:
```bash
go run main.go
```

The server will start on port 8080 by default.

## Project Structure

```
.
├── config/         # Configuration files
│   └── database.go # Database configuration
├── handlers/       # HTTP request handlers
│   ├── auth.go    # Authentication handlers
│   └── quote.go   # Quote handlers
├── middleware/     # Custom middleware
│   └── auth.go    # Authentication middleware
├── models/         # Data models
│   ├── quote.go   # Quote model
│   └── user.go    # User model
├── main.go        # Application entry point
├── go.mod         # Go module file
└── .env           # Environment variables
```

## API Documentation

For detailed API documentation, see [API.md](API.md).

### Quick Start

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

To run tests:
```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

Your Name - [@your_twitter](https://twitter.com/your_twitter)

Project Link: [https://github.com/JustLouispace/Qoute-backend](https://github.com/JustLouispace/Qoute-backend) 