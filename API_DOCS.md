# Innogen Backend API Documentation

Base URL (Production): `https://code.innogenlab.com/api`
*(Local dev: `http://localhost:8080/api`)*

## Authentication (`/auth`)

### 1. Request OTP (Temporarily Disabled)

**[DEPRECATED]** This endpoint is currently disabled.

Sends a 6-digit OTP to the user's email for registration.

- **URL**: `/auth/send-otp`
- **Method**: `POST`
- **Body**:

  ```json
  {
    "email": "user@example.com"
  }
  ```

### 2. Register User (Temporarily Disabled)

**[DEPRECATED]** This endpoint is currently disabled. Users cannot self-register at this time. Use the existing admin account or seed data.

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
        "Slug": "two-sum",
        "AcceptanceRate": 0.0,
        "Title": "Two Sum",
        "Difficulty": "easy",
        "ProblemMd": "Given an array of integers...",
        "TimeLimitMs": 1000,
        "MemoryLimitKb": 256000,
        "IsPublished": false,
        "CreatedBy": 1,
        "CreatedAt": "2026-02-21T10:00:00Z",
        "UpdatedAt": "2026-02-21T10:00:00Z"
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
    "slug": "two-sum",
    "description": "Given an array of integers return indices of the two numbers such that they add up to target.",
    "difficulty": "Easy",
    "time_limit_ms": 1000,
    "memory_limit_kb": 256000
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
    "InputData": "[2,7,11,15]\n9",
    "ExpectedOutput": "[0,1]",
    "IsHidden": true,
    "Role": "hidden"
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
        "ID": "123e4567-e89b-12d3-a456-426614174000",
        "UserID": 1,
        "ProblemID": 1,
        "Code": "def twoSum(nums, target):\n    pass",
        "Language": "python",
        "Status": "pending",
        "RuntimeMs": 0,
        "MemoryKb": 0,
        "ErrorMessage": "",
        "PassCount": 0,
        "TotalTestcases": 0,
        "CreatedAt": "2026-02-21T10:00:00Z"
      }
    }
    ```

  - `500 Server Error`: `{"error": "Failed to queue submission"}`

---

## Piston Proxy (`/piston`)

These endpoints act as a direct reverse-proxy to the internal Piston execution engine. All Piston APIs can be called by prefixing them with `/piston`.

### 1. Execute Code

Directly executes source code using the Piston engine engine without creating a submission record in the database.

- **URL**: `/piston/api/v2/execute`
- **Method**: `POST`
- **Body**:

  ```json
  {
    "language": "python",
    "version": "3.10.0",
    "files": [
      {
        "content": "print('Hello Innogen!')"
      }
    ]
  }
  ```

- **Responses**:
  - `200 OK`: Piston execution result.

    ```json
    {
      "language": "python",
      "version": "3.10.0",
      "run": {
        "stdout": "Hello Innogen!\n",
        "stderr": "",
        "code": 0,
        "signal": null,
        "output": "Hello Innogen!\n"
      }
    }
    ```

---

## Nginx Reverse Proxy Configuration (Production)

If you are running the service on a VPS and want to publicize the APIs with an SSL certificate using Nginx, you can reference the configuration below.

### 1. Piston Engine (`excode.innogenlab.com`)

This points to the internal Piston engine running on port `2000`.

```nginx
server {
    listen 80;
    server_name excode.innogenlab.com www.excode.innogenlab.com;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl;
    server_name excode.innogenlab.com www.excode.innogenlab.com;

    ssl_certificate /etc/nginx/ssl/innogenlab.com.pem;
    ssl_certificate_key /etc/nginx/ssl/innogenlab.com.key;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    location / {
        proxy_pass http://127.0.0.1:2000;

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 2. Frontend & Backend API (`code.innogenlab.com`)

This configuration serves the Frontend application on the root path `/` and proxies the Main Go Backend running via Docker to the `/api`, `/piston`, and `/health` paths.

```nginx
server {
    listen 80;
    server_name code.innogenlab.com www.code.innogenlab.com;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl;
    server_name code.innogenlab.com www.code.innogenlab.com;

    ssl_certificate /etc/nginx/ssl/innogenlab.com.pem;
    ssl_certificate_key /etc/nginx/ssl/innogenlab.com.key;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # 1. Frontend Application
    location / {
        # Example: Proxy to a Next.js/React frontend running on port 3000
        # proxy_pass http://127.0.0.1:3000;
        
        # Or serve static files built by React/Vue
        # root /var/www/innogen-frontend;
        # index index.html;
        # try_files $uri $uri/ /index.html;
    }

    # 2. Backend API
    location /api/ {
        proxy_pass http://127.0.0.1:8080/api/;

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 3. Piston Execution Reverse Proxy (Called by Frontend)
    location /piston/ {
        proxy_pass http://127.0.0.1:8080/piston/;

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 4. Backend Health Check
    location /health {
        proxy_pass http://127.0.0.1:8080/health;
    }
}
```
