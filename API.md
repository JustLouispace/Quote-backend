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

**Query Parameters**
- `author` (string, optional): Filter quotes by author.
- `search` (string, optional): Search for a term in quote content and author.
- `sortBy` (string, optional): Field to sort by (`created_at`, `author`, `content`). Defaults to `created_at`.
- `order` (string, optional): Sort order (`asc` or `desc`). Defaults to `desc`.

**Response (200 OK)**
```json
[
    {
        "id": "number",
        "content": "string",
        "author": "string",
        "vote_count": "number",
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

### Votes

#### Create Vote
```http
POST /quotes/{id}/vote
Authorization: Bearer <token>
```

**Response (201 Created)**
```json
{
    "message": "Vote recorded successfully",
    "voteCount": 1,
    "vote": {
        "id": "number",
        "user_id": "number",
        "quote_id": "number",
        "created_at": "string",
        "updated_at": "string",
        "user": {
            "id": "number",
            "username": "string"
        }
    }
}
```

**Error Responses**
- 400 Bad Request: Invalid quote ID
- 401 Unauthorized: Missing or invalid token
- 404 Not Found: Quote not found
- 409 Conflict: User has already voted for this quote
- 409 Conflict: Voting is only allowed when the quote has 0 votes

#### Delete Vote
```http
DELETE /quotes/{id}/vote
Authorization: Bearer <token>
```

**Response (200 OK)**
```json
{
    "message": "Vote removed successfully",
    "voteCount": 0
}
```

**Error Responses**
- 400 Bad Request: Invalid quote ID
- 401 Unauthorized: Missing or invalid token
- 404 Not Found: Vote not found

#### Get Vote Count
```http
GET /quotes/{id}/vote/count
Authorization: Bearer <token>
```

**Response (200 OK)**
```json
{
    "count": 5
}
```

**Error Responses**
- 400 Bad Request: Invalid quote ID
- 404 Not Found: Quote not found

#### Check User Vote
```http
GET /quotes/{id}/vote/check
Authorization: Bearer <token>
```

**Response (200 OK)**
```json
{
    "has_voted": true
}
```

**Error Responses**
- 400 Bad Request: Invalid quote ID
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
    votes: Vote[];
    vote_count: number;
    created_at: string;
    updated_at: string;
}

interface Vote {
    id: number;
    user_id: number;
    quote_id: number;
    created_at: string;
    updated_at: string;
    user: {
        id: number;
        username: string;
    };
    quote: {
        id: number;
        content: string;
        author: string;
    };
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
3. Include the token in the Authorization header for all protected endpoints:
   ```
   Authorization: Bearer <your_jwt_token>
   ```
4. Token expires after 24 hours

## Rate Limiting
Currently, there is no rate limiting implemented.

## CORS
CORS is enabled for the following origins:
- http://localhost:3000
- http://127.0.0.1:3000
- http://localhost:5173
- http://127.0.0.1:5173
- https://quote-frontend-zeta.vercel.app

## Environment Variables
Required environment variables:
```
PORT=8080
GIN_MODE=debug
JWT_SECRET=your-secret-key-here
DATABASE_DSN=quotes.db
``` 