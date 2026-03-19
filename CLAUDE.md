# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Innogen is a competitive programming platform backend API built in Go. It provides problem management, user authentication, and code execution capabilities using an asynchronous worker architecture with Redis queue.

## High-Level Architecture

The application follows a clean layered architecture:
- **cmd/**: Entry points (`main.go` for API server, `seed/main.go` for database seeding)
- **internal/routes/**: HTTP route registration with Gin
- **internal/controllers/**: Request handlers (auth, problems, submissions, run)
- **internal/services/**: Business logic (authentication, code execution via Piston)
- **internal/models/**: GORM data models (User, Problem, Submission, Testcase, etc.)
- **internal/database/**: PostgreSQL (GORM) and Redis connections
- **internal/judge/**: Background worker for processing code submissions
- **internal/middleware/**: JWT authentication and role-based access control

**Key Workflow**: When users submit code, it's queued in Redis (`judge_queue`) and processed asynchronously by the worker in `internal/judge/worker.go`. The worker calls the Piston service (external code execution API) to run code against test cases.

## Common Commands

### Local Development

```bash
# Start PostgreSQL, Redis, and backend via Docker
docker-compose up -d

# In another terminal, set up environment
cp .env.example .env
# Edit .env with your settings

# Run database migrations (AutoMigrate runs on startup)
# Seed admin user and sample data
go run ./cmd/seed/main.go

# Start the backend server
go run ./cmd/main.go
```

### Docker Development

```bash
# Build and start all services
docker-compose up --build

# Run in detached mode
docker-compose up -d --build

# View logs
docker-compose logs -f backend
```

### Building

```bash
# Build the server binary
go build -o innogen-backend ./cmd/main.go

# Build the seeder binary
go build -o innogen-seed ./cmd/seed/main.go

# Build for production (used in Dockerfile)
CGO_ENABLED=0 GOOS=linux go build -o innogen-backend ./cmd/main.go
```

### Database

```bash
# Seed admin user and sample problem
go run ./cmd/seed/main.go

# Database migrations are automatic on startup via GORM AutoMigrate
```

## Configuration

Required environment variables (see `.env.example`):
- `POSTGRES_HOST`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRES_PORT`
- `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`
- `JWT_SECRET`: Secret key for JWT tokens
- `PISTON_URL`: URL of Piston code execution service (e.g., `http://localhost:2000` or production URL)
- `ADMIN_PASSWORD`: Password for the default admin user
- `PORT`: Server port (default: 8080)

**Production Piston URL**: `https://excode.innogenlab.com`

## API Structure

All routes are prefixed with `/api`:
- `POST /api/auth/login`: User authentication (returns JWT)
- `GET /api/health`: Health check
- `GET /api/me`: Get current user (requires JWT)
- `GET /api/problems`: List all problems
- `GET /api/problems/:id`: Get problem details
- `POST /api/admin/problems`: Create problem (admin/teacher only)
- `POST /api/submit`: Submit code for judging (requires JWT)
- `GET /api/submit/:id`: Get submission status
- `POST /api/run`: Run code directly without test cases (requires JWT)

See `API_DOCS.md` and `swagger.yaml` for detailed API documentation.

## Key Components

### Judge Worker (`internal/judge/worker.go`)
Background worker that:
1. Polls Redis queue (`judge_queue`) for submissions
2. Retrieves submission and test cases from PostgreSQL
3. Executes code via Piston service for each test case
4. Updates submission status (Accepted, Wrong Answer, Runtime Error, etc.)

### Code Execution (`internal/services/piston.go`)
Integrates with Piston API (`/api/v2/execute`) to execute code in multiple languages. The service handles HTTP requests to the Piston server and parses responses.

### Authentication (`internal/middleware/jwt.go`, `internal/utils/jwt.go`)
JWT-based authentication middleware. Roles: `admin`, `teacher`, `student`. Admin/teacher roles can create problems.

### Models (`internal/models/`)
Main entities:
- **User**: User accounts with role-based permissions
- **Problem**: Programming problems with metadata (time/memory limits, difficulty)
- **Testcase**: Input/output pairs for problem validation
- **Submission**: User code submissions with execution results

## Development Notes

- **No test files found**: Testing should be added for controllers, services, and judge worker
- **Binary checked in**: `innogen` binary exists in root (should be in `.gitignore`)
- **Dependencies**: Uses Gin, GORM (PostgreSQL), Redis, JWT, and Piston for code execution
- **Port conflicts**: Backend runs on 8080 by default, docker-compose maps to 8081
- **CORS**: Enabled with `AllowAllOrigins = true` for development

## Production Considerations

- Backend is containerized with `Dockerfile`
- Uses external Piston service for code execution (not running locally)
- PostgreSQL and Redis run as separate containers
- Worker runs in same process as API server (started in goroutine)
- Static file serving not implemented (see `index.html` - appears to be for testing)

## Database Schema

Auto-migrated models include: User, Problem, Testcase, Submission, Subject, SubjectSession, Lesson, LessonProblem, Tag, ProblemTag. Database connection and migration happen in `internal/database/db.go`.