# Innogen Backend API Documentation

Base URL: `http://localhost:8080/api`

## Authentication (`/auth`)

### 1. Request OTP

Sends a 6-digit OTP to the user's email for registration.

- **URL**: `/auth/send-otp`
- **Method**: `POST`
- **Body**:

  ```json
  {
    "email": "user@example.com"
  }
  ```

- **Responses**:
  - `200 OK`: `{"message": "OTP sent"}`
  - `400 Bad Request`: `{"error": "valid email is required"}` or `{"error": "user exists"}`
  - `500 Server Error`: `{"error": "failed to store OTP"}`

### 2. Register User

Verifies the OTP and creates a new user with the `student` role.

- **URL**: `/auth/register`
- **Method**: `POST`
- **Body**:

  ```json
  {
    "email": "user@example.com",
    "password": "securepassword123",
    "otp": "123456"
  }
  ```

- **Responses**:
  - `201 Created`: `{"message": "registered"}`
  - `400 Bad Request`: `{"error": "invalid or expired OTP"}`
  - `500 Server Error`: `{"error": "user exists"}` or `{"error": "hash failed"}`

### 3. Login

Authenticates the user and returns a JWT token.

- **URL**: `/auth/login`
- **Method**: `POST`
- **Body**:

  ```json
  {
    "email": "user@example.com",
    "password": "securepassword123"
  }
  ```

- **Responses**:
  - `200 OK`: `{"token": "eyJhbG..."}`
  - `401 Unauthorized`: `{"error": "invalid credentials"}`

---

## User Information (`/me`)

_Requires `Authorization: Bearer <token>` header._

### 1. Get Current User

- **URL**: `/me`
- **Method**: `GET`
- **Responses**:
  - `200 OK`:

    ```json
    {
      "user_id": 1,
      "role": "student"
    }
    ```

---

## Problems (`/problems`)

_Requires `Authorization: Bearer <token>` header._

### 1. Get All Problems

- **URL**: `/problems`
- **Method**: `GET`
- **Responses**:
  - `200 OK`:

    ```json
    [
      {
        "ID": 1,
        "Title": "Two Sum",
        "Description": "Given an array of integers...",
        "Difficulty": "easy",
        "TimeLimitMs": 1000,
        "MemoryLimitMB": 256,
        "CreatedBy": 1,
        "CreatedAt": "2026-02-21T10:00:00Z"
      }
    ]
    ```

### 2. Get Problem by ID

- **URL**: `/problems/:id`
- **Method**: `GET`
- **Responses**:
  - `200 OK`: Single problem object (same format as above)
  - `404 Not Found`: `{"error": "Problem not found"}`

### 3. Create Problem (Admin/Teacher Only)

- **URL**: `/problems`
- **Method**: `POST`
- **Body**:

  ```json
  {
    "title": "Two Sum",
    "description": "Given an array of integers return indices of the two numbers such that they add up to target.",
    "difficulty": "easy",
    "time_limit_ms": 1000,
    "memory_limit_mb": 256
  }
  ```

- **Responses**:
  - `201 Created`: Returns the newly created problem object.
  - `400 Bad Request`: Validation errors.
  - `403 Forbidden`: Insufficient role permissions.

---

## Testcases (`/testcases`)

_Requires `Authorization: Bearer <token>` header with `admin` or `teacher` role._

### 1. Create Testcase

- **URL**: `/testcases`
- **Method**: `POST`
- **Body**:

  ```json
  {
    "ProblemID": 1,
    "Input": "[2,7,11,15]\n9",
    "Output": "[0,1]"
  }
  ```

- **Responses**:
  - `201 Created`: Returns the newly created testcase object.
  - `400 Bad Request`: Validation errors.
  - `403 Forbidden`: Insufficient role permissions.

---

## Submissions (`/submit`)

_Requires `Authorization: Bearer <token>` header._

### 1. Submit Code

Submits code for a specific problem and queues it for the judge worker.

- **URL**: `/submit`
- **Method**: `POST`
- **Body**:

  ```json
  {
    "problem_id": 1,
    "code": "def twoSum(nums, target):\n    pass",
    "language": "python"
  }
  ```

- **Responses**:
  - `201 Created`:

    ```json
    {
      "message": "Submission queued",
      "submission": {
        "ID": 1,
        "UserID": 1,
        "ProblemID": 1,
        "Code": "def twoSum(nums, target):\n    pass",
        "Language": "python",
        "Status": "pending"
        // ... (other internal submission fields)
      }
    }
    ```

  - `400 Bad Request`: Validation errors.
  - `500 Server Error`: `{"error": "Failed to queue submission"}`
