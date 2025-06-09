# Quote API Documentation

## Base URL
```
http://localhost:8080
```

## Authentication
All protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

## Endpoints

### Authentication

#### Register User
```http
POST /register
Content-Type: application/json

{
    "username": "string",
    "password": "string"
}
```

**Response (201 Created)**
```json
{
    "message": "User registered successfully"
}
```

**Error Responses**
- 400 Bad Request: Invalid input
- 400 Bad Request: Username already exists

#### Login
```http
POST /login
Content-Type: application/json

{
    "username": "string",
    "password": "string"
}
```

**Response (200 OK)**
```json
{
    "token": "string",
    "user": {
        "id": "number",
        "username": "string"
    }
}
```

**Error Responses**
- 400 Bad Request: Invalid input
- 401 Unauthorized: Invalid credentials

### Quotes

#### Create Quote
```http
POST /quotes
Authorization: Bearer <token>
Content-Type: application/json

{
    "content": "string",
    "author": "string"
}
```

**Response (201 Created)**
```json
{
    "id": "number",
    "content": "string",
    "author": "string",
    "created_at": "string",
    "updated_at": "string"
}
```

**Error Responses**
- 400 Bad Request: Invalid input
- 401 Unauthorized: Missing or invalid token

#### Get All Quotes
```http
GET /quotes
Authorization: Bearer <token>
```

**Response (200 OK)**
```json
[
    {
        "id": "number",
        "content": "string",
        "author": "string",
        "created_at": "string",
        "updated_at": "string"
    }
]
```

**Error Responses**
- 401 Unauthorized: Missing or invalid token

#### Get Quote by ID
```http
GET /quotes/{id}
Authorization: Bearer <token>
```

**Response (200 OK)**
```json
{
    "id": "number",
    "content": "string",
    "author": "string",
    "created_at": "string",
    "updated_at": "string"
}
```

**Error Responses**
- 401 Unauthorized: Missing or invalid token
- 404 Not Found: Quote not found

#### Update Quote
```http
PUT /quotes/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
    "content": "string",
    "author": "string"
}
```

**Response (200 OK)**
```json
{
    "id": "number",
    "content": "string",
    "author": "string",
    "created_at": "string",
    "updated_at": "string"
}
```

**Error Responses**
- 400 Bad Request: Invalid input
- 401 Unauthorized: Missing or invalid token
- 404 Not Found: Quote not found

#### Delete Quote
```http
DELETE /quotes/{id}
Authorization: Bearer <token>
```

**Response (200 OK)**
```json
{
    "message": "Quote deleted successfully"
}
```

**Error Responses**
- 401 Unauthorized: Missing or invalid token
- 404 Not Found: Quote not found

### Health Check

#### Check API Status
```http
GET /health
```

**Response (200 OK)**
```json
{
    "status": "ok"
}
```

## Data Models

### User
```typescript
interface User {
    id: number;
    username: string;
    created_at: string;
    updated_at: string;
}
```

### Quote
```typescript
interface Quote {
    id: number;
    content: string;
    author: string;
    created_at: string;
    updated_at: string;
}
```

## Error Responses
All error responses follow this format:
```json
{
    "error": "string"
}
```

## Authentication Flow
1. Register a new user using `/register`
2. Login using `/login` to get a JWT token
3. Include the token in the Authorization header for all protected endpoints
4. Token expires after 24 hours

## Rate Limiting
Currently, there is no rate limiting implemented.

## CORS
CORS is enabled for all origins in development mode.

## Environment Variables
Required environment variables:
```
PORT=8080
JWT_SECRET=your-secret-key
GIN_MODE=debug
``` 