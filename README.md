# Innogen Backend

A competitive programming platform backend built in Go with JWT authentication, code execution via Piston, and Redis-based job queues.

## 🚀 Quick Start with Docker Compose

The easiest way to run everything:

```bash
# 1. Start all services (PostgreSQL, Redis, Backend)
docker-compose up -d --build

# 2. Seed the database with test accounts
docker-compose run --rm seeder

# 3. Access the API
# API: http://localhost:8081
# Health: http://localhost:8081/api/health
# Swagger: http://localhost:8081/swagger/index.html
```

**Test Accounts:**
- Admin: `admin@admin.com` / `admin123`
- Teacher: `teacher@innogen.com` / `teacher123`
- Student: `student1@innogen.com` / `student123`

See [DOCKER_COMPOSE_GUIDE.md](./DOCKER_COMPOSE_GUIDE.md) for detailed Docker instructions.

## 🛠️ Local Development

### Prerequisites
- Go 1.21+
- PostgreSQL 15+
- Redis 7+

### Setup

```bash
# 1. Copy environment file
cp .env.example .env

# 2. Edit .env with your settings
#    Update JWT_SECRET and JWT_REFRESH_SECRET to strong random values

# 3. Start PostgreSQL and Redis
# Option A: Docker
docker run -d --name postgres -p 5433:5432 -e POSTGRES_PASSWORD=maiphuongdangyeu -e POSTGRES_DB=innogendb postgres:15
docker run -d --name redis -p 6380:6379 redis:7

# Option B: Local installation
# Make sure PostgreSQL and Redis are running locally

# 4. Seed database
go run ./cmd/seed/main.go

# 5. Start backend
go run ./cmd/main.go
```

### Configuration

Required environment variables (see `.env.example`):

- `POSTGRES_HOST`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRES_PORT`
- `REDIS_HOST`, `REDIS_PORT`
- `JWT_SECRET`: Secret key for access tokens (32+ chars)
- `JWT_REFRESH_SECRET`: Secret key for refresh tokens (32+ chars, different from JWT_SECRET)
- `ADMIN_PASSWORD`: Password for admin user
- `PISTON_URL`: URL of Piston code execution service

**Production Piston URL**: `https://excode.innogenlab.com`

## 📚 API Documentation

### Authentication

The API uses **access token + refresh token** authentication:

#### 1. Login
```bash
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

Response:
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 2. Use Access Token
```bash
GET /api/me
Authorization: Bearer <accessToken>
```

#### 3. Refresh Token
```bash
POST /api/auth/refresh
Content-Type: application/json

{
  "refreshToken": "your-refresh-token"
}
```

Returns new `accessToken` and `refreshToken`.

### Key Endpoints

- `POST /api/auth/login` - Login
- `POST /api/auth/refresh` - Refresh access token
- `GET /api/me` - Get current user
- `GET /api/problems` - List problems
- `POST /api/admin/problems` - Create problem (admin/teacher)
- `POST /api/submit` - Submit code
- `GET /api/submit/:id` - Get submission status

See [API_QUICK_REFERENCE.md](./API_QUICK_REFERENCE.md) for complete API documentation.

## 🏗️ Architecture

- **API Server**: Gin framework (port 8081)
- **Database**: PostgreSQL with GORM
- **Job Queue**: Redis
- **Code Execution**: External Piston service
- **Authentication**: JWT (access token + refresh token)

### Directory Structure

```
.
├── cmd/
│   ├── main.go          # API server entry point
│   └── seed/main.go     # Database seeding
├── internal/
│   ├── controllers/     # HTTP handlers
│   ├── middleware/      # JWT auth, CORS, etc.
│   ├── models/          # GORM models
│   ├── services/        # Business logic
│   ├── routes/          # Route registration
│   ├── database/        # DB connection
│   ├── judge/           # Background worker
│   └── utils/           # JWT utilities
└── docker-compose.yml   # Docker Compose config
```

## 🔧 Development

### Building

```bash
# Build server
go build -o innogen-backend ./cmd/main.go

# Build seeder
go build -o innogen-seed ./cmd/seed/main.go
```

### Testing

```bash
# Run all tests
go test ./...

# Run specific package
go test ./internal/controllers/
```

### Code Submission Flow

1. User submits code via `POST /api/submit`
2. Submission stored in database with status "pending"
3. Job queued in Redis (`judge_queue`)
4. Background worker processes queue
5. Worker calls Piston to execute code
6. Results stored in database
7. Status updated: "Accepted", "Wrong Answer", etc.

## 📦 Docker

### Build Images

```bash
docker-compose build
```

### Run Services

```bash
# Start all
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all
docker-compose down

# Clean everything
docker-compose down -v
```

### Services

| Service | Port | Description |
|---------|------|-------------|
| backend | 8081 | API Server |
| postgres | 5433 | PostgreSQL database |
| redis | 6380 | Redis job queue |
| seeder | - | Database seeding tool |

See [DOCKER_COMPOSE_GUIDE.md](./DOCKER_COMPOSE_GUIDE.md) for detailed Docker instructions.

## 🐛 Troubleshooting

### Database Connection Failed
```bash
# Check PostgreSQL is running
docker-compose ps postgres

# Check logs
docker-compose logs postgres

# Wait and retry
sleep 5
go run ./cmd/seed/main.go
```

### Port Already in Use
```bash
# Find process using port 8081
lsof -i :8081

# Or change port in docker-compose.yml
```

### JWT Secret Not Set
```bash
# Make sure .env exists and JWT_SECRET is set
cat .env

# Generate secrets
openssl rand -base64 32
```

## 📄 License

This project is proprietary software by Innogen Labs.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Write tests
5. Submit a pull request

## 📚 Documentation

- [JWT Migration Guide](./JWT_MIGRATION_GUIDE.md) - Auth system upgrade
- [JWT Implementation Summary](./JWT_IMPLEMENTATION_SUMMARY.md) - Technical details
- [API Quick Reference](./API_QUICK_REFERENCE.md) - API usage examples
- [Docker Compose Guide](./DOCKER_COMPOSE_GUIDE.md) - Docker setup
- [CLAUDE.md](./CLAUDE.md) - Project guidelines

## 🔐 Security

- Access tokens expire after 15 minutes
- Refresh tokens expire after 30 days
- Refresh tokens rotate on each use
- Tokens stored as hashes in database
- Different secrets for access and refresh tokens

For security issues, please contact the development team.