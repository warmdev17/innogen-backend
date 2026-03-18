# Innogen Backend API Documentation

## Overview

The Innogen Backend API is a RESTful API for a competitive programming platform. It provides endpoints for user authentication, problem management, code submission, and execution.

## Base URL

```
http://localhost:8081/api
```

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. To access protected endpoints, include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Endpoints

### 1. Authentication

#### POST /auth/login

Authenticate a user and receive a JWT token.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Success Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error Response (401):**
```json
{
  "error": "Invalid credentials"
}
```

#### GET /me

Get information about the currently authenticated user.

**Headers:**
- Authorization: Bearer <token>

**Success Response (200):**
```json
{
  "id": 1,
  "email": "user@example.com",
  "username": "johndoe",
  "fullName": "John Doe",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

**Error Response (401):**
```json
{
  "error": "Unauthorized"
}
```

### 2. Problems

#### GET /problems

Get a list of all available problems.

**Success Response (200):**
```json
[
  {
    "id": 1,
    "slug": "two-sum",
    "title": "Two Sum",
    "difficulty": "easy",
    "problemMd": "## Problem Description\n...",
    "timeLimitMs": 1000,
    "memoryLimitKb": 256,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
]
```

#### POST /admin/problems

Create a new problem (admin/teacher only).

**Headers:**
- Authorization: Bearer <token>

**Request Body:**
```json
{
  "slug": "fibonacci",
  "title": "Fibonacci Sequence",
  "difficulty": "medium",
  "problemMd": "## Problem Description\n...",
  "timeLimitMs": 2000,
  "memoryLimitKb": 512,
  "testcases": [
    {
      "input": "5",
      "expectedOutput": "0\n1\n1\n2\n3\n"
    }
  ]
}
```

**Success Response (201):**
```json
{
  "id": 2,
  "slug": "fibonacci",
  "title": "Fibonacci Sequence",
  "difficulty": "medium",
  "problemMd": "## Problem Description\n...",
  "timeLimitMs": 2000,
  "memoryLimitKb": 512,
  "testcases": [...],
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

**Error Responses:**
- 400: Invalid input
- 403: Insufficient permissions

### 3. Submissions

#### POST /submit

Submit code for judging against test cases.

**Headers:**
- Authorization: Bearer <token>

**Request Body:**
```json
{
  "problemId": 1,
  "language": "python",
  "code": "def two_sum(nums, target):\n    # solution here\n    pass"
}
```

**Success Response (201):**
```json
{
  "message": "Submission queued",
  "submission": {
    "id": "uuid-string",
    "userId": 1,
    "problemId": 1,
    "code": "def two_sum(nums, target):\n    # solution here\n    pass",
    "language": "python",
    "status": "pending",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```

**Supported Languages:**
- `python`
- `javascript`
- `go`
- `cpp`

#### GET /submit/{id}

Get the status and results of a specific submission.

**Headers:**
- Authorization: Bearer <token>

**Path Parameters:**
- id: Submission UUID

**Success Response (200):**
```json
{
  "id": "uuid-string",
  "userId": 1,
  "problemId": 1,
  "code": "def two_sum(nums, target):\n    # solution here\n    pass",
  "language": "python",
  "status": "accepted",
  "testResults": [
    {
      "testNumber": 1,
      "status": "passed",
      "input": "2 7",
      "expectedOutput": "0 1",
      "actualOutput": "0 1",
      "executionTime": 12.5
    }
  ],
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

**Status Values:**
- `pending`: Waiting in queue
- `running`: Currently being judged
- `accepted`: All test cases passed
- `wrong_answer`: Some test cases failed
- `time_limit_exceeded`: Program exceeded time limit
- `memory_limit_exceeded`: Program exceeded memory limit
- `runtime_error`: Program crashed
- `compilation_error`: Code failed to compile

**Error Responses:**
- 403: Not authorized to view this submission
- 404: Submission not found

### 4. Code Execution

#### POST /run

Execute code directly without test cases (useful for testing).

**Headers:**
- Authorization: Bearer <token>

**Request Body:**
```json
{
  "language": "python",
  "code": "print('Hello, World!')",
  "input": ""
}
```

**Success Response (200):**
```json
{
  "output": "Hello, World!\n",
  "error": "",
  "executionTime": 5.2,
  "memoryUsed": 1024
}
```

**Error Responses:**
- 400: Invalid input
- 500: Execution failed

## Error Handling

All endpoints return errors in the following format:

```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:
- `200`: Success
- `201`: Created successfully
- `400`: Bad Request (invalid input)
- `401`: Unauthorized (invalid or missing token)
- `403`: Forbidden (insufficient permissions)
- `404`: Not Found
- `500`: Internal Server Error

## Rate Limiting

The API may enforce rate limiting to prevent abuse. If you receive a 429 status code, you have exceeded the allowed number of requests.

## Swagger Documentation

Interactive API documentation is available at:
- Swagger UI: `/swagger/index.html`
- API Docs (local): `/docs/index.html`
- JSON Spec: `/swagger/doc.json`
- YAML Spec: `/swagger/doc.yaml`

## Testing

A test client is available at `/test-client.html` for interactive testing of all endpoints.

## SDKs and Libraries

### JavaScript/Node.js Example

```javascript
const API_BASE = 'http://localhost:8081/api';

// Login
const loginResponse = await fetch(`${API_BASE}/auth/login`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ email, password })
});
const { token } = await loginResponse.json();

// Get problems
const problemsResponse = await fetch(`${API_BASE}/problems`);
const problems = await problemsResponse.json();

// Submit code
const submitResponse = await fetch(`${API_BASE}/submit`, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({ problemId, language, code })
});
const submission = await submitResponse.json();
```

### Python Example

```python
import requests

API_BASE = 'http://localhost:8081/api'

# Login
response = requests.post(f'{API_BASE}/auth/login', json={
    'email': email,
    'password': password
})
token = response.json()['token']

# Get problems
problems = requests.get(f'{API_BASE}/problems').json()

# Submit code
headers = {'Authorization': f'Bearer {token}'}
submission = requests.post(f'{API_BASE}/submit', json={
    'problemId': problem_id,
    'language': language,
    'code': code
}, headers=headers).json()
```

## Support

For issues and questions, please contact the development team or create an issue in the project repository.