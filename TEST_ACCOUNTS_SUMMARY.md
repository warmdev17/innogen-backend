# Test Accounts Summary

## Overview

The register and send-otp routes have been **temporarily disabled**. Only login is available for existing accounts. The database seeder has been updated to create 5 test accounts with different roles.

## Current Authentication Status

### ✅ Available Endpoints
- `POST /api/auth/login` - Login with existing accounts
- `POST /api/auth/refresh` - Refresh access token
- `POST /api/auth/logout` - Logout from current device
- `POST /api/auth/logout-all` - Logout from all devices

### ❌ Disabled Endpoints
- `POST /api/auth/register` - **Disabled** (self-registration not allowed)
- `POST /api/auth/send-otp` - **Disabled** (OTP verification not used)

**Note**: These endpoints are not available and will return 404 if accessed.

## Test Accounts Created

The seeder creates 5 accounts with the following credentials:

### 1. Admin Account
- **Email**: `admin@admin.com`
- **Password**: Value of `ADMIN_PASSWORD` environment variable (default: Check your `.env` file)
- **Role**: `admin`
- **Permissions**: Full access, can create problems

### 2. Teacher Account
- **Email**: `teacher@innogen.com`
- **Password**: `teacher123`
- **Role**: `teacher`
- **Permissions**: Can create problems, manage content

### 3. Student Accounts
All student accounts have the same password: `student123`

#### Student 1
- **Email**: `student1@innogen.com`
- **Username**: `student1`
- **Full Name**: `Alice Student`
- **Role**: `student`

#### Student 2
- **Email**: `student2@innogen.com`
- **Username**: `student2`
- **Full Name**: `Bob Student`
- **Role**: `student`

#### Student 3
- **Email**: `student3@innogen.com`
- **Username**: `student3`
- **Full Name**: `Charlie Student`
- **Role**: `student`

## How to Use Test Accounts

### 1. Run the Seeder

```bash
# Make sure PostgreSQL and Redis are running
docker-compose up -d

# Set ADMIN_PASSWORD in .env file
echo "ADMIN_PASSWORD=your-password-here" > .env

# Run the seeder to create accounts
go run ./cmd/seed/main.go
```

The seeder will output:
```
Created user: admin@admin.com (role: admin)
Created user: teacher@innogen.com (role: teacher)
Created user: student1@innogen.com (role: student)
Created user: student2@innogen.com (role: student)
Created user: student3@innogen.com (role: student)
Successfully created sample problem: Two Sum
Successfully created 3 sample test cases

=== Test Accounts Created ===
Admin: admin@admin.com / your-password-here
Teacher: teacher@innogen.com / teacher123
Student 1: student1@innogen.com / student123
Student 2: student2@innogen.com / student123
Student 3: student3@innogen.com / student123
===============================
```

### 2. Login Example

```bash
curl -X POST http://localhost:8081/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student1@innogen.com",
    "password": "student123"
  }'
```

Response:
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 3. Access Protected Endpoints

```bash
# Get current user info
curl http://localhost:8081/api/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Submit code (students and teachers can do this)
curl -X POST http://localhost:8081/api/submit \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "problemId": 1,
    "language": "python",
    "code": "print(input())"
  }'

# Create problem (admin and teacher only)
curl -X POST http://localhost:8081/api/admin/problems \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{...}'
```

## Role-Based Access Control

### Admin Role (`admin`)
- ✅ Can create, edit, delete problems
- ✅ Can view all submissions
- ✅ Can manage users
- ✅ Can access all endpoints

### Teacher Role (`teacher`)
- ✅ Can create, edit, delete problems
- ✅ Can view submissions from students
- ❌ Cannot manage other teachers or admins

### Student Role (`student`)
- ✅ Can view published problems
- ✅ Can submit code for judging
- ✅ Can run code without test cases
- ❌ Cannot create problems
- ❌ Cannot view other students' submissions

## Re-enabling Registration

To re-enable user registration:

1. **Add Register Controller Function** in `internal/controllers/auth.go`:
```go
func Register(c *gin.Context) {
    // Implementation for user registration
}
```

2. **Add SendOTP Controller Function** for email verification:
```go
func SendOTP(c *gin.Context) {
    // Implementation for sending OTP
}
```

3. **Add Routes** in `internal/routes/routes.go`:
```go
api.POST("/auth/register", controllers.Register)
api.POST("/auth/send-otp", controllers.SendOTP)
```

4. **Update Login Logic** if you want both registration and login available

## Current API Endpoints

### Authentication
```
POST   /api/auth/login         # Login (available)
POST   /api/auth/refresh       # Refresh token (available)
POST   /api/auth/logout        # Logout (available)
POST   /api/auth/logout-all    # Logout all devices (available)
POST   /api/auth/register      # Register (DISABLED)
POST   /api/auth/send-otp      # Send OTP (DISABLED)
```

### User
```
GET    /api/me                 # Get current user (requires auth)
```

### Problems
```
GET    /api/problems           # List all problems
GET    /api/problems/:id       # Get problem details
POST   /api/admin/problems     # Create problem (admin/teacher only)
```

### Submissions
```
POST   /api/submit             # Submit code (requires auth)
GET    /api/submit/:id         # Get submission status (requires auth)
```

### Code Execution
```
POST   /api/run                # Run code directly (requires auth)
```

## Security Notes

1. **No Self-Registration**: Users cannot create accounts themselves
2. **Admin-Only User Creation**: Only admins can create users via database/seeder
3. **Password Security**: All passwords are hashed using bcrypt
4. **Token-Based Auth**: Uses access + refresh token pattern
5. **Role Enforcement**: All endpoints enforce proper role-based access

## Testing Workflow

1. Run seeder to create test accounts
2. Use login endpoint to authenticate
3. Store access token for API requests
4. Test different features with different roles
5. Use logout endpoints when done

## Support

If you need to:
- Create additional users → Edit `cmd/seed/main.go`
- Change passwords → Re-run seeder (will skip existing users)
- Re-enable registration → Follow "Re-enabling Registration" section above
- Reset database → Drop and recreate tables, then run seeder

---

**Status**: ✅ Ready for Testing

**Total Test Accounts**: 5 (1 admin, 1 teacher, 3 students)

**Build Status**: ✅ Compiles Successfully