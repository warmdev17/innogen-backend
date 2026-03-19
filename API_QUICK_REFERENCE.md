# API Quick Reference - Access Token + Refresh Token

## Base URL
```
http://localhost:8081/api
```

## 1. Login

### Request
```bash
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

### Response (200 OK)
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

## 2. Refresh Token

### Request
```bash
POST /auth/refresh
Content-Type: application/json

{
  "refreshToken": "your-refresh-token-here"
}
```

### Response (200 OK)
```json
{
  "accessToken": "new-access-token",
  "refreshToken": "new-refresh-token"
}
```

### Response (401 Unauthorized)
```json
{
  "error": "refresh token expired"
}
```

---

## 3. Get Current User (Protected)

### Request
```bash
GET /me
Authorization: Bearer your-access-token-here
```

### Response (200 OK)
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

### Response (401 Unauthorized)
```json
{
  "error": "invalid token"
}
```

---

## 4. Logout (Single Device)

### Request
```bash
POST /auth/logout
Content-Type: application/json
Authorization: Bearer your-access-token-here

{
  "refreshToken": "token-to-revoke"
}
```

### Response (200 OK)
```json
{
  "message": "logged out successfully"
}
```

---

## 5. Logout All Devices

### Request
```bash
POST /auth/logout-all
Authorization: Bearer your-access-token-here
```

### Response (200 OK)
```json
{
  "message": "logged out from all devices successfully"
}
```

---

## 6. Other Protected Endpoints

All other protected endpoints use the same Authorization header:

```bash
GET /problems
Authorization: Bearer your-access-token-here

POST /submit
Authorization: Bearer your-access-token-here
Content-Type: application/json

{
  "problemId": 1,
  "language": "python",
  "code": "print('Hello, World!')"
}
```

---

## Client-Side Implementation Examples

### JavaScript/Node.js

```javascript
const API_BASE = 'http://localhost:8081/api';

// Store tokens
let accessToken = localStorage.getItem('accessToken');
let refreshToken = localStorage.getItem('refreshToken');

// Login
async function login(email, password) {
  const response = await fetch(`${API_BASE}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  });

  if (!response.ok) throw new Error('Login failed');

  const data = await response.json();
  accessToken = data.accessToken;
  refreshToken = data.refreshToken;
  localStorage.setItem('accessToken', accessToken);
  localStorage.setItem('refreshToken', refreshToken);

  return data;
}

// API request with auto-refresh
async function apiRequest(url, options = {}) {
  const headers = {
    ...options.headers,
    'Authorization': `Bearer ${accessToken}`
  };

  const response = await fetch(url, {
    ...options,
    headers
  });

  // If unauthorized, try to refresh
  if (response.status === 401) {
    const newTokens = await refreshAccessToken();
    if (newTokens) {
      headers['Authorization'] = `Bearer ${accessToken}`;
      return fetch(url, { ...options, headers });
    }
  }

  return response;
}

// Refresh access token
async function refreshAccessToken() {
  if (!refreshToken) return null;

  const response = await fetch(`${API_BASE}/auth/refresh`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ refreshToken })
  });

  if (!response.ok) {
    // Refresh failed, redirect to login
    window.location.href = '/login';
    return null;
  }

  const data = await response.json();
  accessToken = data.accessToken;
  refreshToken = data.refreshToken;
  localStorage.setItem('accessToken', accessToken);
  localStorage.setItem('refreshToken', refreshToken);

  return data;
}

// Logout
async function logout() {
  if (refreshToken) {
    await fetch(`${API_BASE}/auth/logout`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${accessToken}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ refreshToken })
    });
  }

  localStorage.removeItem('accessToken');
  localStorage.removeItem('refreshToken');
  window.location.href = '/login';
}

// Get current user
async function getCurrentUser() {
  const response = await apiRequest(`${API_BASE}/me`);
  return response.json();
}
```

---

### Python

```python
import requests
from datetime import datetime, timedelta

API_BASE = 'http://localhost:8081/api'
access_token = None
refresh_token = None

def login(email, password):
    global access_token, refresh_token

    response = requests.post(f'{API_BASE}/auth/login', json={
        'email': email,
        'password': password
    })

    if response.status_code != 200:
        raise Exception('Login failed')

    data = response.json()
    access_token = data['accessToken']
    refresh_token = data['refreshToken']
    return data

def api_request(method, url, **kwargs):
    global access_token

    headers = kwargs.get('headers', {})
    headers['Authorization'] = f'Bearer {access_token}'
    kwargs['headers'] = headers

    response = requests.request(method, url, **kwargs)

    # If unauthorized, try to refresh
    if response.status_code == 401:
        new_tokens = refresh_access_token()
        if new_tokens:
            headers['Authorization'] = f'Bearer {access_token}'
            response = requests.request(method, url, **kwargs)

    return response

def refresh_access_token():
    global access_token, refresh_token

    if not refresh_token:
        return None

    response = requests.post(f'{API_BASE}/auth/refresh', json={
        'refreshToken': refresh_token
    })

    if response.status_code != 200:
        raise Exception('Refresh failed')

    data = response.json()
    access_token = data['accessToken']
    refresh_token = data['refreshToken']
    return data

def logout():
    global access_token, refresh_token

    if refresh_token:
        requests.post(f'{API_BASE}/auth/logout', json={
            'refreshToken': refresh_token
        }, headers={'Authorization': f'Bearer {access_token}'})

    access_token = None
    refresh_token = None

def get_current_user():
    response = api_request('GET', f'{API_BASE}/me')
    return response.json()
```

---

### cURL

```bash
# Login
curl -X POST http://localhost:8081/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'

# Store tokens
ACCESS_TOKEN="eyJhbGciOi..."
REFRESH_TOKEN="eyJhbGciOi..."

# Refresh token
curl -X POST http://localhost:8081/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refreshToken\":\"$REFRESH_TOKEN\"}"

# Get current user
curl http://localhost:8081/api/me \
  -H "Authorization: Bearer $ACCESS_TOKEN"

# Logout
curl -X POST http://localhost:8081/api/auth/logout \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"refreshToken\":\"$REFRESH_TOKEN\"}"

# Logout all devices
curl -X POST http://localhost:8081/api/auth/logout-all \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

---

## Error Codes

| Code | Message | Description |
|------|---------|-------------|
| 200 | OK | Success |
| 400 | Bad Request | Invalid input |
| 401 | Unauthorized | Invalid/expired token |
| 404 | Not Found | Resource not found |
| 500 | Internal Server Error | Server error |

---

## Token Expiration

- **Access Token**: 15 minutes
- **Refresh Token**: 30 days

Tokens will automatically expire. Use the refresh endpoint to get new tokens.

---

## Security Notes

1. **Never** expose refresh tokens in URLs or logs
2. Use **HTTPS** in production
3. Store tokens **securely** (HTTP-only cookies recommended)
4. **Rotate** refresh tokens on each use
5. **Revoke** tokens on logout
6. Implement **rate limiting** on auth endpoints
7. Use **strong, random secrets** for JWT signing

---

## Testing with Postman

1. **Create a collection** for your API
2. **Add login request** → Save tokens to variables
3. **Add Pre-request Script** for automatic token refresh
4. **Use variables** in Authorization header

**Pre-request Script Example:**
```javascript
if (!pm.environment.get('accessToken') ||
    pm.response.code === 401) {

  const refreshToken = pm.environment.get('refreshToken');
  if (refreshToken) {
    const refreshResponse = pm.sendRequest({
      url: 'http://localhost:8081/api/auth/refresh',
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: { refreshToken: refreshToken }
    });

    refreshResponse.then((response) => {
      const data = response.json();
      pm.environment.set('accessToken', data.accessToken);
      pm.environment.set('refreshToken', data.refreshToken);
    });
  }
}
```

---

## Need Help?

- 📖 Read: [JWT_MIGRATION_GUIDE.md](./JWT_MIGRATION_GUIDE.md)
- 📊 Full Summary: [JWT_IMPLEMENTATION_SUMMARY.md](./JWT_IMPLEMENTATION_SUMMARY.md)
- 🔍 Debug tokens: https://jwt.io/
- 📝 Test with Swagger: http://localhost:8081/swagger/index.html