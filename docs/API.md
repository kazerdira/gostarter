# API Documentation

Base URL: `http://localhost:8080/api/v1`

## Table of Contents
- [Authentication](#authentication)
- [Users](#users)
- [Error Responses](#error-responses)

---

## Authentication

### Register

Create a new user account.

**Endpoint:** `POST /auth/register`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123",
  "full_name": "John Doe"
}
```

**Response:** `201 Created`
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_at": "2024-01-01T12:15:00Z",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "full_name": "John Doe",
    "is_admin": false
  }
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123",
    "full_name": "John Doe"
  }'
```

### Login

Authenticate an existing user.

**Endpoint:** `POST /auth/login`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response:** `200 OK`
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_at": "2024-01-01T12:15:00Z",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "full_name": "John Doe",
    "is_admin": false
  }
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123"
  }'
```

### Refresh Token

Get a new access token using a refresh token.

**Endpoint:** `POST /auth/refresh`

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response:** `200 OK`
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_at": "2024-01-01T12:15:00Z",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "full_name": "John Doe",
    "is_admin": false
  }
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  }'
```

### Logout

Invalidate a refresh token.

**Endpoint:** `POST /auth/logout`

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response:** `200 OK`
```json
{
  "message": "logged out successfully"
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  }'
```

---

## Users

All user endpoints require authentication via Bearer token.

### Get Current User

Get the authenticated user's information.

**Endpoint:** `GET /users/me`

**Headers:**
```
Authorization: Bearer {access_token}
```

**Response:** `200 OK`
```json
{
  "id": 1,
  "email": "user@example.com",
  "full_name": "John Doe",
  "is_admin": false,
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z"
}
```

**cURL Example:**
```bash
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

### Update Current User

Update the authenticated user's information.

**Endpoint:** `PUT /users/me`

**Headers:**
```
Authorization: Bearer {access_token}
```

**Request Body:**
```json
{
  "full_name": "Jane Doe",
  "email": "jane@example.com"
}
```

**Note:** All fields are optional. Only send fields you want to update.

**Response:** `200 OK`
```json
{
  "id": 1,
  "email": "jane@example.com",
  "full_name": "Jane Doe",
  "is_admin": false,
  "updated_at": "2024-01-01T12:00:00Z"
}
```

**cURL Example:**
```bash
curl -X PUT http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Jane Doe"
  }'
```

### Delete Current User

Soft-delete the authenticated user's account.

**Endpoint:** `DELETE /users/me`

**Headers:**
```
Authorization: Bearer {access_token}
```

**Response:** `200 OK`
```json
{
  "message": "account deleted successfully"
}
```

**cURL Example:**
```bash
curl -X DELETE http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

### Get User by ID (Admin Only)

Get a specific user's information.

**Endpoint:** `GET /users/:id`

**Headers:**
```
Authorization: Bearer {access_token}
```

**Response:** `200 OK`
```json
{
  "id": 2,
  "email": "other@example.com",
  "full_name": "Other User",
  "is_admin": false,
  "is_active": true,
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z"
}
```

**cURL Example:**
```bash
curl -X GET http://localhost:8080/api/v1/users/2 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

### List Users (Admin Only)

Get a paginated list of users.

**Endpoint:** `GET /users?page={page}&limit={limit}`

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20, max: 100)

**Headers:**
```
Authorization: Bearer {access_token}
```

**Response:** `200 OK`
```json
{
  "users": [
    {
      "id": 1,
      "email": "user1@example.com",
      "full_name": "User One",
      "is_admin": false,
      "is_active": true,
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    },
    {
      "id": 2,
      "email": "user2@example.com",
      "full_name": "User Two",
      "is_admin": false,
      "is_active": true,
      "created_at": "2024-01-01T11:00:00Z",
      "updated_at": "2024-01-01T11:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 42,
    "total_pages": 3
  }
}
```

**cURL Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/users?page=1&limit=20" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

---

## Health Checks

### Health Check

Basic health check endpoint (no auth required).

**Endpoint:** `GET /health`

**Response:** `200 OK`
```json
{
  "status": "healthy",
  "env": "development"
}
```

**cURL Example:**
```bash
curl -X GET http://localhost:8080/health
```

### Readiness Check

Check if the service is ready to accept traffic (includes database check).

**Endpoint:** `GET /ready`

**Response:** `200 OK`
```json
{
  "status": "ready"
}
```

**cURL Example:**
```bash
curl -X GET http://localhost:8080/ready
```

---

## Error Responses

All error responses follow this format:

```json
{
  "error": "error message here"
}
```

### Common Status Codes

- `400 Bad Request` - Invalid request format or missing required fields
- `401 Unauthorized` - Missing or invalid authentication token
- `403 Forbidden` - Insufficient permissions (e.g., not an admin)
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource already exists (e.g., email already registered)
- `500 Internal Server Error` - Server error

### Examples

**Invalid Email Format:**
```json
{
  "error": "Key: 'RegisterRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"
}
```

**Unauthorized Access:**
```json
{
  "error": "authorization header required"
}
```

**User Not Found:**
```json
{
  "error": "user not found"
}
```

**Email Already Registered:**
```json
{
  "error": "email already registered"
}
```

---

## Authentication Flow

### Typical Flow

1. **Register or Login** to get tokens
2. **Use Access Token** for API requests (valid for 15 minutes)
3. **Refresh Token** when access token expires
4. **Logout** to invalidate refresh token

### Token Usage

Include the access token in the `Authorization` header:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Token Expiry

- **Access Token**: 15 minutes (configurable)
- **Refresh Token**: 7 days (configurable)

When the access token expires, use the `/auth/refresh` endpoint to get a new one.

---

## Rate Limiting

Default rate limits:
- 100 requests per minute per IP address

Exceeding the rate limit returns `429 Too Many Requests`.

---

## Testing with Postman

1. Import the API endpoints into Postman
2. Create an environment with:
   - `base_url`: http://localhost:8080/api/v1
   - `access_token`: (set after login)
3. Register/login to get tokens
4. Add `{{access_token}}` to Authorization header in subsequent requests

---

Need help? Check the main README or open an issue!
