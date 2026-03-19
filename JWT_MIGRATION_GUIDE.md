# JWT Authentication Migration Guide

## Overview

This guide explains how to upgrade from a single JWT token (72-hour expiration) to a modern access token + refresh token authentication system.

## What Changed

### Old System (Single Token)
- One JWT token with 72-hour expiration
- Simple but insecure for long sessions
- Difficult to revoke once issued
- No token rotation

### New System (Access Token + Refresh Token)
- **Access Token**: Short-lived (15 minutes), used for API requests
- **Refresh Token**: Long-lived (30 days), used to get new access tokens
- Token rotation on refresh for security
- Database-backed refresh token management
- Ability to revoke tokens individually or in bulk

## Changes Required

### 1. Environment Variables

Add a new environment variable for the refresh token secret:

```bash
# .env file
JWT_SECRET=your-jwt-secret-for-access-tokens
JWT_REFRESH_SECRET=your-separate-secret-for-refresh-tokens
PORT=8081
```

**Important**: Use different secrets for access and refresh tokens!

### 2. Database Migration

Create the refresh tokens table:

```sql
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id BIGINT NOT NULL,
    token_hash TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_expires_at (expires_at)
);
```

For GORM users, simply run the application and it will auto-migrate based on the model.

## API Changes

### 1. Login Response

**Old Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**New Response:**
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 2. New Endpoints

#### POST /api/auth/refresh

Request:
```json
{
  "refreshToken": "your-refresh-token"
}
```

Response:
```json
{
  "accessToken": "new-access-token",
  "refreshToken": "new-refresh-token"
}
```

#### POST /api/auth/logout

Request:
```json
{
  "refreshToken": "token-to-revoke"
}
```

Response:
```json
{
  "message": "logged out successfully"
}
```

#### POST /api/auth/logout-all

No request body required.

Response:
```json
{
  "message": "logged out from all devices successfully"
}
```

## Frontend Integration

### 1. Token Storage

Store tokens securely:
- Access Token: Can be stored in memory (changes frequently)
- Refresh Token: Should be in HTTP-only cookie OR secure storage

**Example with HTTP-only cookies:**
```javascript
// Set cookies (server-side)
res.cookie('accessToken', accessToken, { httpOnly: true, secure: true, sameSite: 'strict', maxAge: 15 * 60 * 1000 });
res.cookie('refreshToken', refreshToken, { httpOnly: true, secure: true, sameSite: 'strict', maxAge: 30 * 24 * 60 * 60 * 1000 });
```

**Example with localStorage (less secure):**
```javascript
localStorage.setItem('accessToken', accessToken);
localStorage.setItem('refreshToken', refreshToken);
```

### 2. Request Flow

```javascript
// 1. Login
const loginResponse = await fetch('/api/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ email, password })
});
const { accessToken, refreshToken } = await loginResponse.json();

// 2. Make API request with access token
const protectedResponse = await fetch('/api/me', {
  headers: { 'Authorization': `Bearer ${accessToken}` }
});

// 3. If 401 error, try to refresh
if (protectedResponse.status === 401) {
  const refreshResponse = await fetch('/api/auth/refresh', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ refreshToken })
  });

  if (refreshResponse.ok) {
    const { accessToken: newAccessToken, refreshToken: newRefreshToken } = await refreshResponse.json();
    // Retry original request with new access token
    // Update stored tokens
  } else {
    // Refresh failed, redirect to login
    window.location.href = '/login';
  }
}
```

### 3. Auto-Refresh Logic

```javascript
// Set up automatic token refresh
let accessToken = localStorage.getItem('accessToken');
let refreshToken = localStorage.getItem('refreshToken');

function scheduleTokenRefresh() {
  const expiresIn = 15 * 60 * 1000; // 15 minutes
  const refreshTime = expiresIn - 60 * 1000; // Refresh 1 minute before expiry

  setTimeout(async () => {
    try {
      const response = await fetch('/api/auth/refresh', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refreshToken })
      });

      if (response.ok) {
        const data = await response.json();
        accessToken = data.accessToken;
        refreshToken = data.refreshToken;
        localStorage.setItem('accessToken', accessToken);
        localStorage.setItem('refreshToken', refreshToken);
        scheduleTokenRefresh(); // Schedule next refresh
      }
    } catch (error) {
      console.error('Token refresh failed:', error);
    }
  }, refreshTime);
}

// Start the refresh timer after login
scheduleTokenRefresh();
```

### 4. Logout Handling

```javascript
// Logout current device
async function logout() {
  const refreshToken = localStorage.getItem('refreshToken');
  await fetch('/api/auth/logout', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${accessToken}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ refreshToken })
  });

  localStorage.removeItem('accessToken');
  localStorage.removeItem('refreshToken');
  window.location.href = '/login';
}

// Logout all devices
async function logoutAll() {
  await fetch('/api/auth/logout-all', {
    method: 'POST',
    headers: { 'Authorization': `Bearer ${accessToken}` }
  });

  localStorage.removeItem('accessToken');
  localStorage.removeItem('refreshToken');
  window.location.href = '/login';
}
```

## Security Best Practices Implemented

### 1. Token Separation
- Different secrets for access and refresh tokens
- Different claims and expiration times
- Separate validation logic

### 2. Token Rotation
- Refresh endpoint generates new tokens
- Old refresh tokens are revoked
- Prevents token reuse attacks

### 3. Database Storage
- Refresh tokens hashed before storage
- Tokens tracked with expiration
- Revocation support

### 4. Security Features
- Short-lived access tokens (15 min)
- Secure random token generation
- Token type validation
- Expiration checks

## Migration Checklist

### Backend
- [ ] Update environment variables (add JWT_REFRESH_SECRET)
- [ ] Run database migration for refresh_tokens table
- [ ] Test all endpoints with new authentication flow

### Frontend
- [ ] Update login handler to store both tokens
- [ ] Add automatic token refresh logic
- [ ] Update API request headers to use accessToken
- [ ] Handle 401 responses with refresh flow
- [ ] Add logout handlers
- [ ] Test complete authentication flow

### Testing
- [ ] Login → Get tokens
- [ ] Make protected API requests
- [ ] Wait for token expiration → Auto-refresh
- [ ] Logout → Token revoked
- [ ] Multiple devices → Independent tokens

## Troubleshooting

### "invalid refresh token" error
- Token not found in database
- Token has been revoked
- Token expired
- Token type mismatch

### "failed to rotate refresh token" error
- Database operation failed
- Old token verification failed
- Database connection issue

### 401 on protected endpoint
- Access token expired
- Access token missing
- Access token invalid

## Rollback Plan

To rollback to single token system:

1. Remove refresh token code
2. Update JWT utils to use single token generation
3. Remove refresh token table (optional)
4. Update API responses to return single token
5. Update frontend to use single token

## Support

For issues or questions:
- Check logs for detailed error messages
- Verify environment variables are set
- Ensure database connection is working
- Check token format and claims

## Additional Resources

- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)
- [OWASP Authentication Cheatsheet](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html)
- [Token Storage Recommendations](https://stackoverflow.com/questions/35226774/where-to-store-jwt-in-browser-and-how-to-deal-with-csrf)