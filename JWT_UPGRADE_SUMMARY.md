# JWT Access Token + Refresh Token Implementation Summary

## Implementation Overview

This document summarizes the complete upgrade from a single JWT token to a secure access token + refresh token authentication system.

## Stack Information

- **Language:** Go
- **Framework:** Gin
- **JWT Package:** golang-jwt/jwt/v5
- **Database:** PostgreSQL with GORM
- **Current Token:** Single JWT with 72-hour expiration → **NEW:** Access (15 min) + Refresh (30 days)

## Architecture

### Token Types

1. **Access Token**
   - Lifetime: 15 minutes
   - Purpose: Authenticate API requests
   - Sent via: Authorization header (Bearer token)
   - Validation: JWT signature verification

2. **Refresh Token**
   - Lifetime: 30 days
   - Purpose: Get new access tokens
   - Storage: Database (hashed)
   - Security: Revocable, rotatable

### Security Features

✅ **Separate Secrets:** Different keys for access and refresh tokens
✅ **Token Rotation:** New refresh token on each use
✅ **Database Hashing:** SHA256 hashing of refresh tokens
✅ **Expiration Checks:** Automatic validation on both tokens
✅ **Token Type Validation:** Prevents token misuse
✅ **Revocation Support:** Individual and bulk revocation
✅ **Short-Lived Access:** Minimizes exposure window

## Files Modified

### 1. `/internal/utils/jwt.go`

**Changes:**
- Added `GenerateAccessToken()` - creates 15-minute access tokens
- Added `GenerateRefreshToken()` - creates 30-day refresh JWT tokens
- Added `ParseAccessToken()` - validates access tokens with type check
- Added `ParseRefreshToken()` - validates refresh tokens with type check
- Added `GetTokenClaims()` - extracts claims from parsed tokens
- Deprecated: `GenerateToken()` (kept for compatibility)

**New Functions:**
```go
GenerateAccessToken(userID uint, role string) (string, error)
GenerateRefreshToken(userID uint, role string) (string, error)
ParseAccessToken(tokenStr string) (*jwt.Token, error)
ParseRefreshToken(tokenStr string) (*jwt.Token, error)
GetTokenClaims(token *jwt.Token) (jwt.MapClaims, error)
```

### 2. `/internal/services/refresh_token.go`

**New Service:** Complete refresh token management

**Functions:**
```go
CreateRefreshToken(userID uint) (string, error)
    // Creates new refresh token, stores hashed version in DB

VerifyRefreshToken(rawToken string, userID uint) (*RefreshToken, error)
    // Validates refresh token exists and is not revoked

RevokeRefreshToken(rawToken string, userID uint) error
    // Marks specific token as revoked

RevokeAllUserTokens(userID uint) error
    // Revokes all tokens for a user (logout all devices)

RotateRefreshToken(oldToken string, userID uint) (string, error)
    // Creates new token and revokes old one (security best practice)
```

**Data Model:**
```go
type RefreshToken struct {
    ID         uuid.UUID  // Primary key
    UserID     uint       // Foreign key to User
    TokenHash  string     // Hashed token (SHA256)
    ExpiresAt  time.Time  // Expiration date
    Revoked    bool       // Revocation status
    CreatedAt  time.Time  // Creation timestamp
    UpdatedAt  time.Time  // Last update timestamp
}
```

### 3. `/internal/controllers/auth.go`

**Updated Endpoints:**

#### POST /auth/login
**Before:**
```go
c.JSON(200, gin.H{"token": token})
```

**After:**
```go
c.JSON(200, gin.H{
    "accessToken":  accessToken,
    "refreshToken": refreshTokenStr,
})
```

#### NEW: POST /auth/refresh
Generates new access + refresh tokens using valid refresh token

```go
func RefreshToken(c *gin.Context) {
    // 1. Validate refresh token format
    // 2. Extract claims
    // 3. Check expiration
    // 4. Verify in database
    // 5. Generate new access token
    // 6. Rotate refresh token
    // 7. Return both tokens
}
```

#### NEW: POST /auth/logout
Revokes specific refresh token

```go
func Logout(c *gin.Context) {
    // 1. Get refresh token from request
    // 2. Get user from access token
    // 3. Revoke refresh token
    // 4. Return success
}
```

#### NEW: POST /auth/logout-all
Revokes all refresh tokens for user

```go
func LogoutAll(c *gin.Context) {
    // 1. Get user from access token
    // 2. Revoke all tokens for user
    // 3. Return success
}
```

**Response Types Added:**
```go
type RefreshTokenRequest struct {
    RefreshToken string `json:"refreshToken"`
}

type RefreshTokenResponse struct {
    AccessToken  string `json:"accessToken"`
    RefreshToken string `json:"refreshToken"`
}
```

### 4. `/internal/middleware/jwt.go`

**Updated:**
- Changed `ParseToken()` to `ParseAccessToken()`
- Added claims to context: `c.Set("claims", claims)`
- Improved token type validation

### 5. `/internal/routes/routes.go`

**Updated Route Structure:**
```go
api := r.Group("/api")
{
    // Public endpoints
    api.POST("/auth/login", controllers.Login)
    api.POST("/auth/refresh", controllers.RefreshToken)  // NEW

    // Protected endpoints (require access token)
    protected := api.Group("")
    protected.Use(middleware.JWTAuth())
    {
        // User info
        me := protected.Group("/me")
        {
            me.GET("", controllers.GetCurrentUser)
        }

        // Auth
        auth := protected.Group("/auth")  // NEW
        {
            auth.POST("/logout", controllers.Logout)        // NEW
            auth.POST("/logout-all", controllers.LogoutAll) // NEW
        }

        // ... other protected routes
    }
}
```

## Environment Variables Required

Add to `.env`:
```bash
JWT_SECRET=your-access-token-secret
JWT_REFRESH_SECRET=your-refresh-token-secret  # NEW
```

**Important:** Use different, strong secrets for each token type!

## Database Migration

Run this SQL to create the refresh_tokens table:

```sql
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id BIGINT NOT NULL,
    token_hash TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
CREATE INDEX idx_refresh_tokens_revoked ON refresh_tokens(revoked);
```

## Authentication Flow

### 1. Login Flow

```
User → POST /api/auth/login
        email + password
        ↓
        Validate credentials
        ↓
        Generate access token (15 min)
        ↓
        Generate refresh token (30 days)
        ↓
        Store hashed refresh token in DB
        ↓
        Return: accessToken + refreshToken
```

### 2. API Request Flow

```
Client → GET /api/me
         Authorization: Bearer <accessToken>
         ↓
         Validate access token
         ↓
         Extract user_id from claims
         ↓
         Return user data
```

### 3. Token Refresh Flow

```
Client → GET /api/me (401 error)
         ↓
         POST /api/auth/refresh
         { refreshToken: "<refresh_token>" }
         ↓
         Validate refresh token format & claims
         ↓
         Check token exists & not revoked in DB
         ↓
         Generate new access token (15 min)
         ↓
         Rotate refresh token (revoke old, create new)
         ↓
         Return: new accessToken + new refreshToken
         ↓
         Retry original API request
```

### 4. Logout Flow

```
Client → POST /api/auth/logout
         Authorization: Bearer <accessToken>
         { refreshToken: "<token_to_revoke>" }
         ↓
         Get user from access token
         ↓
         Mark refresh token as revoked in DB
         ↓
         Return: success message
```

## Security Considerations

### Token Storage

**Access Token:**
- Store in memory or short-term storage
- Transmits with every request
- Low security risk (15 min lifetime)

**Refresh Token:**
- Store in HTTP-only cookie (recommended) OR secure storage
- Only used for token refresh
- High security risk (30 day lifetime)

### Recommended Frontend Implementation

```javascript
// Option 1: HTTP-only cookies (most secure)
// Server sets cookies with httpOnly flag
// Tokens automatically sent with requests
// Protected from XSS attacks

// Option 2: Memory + localStorage (convenient)
// Store in localStorage
// Add to Authorization header
// Implement automatic refresh
```

### Token Reuse Prevention

The system implements **token rotation**:
- Each refresh generates NEW refresh token
- Old refresh token is revoked
- Prevents replay attacks
- Limits token lifetime

### Database Security

- Tokens hashed with SHA256 before storage
- Only hash stored, never raw tokens
- Automatic cleanup of expired tokens (implement scheduled job)

## Testing

### Test Cases

1. **Login**
   ```bash
   curl -X POST http://localhost:8081/api/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"user@example.com","password":"password123"}'
   ```

2. **Access Protected Endpoint**
   ```bash
   curl -X GET http://localhost:8081/api/me \
     -H "Authorization: Bearer <accessToken>"
   ```

3. **Refresh Token**
   ```bash
   curl -X POST http://localhost:8081/api/auth/refresh \
     -H "Content-Type: application/json" \
     -d '{"refreshToken":"<refreshToken>"}'
   ```

4. **Logout**
   ```bash
   curl -X POST http://localhost:8081/api/auth/logout \
     -H "Authorization: Bearer <accessToken>" \
     -H "Content-Type: application/json" \
     -d '{"refreshToken":"<refreshToken>"}'
   ```

## Benefits of This Implementation

✅ **Improved Security:**
- Short-lived access tokens reduce exposure window
- Refresh token rotation prevents reuse
- Database tracking enables revocation

✅ **Better User Experience:**
- Users stay logged in for 30 days
- Seamless token refresh (no re-login required)
- Multiple device support

✅ **Compliance:**
- Follows OWASP recommendations
- Aligns with OAuth 2.0 best practices
- Industry-standard approach

✅ **Scalability:**
- Token revocation at scale
- Per-device session management
- Database-backed for reliability

## Migration Path

### For Existing Users

1. **Phase 1: Deploy Backend**
   - Add new environment variable
   - Run migration
   - Deploy with dual token support

2. **Phase 2: Update Frontend**
   - Update login to handle both tokens
   - Add refresh logic
   - Test thoroughly

3. **Phase 3: Cleanup**
   - Remove old single-token code
   - Update documentation

### Backward Compatibility

The old `utils.GenerateToken()` function is kept for compatibility but deprecated. No breaking changes for existing code during migration period.

## Performance Impact

- **Database Queries:** One extra query per refresh (verify token)
- **Token Generation:** Negligible overhead
- **Storage:** ~200 bytes per refresh token in DB

## Monitoring & Observability

Recommended metrics to track:
- Token refresh success/failure rate
- Token expiration frequency
- Logout events (single vs all)
- Active refresh tokens per user

## Conclusion

This implementation provides a production-ready, secure authentication system following industry best practices. The access token + refresh token pattern is widely adopted and provides the right balance between security and user experience.

For questions or issues, refer to the JWT_MIGRATION_GUIDE.md or check the application logs for detailed error messages.